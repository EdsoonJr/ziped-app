package zipformat

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"ziped-app/backend/models"
)

// Reader implementa a leitura de pacotes ZIP
type Reader struct {
	r *zip.ReadCloser
}

// OpenReader abre um arquivo ZIP para leitura
func OpenReader(path string) (*Reader, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return &Reader{r: r}, nil
}

// Close fecha o arquivo
func (r *Reader) Close() error {
	return r.r.Close()
}

// GetFiles retorna a lista de arquivos para a interface
func (r *Reader) GetFiles() []models.FileInfo {
	var files []models.FileInfo
	for _, f := range r.r.File {
		tipo := "Arquivo"
		if f.FileInfo().IsDir() {
			tipo = "Pasta"
		}
		files = append(files, models.FileInfo{
			Name:           filepath.Base(f.Name),
			Path:           f.Name,
			OriginalSize:   int64(f.UncompressedSize64),
			CompressedSize: int64(f.CompressedSize64),
			Type:           tipo,
			Modified:       f.Modified.Format("02/01/2006 15:04"),
			IsDir:          f.FileInfo().IsDir(),
		})
	}
	return files
}

// ExtractFile extrai um arquivo específico
func (r *Reader) ExtractFile(internalPath string, destPath string, cb func(int)) error {
	var target *zip.File
	for _, f := range r.r.File {
		if f.Name == internalPath {
			target = f
			break
		}
	}
	if target == nil {
		return fmt.Errorf("arquivo não encontrado no zip")
	}

	if target.FileInfo().IsDir() {
		return os.MkdirAll(destPath, 0755)
	}

	src, err := target.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Garante que a pasta destino exista
	os.MkdirAll(filepath.Dir(destPath), 0755)

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if cb != nil {
		_, err = io.Copy(dest, &progressReader{r: src, cb: cb})
	} else {
		_, err = io.Copy(dest, src)
	}
	return err
}

// Writer implementa a criação de pacotes ZIP
type Writer struct {
	f *os.File
	w *zip.Writer
}

// NewWriter cria um novo pacote ZIP
func NewWriter(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &Writer{
		f: f,
		w: zip.NewWriter(f),
	}, nil
}

// AddFile adiciona um arquivo do disco para o ZIP
func (w *Writer) AddFile(diskPath string, internalPath string, cb func(int)) error {
	info, err := os.Stat(diskPath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// O ZIP usa barras normais para caminhos internos
	header.Name = filepath.ToSlash(internalPath)
	
	if !info.IsDir() {
		header.Method = zip.Deflate // Comprimir apenas arquivos
		header.Flags |= 0x8         // OBRIGATÓRIO: informa que o tamanho comprimido será escrito após o arquivo
	} else {
		if !strings.HasSuffix(header.Name, "/") {
			header.Name += "/" // Opcional no ZIP para marcar como pasta
		}
	}

	dest, err := w.w.CreateHeader(header)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	src, err := os.Open(diskPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if cb != nil {
		_, err = io.Copy(dest, &progressReader{r: src, cb: cb})
	} else {
		_, err = io.Copy(dest, src)
	}
	return err
}

// Close finaliza e fecha o ZIP
func (w *Writer) Close() error {
	var err1 error
	if w.w != nil {
		err1 = w.w.Close()
	}
	err2 := w.f.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// Helper para emitir progresso
type progressReader struct {
	r  io.Reader
	cb func(int)
}

func (p *progressReader) Read(buf []byte) (int, error) {
	n, err := p.r.Read(buf)
	if n > 0 && p.cb != nil {
		p.cb(n)
	}
	return n, err
}

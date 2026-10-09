package sevenzformat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"ziped-app/backend/models"

	"github.com/bodgit/sevenzip"
)

// Reader implementa a leitura de pacotes 7z
type Reader struct {
	r *sevenzip.ReadCloser
}

// OpenReader abre o pacote 7z
func OpenReader(path string) (*Reader, error) {
	r, err := sevenzip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return &Reader{r: r}, nil
}

// Close fecha o pacote
func (r *Reader) Close() error {
	if r.r != nil {
		return r.r.Close()
	}
	return nil
}

// GetFiles lista o conteúdo do 7z
func (r *Reader) GetFiles() []models.FileInfo {
	var files []models.FileInfo
	for _, f := range r.r.File {
		tipo := "Arquivo"
		if f.FileInfo().IsDir() {
			tipo = "Pasta"
		}

		files = append(files, models.FileInfo{
			Name:           filepath.Base(f.Name),
			Path:           filepath.ToSlash(f.Name),
			OriginalSize:   f.FileInfo().Size(),
			CompressedSize: 0, // A biblioteca não expõe tamanho comprimido por arquivo facilmente
			Type:           tipo,
			Modified:       f.FileInfo().ModTime().Format("02/01/2006 15:04"),
			IsDir:          f.FileInfo().IsDir(),
		})
	}
	return files
}

// ExtractFile encontra o arquivo correto no 7z e o extrai
func (r *Reader) ExtractFile(internalPath string, destPath string, cb func(int)) error {
	var target *sevenzip.File
	for _, f := range r.r.File {
		if filepath.ToSlash(f.Name) == internalPath {
			target = f
			break
		}
	}

	if target == nil {
		return fmt.Errorf("arquivo não encontrado: %s", internalPath)
	}

	if target.FileInfo().IsDir() {
		return os.MkdirAll(destPath, 0755)
	}

	src, err := target.Open()
	if err != nil {
		return err
	}
	defer src.Close()

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

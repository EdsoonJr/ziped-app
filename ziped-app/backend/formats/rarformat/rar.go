package rarformat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"ziped-app/backend/models"

	"github.com/nwaples/rardecode"
)

// Reader implementa a leitura de pacotes RAR
type Reader struct {
	path string
}

// OpenReader apenas verifica se o arquivo é um RAR válido
func OpenReader(path string) (*Reader, error) {
	r, err := rardecode.OpenReader(path, "")
	if err != nil {
		return nil, err
	}
	r.Close()
	return &Reader{path: path}, nil
}

// Close no caso do RAR (fechamento lazy) não faz nada no nível superior
func (r *Reader) Close() error {
	return nil
}

// GetFiles varre o pacote RAR e retorna os metadados dos arquivos
func (r *Reader) GetFiles() []models.FileInfo {
	var files []models.FileInfo
	rr, err := rardecode.OpenReader(r.path, "")
	if err != nil {
		return files
	}
	defer rr.Close()

	for {
		header, err := rr.Next()
		if err == io.EOF || err != nil {
			break
		}

		tipo := "Arquivo"
		if header.IsDir {
			tipo = "Pasta"
		}

		files = append(files, models.FileInfo{
			Name:           filepath.Base(header.Name),
			Path:           header.Name,
			OriginalSize:   header.UnPackedSize,
			CompressedSize: header.PackedSize,
			Type:           tipo,
			Modified:       header.ModificationTime.Format("02/01/2006 15:04"),
			IsDir:          header.IsDir,
		})
	}
	return files
}

// ExtractFile varre o RAR até encontrar o arquivo correto e o extrai para o destino
func (r *Reader) ExtractFile(internalPath string, destPath string, cb func(int)) error {
	rr, err := rardecode.OpenReader(r.path, "")
	if err != nil {
		return err
	}
	defer rr.Close()

	for {
		header, err := rr.Next()
		if err == io.EOF {
			return fmt.Errorf("arquivo não encontrado no RAR: %s", internalPath)
		}
		if err != nil {
			return err
		}

		if header.Name == internalPath {
			if header.IsDir {
				return os.MkdirAll(destPath, 0755)
			}

			// Garante que a pasta pai exista
			os.MkdirAll(filepath.Dir(destPath), 0755)
			dest, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer dest.Close()

			if cb != nil {
				_, err = io.Copy(dest, &progressReader{r: rr, cb: cb})
			} else {
				_, err = io.Copy(dest, rr)
			}
			return err
		}
	}
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

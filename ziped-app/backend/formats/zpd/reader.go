package zpd

import (
	"compress/flate"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"time"
	"ziped-app/backend/models"
)

// Reader é responsável por abrir um arquivo .zpd e extrair seus conteúdos
type Reader struct {
	f       *os.File
	Entries []FileEntry
}

// OpenReader abre um arquivo, valida os Magic Bytes e extrai a tabela de índices (File Table)
func OpenReader(filename string) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if err := ReadHeader(f); err != nil {
		f.Close()
		return nil, err
	}

	// O último dado do arquivo (8 bytes) é o offset da Tabela de Índices
	if _, err := f.Seek(-8, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}

	var tableOffset uint64
	if err := binary.Read(f, binary.LittleEndian, &tableOffset); err != nil {
		f.Close()
		return nil, err
	}

	// Vai até a tabela de índices
	if _, err := f.Seek(int64(tableOffset), io.SeekStart); err != nil {
		f.Close()
		return nil, err
	}

	var count uint32
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		f.Close()
		return nil, err
	}

	entries := make([]FileEntry, count)
	for i := uint32(0); i < count; i++ {
		var entry FileEntry
		if err := entry.Decode(f); err != nil {
			f.Close()
			return nil, err
		}
		entries[i] = entry
	}

	return &Reader{
		f:       f,
		Entries: entries,
	}, nil
}

// GetFiles retorna as entradas formatadas para a interface do frontend
func (r *Reader) GetFiles() []models.FileInfo {
	var files []models.FileInfo
	for _, entry := range r.Entries {
		t := time.Unix(entry.ModTime, 0).Format("02/01/2006 15:04")
		tipo := "Arquivo"
		if entry.IsDir {
			tipo = "Pasta"
		}
		files = append(files, models.FileInfo{
			Name:           filepath.Base(entry.Name),
			Path:           entry.Name,
			OriginalSize:   int64(entry.OriginalSize),
			CompressedSize: int64(entry.CompressedSize),
			Type:           tipo,
			Modified:       t,
			IsDir:          entry.IsDir,
		})
	}
	return files
}

// ExtractFile descompacta uma entrada específica para o destino usando streams
func (r *Reader) ExtractFile(internalPath string, destPath string, progressCallback func(bytesProcessed int)) error {
	var targetEntry *FileEntry
	for _, e := range r.Entries {
		if e.Name == internalPath {
			targetEntry = &e
			break
		}
	}
	if targetEntry == nil {
		return fmt.Errorf("arquivo %s não encontrado no zpd", internalPath)
	}

	entry := *targetEntry

	if entry.IsDir {
		return os.MkdirAll(destPath, 0755)
	}

	// Vai até onde começam os bytes comprimidos deste arquivo
	if _, err := r.f.Seek(int64(entry.DataOffset), io.SeekStart); err != nil {
		return err
	}

	// Limita a leitura ao tamanho comprimido exato deste arquivo
	limitReader := io.LimitReader(r.f, int64(entry.CompressedSize))
	flateReader := flate.NewReader(limitReader)
	defer flateReader.Close()

	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()

	hasher := crc32.NewIEEE()
	multiWriter := io.MultiWriter(destFile, hasher)

	var writer io.Writer = multiWriter
	if progressCallback != nil {
		writer = &progressWriter{w: multiWriter, cb: progressCallback}
	}

	// Extração eficiente processada em chunks
	if _, err := io.Copy(writer, flateReader); err != nil {
		return err
	}

	// Validação de integridade do arquivo
	if hasher.Sum32() != entry.Checksum {
		return errors.New("falha de integridade: CRC32 não confere. O arquivo pode estar corrompido")
	}

	return nil
}

// Close fecha o arquivo .zpd
func (r *Reader) Close() error {
	return r.f.Close()
}

type progressWriter struct {
	w  io.Writer
	cb func(int)
}

func (p *progressWriter) Write(buf []byte) (int, error) {
	n, err := p.w.Write(buf)
	if n > 0 && p.cb != nil {
		p.cb(n)
	}
	return n, err
}

package zpd

import (
	"compress/flate"
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
)

// Writer é responsável por criar o arquivo .zpd e adicionar arquivos nele sob demanda
type Writer struct {
	f          *os.File
	offset     uint64
	entries    []FileEntry
	isClosed   bool
}

// NewWriter cria ou sobreescreve um arquivo .zpd e prepara para escrever
func NewWriter(filename string) (*Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	if err := WriteHeader(f); err != nil {
		f.Close()
		return nil, err
	}

	return &Writer{
		f:      f,
		offset: uint64(len(MagicBytes) + 2), // 4 bytes do magic + 2 bytes da versão
	}, nil
}

// AddFile adiciona um arquivo ao compactador processando em streams (pedaços)
// Não carrega arquivos grandes na RAM.
func (w *Writer) AddFile(filePath string, internalPath string, progressCallback func(bytesProcessed int)) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	entry := FileEntry{
		Name:         internalPath,
		IsDir:        info.IsDir(),
		OriginalSize: uint64(info.Size()),
		ModTime:      info.ModTime().Unix(),
		DataOffset:   w.offset,
	}

	// Se for diretório, não precisamos comprimir nada, apenas salvar metadados
	if entry.IsDir {
		w.entries = append(w.entries, entry)
		return nil
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Prepara compressor e checksum
	hasher := crc32.NewIEEE()
	flateWriter, err := flate.NewWriter(w.f, flate.BestCompression)
	if err != nil {
		return err
	}

	// Lê do arquivo, manda para o hasher e depois para o compressor via TeeReader
	teeReader := io.TeeReader(srcFile, hasher)

	var reader io.Reader = teeReader
	if progressCallback != nil {
		reader = &progressReader{r: teeReader, cb: progressCallback}
	}

	// io.Copy lida com os chunks automaticamente, usando pouca memória
	copied, err := io.Copy(flateWriter, reader)
	if err != nil {
		flateWriter.Close()
		return err
	}

	if err := flateWriter.Close(); err != nil {
		return err
	}

	newOffset, err := w.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	entry.CompressedSize = uint64(newOffset) - w.offset
	entry.Checksum = hasher.Sum32()
	entry.OriginalSize = uint64(copied)
	w.offset = uint64(newOffset)

	w.entries = append(w.entries, entry)
	return nil
}

// Close finaliza o pacote escrevendo a tabela de índices no final do arquivo
func (w *Writer) Close() error {
	if w.isClosed {
		return nil
	}
	w.isClosed = true

	tableOffset := w.offset

	// Escreve a quantidade de arquivos (uint32)
	if err := binary.Write(w.f, binary.LittleEndian, uint32(len(w.entries))); err != nil {
		return err
	}

	// Escreve cada entrada do índice
	for _, entry := range w.entries {
		if err := entry.Encode(w.f); err != nil {
			return err
		}
	}

	// Escreve a posição onde a tabela começou (nos últimos 8 bytes do arquivo)
	if err := binary.Write(w.f, binary.LittleEndian, tableOffset); err != nil {
		return err
	}

	return w.f.Close()
}

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

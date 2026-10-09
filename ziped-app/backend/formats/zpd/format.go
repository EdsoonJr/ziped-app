package zpd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	// MagicBytes identifica que este é um arquivo do nosso compactador
	MagicBytes = []byte{'Z', 'P', 'D', 0x01}
	// FormatVersion é a versão atual do nosso formato
	FormatVersion = uint16(1)

	ErrInvalidFormat = errors.New("arquivo corrompido ou formato inválido: magic bytes não conferem")
	ErrUnsupported   = errors.New("versão do formato não suportada")
)

// FileEntry representa as informações de um único arquivo/pasta dentro do pacote
type FileEntry struct {
	Name           string // Suporta Unicode nativamente em Go (UTF-8)
	IsDir          bool
	OriginalSize   uint64
	CompressedSize uint64
	ModTime        int64  // Unix timestamp
	DataOffset     uint64 // Posição (byte) onde os dados comprimidos deste arquivo começam
	Checksum       uint32 // CRC32 dos dados comprimidos para validar corrupção
}

// Encode escreve a entrada do arquivo em um formato binário
func (e *FileEntry) Encode(w io.Writer) error {
	nameBytes := []byte(e.Name)
	if err := binary.Write(w, binary.LittleEndian, uint16(len(nameBytes))); err != nil {
		return err
	}
	if _, err := w.Write(nameBytes); err != nil {
		return err
	}
	
	isDirByte := byte(0)
	if e.IsDir {
		isDirByte = 1
	}
	if err := binary.Write(w, binary.LittleEndian, isDirByte); err != nil {
		return err
	}
	
	if err := binary.Write(w, binary.LittleEndian, e.OriginalSize); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, e.CompressedSize); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, e.ModTime); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, e.DataOffset); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, e.Checksum); err != nil {
		return err
	}
	
	return nil
}

// Decode lê a entrada do arquivo a partir de um formato binário
func (e *FileEntry) Decode(r io.Reader) error {
	var nameLen uint16
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return err
	}
	
	nameBytes := make([]byte, nameLen)
	if _, err := io.ReadFull(r, nameBytes); err != nil {
		return err
	}
	e.Name = string(nameBytes)
	
	var isDirByte byte
	if err := binary.Read(r, binary.LittleEndian, &isDirByte); err != nil {
		return err
	}
	e.IsDir = isDirByte == 1
	
	if err := binary.Read(r, binary.LittleEndian, &e.OriginalSize); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.CompressedSize); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.ModTime); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.DataOffset); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.Checksum); err != nil {
		return err
	}
	
	return nil
}

// WriteHeader escreve a assinatura e a versão inicial do arquivo
func WriteHeader(w io.Writer) error {
	if _, err := w.Write(MagicBytes); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, FormatVersion)
}

// ReadHeader lê a assinatura e valida a versão do formato
func ReadHeader(r io.Reader) error {
	magic := make([]byte, 4)
	if _, err := io.ReadFull(r, magic); err != nil {
		return err
	}
	if !bytes.Equal(magic, MagicBytes) {
		return ErrInvalidFormat
	}

	var version uint16
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != FormatVersion {
		return ErrUnsupported
	}
	return nil
}

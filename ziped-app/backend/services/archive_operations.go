package services

import (
	"fmt"
	"os"
	"path/filepath"
	"ziped-app/backend/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (s *ArchiveService) OpenArchive(path string) ([]models.FileInfo, error) {
	reader, err := getReader(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %v", err)
	}
	defer reader.Close()

	return reader.GetFiles(), nil
}

func (s *ArchiveService) CreateNewArchive(destPath string, filePaths []string) (models.ArchiveResult, error) {
	writer, err := getWriter(destPath)
	if err != nil {
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}
	defer writer.Close()

	var totalSize int64
	for _, p := range filePaths {
		info, err := os.Stat(p)
		if err == nil && !info.IsDir() {
			totalSize += info.Size()
		}
	}

	var processedSize int64
	for _, p := range filePaths {
		internalPath := filepath.Base(p)
		
		err := writer.AddFile(p, internalPath, func(bytesProcessed int) {
			processedSize += int64(bytesProcessed)
			if totalSize > 0 {
				percent := int(float64(processedSize) / float64(totalSize) * 100)
				runtime.EventsEmit(s.ctx, "progress", models.ProgressEvent{
					Percent:     percent,
					CurrentFile: internalPath,
					Status:      "Compactando...",
				})
			}
		})
		
		if err != nil {
			return models.ArchiveResult{Success: false, Error: "Erro ao adicionar " + internalPath + ": " + formatErrorMsg(err)}, err
		}
	}

	return models.ArchiveResult{Success: true, Message: "Arquivo criado com sucesso!"}, nil
}

func (s *ArchiveService) AddFiles(archivePath string, files []string) (models.ArchiveResult, error) {
	if len(files) == 0 {
		return models.ArchiveResult{Success: false, Error: "Nenhum arquivo selecionado"}, nil
	}

	tempPath := archivePath + ".tmp"
	writer, err := getWriter(tempPath)
	if err != nil {
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}

	reader, err := getReader(archivePath)
	if err == nil {
		tempDir, err := os.MkdirTemp("", "archive_add_*")
		if err != nil {
			writer.Close()
			reader.Close()
			return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
		}
		defer os.RemoveAll(tempDir)

		for _, entry := range reader.GetFiles() {
			fullDestPath := filepath.Join(tempDir, entry.Path)
			os.MkdirAll(filepath.Dir(fullDestPath), 0755)
			if err := reader.ExtractFile(entry.Path, fullDestPath, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir arquivo"}, err
			}
			if err := writer.AddFile(fullDestPath, entry.Path, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir pacote"}, err
			}
		}
		reader.Close()
	}

	for _, p := range files {
		internalPath := filepath.Base(p)
		if err := writer.AddFile(p, internalPath, nil); err != nil {
			writer.Close()
			return models.ArchiveResult{Success: false, Error: "Erro ao adicionar novo arquivo: " + formatErrorMsg(err)}, err
		}
	}

	if err := writer.Close(); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao finalizar novo pacote: " + formatErrorMsg(err)}, err
	}
	os.Remove(archivePath)
	if err := os.Rename(tempPath, archivePath); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao salvar pacote atualizado"}, err
	}

	return models.ArchiveResult{Success: true, Message: "Arquivos adicionados com sucesso!"}, nil
}

func (s *ArchiveService) AddFolder(archivePath string, folder string) (models.ArchiveResult, error) {
	if folder == "" {
		return models.ArchiveResult{Success: false, Error: "Nenhuma pasta selecionada"}, nil
	}

	tempPath := archivePath + ".tmp"
	writer, err := getWriter(tempPath)
	if err != nil {
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}

	reader, err := getReader(archivePath)
	if err == nil {
		tempDir, err := os.MkdirTemp("", "archive_add_*")
		if err != nil {
			writer.Close()
			reader.Close()
			return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
		}
		defer os.RemoveAll(tempDir)

		for _, entry := range reader.GetFiles() {
			fullDestPath := filepath.Join(tempDir, entry.Path)
			os.MkdirAll(filepath.Dir(fullDestPath), 0755)
			if err := reader.ExtractFile(entry.Path, fullDestPath, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir arquivo"}, err
			}
			if err := writer.AddFile(fullDestPath, entry.Path, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir pacote"}, err
			}
		}
		reader.Close()
	}

	baseFolderName := filepath.Base(folder)
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		
		relPath, _ := filepath.Rel(folder, path)
		internalPath := filepath.Join(baseFolderName, relPath)
		internalPath = filepath.ToSlash(internalPath)
		
		return writer.AddFile(path, internalPath, nil)
	})

	if err != nil {
		writer.Close()
		return models.ArchiveResult{Success: false, Error: "Erro ao adicionar pasta: " + formatErrorMsg(err)}, err
	}

	if err := writer.Close(); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao finalizar novo pacote: " + formatErrorMsg(err)}, err
	}
	os.Remove(archivePath)
	if err := os.Rename(tempPath, archivePath); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao salvar pacote atualizado"}, err
	}

	return models.ArchiveResult{Success: true, Message: "Pasta adicionada com sucesso!"}, nil
}

func (s *ArchiveService) ExtractFiles(archivePath string, destPath string, files []string) (models.ArchiveResult, error) {
	reader, err := getReader(archivePath)
	if err != nil {
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}
	defer reader.Close()

	extractAll := len(files) == 0
	filesMap := make(map[string]bool)
	for _, f := range files {
		filesMap[f] = true
	}

	archiveFiles := reader.GetFiles()
	var totalSize uint64
	for _, entry := range archiveFiles {
		if extractAll || filesMap[entry.Path] {
			totalSize += uint64(entry.OriginalSize)
		}
	}

	var processedSize uint64
	for _, entry := range archiveFiles {
		if extractAll || filesMap[entry.Path] {
			fullDestPath := filepath.Join(destPath, entry.Path)
			os.MkdirAll(filepath.Dir(fullDestPath), 0755)

			err := reader.ExtractFile(entry.Path, fullDestPath, func(n int) {
				processedSize += uint64(n)
				if totalSize > 0 {
					percent := int(float64(processedSize) / float64(totalSize) * 100)
					runtime.EventsEmit(s.ctx, "progress", models.ProgressEvent{
						Percent:     percent,
						CurrentFile: entry.Path,
						Status:      "Extraindo...",
					})
				}
			})
			if err != nil {
				return models.ArchiveResult{Success: false, Error: "Erro ao extrair " + entry.Path + ": " + formatErrorMsg(err)}, err
			}
		}
	}

	return models.ArchiveResult{Success: true, Message: "Arquivos extraídos com sucesso!"}, nil
}

func (s *ArchiveService) DeleteEntries(archivePath string, entryNames []string) (models.ArchiveResult, error) {
	if len(entryNames) == 0 {
		return models.ArchiveResult{Success: false, Error: "Nenhum arquivo selecionado"}, nil
	}

	reader, err := getReader(archivePath)
	if err != nil {
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}

	deleteMap := make(map[string]bool)
	for _, f := range entryNames {
		deleteMap[f] = true
	}

	tempPath := archivePath + ".tmp"
	writer, err := getWriter(tempPath)
	if err != nil {
		reader.Close()
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}

	tempDir, err := os.MkdirTemp("", "archive_delete_*")
	if err != nil {
		writer.Close()
		reader.Close()
		return models.ArchiveResult{Success: false, Error: formatErrorMsg(err)}, err
	}
	defer os.RemoveAll(tempDir)

	for _, entry := range reader.GetFiles() {
		if !deleteMap[entry.Path] {
			fullDestPath := filepath.Join(tempDir, entry.Path)
			os.MkdirAll(filepath.Dir(fullDestPath), 0755)
			if err := reader.ExtractFile(entry.Path, fullDestPath, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir arquivo"}, err
			}
			if err := writer.AddFile(fullDestPath, entry.Path, nil); err != nil {
				writer.Close()
				reader.Close()
				return models.ArchiveResult{Success: false, Error: "Erro ao reconstruir pacote"}, err
			}
		}
	}

	if err := writer.Close(); err != nil {
		reader.Close()
		return models.ArchiveResult{Success: false, Error: "Erro ao finalizar novo pacote: " + formatErrorMsg(err)}, err
	}
	reader.Close()

	if err := os.Remove(archivePath); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao deletar original"}, err
	}
	if err := os.Rename(tempPath, archivePath); err != nil {
		return models.ArchiveResult{Success: false, Error: "Erro ao renomear novo arquivo"}, err
	}

	return models.ArchiveResult{Success: true, Message: "Itens excluídos com sucesso!"}, nil
}

func (s *ArchiveService) TestArchive(archivePath string) (models.ArchiveResult, error) {
	reader, err := getReader(archivePath)
	if err != nil {
		return models.ArchiveResult{Success: false, Message: "Arquivo corrompido ou inválido: " + formatErrorMsg(err)}, nil
	}
	defer reader.Close()
	return models.ArchiveResult{Success: true, Message: "Arquivo está intacto!"}, nil
}

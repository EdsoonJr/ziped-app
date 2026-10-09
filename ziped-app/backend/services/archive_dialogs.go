package services

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (s *ArchiveService) SelectArchiveFile() (string, error) {
	return runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Abrir pacote",
		Filters: []runtime.FileFilter{{DisplayName: "Pacotes Suportados", Pattern: "*.zpd;*.zip;*.garc;*.rar;*.7z"}},
	})
}

func (s *ArchiveService) SelectSaveFile() (string, error) {
	return runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		Title:           "Criar pacote",
		DefaultFilename: "novo_pacote.zip",
		Filters: []runtime.FileFilter{
			{DisplayName: "Arquivo ZIP Padrão", Pattern: "*.zip"},
			{DisplayName: "Formato ZPD", Pattern: "*.zpd"},
		},
	})
}

func (s *ArchiveService) SelectExtractFolder() (string, error) {
	return runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Selecione a pasta de destino para extração",
	})
}

func (s *ArchiveService) SelectFilesToCompact() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Selecione os arquivos",
	})
}

func (s *ArchiveService) SelectFolderToCompact() (string, error) {
	return runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Selecione a pasta",
	})
}

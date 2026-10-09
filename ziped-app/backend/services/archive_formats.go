package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"ziped-app/backend/formats/rarformat"
	"ziped-app/backend/formats/sevenzformat"
	"ziped-app/backend/formats/zipformat"
	"ziped-app/backend/formats/zpd"
	"ziped-app/backend/models"
)

type ArchiveReader interface {
	GetFiles() []models.FileInfo
	ExtractFile(internalPath string, destPath string, cb func(int)) error
	Close() error
}

type ArchiveWriter interface {
	AddFile(diskPath string, internalPath string, cb func(int)) error
	Close() error
}

func getReader(path string) (ArchiveReader, error) {
	cleanPath := strings.TrimSuffix(path, ".tmp")
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext == ".zip" {
		return zipformat.OpenReader(path)
	}
	if ext == ".rar" {
		return rarformat.OpenReader(path)
	}
	if ext == ".7z" {
		return sevenzformat.OpenReader(path)
	}
	return zpd.OpenReader(path)
}

func getWriter(path string) (ArchiveWriter, error) {
	cleanPath := strings.TrimSuffix(path, ".tmp")
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext == ".rar" {
		return nil, fmt.Errorf("o formato RAR é suportado apenas para leitura")
	}
	if ext == ".7z" {
		return nil, fmt.Errorf("o formato 7z é suportado apenas para leitura")
	}
	if ext == ".zip" {
		return zipformat.NewWriter(path)
	}
	return zpd.NewWriter(path)
}

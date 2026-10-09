package services

import (
	"os"
	"strings"
)

func formatErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()

	if os.IsNotExist(err) || strings.Contains(msg, "no such file or directory") || strings.Contains(msg, "não foi possível encontrar") {
		return "Arquivo ou diretório inexistente."
	}
	if os.IsPermission(err) || strings.Contains(msg, "permission denied") || strings.Contains(msg, "Acesso negado") || strings.Contains(msg, "Access is denied") {
		return "Acesso negado: falta de permissão ou arquivo bloqueado por outro processo."
	}
	if strings.Contains(msg, "being used by another process") || strings.Contains(msg, "está sendo usado por outro processo") {
		return "Arquivo bloqueado: está sendo usado por outro processo."
	}
	if strings.Contains(msg, "no space left on device") || strings.Contains(msg, "não há espaço") || strings.Contains(msg, "disk full") {
		return "Erro: disco sem espaço suficiente para concluir a operação."
	}
	if strings.Contains(msg, "EOF") || strings.Contains(msg, "corrupt") || strings.Contains(msg, "invalid format") || strings.Contains(msg, "not a supported archive") || strings.Contains(msg, "unexpected EOF") {
		return "Arquivo corrompido ou formato desconhecido."
	}
	

	if strings.Contains(msg, "suportado apenas para leitura") {
		return msg
	}

	return "Ocorreu um erro: " + msg
}

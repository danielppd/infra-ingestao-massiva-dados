package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ponderada/back/queue"
)

// TelemetryHandler é o nosso controlador para as rotas de telemetria
type TelemetryHandler struct {
	filaDeMensagens *queue.RabbitMQ
}

// Construtorzinho que o professor mostrou na aula de Padrões
func NewTelemetryHandler(mq *queue.RabbitMQ) *TelemetryHandler {
	return &TelemetryHandler{filaDeMensagens: mq}
}

// PostTelemetry é a função que o router do Gin vai chamar quando bater o POST
func (h *TelemetryHandler) PostTelemetry(contexto *gin.Context) {
	// Vou usar um map pro JSON livre pra não precisar criar structs aqui no backend atoa
	var corpoDaRequisicao map[string]interface{}
	
	// Tenta ler o JSON do corpo. Se for um JSON inválido, já dá erro HTTP 400
	if erro := contexto.ShouldBindJSON(&corpoDaRequisicao); erro != nil {
		contexto.JSON(http.StatusBadRequest, gin.H{"erro": "Opa, o formato do JSON tá inválido!"})
		return
	}

	// Pede pra classe de fila publicar a nossa mensagem
	if erro := h.filaDeMensagens.Publish(corpoDaRequisicao); erro != nil {
		contexto.JSON(http.StatusInternalServerError, gin.H{"erro": "Eita, falhou na hora de mandar a mensagem pra fila do rabbit"})
		return
	}

	// Devolvemos o Status 202 (Accepted) que configuramos no teste de carga (k6)
	contexto.Status(http.StatusAccepted)
}

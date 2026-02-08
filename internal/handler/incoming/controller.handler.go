package incoming

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/rabbitmq"
	incomingService "pannypal/internal/service/incoming"
)

type Handler struct {
	ctx      context.Context
	rabbitmq *rabbitmq.ConnectionManager
	incoming incomingService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	WebhookEventBaileys(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, incoming incomingService.IService) IHandler {
	return &Handler{
		ctx:      ctx,
		rabbitmq: rabbitmq,
		incoming: incoming,
	}
}

// WebhookEventBaileys godoc
// @Summary Webhook Event Baileys
// @Description Handles webhook events from Baileys WhatsApp service
// @Tags Webhook
// @Accept json
// @Produce json
// @Param webhook body interface{} true "Baileys webhook event payload"
// @Success 200 {object} types.Response "Webhook event processed successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /incoming/baileys [post]
func (h *Handler) WebhookEventBaileys(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))

	var payload interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.incoming.HandleWebhookEventBaileys(payload))
}

package webhook

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/rabbitmq"
	webhookService "pannypal/internal/service/webhook"
)

type Handler struct {
	ctx            context.Context
	rabbitmq       *rabbitmq.ConnectionManager
	webhookService webhookService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	WebhookEventWaha(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, webhookService webhookService.IService) IHandler {
	return &Handler{
		ctx:            ctx,
		rabbitmq:       rabbitmq,
		webhookService: webhookService,
	}
}

// WebhookEventWaha godoc
// @Summary Webhook Event Waha
// @Description Handles webhook events from Waha
// @Tags Webhook
// @Accept json
// @Produce json
// @Param webhook body interface{} true "Webhook event payload"
// @Success 201 {object} types.Response "Webhook event processed successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /webhook/waha [post]
func (h *Handler) WebhookEventWaha(c *gin.Context) {
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

	send(h.webhookService.HandleWebhookEventWaha(payload))
}

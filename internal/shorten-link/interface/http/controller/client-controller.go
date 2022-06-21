package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"trueid-shorten-link/internal/shorten-link/entity"
	"trueid-shorten-link/internal/shorten-link/service"
)

type ClientController struct {
	Controller
	client *service.ClientService
}

func NewClientController(service *service.ClientService) *ClientController {
	return &ClientController{
		client: service,
	}
}

func (c *ClientController) UpdateClient(ctx *fiber.Ctx) error {
	var body entity.UpdateClientRequest
	if err := json.Unmarshal(ctx.Body(), &body); err != nil {
		return c.Failure(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
	}

	err := c.client.UpdateClient(&body)
	if err != nil {
		return c.Failure(ctx, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
	}

	return c.Success(ctx, http.StatusOK, http.StatusOK, "success", nil)
}

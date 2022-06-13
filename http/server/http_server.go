package server

import (
	"encoding/json"
	"fmt"
	"go-forex/model"
	"go-forex/usecase"

	"github.com/gofiber/fiber/v2"
)

type ServerHttp struct {
	getRateUsecase usecase.GetRateUsecase
}

func NewServerHttp(getRateUsecase usecase.GetRateUsecase) ServerHttp {
	return ServerHttp{getRateUsecase}
}

func (h *ServerHttp) Init() {
	fiberApp := fiber.New()
	fiberApp.Use(DefaultMiddleware())

	fiberApp.Get("v1/get-rate", h.GetRate)

	fiberApp.Post("v1/book-rate", h.BookRate)

	fmt.Printf("HTTP listen on port :80 \n")
	fiberApp.Listen(fmt.Sprintf(":80"))
}

func (h *ServerHttp) GetRate(ctx *fiber.Ctx) error {
	traceId := fmt.Sprint(ctx.Locals(model.HEADER_TRACE_ID))
	originCurrency := ctx.Query("origin")
	destinationCurrency := ctx.Query("destination")
	response := h.getRateUsecase.GetRate(traceId, originCurrency, destinationCurrency, "")
	if !response.Success {
		return ctx.Status(500).JSON(response)
	}
	return ctx.Status(200).JSON(response)
}

func (h *ServerHttp) BookRate(ctx *fiber.Ctx) error {
	traceId := fmt.Sprint(ctx.Locals(model.HEADER_TRACE_ID))
	tokenAccess := fmt.Sprint(ctx.Request().Header.Peek("token_access"))
	var request model.RequestBook
	json.Unmarshal(ctx.Request().Body(), &request)
	response := h.getRateUsecase.LockOrBookRate(traceId, tokenAccess, request)
	if !response.Success {
		return ctx.Status(500).JSON(response)
	}
	return ctx.Status(200).JSON(response)
}

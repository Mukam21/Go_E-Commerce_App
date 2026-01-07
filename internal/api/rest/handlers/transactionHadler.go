package handlers

import (
	"net/http"

	"github.com/Mukam21/Go_E-Commerce_App/internal/api/rest"
	"github.com/Mukam21/Go_E-Commerce_App/internal/helper"
	"github.com/Mukam21/Go_E-Commerce_App/internal/repository"
	"github.com/Mukam21/Go_E-Commerce_App/internal/service"
	"github.com/Mukam21/Go_E-Commerce_App/pkg/payment"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc           service.TransactionService
	paymentClient payment.PaymentClient
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(as *rest.RestHandler) {

	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)

	handler := TransactionHandler{
		svc:           svc,
		paymentClient: as.Pc,
	}

	secRouter := app.Group("/", as.Auth.Authorize)
	secRouter.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("orders", handler.GetOrders)
	sellerRoute.Get("orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	// 1. call user service get cart data to aggregate the total amount and collect payment

	// 2. check if payment session active or Create a new payment session
	sessionResult, err := h.paymentClient.CreatePayment(2, 123, 456)

	// 3. Store payment session in db to create and validate order
	if err != nil {
		return ctx.Status(400).JSON(err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":     "login",
		"result":      sessionResult,
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

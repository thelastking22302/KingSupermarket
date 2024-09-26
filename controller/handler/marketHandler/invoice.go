package markethandler

import (
	"net/http"
	"time"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/repository"
	repomarketiml "github.com/KingSupermarket/repository/repo_market_iml"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateInvoiceHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dataInvoice marketmodels.Invoice
		if err := c.BodyParser(&dataInvoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind invoice",
			})
		}

		validate := validator.New()
		if err := validate.Struct(dataInvoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not validate invoice",
			})
		}

		idInvoice, err := uuid.NewUUID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not create UUID for invoice",
			})
		}

		id := idInvoice.String()
		now := time.Now().UTC()
		status := "PENDING"
		if dataInvoice.Payment_status == nil {
			dataInvoice.Payment_status = &status
		}

		invoice := &marketmodels.Invoice{
			ID:               primitive.NewObjectID(),
			Invoice_id:       id,
			Order_id:         dataInvoice.Order_id,
			Payment_method:   dataInvoice.Payment_method,
			Payment_status:   dataInvoice.Payment_status,
			Payment_due_date: dataInvoice.Payment_due_date,
			Created_at:       now,
			Updated_at:       now,
		}

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateInvoice(c.Context(), invoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not business invoice",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "added invoice successfully!",
		})
	}
}
func HandlerGetInvoice(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idInvoice := c.Params("invoice_id")

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetInvoice(c.Context(), idInvoice)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "getInvoice failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"invoice": data,
		})
	}
}
func HandlerUpdateInvoices(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idInvoice := c.Params("invoice_id")
		var data marketmodels.Invoice
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bind data failed",
			})
		}

		data.Updated_at = time.Now()

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateInvoice(c.Context(), idInvoice, &data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "data update failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}
func HandlerDeleteInvoice(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idInvoice := c.Params("invoice_id")
		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteInvoice(c.Context(), idInvoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "delete invoice failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "Deleted successfully!",
		})
	}
}
func GetListInvoice(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagging common.Pagging
		if err := c.BodyParser(&pagging); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind paging",
			})
		}

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		listInvoice, err := bus.NewGetListInvoice(c.Context(), bson.M{}, &pagging)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not get list business invoice",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"listInvoice": listInvoice,
			"total":       pagging.Total,
		})
	}
}

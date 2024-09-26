package markethandler

import (
	"net/http"
	"time"

	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/repository"
	repomarketiml "github.com/KingSupermarket/repository/repo_market_iml"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrderHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData, ok := c.Locals("userId").(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		var dataOrder marketmodels.Order
		if err := c.BodyParser(&dataOrder); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind order",
			})
		}

		validate := validator.New()
		if err := validate.Struct(dataOrder); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not validate order",
			})
		}

		idOrder, err := uuid.NewUUID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not create UUID for order",
			})
		}

		id := idOrder.String()
		now := time.Now().UTC()

		order := &marketmodels.Order{
			Id:           primitive.NewObjectID(),
			Order_Id:     &id,
			User_Id:      &tokenData,
			Address:      dataOrder.Address,
			Phone_Number: dataOrder.Phone_Number,
			Total_amount: dataOrder.Total_amount,
			Status:       dataOrder.Status,
			Notes:        dataOrder.Notes,
			Order_day:    &now,
			Created_at:   &now,
			Updated_at:   &now,
		}

		bus := repository.NewOrderRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateOrder(c.Context(), order, tokenData); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not business order",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "added order successfully!",
		})
	}
}
func HandlerGetOrder(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData, ok := c.Locals("userId").(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		idOrder := c.Params("order_id")

		bus := repository.NewOrderRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetOrder(c.Context(), idOrder, tokenData)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "getOrder failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"order": data,
		})
	}
}
func HandlerUpdateOrder(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData, ok := c.Locals("userId").(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		idOrder := c.Params("order_id")
		var data marketmodels.Order
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bind data failed",
			})
		}

		now := time.Now().UTC()
		data.Updated_at = &now
		data.Order_day = &now

		bus := repository.NewOrderRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateOrder(c.Context(), idOrder, &data, tokenData); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "data update failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}
func HandlerDeleteOrder(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData, ok := c.Locals("userId").(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		idOrder := c.Params("order_id")
		bus := repository.NewOrderRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteOrder(c.Context(), idOrder, tokenData); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "delete order failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "Deleted successfully!",
		})
	}
}

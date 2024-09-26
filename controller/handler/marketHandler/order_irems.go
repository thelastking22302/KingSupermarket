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

func CreateOrderItemsHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dataOrderItems marketmodels.OrderItems
		if err := c.BodyParser(&dataOrderItems); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind order items",
			})
		}

		validate := validator.New()
		if err := validate.Struct(dataOrderItems); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not validate order items",
			})
		}

		idOrderItems, err := uuid.NewUUID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not create UUID for order items",
			})
		}

		id := idOrderItems.String()
		now := time.Now().UTC()

		orderItems := &marketmodels.OrderItems{
			Id:            primitive.NewObjectID(),
			Order_Item_Id: &id,
			Order_Id:      dataOrderItems.Order_Id,
			Product_id:    dataOrderItems.Product_id,
			Quantity:      dataOrderItems.Quantity,
			Price:         dataOrderItems.Price,
			Created_at:    &now,
			Updated_at:    &now,
		}

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateOrderItems(c.Context(), orderItems); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not business order items",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "added order items successfully!",
		})
	}
}
func HandlerGetOrderItems(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrderItems := c.Params("order_item_id")

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetOrderItems(c.Context(), idOrderItems)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "getOrderItems failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"order item": data,
		})
	}
}
func HandlerUpdateOrderItems(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrderItems := c.Params("order_item_id")
		var data marketmodels.OrderItems
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bind data failed",
			})
		}

		now := time.Now().UTC()
		data.Updated_at = &now

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateOrderItems(c.Context(), idOrderItems, &data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "data update failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}
func HandlerDeleteOrderItems(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrderItems := c.Params("order_item_id")
		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteOrderItems(c.Context(), idOrderItems); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "delete order items failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "Deleted successfully!",
		})
	}
}
func GetListOrderItems(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagging common.Pagging
		if err := c.BodyParser(&pagging); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind paging",
			})
		}

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		listOrderItems, err := bus.NewGetListOrderItems(c.Context(), bson.M{}, &pagging)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not get list business order items",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"listOrderItems": listOrderItems,
			"total":          pagging.Total,
		})
	}
}

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCategoryHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dataCategory marketmodels.Category
		if err := c.BodyParser(&dataCategory); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind category",
			})
		}

		validate := validator.New()
		if err := validate.Struct(dataCategory); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not validate category",
			})
		}

		idCategory, err := uuid.NewUUID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not create UUID for category",
			})
		}

		id := idCategory.String()
		now := time.Now().UTC()

		category := &marketmodels.Category{
			Id:          primitive.NewObjectID(),
			Category_id: &id,
			Name:        dataCategory.Name,
			Created_at:  &now,
			Updated_at:  &now,
		}

		bus := repository.NewCategoryRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateCategory(c.Context(), category); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not business category",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "added category successfully!",
		})
	}
}
func HandlerGetCategory(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idCategory := c.Params("category_id")

		bus := repository.NewCategoryRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetCategory(c.Context(), idCategory)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "getCategory failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"category": data,
		})
	}
}
func HandlerUpdateCategory(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idCategory := c.Params("category_id")
		var data marketmodels.Category
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bind data failed",
			})
		}

		now := time.Now().UTC()
		data.Updated_at = &now

		bus := repository.NewCategoryRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateCategory(c.Context(), idCategory, &data); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "data update failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	}
}
func HandlerDeleteCategory(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idCategory := c.Params("category_id")
		bus := repository.NewCategoryRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteCategory(c.Context(), idCategory); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "delete category failed",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "Deleted successfully!",
		})
	}
}
func GetListCategory(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagging common.Pagging
		if err := c.BodyParser(&pagging); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not bind paging",
			})
		}

		bus := repository.NewCategoryRepoImpl(repomarketiml.NewDb(db))
		listCategory, err := bus.NewGetListCategory(c.Context(), &pagging)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Can't not get list business category",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"listProduct": listCategory,
			"total":       pagging.Total,
		})
	}
}

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

func CreateProductHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dataProduct marketmodels.Product
		if err := c.BodyParser(&dataProduct); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Can't bind product data"})
		}

		validate := validator.New()
		if err := validate.Struct(dataProduct); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Product validation failed"})
		}

		if dataProduct.Stock == 0 {
			dataProduct.Status = "out of stock"
		}

		id := uuid.New().String()
		currentTime := time.Now().UTC()
		product := &marketmodels.Product{
			Id:          primitive.NewObjectID(),
			Product_id:  &id,
			Title:       dataProduct.Title,
			Image:       dataProduct.Image,
			Description: dataProduct.Description,
			Price:       dataProduct.Price,
			Stock:       dataProduct.Stock,
			Status:      dataProduct.Status,
			Category_id: dataProduct.Category_id,
			Created_at:  &currentTime,
			Updated_at:  &currentTime,
		}

		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))
		if err := repo.NewCreateProduct(c.Context(), product); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"comment": "Product added successfully!"})
	}
}
func UpdateProductHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var updateProduct marketmodels.Product
		if err := c.BodyParser(&updateProduct); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Can't bind product data"})
		}

		idProductUpd := c.Params("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		if err := repo.NewUpdateProduct(c.Context(), idProductUpd, &updateProduct); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"comment": "Product updated successfully!"})
	}
}
func DeleteProductHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idProduct := c.Params("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		if err := repo.NewDeleteProduct(c.Context(), idProduct); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"comment": "Product deleted successfully!"})
	}
}
func GetProductHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idProduct := c.Params("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		dataProduct, err := repo.NewGetProduct(c.Context(), idProduct)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"dataProduct": dataProduct})
	}
}
func GetListProduct(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagging common.Pagging
		if err := c.BodyParser(&pagging); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Can't bind pagination data"})
		}

		filter := bson.M{"status": "stocking"}
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		listProduct, err := repo.NewGetListProduct(c.Context(), filter, &pagging)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product list"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"listProduct": listProduct, "total": pagging.Total})
	}
}

package markethandler

import (
	"net/http"
	"time"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/repository"
	repomarketiml "github.com/KingSupermarket/repository/repo_market_iml"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateProductHandler xử lý việc tạo sản phẩm mới.
func CreateProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataProduct marketmodels.Product
		if err := c.ShouldBind(&dataProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't bind product data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(dataProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product validation failed"})
			return
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
		if err := repo.NewCreateProduct(c.Request.Context(), product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"comment": "Product added successfully!"})
	}
}

// UpdateProductHandler xử lý việc cập nhật sản phẩm.
func UpdateProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateProduct marketmodels.Product
		if err := c.ShouldBind(&updateProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't bind product data"})
			return
		}

		idProductUpd := c.Param("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		if err := repo.NewUpdateProduct(c.Request.Context(), idProductUpd, &updateProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"comment": "Product updated successfully!"})
	}
}

// DeleteProductHandler xử lý việc xóa sản phẩm.
func DeleteProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idProduct := c.Param("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		if err := repo.NewDeleteProduct(c.Request.Context(), idProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"comment": "Product deleted successfully!"})
	}
}

// GetProductHandler xử lý việc lấy thông tin sản phẩm theo ID.
func GetProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idProduct := c.Param("product_id")
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		dataProduct, err := repo.NewGetProduct(c.Request.Context(), idProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"dataProduct": dataProduct})
	}
}

// GetListProduct xử lý việc lấy danh sách sản phẩm.
func GetListProduct(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagging common.Pagging
		if err := c.ShouldBind(&pagging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't bind pagination data"})
			return
		}

		filter := bson.M{"status": "stocking"}
		repo := repository.NewProductsRepoImpl(repomarketiml.NewDb(db))

		listProduct, err := repo.NewGetListProduct(c.Request.Context(), filter, &pagging)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product list"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"listProduct": listProduct, "total": pagging.Total})
	}
}

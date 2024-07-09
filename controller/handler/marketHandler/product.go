package markethandler

import (
	"net/http"
	"time"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	repomarketiml "github.com/KingSupermarket/repository/repo_market_iml"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataProduct *marketmodels.Product
		if err := c.ShouldBind(&dataProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind product",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not validate product",
			})
			return
		}
		idProduct, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not uuid product",
			})
			return
		}
		if dataProduct.Stock == 0 {
			dataProduct.Status = "out of stock"
		}
		id := idProduct.String()
		time := time.Now().UTC()
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
			Created_at:  &time,
			Updated_at:  &time,
		}
		bus := repomarketiml.NewDb(db)
		if err := bus.CreateProduct(c.Request.Context(), product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not bussiness product",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "added product successfully!",
		})
	}
}
func UpdateProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateProduct *marketmodels.Product
		if err := c.ShouldBind(&updateProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind product",
			})
			return
		}
		idProductUpd := c.Param("product_id")

		bus := repomarketiml.NewDb(db)
		if err := bus.UpdateProduct(c.Request.Context(), idProductUpd, updateProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not update bussiness product",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Updated bussiness product successfully!",
		})
	}
}
func DeleteProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idProduct := c.Param("product_id")
		bus := repomarketiml.NewDb(db)
		if err := bus.DeleteProduct(c.Request.Context(), idProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not delete bussiness product",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Delete bussiness product successfully!",
		})
	}
}
func GetProductHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idProduct := c.Param("product_id")
		bus := repomarketiml.NewDb(db)
		dataProduct, err := bus.GetProduct(c.Request.Context(), idProduct)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not get bussiness product",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"dataProduct": dataProduct,
		})
	}
}
func GetListProduct(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagging *common.Pagging
		if err := c.ShouldBind(&pagging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind pagging",
			})
			return
		}
		filter := map[string]interface{}{"status": "stocking"}
		bus := repomarketiml.NewDb(db)
		listProduct, err := bus.GetListProduct(c.Request.Context(), filter, pagging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not getList bussiness product",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"listProduct": listProduct,
			"total":       pagging.Total,
		})
	}
}

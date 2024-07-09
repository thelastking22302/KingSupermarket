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

func CreateCategoryHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataCategory *marketmodels.Category
		if err := c.ShouldBind(&dataCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind category",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not validate category",
			})
			return
		}
		idCategory, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not uuid category",
			})
			return
		}

		id := idCategory.String()
		time := time.Now().UTC()

		catgory := &marketmodels.Category{
			Id:          primitive.NewObjectID(),
			Category_id: &id,
			Name:        dataCategory.Name,
			Created_at:  &time,
			Updated_at:  &time,
		}
		bus := repomarketiml.NewDBCategory(db)
		if err := bus.CreateCategory(c.Request.Context(), catgory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not bussiness category",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "added category successfully!",
		})
	}
}

func HandlerGetCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idCategory := c.Param("category_id")

		biz := repomarketiml.NewDBCategory(db)
		data, err := biz.GetCategory(c.Request.Context(), idCategory)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "getCategory faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"category": data,
		})
	}
}
func HandlerUpdateCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idCategory := c.Param("category_id")
		var data marketmodels.Category
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bind data faild",
			})
			return
		}
		time := time.Now().UTC()
		data.Updated_at = &time

		biz := repomarketiml.NewDBCategory(db)
		if err := biz.UpdateCategory(c.Request.Context(), idCategory, &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " data faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
func HandlerDeleteCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idCategory := c.Param("category_id")
		biz := repomarketiml.NewDBCategory(db)
		if err := biz.DeleteCategory(c.Request.Context(), idCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " delete category faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Deleted success!",
		})
	}
}
func GetListCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagging *common.Pagging
		if err := c.ShouldBind(&pagging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind pagging",
			})
			return
		}

		bus := repomarketiml.NewDBCategory(db)
		listCategory, err := bus.GetListCategory(c.Request.Context(), pagging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not getList bussiness category",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"listProduct": listCategory,
			"total":       pagging.Total,
		})
	}
}

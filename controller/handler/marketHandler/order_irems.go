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

func CreateOrderItemsHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataOrderItems *marketmodels.OrderItems
		if err := c.ShouldBind(&dataOrderItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind order items",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataOrderItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not validate order items",
			})
			return
		}
		idOrderItems, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not uuid order items",
			})
			return
		}

		id := idOrderItems.String()
		time := time.Now().UTC()

		orderItems := &marketmodels.OrderItems{
			Id:            primitive.NewObjectID(),
			Order_Item_Id: &id,
			Order_Id:      dataOrderItems.Order_Id,
			Product_id:    dataOrderItems.Product_id,
			Quantity:      dataOrderItems.Quantity,
			Price:         dataOrderItems.Price,
			Created_at:    &time,
			Updated_at:    &time,
		}
		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateOrderItems(c.Request.Context(), orderItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not bussiness order items",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "added order items successfully!",
		})
	}
}

func HandlerGetOrderItems(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idOrderItems := c.Param("order_item_id")

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetOrderItems(c.Request.Context(), idOrderItems)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "getOrderItems faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"order item": data,
		})
	}
}
func HandlerUpdateOrderItems(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idOrderItems := c.Param("order_item_id")
		var data marketmodels.OrderItems
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bind data faild",
			})
			return
		}
		time := time.Now().UTC()
		data.Updated_at = &time
		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateOrderItems(c.Request.Context(), idOrderItems, &data); err != nil {
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
func HandlerDeleteOrderItems(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idOrderItems := c.Param("order_item_id")
		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteOrderItems(c.Request.Context(), idOrderItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " delete order items faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Deleted success!",
		})
	}
}
func GetListOrderItems(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagging *common.Pagging
		if err := c.ShouldBind(&pagging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind pagging",
			})
			return
		}

		bus := repository.NewOrderItemsRepoImpl(repomarketiml.NewDb(db))
		listOrderItems, err := bus.NewGetListOrderItems(c.Request.Context(), bson.M{}, pagging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not getList bussiness order items",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"listOrderItems": listOrderItems,
			"total":          pagging.Total,
		})
	}
}

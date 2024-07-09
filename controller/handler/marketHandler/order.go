package markethandler

import (
	"net/http"
	"time"

	marketmodels "github.com/KingSupermarket/model/marketModels"
	repomarketiml "github.com/KingSupermarket/repository/repo_market_iml"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrderHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		dataClaims := tokenData.(string)
		var dataOrder *marketmodels.Order
		if err := c.ShouldBind(&dataOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind order",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not validate order",
			})
			return
		}
		idOrder, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not uuid order",
			})
			return
		}

		id := idOrder.String()
		time := time.Now().UTC()

		order := &marketmodels.Order{
			Id:           primitive.NewObjectID(),
			Order_Id:     &id,
			User_Id:      &dataClaims,
			Address:      dataOrder.Address,
			Phone_Number: dataOrder.Phone_Number,
			Total_amount: dataOrder.Total_amount,
			Status:       dataOrder.Status,
			Notes:        dataOrder.Notes,
			Order_day:    &time,
			Created_at:   &time,
			Updated_at:   &time,
		}
		bus := repomarketiml.NewDbOrder(db)
		if err := bus.CreateOrder(c.Request.Context(), order, dataClaims); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not bussiness order",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "added order successfully!",
		})
	}
}

func HandlerGetOrder(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		dataClaims := tokenData.(string)
		idOrder := c.Param("order_id")

		biz := repomarketiml.NewDbOrder(db)
		data, err := biz.GetOrder(c.Request.Context(), idOrder, dataClaims)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "getOrder faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"order": data,
		})
	}
}
func HandlerUpdateOrder(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		dataClaims := tokenData.(string)
		idOrder := c.Param("Order_id")
		var data marketmodels.Order
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bind data faild",
			})
			return
		}
		time := time.Now().UTC()
		data.Updated_at = &time
		data.Order_day = &time
		biz := repomarketiml.NewDbOrder(db)
		if err := biz.UpdateOrder(c.Request.Context(), idOrder, &data, dataClaims); err != nil {
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
func HandlerDeleteOrder(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		dataClaims := tokenData.(string)
		idOrder := c.Param("order_id")
		biz := repomarketiml.NewDbOrder(db)
		if err := biz.DeleteOrder(c.Request.Context(), idOrder, dataClaims); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " delete order faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Deleted success!",
		})
	}
}

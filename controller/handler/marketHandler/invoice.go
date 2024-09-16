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

func CreateInvoiceHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataInvoice *marketmodels.Invoice
		if err := c.ShouldBind(&dataInvoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind invoice",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataInvoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not validate invoice",
			})
			return
		}
		idInvoice, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not uuid invoice",
			})
			return
		}

		id := idInvoice.String()
		time := time.Now().UTC()
		var status string = "PENDING"
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
			Created_at:       time,
			Updated_at:       time,
		}
		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewCreateInvoice(c.Request.Context(), invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not bussiness invoice",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "added invoice successfully!",
		})
	}
}

func HandlerGetInvoice(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idInvoice := c.Param("invoice_id")

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		data, err := bus.NewGetInvoice(c.Request.Context(), idInvoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "getInvoice faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"invoice": data,
		})
	}
}
func HandlerUpdateInvoices(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idInvoice := c.Param("invoice_id")
		var data marketmodels.Invoice
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bind data faild",
			})
			return
		}
		data.Updated_at = time.Now()

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewUpdateInvoice(c.Request.Context(), idInvoice, &data); err != nil {
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
func HandlerDeleteInvoice(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idInvoice := c.Param("invoice_id")
		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		if err := bus.NewDeleteInvoice(c.Request.Context(), idInvoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": " delete invocie faild",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "Deleted success!",
		})
	}
}
func GetListInvoice(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagging *common.Pagging
		if err := c.ShouldBind(&pagging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not shouldBind pagging",
			})
			return
		}

		bus := repository.NewInvoiceRepoImpl(repomarketiml.NewDb(db))
		listInvoice, err := bus.NewGetListInvoice(c.Request.Context(), bson.M{}, pagging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't not getList bussiness invoice",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"listInvoice": listInvoice,
			"total":       pagging.Total,
		})
	}
}

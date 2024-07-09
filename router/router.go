package router

import (
	markethandler "github.com/KingSupermarket/controller/handler/marketHandler"
	"github.com/KingSupermarket/controller/handler/userHandler"
	"github.com/KingSupermarket/middleware"
	"github.com/KingSupermarket/server"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	newRouter := r.Group("kingsupermarket")
	{
		auth := newRouter.Group("/auth")
		{
			auth.POST("/sign-up", userHandler.SignUpHandler(server.NewConnectionMongo()))
			auth.POST("/sign-in", userHandler.SignInHandler(server.NewConnectionMongo()))
		}
		user := newRouter.Group("/user", middleware.JwtMiddleware())
		{
			user.GET("/profile", userHandler.ProfileUserHandler(server.NewConnectionMongo()))
			user.DELETE("/delete", middleware.CheckAdmin(), userHandler.DeleteUserHandler(server.NewConnectionMongo()))
			user.PATCH("/update", userHandler.UpdateUserHandler(server.NewConnectionMongo()))
		}
		client := user.Group("/client")
		{
			market := client.Group("/market")
			{
				market.GET("product/:product_id", markethandler.GetProductHandler(server.NewConnectionMongo()))
				market.GET("category/:category_id", markethandler.HandlerGetCategory(server.NewConnectionMongo()))
				market.GET("invoice/:invoice_id", markethandler.HandlerGetInvoice(server.NewConnectionMongo()))
				market.GET("order/:order_id", markethandler.HandlerGetOrder(server.NewConnectionMongo()))
				market.GET("order-item/:order_item_id", markethandler.HandlerGetOrderItems(server.NewConnectionMongo()))
			}
		}
		admin := newRouter.Group("/admin")
		{
			product := admin.Group("/product")
			{
				product.POST("/", markethandler.CreateProductHandler(server.NewConnectionMongo()))
				product.GET("/:product_id", markethandler.GetProductHandler(server.NewConnectionMongo()))
				product.GET("/list", markethandler.GetListProduct(server.NewConnectionMongo()))
				product.PATCH("/:product_id", markethandler.UpdateProductHandler(server.NewConnectionMongo()))
				product.DELETE("/:product_id", markethandler.DeleteProductHandler(server.NewConnectionMongo()))
			}
			category := admin.Group("/category")
			{
				category.POST("", markethandler.CreateCategoryHandler(server.NewConnectionMongo()))
				category.GET(":category_id", markethandler.HandlerGetCategory(server.NewConnectionMongo()))
				category.GET("/list", markethandler.GetListCategory(server.NewConnectionMongo()))
				category.PATCH("/:category_id", markethandler.HandlerUpdateCategory(server.NewConnectionMongo()))
				category.DELETE("/:category_id", markethandler.HandlerDeleteCategory(server.NewConnectionMongo()))
			}
			invoice := admin.Group("/invoice")
			{
				invoice.POST("/", markethandler.CreateInvoiceHandler(server.NewConnectionMongo()))
				invoice.GET(":invoice_id", markethandler.HandlerGetInvoice(server.NewConnectionMongo()))
				invoice.GET("/list", markethandler.GetListInvoice(server.NewConnectionMongo()))
				invoice.PATCH("/:invoice_id", markethandler.HandlerUpdateInvoices(server.NewConnectionMongo()))
				invoice.DELETE("/:invoice_id", markethandler.HandlerDeleteInvoice(server.NewConnectionMongo()))
			}
			order := admin.Group("/order", middleware.JwtMiddleware())
			{
				order.POST("/", markethandler.CreateOrderHandler(server.NewConnectionMongo()))
				order.GET(":order_id", markethandler.HandlerGetOrder(server.NewConnectionMongo()))
				order.PATCH("/:order_id", markethandler.HandlerUpdateOrder(server.NewConnectionMongo()))
				order.DELETE("/:order_id", markethandler.HandlerDeleteOrder(server.NewConnectionMongo()))
			}
			orderItems := admin.Group("/order-item")
			{
				orderItems.POST("/", markethandler.CreateOrderItemsHandler(server.NewConnectionMongo()))
				orderItems.GET(":order_item_id", markethandler.HandlerGetOrderItems(server.NewConnectionMongo()))
				orderItems.GET("/list", markethandler.GetListOrderItems(server.NewConnectionMongo()))
				orderItems.PATCH("/:order_item_id", markethandler.HandlerUpdateOrderItems(server.NewConnectionMongo()))
				orderItems.DELETE("/:order_item_id", markethandler.HandlerDeleteOrderItems(server.NewConnectionMongo()))
			}
		}
	}
}

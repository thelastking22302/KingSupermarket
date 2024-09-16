package router

import (
	markethandler "github.com/KingSupermarket/controller/handler/marketHandler"
	"github.com/KingSupermarket/controller/handler/userHandler"
	"github.com/KingSupermarket/middleware"
	"github.com/KingSupermarket/server"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	conn := server.GetInstance()
	r := router.Group("/kingsupermarket")

	setupAuthRoutes(r.Group("/auth"), conn)
	setupUserRoutes(r.Group("/user"), conn)
	setupProductRoutes(r.Group("/product"), conn)
	setupCategoryRoutes(r.Group("/category"), conn)
	setupInvoiceRoutes(r.Group("/invoice"), conn)
	setupOrderRoutes(r.Group("/order"), conn)
	setupOrderItemsRoutes(r.Group("/order-item"), conn)

}
func setupAuthRoutes(auth *gin.RouterGroup, conn *server.SingletonMongo) {
	auth.POST("/sign-up", userHandler.SignUpHandler(conn.NewConnectionMongo()))
	auth.POST("/sign-in", userHandler.SignInHandler(conn.NewConnectionMongo()))
}

// User Routes
func setupUserRoutes(user *gin.RouterGroup, conn *server.SingletonMongo) {
	user.Use(middleware.JwtMiddleware())
	user.GET("/profile", userHandler.ProfileUserHandler(conn.NewConnectionMongo()))
	user.DELETE("/delete", middleware.CheckAdmin(), userHandler.DeleteUserHandler(conn.NewConnectionMongo()))
	user.PATCH("/update", userHandler.UpdateUserHandler(conn.NewConnectionMongo()))
}
func setupProductRoutes(product *gin.RouterGroup, conn *server.SingletonMongo) {
	product.POST("/", markethandler.CreateProductHandler(conn.NewConnectionMongo()))
	product.GET("/:product_id", markethandler.GetProductHandler(conn.NewConnectionMongo()))
	product.GET("/list", markethandler.GetListProduct(conn.NewConnectionMongo()))
	product.PATCH("/:product_id", markethandler.UpdateProductHandler(conn.NewConnectionMongo()))
	product.DELETE("/:product_id", markethandler.DeleteProductHandler(conn.NewConnectionMongo()))
}

// Category Routes
func setupCategoryRoutes(category *gin.RouterGroup, conn *server.SingletonMongo) {
	category.POST("/", markethandler.CreateCategoryHandler(conn.NewConnectionMongo()))
	category.GET("/:category_id", markethandler.HandlerGetCategory(conn.NewConnectionMongo()))
	category.GET("/list", markethandler.GetListCategory(conn.NewConnectionMongo()))
	category.PATCH("/:category_id", markethandler.HandlerUpdateCategory(conn.NewConnectionMongo()))
	category.DELETE("/:category_id", markethandler.HandlerDeleteCategory(conn.NewConnectionMongo()))
}

// Invoice Routes
func setupInvoiceRoutes(invoice *gin.RouterGroup, conn *server.SingletonMongo) {
	invoice.POST("/", markethandler.CreateCategoryHandler(conn.NewConnectionMongo()))
	invoice.GET("/:invoice_id", markethandler.HandlerGetInvoice(conn.NewConnectionMongo()))
	invoice.GET("/list", markethandler.GetListInvoice(conn.NewConnectionMongo()))
	invoice.PATCH("/:invoice_id", markethandler.HandlerUpdateInvoices(conn.NewConnectionMongo()))
	invoice.DELETE("/:invoice_id", markethandler.HandlerDeleteInvoice(conn.NewConnectionMongo()))
}

// Order Routes
func setupOrderRoutes(order *gin.RouterGroup, conn *server.SingletonMongo) {
	order.POST("/", markethandler.CreateCategoryHandler(conn.NewConnectionMongo()))
	order.GET("/:order_id", markethandler.HandlerGetOrder(conn.NewConnectionMongo()))
	order.PATCH("/:order_id", markethandler.HandlerUpdateOrder(conn.NewConnectionMongo()))
	order.DELETE("/:order_id", markethandler.HandlerDeleteOrder(conn.NewConnectionMongo()))
}

// Order Items Routes
func setupOrderItemsRoutes(orderItems *gin.RouterGroup, conn *server.SingletonMongo) {
	orderItems.POST("/", markethandler.CreateOrderItemsHandler(conn.NewConnectionMongo()))
	orderItems.GET("/:order_item_id", markethandler.HandlerGetOrderItems(conn.NewConnectionMongo()))
	orderItems.GET("/list", markethandler.GetListOrderItems(conn.NewConnectionMongo()))
	orderItems.PATCH("/:order_item_id", markethandler.HandlerUpdateOrderItems(conn.NewConnectionMongo()))
	orderItems.DELETE("/:order_item_id", markethandler.HandlerDeleteOrderItems(conn.NewConnectionMongo()))
}

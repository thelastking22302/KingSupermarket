package router

import (
	markethandler "github.com/KingSupermarket/controller/handler/marketHandler"
	"github.com/KingSupermarket/controller/handler/userHandler"
	"github.com/KingSupermarket/middleware"
	"github.com/KingSupermarket/server"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App) {
	conn := server.GetInstance()
	r := app.Group("/kingsupermarket")

	setupAuthRoutes(r.Group("/auth").(*fiber.Group), conn)
	setupUserRoutes(r.Group("/user").(*fiber.Group), conn)
	setupProductRoutes(r.Group("/product").(*fiber.Group), conn)
	setupCategoryRoutes(r.Group("/category").(*fiber.Group), conn)
	setupInvoiceRoutes(r.Group("/invoice").(*fiber.Group), conn)
	setupOrderRoutes(r.Group("/order").(*fiber.Group), conn)
	setupOrderItemsRoutes(r.Group("/order-items").(*fiber.Group), conn)
}

// Hàm thiết lập Auth Routes
func setupAuthRoutes(auth *fiber.Group, conn *server.SingletonMongo) {
	auth.Post("/sign-up", userHandler.SignUpHandler(conn.NewConnectionMongo()))
	auth.Post("/sign-in", userHandler.SignInHandler(conn.NewConnectionMongo()))
}

// Hàm thiết lập User Routes
func setupUserRoutes(user *fiber.Group, conn *server.SingletonMongo) {
	user.Use(middleware.JwtMiddleware())
	user.Get("/profile", userHandler.ProfileUserHandler(conn.NewConnectionMongo()))
	user.Delete("/delete", middleware.CheckAdmin(), userHandler.DeleteUserHandler(conn.NewConnectionMongo()))
	user.Patch("/update", userHandler.UpdateUserHandler(conn.NewConnectionMongo()))
}

// Hàm thiết lập Product Routes
func setupProductRoutes(product *fiber.Group, conn *server.SingletonMongo) {
	product.Post("/", markethandler.CreateProductHandler(conn.NewConnectionMongo()))
	product.Get("/:product_id", markethandler.GetProductHandler(conn.NewConnectionMongo()))
	product.Get("/list", markethandler.GetListProduct(conn.NewConnectionMongo()))
	product.Patch("/:product_id", markethandler.UpdateProductHandler(conn.NewConnectionMongo()))
	product.Delete("/:product_id", markethandler.DeleteProductHandler(conn.NewConnectionMongo()))
}

// Hàm thiết lập Category Routes
func setupCategoryRoutes(category *fiber.Group, conn *server.SingletonMongo) {
	category.Post("/", markethandler.CreateCategoryHandler(conn.NewConnectionMongo()))
	category.Get("/:category_id", markethandler.HandlerGetCategory(conn.NewConnectionMongo()))
	category.Get("/list", markethandler.GetListCategory(conn.NewConnectionMongo()))
	category.Patch("/:category_id", markethandler.HandlerUpdateCategory(conn.NewConnectionMongo()))
	category.Delete("/:category_id", markethandler.HandlerDeleteCategory(conn.NewConnectionMongo()))
}

// Hàm thiết lập Invoice Routes
func setupInvoiceRoutes(invoice *fiber.Group, conn *server.SingletonMongo) {
	invoice.Post("/", markethandler.CreateInvoiceHandler(conn.NewConnectionMongo()))
	invoice.Get("/:invoice_id", markethandler.HandlerGetInvoice(conn.NewConnectionMongo()))
	invoice.Get("/list", markethandler.GetListInvoice(conn.NewConnectionMongo()))
	invoice.Patch("/:invoice_id", markethandler.HandlerUpdateInvoices(conn.NewConnectionMongo()))
	invoice.Delete("/:invoice_id", markethandler.HandlerDeleteInvoice(conn.NewConnectionMongo()))
}

// Hàm thiết lập Order Routes
func setupOrderRoutes(order *fiber.Group, conn *server.SingletonMongo) {
	order.Post("/", markethandler.CreateOrderHandler(conn.NewConnectionMongo()))
	order.Get("/:order_id", markethandler.HandlerGetOrder(conn.NewConnectionMongo()))
	order.Patch("/:order_id", markethandler.HandlerUpdateOrder(conn.NewConnectionMongo()))
	order.Delete("/:order_id", markethandler.HandlerDeleteOrder(conn.NewConnectionMongo()))
}

// Hàm thiết lập Order Items Routes
func setupOrderItemsRoutes(orderItems *fiber.Group, conn *server.SingletonMongo) {
	orderItems.Post("/", markethandler.CreateOrderItemsHandler(conn.NewConnectionMongo()))
	orderItems.Get("/:order_item_id", markethandler.HandlerGetOrderItems(conn.NewConnectionMongo()))
	orderItems.Get("/list", markethandler.GetListOrderItems(conn.NewConnectionMongo()))
	orderItems.Patch("/:order_item_id", markethandler.HandlerUpdateOrderItems(conn.NewConnectionMongo()))
	orderItems.Delete("/:order_item_id", markethandler.HandlerDeleteOrderItems(conn.NewConnectionMongo()))
}

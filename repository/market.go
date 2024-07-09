package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryRepo interface {
	CreateCategory(ctx context.Context, data *marketmodels.Category) error
	GetCategory(ctx context.Context, id string) (*marketmodels.Category, error)
	GetListCategory(ctx context.Context, pagging *common.Pagging) ([]marketmodels.Category, error)
	UpdateCategory(ctx context.Context, id string, data *marketmodels.Category) error
	DeleteCategory(ctx context.Context, id string) error
}
type ProductsRepo interface {
	CreateProduct(ctx context.Context, data *marketmodels.Product) error
	GetProduct(ctx context.Context, id string) (*marketmodels.Product, error)
	GetListProduct(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Product, error)
	UpdateProduct(ctx context.Context, id string, data *marketmodels.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type InvoiceRepo interface {
	CreateInvoice(ctx context.Context, data *marketmodels.Invoice) error
	GetInvoice(ctx context.Context, id string) (*marketmodels.Invoice, error)
	GetListInvoice(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Invoice, error)
	UpdateInvoice(ctx context.Context, id string, data *marketmodels.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
}
type OrderRepo interface {
	CreateOrder(ctx context.Context, data *marketmodels.Order, idUser string) error
	GetOrder(ctx context.Context, id string, idUser string) (*marketmodels.Order, error)
	GetListOrder(ctx context.Context, filter bson.M, pagging *common.Pagging, idUser string) ([]marketmodels.Order, error)
	UpdateOrder(ctx context.Context, id string, data *marketmodels.Order, idUser string) error
	DeleteOrder(ctx context.Context, id string, idUser string) error
}
type OrderItemsRepo interface {
	CreateOrderItems(ctx context.Context, data *marketmodels.OrderItems) error
	GetOrderItems(ctx context.Context, id string) (*marketmodels.OrderItems, error)
	GetListOrderItems(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.OrderItems, error)
	UpdateOrderItems(ctx context.Context, id string, data *marketmodels.OrderItems) error
	DeleteOrderItems(ctx context.Context, id string) error
}

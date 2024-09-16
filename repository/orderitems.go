package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderItemsRepo interface {
	CreateOrderItems(ctx context.Context, data *marketmodels.OrderItems) error
	GetOrderItems(ctx context.Context, id string) (*marketmodels.OrderItems, error)
	GetListOrderItems(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.OrderItems, error)
	UpdateOrderItems(ctx context.Context, id string, data *marketmodels.OrderItems) error
	DeleteOrderItems(ctx context.Context, id string) error
}

type orderItemsRepoImpl struct {
	repo   OrderItemsRepo
	logger logger.Logger
}

func NewOrderItemsRepoImpl(repo OrderItemsRepo) *orderItemsRepoImpl {
	return &orderItemsRepoImpl{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

func (o *orderItemsRepoImpl) NewCreateOrderItems(ctx context.Context, data *marketmodels.OrderItems) error {
	if err := o.repo.CreateOrderItems(ctx, data); err != nil {
		o.logger.Errorf("Failed to create order item: %v", err)
		return err
	}
	o.logger.Infof("Order item created successfully: %s", data.Id)
	return nil
}

func (o *orderItemsRepoImpl) NewGetOrderItems(ctx context.Context, id string) (*marketmodels.OrderItems, error) {
	orderItem, err := o.repo.GetOrderItems(ctx, id)
	if err != nil {
		o.logger.Errorf("Error retrieving order item with ID %s: %v", id, err)
		return nil, err
	}
	o.logger.Infof("Retrieved order item with ID %s", id)
	return orderItem, nil
}

func (o *orderItemsRepoImpl) NewGetListOrderItems(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.OrderItems, error) {
	orderItems, err := o.repo.GetListOrderItems(ctx, filter, pagging)
	if err != nil {
		o.logger.Errorf("Error retrieving order items list: %v", err)
		return nil, err
	}
	o.logger.Infof("Retrieved list of order items")
	return orderItems, nil
}

func (o *orderItemsRepoImpl) NewUpdateOrderItems(ctx context.Context, id string, data *marketmodels.OrderItems) error {
	if err := o.repo.UpdateOrderItems(ctx, id, data); err != nil {
		o.logger.Errorf("Update order item failed for ID %s: %v", id, err)
		return err
	}
	o.logger.Infof("Order item updated successfully: %s", id)
	return nil
}

func (o *orderItemsRepoImpl) NewDeleteOrderItems(ctx context.Context, id string) error {
	if err := o.repo.DeleteOrderItems(ctx, id); err != nil {
		o.logger.Errorf("Failed to delete order item with ID %s: %v", id, err)
		return err
	}
	o.logger.Infof("Deleted order item with ID %s successfully", id)
	return nil
}

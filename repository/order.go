package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, data *marketmodels.Order, idUser string) error
	GetOrder(ctx context.Context, id string, idUser string) (*marketmodels.Order, error)
	GetListOrder(ctx context.Context, filter bson.M, pagging *common.Pagging, idUser string) ([]marketmodels.Order, error)
	UpdateOrder(ctx context.Context, id string, data *marketmodels.Order, idUser string) error
	DeleteOrder(ctx context.Context, id string, idUser string) error
}

type orderRepoImpl struct {
	repo   OrderRepo
	logger logger.Logger
}

func NewOrderRepoImpl(repo OrderRepo) *orderRepoImpl {
	return &orderRepoImpl{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

func (o *orderRepoImpl) NewCreateOrder(ctx context.Context, data *marketmodels.Order, idUser string) error {
	if err := o.repo.CreateOrder(ctx, data, idUser); err != nil {
		o.logger.Errorf("Failed to create order for user ID %s: %v", idUser, err)
		return err
	}
	o.logger.Infof("Order created successfully for user ID %s: %s", idUser, data.Id)
	return nil
}

func (o *orderRepoImpl) NewGetOrder(ctx context.Context, id string, idUser string) (*marketmodels.Order, error) {
	order, err := o.repo.GetOrder(ctx, id, idUser)
	if err != nil {
		o.logger.Errorf("Error retrieving order with ID %s for user ID %s: %v", id, idUser, err)
		return nil, err
	}
	o.logger.Infof("Retrieved order with ID %s for user ID %s", id, idUser)
	return order, nil
}

func (o *orderRepoImpl) NewGetListOrder(ctx context.Context, filter bson.M, pagging *common.Pagging, idUser string) ([]marketmodels.Order, error) {
	orders, err := o.repo.GetListOrder(ctx, filter, pagging, idUser)
	if err != nil {
		o.logger.Errorf("Error retrieving order list for user ID %s: %v", idUser, err)
		return nil, err
	}
	o.logger.Infof("Retrieved list of orders for user ID %s", idUser)
	return orders, nil
}

func (o *orderRepoImpl) NewUpdateOrder(ctx context.Context, id string, data *marketmodels.Order, idUser string) error {
	if err := o.repo.UpdateOrder(ctx, id, data, idUser); err != nil {
		o.logger.Errorf("Update order failed for ID %s and user ID %s: %v", id, idUser, err)
		return err
	}
	o.logger.Infof("Order updated successfully for ID %s and user ID %s", id, idUser)
	return nil
}

func (o *orderRepoImpl) NewDeleteOrder(ctx context.Context, id string, idUser string) error {
	if err := o.repo.DeleteOrder(ctx, id, idUser); err != nil {
		o.logger.Errorf("Failed to delete order with ID %s for user ID %s: %v", id, idUser, err)
		return err
	}
	o.logger.Infof("Deleted order with ID %s successfully for user ID %s", id, idUser)
	return nil
}

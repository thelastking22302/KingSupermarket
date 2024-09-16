package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductsRepo interface {
	CreateProduct(ctx context.Context, data *marketmodels.Product) error
	GetProduct(ctx context.Context, id string) (*marketmodels.Product, error)
	GetListProduct(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Product, error)
	UpdateProduct(ctx context.Context, id string, data *marketmodels.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type productsRepoImpl struct {
	repo   ProductsRepo
	logger logger.Logger
}

func NewProductsRepoImpl(repo ProductsRepo) *productsRepoImpl {
	return &productsRepoImpl{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

func (p *productsRepoImpl) NewCreateProduct(ctx context.Context, data *marketmodels.Product) error {
	if err := p.repo.CreateProduct(ctx, data); err != nil {
		p.logger.Errorf("Failed to create product: %v", err)
		return err
	}
	return nil
}

func (p *productsRepoImpl) NewGetProduct(ctx context.Context, id string) (*marketmodels.Product, error) {
	product, err := p.repo.GetProduct(ctx, id)
	if err != nil {
		p.logger.Errorf("Error retrieving product with ID %s: %v", id, err)
		return nil, err
	}
	p.logger.Infof("Retrieved product with ID %s", id)
	return product, nil
}

func (p *productsRepoImpl) NewGetListProduct(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Product, error) {
	products, err := p.repo.GetListProduct(ctx, filter, pagging)
	if err != nil {
		p.logger.Errorf("Error retrieving product list: %v", err)
		return nil, err
	}
	p.logger.Infof("Retrieved list of products")
	return products, nil
}

func (p *productsRepoImpl) NewUpdateProduct(ctx context.Context, id string, data *marketmodels.Product) error {
	if err := p.repo.UpdateProduct(ctx, id, data); err != nil {
		p.logger.Errorf("Update product failed for ID %s: %v", id, err)
		return err
	}
	p.logger.Infof("Product updated successfully: %s", id)
	return nil
}

func (p *productsRepoImpl) NewDeleteProduct(ctx context.Context, id string) error {
	if err := p.repo.DeleteProduct(ctx, id); err != nil {
		p.logger.Errorf("Failed to delete product with ID %s: %v", id, err)
		return err
	}
	p.logger.Infof("Deleted product with ID %s successfully", id)
	return nil
}

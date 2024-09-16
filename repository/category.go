package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/pkg/logger"
)

type CategoryRepo interface {
	CreateCategory(ctx context.Context, data *marketmodels.Category) error
	GetCategory(ctx context.Context, id string) (*marketmodels.Category, error)
	GetListCategory(ctx context.Context, pagging *common.Pagging) ([]marketmodels.Category, error)
	UpdateCategory(ctx context.Context, id string, data *marketmodels.Category) error
	DeleteCategory(ctx context.Context, id string) error
}

type categoryRepoImpl struct {
	repo   CategoryRepo
	logger logger.Logger
}

func NewCategoryRepoImpl(repo CategoryRepo) *categoryRepoImpl {
	return &categoryRepoImpl{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

func (c *categoryRepoImpl) NewCreateCategory(ctx context.Context, data *marketmodels.Category) error {
	if err := c.repo.CreateCategory(ctx, data); err != nil {
		c.logger.Errorf("Failed to create category: %v", err)
		return err
	}
	c.logger.Infof("Category created successfully: %s", data.Name)
	return nil
}

func (c *categoryRepoImpl) NewGetCategory(ctx context.Context, id string) (*marketmodels.Category, error) {
	category, err := c.repo.GetCategory(ctx, id)
	if err != nil {
		c.logger.Errorf("Error retrieving category with ID %s: %v", id, err)
		return nil, err
	}
	c.logger.Infof("Retrieved category with ID %s", id)
	return category, nil
}

func (c *categoryRepoImpl) NewGetListCategory(ctx context.Context, pagging *common.Pagging) ([]marketmodels.Category, error) {
	categories, err := c.repo.GetListCategory(ctx, pagging)
	if err != nil {
		c.logger.Errorf("Error retrieving categories: %v", err)
		return nil, err
	}
	c.logger.Infof("Retrieved list of categories")
	return categories, nil
}

func (c *categoryRepoImpl) NewUpdateCategory(ctx context.Context, id string, data *marketmodels.Category) error {
	if err := c.repo.UpdateCategory(ctx, id, data); err != nil {
		c.logger.Errorf("Update category failed for ID %s: %v", id, err)
		return err
	}
	c.logger.Infof("Category updated successfully: %s", id)
	return nil
}

func (c *categoryRepoImpl) NewDeleteCategory(ctx context.Context, id string) error {
	if err := c.repo.DeleteCategory(ctx, id); err != nil {
		c.logger.Errorf("Failed to delete category with ID %s: %v", id, err)
		return err
	}
	c.logger.Infof("Deleted category with ID %s successfully", id)
	return nil
}

package repository

import (
	"context"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"github.com/KingSupermarket/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type InvoiceRepo interface {
	CreateInvoice(ctx context.Context, data *marketmodels.Invoice) error
	GetInvoice(ctx context.Context, id string) (*marketmodels.Invoice, error)
	GetListInvoice(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Invoice, error)
	UpdateInvoice(ctx context.Context, id string, data *marketmodels.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
}

type invoiceRepoImpl struct {
	repo   InvoiceRepo
	logger logger.Logger
}

func NewInvoiceRepoImpl(repo InvoiceRepo) *invoiceRepoImpl {
	return &invoiceRepoImpl{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

func (i *invoiceRepoImpl) NewCreateInvoice(ctx context.Context, data *marketmodels.Invoice) error {
	if err := i.repo.CreateInvoice(ctx, data); err != nil {
		i.logger.Errorf("Failed to create invoice: %v", err)
		return err
	}
	i.logger.Infof("Invoice created successfully: %s", data.ID)
	return nil
}

func (i *invoiceRepoImpl) NewGetInvoice(ctx context.Context, id string) (*marketmodels.Invoice, error) {
	invoice, err := i.repo.GetInvoice(ctx, id)
	if err != nil {
		i.logger.Errorf("Error retrieving invoice with ID %s: %v", id, err)
		return nil, err
	}
	i.logger.Infof("Retrieved invoice with ID %s", id)
	return invoice, nil
}

func (i *invoiceRepoImpl) NewGetListInvoice(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Invoice, error) {
	invoices, err := i.repo.GetListInvoice(ctx, filter, pagging)
	if err != nil {
		i.logger.Errorf("Error retrieving invoice list: %v", err)
		return nil, err
	}
	i.logger.Infof("Retrieved list of invoices")
	return invoices, nil
}

func (i *invoiceRepoImpl) NewUpdateInvoice(ctx context.Context, id string, data *marketmodels.Invoice) error {
	if err := i.repo.UpdateInvoice(ctx, id, data); err != nil {
		i.logger.Errorf("Update invoice failed for ID %s: %v", id, err)
		return err
	}
	i.logger.Infof("Invoice updated successfully: %s", id)
	return nil
}

func (i *invoiceRepoImpl) NewDeleteInvoice(ctx context.Context, id string) error {
	if err := i.repo.DeleteInvoice(ctx, id); err != nil {
		i.logger.Errorf("Failed to delete invoice with ID %s: %v", id, err)
		return err
	}
	i.logger.Infof("Deleted invoice with ID %s successfully", id)
	return nil
}

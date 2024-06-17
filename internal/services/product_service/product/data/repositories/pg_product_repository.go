package repositories

import (
	"context"
	"fmt"

	gormpsql "github.com/Drack112/Golang-GQRS-Example/internal/pkg/gorm_psql"
	"github.com/Drack112/Golang-GQRS-Example/internal/pkg/logger"
	"github.com/Drack112/Golang-GQRS-Example/internal/pkg/utils"
	"github.com/Drack112/Golang-GQRS-Example/internal/services/product_service/product/data/contracts"
	"github.com/Drack112/Golang-GQRS-Example/internal/services/product_service/product/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PostgresProductRepository struct {
	log  logger.ILogger
	cfg  *gormpsql.GormPostgresConfig
	db   *pgxpool.Pool
	gorm *gorm.DB
}

func NewPostgresProductRepository(log logger.ILogger, cfg *gormpsql.GormPostgresConfig, gorm *gorm.DB) contracts.ProductRepository {
	return &PostgresProductRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *PostgresProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	result, err := gormpsql.Paginate[*models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgresProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gormpsql.Paginate[*models.Product](ctx, listQuery, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	return &product, nil
}

func (p *PostgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return product, nil
}

func (p *PostgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error in updating product with id %s into the database.", updateProduct.ProductId))
	}

	return updateProduct, nil
}

func (p *PostgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return errors.Wrap(err, "error in the deleting product into the database.")
	}

	return nil
}

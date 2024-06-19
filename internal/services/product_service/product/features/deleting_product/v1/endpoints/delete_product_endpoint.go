package endpoints

import (
	"context"
	"net/http"

	echomiddleware "github.com/Drack112/Golang-GQRS-Example/internal/pkg/http/echo/middleware"
	"github.com/Drack112/Golang-GQRS-Example/internal/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	commandsv1 "github.com/Drack112/Golang-GQRS-Example/internal/services/product_service/product/features/deleting_product/v1/commands"
	dtosv1 "github.com/Drack112/Golang-GQRS-Example/internal/services/product_service/product/features/deleting_product/v1/dtos"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.DELETE("/:id", deleteProduct(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

// DeleteProduct
// @Tags        Products
// @Summary     Delete product
// @Description Delete existing product
// @Accept      json
// @Produce     json
// @Success     204
// @Param       id path string true "Product ID"
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [delete]
func deleteProduct(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &dtosv1.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commandsv1.NewDeleteProduct(request.ProductID)
		if err := validator.StructCtx(ctx, command); err != nil {
			log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*commandsv1.DeleteProduct, mediatr.Unit](ctx, command)
		if err != nil {
			log.Warn("DeleteProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

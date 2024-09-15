package delivery

import (
	"errors"
	tenderUsecase "mymodule/app/internal/tender/usecase"
	"mymodule/app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Delivery struct {
	TenderUC tenderUsecase.UseCaseI
}

func (delivery *Delivery) createTender(c echo.Context) error {
	var tender models.Tender
	if err := c.Bind(&tender); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	err := delivery.TenderUC.CreateTender(&tender)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotPermission):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusForbidden, models.ErrUserNotPermission.Error())
		case errors.Is(err, models.ErrUserInvalid):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUserInvalid.Error())
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, tender)

}

func (delivery *Delivery) getStatusTender(c echo.Context) error {
	var tender models.Tender
	tender.Id = c.Param("tenderId")
	tender.CreatorUsername = c.QueryParam("username")

	status, err := delivery.TenderUC.GetStatusTender(tender)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotPermission):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusForbidden, models.ErrUserNotPermission.Error())
		case errors.Is(err, models.ErrUserInvalid):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUserInvalid.Error())
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		case errors.Is(err, models.ErrTenderNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrTenderNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, status)
}

func (delivery *Delivery) updateStatusTender(c echo.Context) error {
	var tender models.Tender
	tender.Id = c.Param("tenderId")
	tender.CreatorUsername = c.QueryParam("username")
	tender.Status = c.QueryParam("status")

	if tender.Status != models.CREATEDTEN && tender.Status != models.CLOSEDTEN && tender.Status != models.PUBLISHEDTEN {
		c.Logger().Error(models.ErrBadData)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.TenderUC.UpdateStatusTender(&tender)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotPermission):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusForbidden, models.ErrUserNotPermission.Error())
		case errors.Is(err, models.ErrUserInvalid):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUserInvalid.Error())
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		case errors.Is(err, models.ErrTenderNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrTenderNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, tender)
}

func (delivery *Delivery) updateTender(c echo.Context) error {
	var tender models.Tender
	if err := c.Bind(&tender); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	tender.Id = c.Param("tenderId")
	tender.CreatorUsername = c.QueryParam("username")

	if tender.Status != models.CREATEDTEN && tender.Status != models.CLOSEDTEN && tender.Status != models.PUBLISHEDTEN {
		c.Logger().Error(models.ErrBadData)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.TenderUC.UpdateTender(&tender)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotPermission):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusForbidden, models.ErrUserNotPermission.Error())
		case errors.Is(err, models.ErrUserInvalid):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUserInvalid.Error())
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		case errors.Is(err, models.ErrTenderNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrTenderNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, tender)
}

func (delivery *Delivery) selectTenders(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	serviceType := c.QueryParam("service_type")

	limit, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	} else if limitStr == "" {
		limit = 5
	}
	if limit < 1 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetStr)
	if offsetStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	} else if offsetStr == "" {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}

	if serviceType != "" && serviceType != models.CONSTRUCTION && serviceType != models.DELIVERY && serviceType != models.MANUFACTURE {
		c.Logger().Error(models.ErrBadData)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	tenders, err := delivery.TenderUC.SelectTenders(limit, offset, serviceType)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, tenders)
}

func (delivery *Delivery) selectTendersByUsername(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	username := c.QueryParam("username")

	limit, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	} else if limitStr == "" {
		limit = 5
	}
	if limit < 1 {
		limit = 5
	}

	offset, err := strconv.Atoi(offsetStr)
	if offsetStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	} else if offsetStr == "" {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}

	tenders, err := delivery.TenderUC.SelectTendersByUsername(limit, offset, username)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrBadData):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
		case errors.Is(err, models.ErrUserInvalid):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUserInvalid.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, tenders)
}

func NewDelivery(e *echo.Group, tenderUC tenderUsecase.UseCaseI) {
	handler := &Delivery{
		TenderUC: tenderUC,
	}
	e.POST("/tenders/new", handler.createTender)
	e.GET("/tenders/:tenderId/status", handler.getStatusTender)
	e.PUT("/tenders/:tenderId/status", handler.updateStatusTender)
	e.PATCH("/tenders/:tenderId/edit", handler.updateTender)
	e.GET("/tenders", handler.selectTenders)
	e.GET("/tenders/my", handler.selectTendersByUsername)
}

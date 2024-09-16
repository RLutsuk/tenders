package delivery

import (
	"errors"
	bidUsecase "mymodule/app/internal/bid/usecase"
	"mymodule/app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Delivery struct {
	BidUC bidUsecase.UseCaseI
}

func (delivery *Delivery) createBid(c echo.Context) error {
	var bid models.Bid
	if err := c.Bind(&bid); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	err := delivery.BidUC.CreateBid(&bid)
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
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrTenderNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, bid)

}

func (delivery *Delivery) getStatusBid(c echo.Context) error {
	bidId := c.Param("bidId")
	username := c.QueryParam("username")

	status, err := delivery.BidUC.GetStatusBid(bidId, username)

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
		case errors.Is(err, models.ErrBidNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrBidNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, status)
}

func (delivery *Delivery) updateStatusBid(c echo.Context) error {
	var bid models.Bid
	bid.Id = c.Param("bidId")
	username := c.QueryParam("username")
	bid.Status = c.QueryParam("status")

	if bid.Status != models.CANCELEDBID && bid.Status != models.CREATEDBID && bid.Status != models.PUBLISHEDBID {
		c.Logger().Error(models.ErrBadData)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.BidUC.UpdateStatusBid(&bid, username)

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
		case errors.Is(err, models.ErrBidNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrBidNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, bid)
}

func (delivery *Delivery) updateBid(c echo.Context) error {
	var bid models.Bid
	if err := c.Bind(&bid); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	bid.Id = c.Param("bidId")
	username := c.QueryParam("username")

	err := delivery.BidUC.UpdateBid(&bid, username)

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
		case errors.Is(err, models.ErrBidNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrBidNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, bid)
}

func (delivery *Delivery) selectBidsByUsername(c echo.Context) error {
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

	bids, err := delivery.BidUC.SelectBidsByUsername(limit, offset, username)

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
	return c.JSON(http.StatusOK, bids)
}

func (delivery *Delivery) selectBidsByTender(c echo.Context) error {
	tenderId := c.Param("tenderId")
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

	bids, err := delivery.BidUC.SelectBidsByTenderId(limit, offset, username, tenderId)

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
	return c.JSON(http.StatusOK, bids)
}

func (delivery *Delivery) submitDecision(c echo.Context) error {
	var bid models.Bid
	bid.Id = c.Param("bidId")
	decision := c.QueryParam("decision")
	username := c.QueryParam("username")

	if decision != "Approved" && decision != "Rejected" {
		c.Logger().Error(models.ErrBadData)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadData.Error())
	}

	err := delivery.BidUC.SubmitDecision(&bid, username, decision)
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
		case errors.Is(err, models.ErrBidNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrBidNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, bid)
}

func NewDelivery(e *echo.Group, bidUC bidUsecase.UseCaseI) {
	handler := &Delivery{
		BidUC: bidUC,
	}
	e.POST("/bids/new", handler.createBid)
	e.GET("/bids/:bidId/status", handler.getStatusBid)
	e.PUT("/bids/:bidId/status", handler.updateStatusBid)
	e.PATCH("/bids/:bidId/edit", handler.updateBid)
	e.GET("/bids/my", handler.selectBidsByUsername)
	e.GET("/bids/:tenderId/list", handler.selectBidsByTender)
	e.PUT("/bids/:bidId/submit_decision", handler.submitDecision)
}

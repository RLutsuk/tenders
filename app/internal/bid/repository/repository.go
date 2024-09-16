package repository

import (
	"mymodule/app/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RepositoryI interface {
	CreateBid(bid *models.Bid) error
	SelectBidById(bidId string) (*models.Bid, error)
	SelectBidsByUserId(limit, offset int, userId string) ([]*models.Bid, error)
	SelectBidsByTederId(limit, offset int, username, tenderId string) ([]*models.Bid, error)
	UpdateBid(bid models.Bid) error
	UpdateStatusBid(bid models.Bid) error
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbBid *dataBase) CreateBid(bid *models.Bid) error {
	tx := dbBid.db.Table("bid").Create(bid)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table bid)")
	}
	return nil
}

func (dbBid *dataBase) SelectBidById(bidId string) (*models.Bid, error) {
	bid := models.Bid{}

	tx := dbBid.db.Table("bid").Where("id = ?", bidId).Take(&bid)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table bid)")
	}
	return &bid, nil
}

func (dbBid *dataBase) UpdateStatusBid(bid models.Bid) error {
	tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("status", bid.Status)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table bid)")
	}
	tx = dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("version", bid.Version)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table bid)")
	}
	return nil
}

func (dbBid *dataBase) UpdateBid(bid models.Bid) error {
	if bid.Name != "" {
		tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("name", bid.Name)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table bid)")
		}
	}

	if bid.Description != "" {
		tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("description", bid.Description)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table bid)")
		}
	}

	if bid.Status != "" {
		tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("status", bid.Status)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table bid)")
		}
	}

	if bid.AuthorType != "" {
		tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("author_type", bid.AuthorType)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table bid)")
		}
	}

	tx := dbBid.db.Table("bid").Where("id = ?", bid.Id).Update("version", bid.Version)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table bid)")
	}

	return nil
}

func (dbBid *dataBase) SelectBidsByUserId(limit, offset int, userId string) ([]*models.Bid, error) {
	bids := make([]*models.Bid, 0, limit)
	if offset != 0 {
		tx := dbBid.db.Table("bid").Offset(offset).Limit(limit).Where("user_id = ?", userId).Order("name asc").Find(&bids)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table bid)")
		}
	} else {
		tx := dbBid.db.Table("bid").Limit(limit).Where("user_id = ?", userId).Order("name asc").Find(&bids)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table bid)")
		}
	}
	return bids, nil
}

func (dbBid *dataBase) SelectBidsByTederId(limit, offset int, username, tenderId string) ([]*models.Bid, error) {
	bids := make([]*models.Bid, 0, limit)
	if offset != 0 {
		tx := dbBid.db.Table("bid").Offset(offset).Limit(limit).Where("tender_id = ?", tenderId).Order("name asc").Find(&bids)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table bid)")
		}
	} else {
		tx := dbBid.db.Table("bid").Limit(limit).Where("tender_id = ?", tenderId).Order("name asc").Find(&bids)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table bid)")
		}
	}
	return bids, nil
}

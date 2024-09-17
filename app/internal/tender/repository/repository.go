package repository

import (
	"mymodule/app/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RepositoryI interface {
	CreateTender(tender *models.Tender) error
	SelectTenderById(tenderId string) (*models.Tender, error)
	UpdateTenderStatus(tender models.Tender) error
	UpdateTender(tender models.Tender) error
	SelectTenders(limit, offset int, service_type string) ([]*models.Tender, error)
	SelectTendersByUsername(limit, offset int, username string) ([]*models.Tender, error)
	ClosedСonfirmedTender(tendetId string) error
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbTender *dataBase) CreateTender(tender *models.Tender) error {
	tx := dbTender.db.Table("tender").Create(tender)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tenders)")
	}
	return nil
}

func (dbTender *dataBase) SelectTenderById(tenderId string) (*models.Tender, error) {
	tender := models.Tender{}
	tx := dbTender.db.Table("tender").Where("id = ?", tenderId).Take(&tender)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table tender)")
	}
	return &tender, nil
}

func (dbTender *dataBase) UpdateTenderStatus(tender models.Tender) error {
	tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("status", tender.Status)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tender)")
	}
	return nil
}

func (dbTender *dataBase) UpdateTender(tender models.Tender) error {
	if tender.Name != "" {
		tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("name", tender.Name)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table tender)")
		}
	}

	if tender.Description != "" {
		tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("description", tender.Description)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table tender)")
		}
	}

	if tender.Status != "" {
		tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("status", tender.Status)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table tender)")
		}
	}

	if tender.ServiceType != "" {
		tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("service_type", tender.ServiceType)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table tender)")
		}
	}

	tx := dbTender.db.Table("tender").Where("id = ?", tender.Id).Update("version", tender.Version+1)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tender)")
	}

	return nil
}

func (dbTender *dataBase) SelectTenders(limit, offset int, service_type string) ([]*models.Tender, error) {
	tenders := make([]*models.Tender, 0, limit) //почему 10

	if service_type == "" {
		if offset != 0 {
			tx := dbTender.db.Table("tender").Offset(offset).Limit(limit).Where("status = ?", models.PUBLISHEDTEN).Order("name asc").Find(&tenders)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table tender)")
			}
		} else {
			tx := dbTender.db.Table("tender").Limit(limit).Where("status = ?", models.PUBLISHEDTEN).Order("name asc").Find(&tenders)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table tender)")
			}
		}
	} else {
		if offset != 0 {
			tx := dbTender.db.Table("tender").Offset(offset).Limit(limit).Where("service_type = ? AND status = ?", service_type, models.PUBLISHEDTEN).Order("name asc").Find(&tenders)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table tender)")
			}
		} else {
			tx := dbTender.db.Table("tender").Limit(limit).Where("service_type = ? AND status = ?", service_type, models.PUBLISHEDTEN).Order("name asc").Find(&tenders)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table tender)")
			}
		}
	}
	return tenders, nil
}

func (dbTender *dataBase) SelectTendersByUsername(limit, offset int, username string) ([]*models.Tender, error) {
	tenders := make([]*models.Tender, 0, limit)
	if offset != 0 {
		tx := dbTender.db.Table("tender").Offset(offset).Limit(limit).Where("user_name = ?", username).Order("name asc").Find(&tenders)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table tenders)")
		}
	} else {
		tx := dbTender.db.Table("tender").Limit(limit).Where("user_name = ?", username).Order("name asc").Find(&tenders)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table tenders)")
		}
	}
	return tenders, nil
}

func (dbTender *dataBase) ClosedСonfirmedTender(tenderId string) error {
	tx := dbTender.db.Table("tender").Where("id = ?", tenderId).Update("status", "Closed")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tender)")
	}
	tx = dbTender.db.Table("tender").Where("id = ?", tenderId).UpdateColumn("version", gorm.Expr("version + ?", 1))
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tender)")
	}
	return nil
}

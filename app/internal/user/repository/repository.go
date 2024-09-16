package repository

import (
	"mymodule/app/models"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type RepositoryI interface {
	SelectUserByUsername(nickname string) (*models.User, error)
	SelectUserById(userId string) (*models.User, error)
	CheckUserOrganization(userId, organizationId string) (bool, error)
	CheckUserIsWorkerOrganization(userId string) (bool, error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbUser *dataBase) SelectUserByUsername(username string) (*models.User, error) {
	user := models.User{}

	tx := dbUser.db.Table("employee").Where("username = ?", username).Take(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}
	return &user, nil
}

func (dbUser *dataBase) SelectUserById(userId string) (*models.User, error) {
	user := models.User{}

	tx := dbUser.db.Table("employee").Where("id = ?", userId).Take(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}
	return &user, nil
}

func (dbUser *dataBase) CheckUserOrganization(userId, organizationId string) (bool, error) {
	exist := false
	tx := dbUser.db.Table("organization_responsible").Where("user_id = ? AND organization_id = ?", userId, organizationId).First(&models.OrganizationResponsible{})
	if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table organization_responsible)")
	}
	if tx.RowsAffected > 0 {
		exist = true
	}

	return exist, nil
}

func (dbUser *dataBase) CheckUserIsWorkerOrganization(userId string) (bool, error) {
	exist := false
	tx := dbUser.db.Table("organization_responsible").Where("user_id = ?", userId).First(&models.OrganizationResponsible{})
	if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table organization_responsible)")
	}
	if tx.RowsAffected > 0 {
		exist = true
	}

	return exist, nil
}

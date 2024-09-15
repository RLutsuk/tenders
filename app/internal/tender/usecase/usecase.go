package usecase

import (
	tenderRep "mymodule/app/internal/tender/repository"
	userRep "mymodule/app/internal/user/repository"
	"mymodule/app/models"
	"time"

	"github.com/go-openapi/strfmt"
)

type UseCaseI interface {
	CreateTender(tender *models.Tender) error
	GetStatusTender(tender models.Tender) (string, error)
	UpdateStatusTender(tender *models.Tender) error
	UpdateTender(tender *models.Tender) error
	SelectTenders(limit, offset int, service_type string) ([]*models.Tender, error)
	SelectTendersByUsername(limit, offset int, username string) ([]*models.Tender, error)
}

type useCase struct {
	tenderRepository tenderRep.RepositoryI
	userRepository   userRep.RepositoryI
}

func New(tenderRepository tenderRep.RepositoryI, userRepository userRep.RepositoryI) UseCaseI {
	return &useCase{
		tenderRepository: tenderRepository,
		userRepository:   userRepository,
	}
}

func (uc *useCase) CreateTender(tender *models.Tender) error {
	err := uc.validateTender(tender)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.SelectUserByUsername(tender.CreatorUsername)
	if err != nil {
		return models.ErrUserInvalid
	}

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, tender.OrganizationId)
	if err != nil || !exist {
		return models.ErrUserNotPermission
	}

	tender.Version = 1 // потому что создание
	tender.Status = models.CREATEDTEN
	tender.Created = strfmt.DateTime(time.Now()) // по хорошему спросить у чата гпт как работает тут горм и зачем тут strfmt.DateTime

	err = uc.tenderRepository.CreateTender(tender)
	if err != nil {
		return err
	}
	tender.CreatorUsername = ""
	tender.OrganizationId = ""
	return nil
}

func (uc *useCase) validateTender(tender *models.Tender) error {
	if len([]rune(tender.CreatorUsername)) > 100 {
		return models.ErrUserInvalid
	}

	if len([]rune(tender.Name)) > 100 {
		return models.ErrBadData
	}

	if len([]rune(tender.Description)) > 500 {
		return models.ErrBadData
	}

	if tender.ServiceType != models.DELIVERY && tender.ServiceType != models.CONSTRUCTION && tender.ServiceType != models.MANUFACTURE {
		return models.ErrBadData
	}

	return nil
}

func (uc *useCase) GetStatusTender(tender models.Tender) (string, error) {
	err := uc.validateTender(&tender)
	if err != nil {
		return tender.Status, err
	}

	user, err := uc.userRepository.SelectUserByUsername(tender.CreatorUsername)
	if err != nil {
		return tender.Status, models.ErrUserInvalid
	}

	existTender, err := uc.tenderRepository.SelectTenderById(tender.Id)
	if err != nil {
		return tender.Status, models.ErrTenderNotFound
	}
	tender.Status, tender.OrganizationId = existTender.Status, existTender.OrganizationId

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, tender.OrganizationId)
	if err != nil || !exist {
		return tender.Status, models.ErrUserNotPermission
	}

	return tender.Status, nil
}

func (uc *useCase) UpdateStatusTender(tender *models.Tender) error {
	err := uc.validateTender(tender)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.SelectUserByUsername(tender.CreatorUsername)
	if err != nil {
		return models.ErrUserInvalid
	}

	existTender, err := uc.tenderRepository.SelectTenderById(tender.Id)
	if err != nil {
		return models.ErrTenderNotFound
	}

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, existTender.OrganizationId)
	if err != nil || !exist {
		return models.ErrUserNotPermission
	}

	err = uc.tenderRepository.UpdateTenderStatus(*tender)
	if err != nil {
		return models.ErrBadData
	}

	tender.Name = existTender.Name
	tender.Description = existTender.Description
	tender.ServiceType = existTender.ServiceType
	tender.Version = existTender.Version
	tender.Created = existTender.Created
	tender.CreatorUsername = ""
	return nil
}

func (uc *useCase) UpdateTender(tender *models.Tender) error {
	err := uc.validateTender(tender)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.SelectUserByUsername(tender.CreatorUsername)
	if err != nil {
		return models.ErrUserInvalid
	}

	existTender, err := uc.tenderRepository.SelectTenderById(tender.Id)
	if err != nil {
		return models.ErrTenderNotFound
	}

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, existTender.OrganizationId)
	if err != nil || !exist {
		return models.ErrUserNotPermission
	}

	tender.Version = existTender.Version

	err = uc.tenderRepository.UpdateTender(*tender)
	if err != nil {
		return models.ErrBadData
	}
	if tender.Name == "" {
		tender.Name = existTender.Name
	}
	if tender.Description == "" {
		tender.Description = existTender.Description
	}
	if tender.ServiceType == "" {
		tender.ServiceType = existTender.ServiceType
	}
	if tender.Status == "" {
		tender.Status = existTender.Status
	}

	tender.Version = tender.Version + 1
	tender.Created = existTender.Created
	tender.CreatorUsername = ""

	return nil
}

func (uc *useCase) SelectTenders(limit, offset int, service_type string) ([]*models.Tender, error) {
	tenders, err := uc.tenderRepository.SelectTenders(limit, offset, service_type)
	if err != nil {
		return nil, err
	}

	for i := range tenders {
		tenders[i].CreatorUsername = ""
		tenders[i].OrganizationId = ""
	}

	return tenders, nil
}

func (uc *useCase) SelectTendersByUsername(limit, offset int, username string) ([]*models.Tender, error) {
	_, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return nil, models.ErrUserInvalid
	}

	tenders, err := uc.tenderRepository.SelectTendersByUsername(limit, offset, username)
	if err != nil {
		return nil, err
	}
	for i := range tenders {
		tenders[i].CreatorUsername = ""
		tenders[i].OrganizationId = ""
	}

	return tenders, nil
}

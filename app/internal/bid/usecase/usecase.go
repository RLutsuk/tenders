package usecase

import (
	bidRep "mymodule/app/internal/bid/repository"
	tenderRep "mymodule/app/internal/tender/repository"
	userRep "mymodule/app/internal/user/repository"
	"mymodule/app/models"
	"time"

	"github.com/go-openapi/strfmt"
)

type UseCaseI interface {
	CreateBid(bid *models.Bid) error
	GetStatusBid(bidId string, username string) (string, error)
	SelectBidsByUsername(limit, offset int, username string) ([]*models.Bid, error)
	SelectBidsByTenderId(limit, offset int, username string, tenderId string) ([]*models.Bid, error)
	SubmitDecision(bid *models.Bid, username, decision string) error
	UpdateBid(bid *models.Bid, username string) error
	UpdateStatusBid(bid *models.Bid, username string) error
}

type useCase struct {
	bidRepository    bidRep.RepositoryI
	tenderRepository tenderRep.RepositoryI
	userRepository   userRep.RepositoryI
}

func New(bidRepository bidRep.RepositoryI, tenderRepository tenderRep.RepositoryI, userRepository userRep.RepositoryI) UseCaseI {
	return &useCase{
		bidRepository:    bidRepository,
		tenderRepository: tenderRepository,
		userRepository:   userRepository,
	}
}

func (uc *useCase) CreateBid(bid *models.Bid) error {
	err := uc.validateBid(bid)
	if err != nil {
		return err
	}
	existTender, err := uc.tenderRepository.SelectTenderById(bid.TenderId)
	if err != nil {
		return models.ErrTenderNotFound
	}

	if existTender.Status != models.PUBLISHEDTEN {
		return models.ErrUserNotPermission
	}

	_, err = uc.userRepository.SelectUserById(bid.AuthorID)
	if err != nil {
		return models.ErrUserInvalid
	}

	bid.Status = models.CREATEDBID
	bid.Version = 1
	bid.Created = strfmt.DateTime(time.Now())
	err = uc.bidRepository.CreateBid(bid)
	if err != nil {
		return err
	}
	bid.Description = ""
	bid.TenderId = ""
	return nil
}

// modify
func (uc *useCase) GetStatusBid(bidId, username string) (string, error) {
	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return "", models.ErrUserInvalid
	}

	existBid, err := uc.bidRepository.SelectBidById(bidId)
	if err != nil {
		return "", models.ErrBidNotFound
	}

	if existBid.AuthorID != user.Id {
		return "", models.ErrUserNotPermission
	}
	return existBid.Status, nil
}

func (uc *useCase) UpdateStatusBid(bid *models.Bid, username string) error {
	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return models.ErrUserInvalid
	}

	existBid, err := uc.bidRepository.SelectBidById(bid.Id)
	if err != nil {
		return models.ErrBidNotFound
	}

	if existBid.AuthorID != user.Id {
		return models.ErrUserNotPermission
	}

	bid.Version = existBid.Version + 1
	err = uc.bidRepository.UpdateStatusBid(*bid)
	if err != nil {
		return models.ErrBadData
	}
	bid.Name = existBid.Name
	bid.AuthorType = existBid.AuthorType
	bid.AuthorID = existBid.AuthorID
	bid.Version = existBid.Version
	return nil
}

func (uc *useCase) UpdateBid(bid *models.Bid, username string) error {
	err := uc.validateBid(bid)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return models.ErrUserInvalid
	}

	existBid, err := uc.bidRepository.SelectBidById(bid.Id)
	if err != nil {
		return models.ErrBidNotFound
	}

	if user.Id != existBid.AuthorID {
		return models.ErrUserNotPermission
	}

	bid.Version = existBid.Version + 1

	err = uc.bidRepository.UpdateBid(*bid)

	if err != nil {
		return models.ErrBadData
	}
	if bid.Name == "" {
		bid.Name = existBid.Name
	}
	if bid.AuthorType == "" {
		bid.AuthorType = existBid.AuthorType
	}
	if bid.Status == "" {
		bid.Status = existBid.Status
	}
	bid.AuthorID = existBid.AuthorID
	bid.Created = existBid.Created
	bid.TenderId = ""
	bid.Description = ""
	return nil
}

func (uc *useCase) SelectBidsByUsername(limit, offset int, username string) ([]*models.Bid, error) {
	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return nil, models.ErrUserInvalid
	}

	bids, err := uc.bidRepository.SelectBidsByUserId(limit, offset, user.Id)
	if err != nil {
		return nil, err
	}
	for i := range bids {
		bids[i].Description = ""
		bids[i].TenderId = ""
	}

	return bids, nil
}

func (uc *useCase) SelectBidsByTenderId(limit, offset int, username string, tenderId string) ([]*models.Bid, error) {

	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return nil, models.ErrUserInvalid
	}

	existTender, err := uc.tenderRepository.SelectTenderById(tenderId)
	if err != nil {
		return nil, models.ErrTenderNotFound
	}

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, existTender.OrganizationId) //проверить, что достаются нужные данные
	if err != nil || !exist {
		return nil, models.ErrUserNotPermission
	}

	bids, err := uc.bidRepository.SelectBidsByTederId(limit, offset, username, tenderId)
	if err != nil {
		return nil, err
	}
	for i := range bids {
		bids[i].Description = ""
		bids[i].TenderId = ""
	}

	return bids, nil
}

func (uc *useCase) SubmitDecision(bid *models.Bid, username, decision string) error {
	user, err := uc.userRepository.SelectUserByUsername(username)
	if err != nil {
		return models.ErrUserInvalid
	}

	existBid, err := uc.bidRepository.SelectBidById(bid.Id)
	if err != nil {
		return models.ErrBidNotFound
	}

	existTender, err := uc.tenderRepository.SelectTenderById(existBid.TenderId)
	if err != nil {
		return models.ErrTenderNotFound
	}

	exist, err := uc.userRepository.CheckUserOrganization(user.Id, existTender.OrganizationId)
	if err != nil || !exist {
		return models.ErrUserNotPermission
	}

	if decision == "Approved" {
		err := uc.tenderRepository.ClosedСonfirmedTender(existTender.Id)
		if err != nil {
			return err
		}
	}

	bid.Name = existBid.Name
	bid.Status = existBid.Name
	bid.AuthorType = existBid.AuthorType
	bid.AuthorID = existBid.AuthorID
	bid.Version = existBid.Version
	return nil

}

func (uc *useCase) validateBid(bid *models.Bid) error {
	if len([]rune(bid.Name)) > 100 {
		return models.ErrBadData
	}

	if len([]rune(bid.Description)) > 500 {
		return models.ErrBadData
	}

	if bid.AuthorType != "User" && bid.AuthorType != "Organization" && bid.AuthorType != "" {
		return models.ErrBadData
	}

	return nil
}

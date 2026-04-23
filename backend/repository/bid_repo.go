package repository

import (
	"nft-backend/models"

	"gorm.io/gorm"
)

type BidRepo struct {
	DB *gorm.DB
}

func NewBidRepo(db *gorm.DB) *BidRepo {
	return &BidRepo{DB: db}
}

func (r *BidRepo) CreateBid(b *models.Bid) error {
	return r.DB.Create(b).Error
}

func (r *BidRepo) GetBids(auctionID uint64) ([]models.Bid, error) {
	var bids []models.Bid

	err := r.DB.
		Where("auction_id = ?", auctionID).
		Order("amount DESC"). // highest first
		Find(&bids).Error

	return bids, err
}

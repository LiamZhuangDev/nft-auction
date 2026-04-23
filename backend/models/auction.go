package models

type Auction struct {
	ID          uint64 `gorm:"primaryKey"`
	Seller      string
	NftContract string
	TokenId     string
	StartPrice  string
	EndTime     uint64
	Active      bool
	Bids        []Bid `gorm:"foreignKey:AuctionId"` // it's not required, gorm can actually figure foreign key out if following the convention
}

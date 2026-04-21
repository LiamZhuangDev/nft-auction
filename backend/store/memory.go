package store

type Listing struct {
	ID          uint64
	Seller      string
	NftContract string
	TokenId     string
	Active      bool
}

var Listings = make(map[uint64]Listing)

type Auction struct {
	ID          uint64
	Seller      string
	NftContract string
	TokenId     string
	StartPrice  string
	EndTime     uint64
	Active      bool
}

var Auctions = make(map[uint64]Auction)

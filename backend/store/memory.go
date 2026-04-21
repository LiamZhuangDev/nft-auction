package store

type Listing struct {
	ID          uint64
	Seller      string
	NftContract string
	TokenId     string
	Active      bool
}

var Listings = make(map[uint64]Listing)

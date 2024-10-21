package models

type GarantexRates struct {
	Timestamp int32 `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Bid `json:"bids"`
}

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Bid struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type RatesToDB struct {
	Timestamp int32 `json:"timestamp" db:"timestamp"`
	Ask       Ask   `json:"asks" db:"asks"`
	Bid       Bid   `json:"bids" db:"bids"`
}

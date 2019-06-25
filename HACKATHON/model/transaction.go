package model

type Transaction struct {
	Uuid              string  `json:"uuid"`
	Description       string  `json:"description"`
	Concept           string  `json:"concept"`
	Mcc               string  `json:"mcc"`
	Type              string  `json:"type"`
	Amount            float32 `json:"amount"`
	Currency          string  `json:"currency"`
	AmountFee         float32 `json:"amount_fee"`
	CurrencyFee       string  `json:"currency_fee"`
	AmountInTransit   float32 `json:"amount_in_transit"`
	CurrencyInTransit string  `json:"currency_in_transit"`
	Date              string  `json:"date"`
	Cat               string  `json:"cat"`
	SubCat            string  `json:"subcat"`
	Folio             string  `json:"folio"`
}

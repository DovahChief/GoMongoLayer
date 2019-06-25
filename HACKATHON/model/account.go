package model

type Account struct {
	Clabe     string `json:"clabe"`
	Reference string `json:"reference"`
	Facade_id string `json:"facade_id"`
	Balance   float32 `json:"balance"`
	Card	 Card   `json:"card"`
}


type Card struct {
	Current Current   `json:"current"`
	Historical Current[]  `json:"hisorical"`
}

type Current struct {
	Number string `json:"number"`
	Expmmyy string `json:"exp_mmyy"`
	Code string `json:"code"`
	Status string `json:"status"`
	Adate string `json:"adate"` 
	Cdate string `json:"cdate"`
}
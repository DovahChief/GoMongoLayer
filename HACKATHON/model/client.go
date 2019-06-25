package model

type User struct {
	Uuid    string  `json:"uuid"`
	Fname   string  `json:"fname"`
	Gname   string  `json:"gname"`
	Gender  string  `json:"gender"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
	Bdate   string  `json:"bdate"`
	Cdate   string  `json:"cdate"`
	Address Address `json:"address"`
}

type Address struct {
	Shipping  Adr `json:"shipping"`
	Invoicing Adr `json:"invoicing"`
}

type Adr struct {
	Country        string `json:"country"`
	Locality       string `json:"locality"`
	Street         string `json:"street"`
	Municipality   string `json:"municipality"`
	Town           string `json:"town"`
	OutsideNumber  string `json:"outside_number"`
	InteriorNumber string `json:"interior_number"`
	Zip            string `json:"zip"`
}

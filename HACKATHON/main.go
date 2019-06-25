package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

var database *sql.DB
var redisClient *redis.Client

func main() {
	fmt.Println("Iniciando server")
	db, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/hackathon")
	defer db.Close()
	if err != nil {
		fmt.Println("Error al conectar con base de datos")
	}

	database = db
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisClient = client

	ruta := mux.NewRouter()
	ruta.HandleFunc("/albo/db/user/new", createuser).Methods("POST")
	ruta.HandleFunc("/albo/db/account/assign", createAccount).Methods("POST")
	ruta.HandleFunc("/albo/db/transaction", createTransaction).Methods("POST")
	ruta.HandleFunc("/albo/db/account/{uuid}", getAccount).Methods("GET")

	http.Handle("/", ruta)
	_ = http.ListenAndServe(":8081", nil)
}

func createTransaction(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Println("Guardando transaccion")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(request.Body)
	var body transaction
	err := decoder.Decode(&body)
	handleError(err)

	tx, err := database.Begin()
	handleError(err)
	amount := fmt.Sprintf("%f",body.User.Account.Transaction.Amount)
	amountFee := fmt.Sprintf("%f",body.User.Account.Transaction.AmountFee)
	amountInTransit := fmt.Sprintf("%f",body.User.Account.Transaction.AmountInTransit)
	txid:= genXid()
	sqlq := "INSERT INTO transactions set accountId= '" + body.User.UUID +
		"', description= '" + body.User.Account.Transaction.Description + "', concepto = '" +
		body.User.Account.Transaction.Concept + "', mcc ='" + body.User.Account.Transaction.Mcc + "', tipo= '" + body.User.Account.Transaction.Type +
		"', amount = '" + amount + "', divisa ='" + body.User.Account.Transaction.Currency +
		"', amountFee ='" + amountFee + "', currencyFee= '" + body.User.Account.Transaction.CurrencyFee +
		"', amountOnTransit = '" + amountInTransit + "', currencyinTransit ='" + body.User.Account.Transaction.CurrencyInTransit +
		"', date ='" + body.User.Account.Transaction.Date + "', cat= '" + body.User.Account.Transaction.Cat +
		"', subcat = '" + body.User.Account.Transaction.SubCat + "', folio ='" + body.User.Account.Transaction.Folio + "', txid='"+txid+"'"
	fmt.Println(sqlq)
	_, err3 := database.Exec(sqlq)
	if err3 != nil {
		_ = tx.Rollback()
		log.Fatal(err3)
	}

	queryAccount := "SELECT * FROM account WHERE user_uuid='" + body.User.UUID+"';"
	fmt.Println(queryAccount)
	filas := database.QueryRow(queryAccount)
	account := dbaccount{}
	_ = filas.Scan(&account.AccountUUID, &account.Balance, &account.UserUUID,
		&account.Clabe, &account.Reference, &account.FacadeID)

	var newamount float32
	if body.User.Account.Transaction.Type == "tx.ab" {
		newamount = account.Balance + body.User.Account.Transaction.Amount + body.User.Account.Transaction.AmountFee
	} else {
		newamount = account.Balance - body.User.Account.Transaction.Amount - body.User.Account.Transaction.AmountFee
	}
	newAmountStr := fmt.Sprintf("%f",newamount)
	strq := "UPDATE account SET balance = " + newAmountStr + " where user_uuid = '" + body.User.UUID+"';"
	fmt.Println(strq)
	_, err5 := database.Exec(strq)
	if err5 != nil {
		_ = tx.Rollback()
		log.Fatal(err5)
	}
	handleError(tx.Commit())

	out := `{
		"response" : "OK",
		"user" : {
		"uuid" : "` + body.User.UUID + `",
		"account" : {
		"transaction" : {
		"uuid" : "` + txid + `"
		}
		}
		}
		}`

	_, _ = fmt.Fprintf(writer, out)
}

func genXid() string {
	id := xid.New()
	return id.String()
}

func createuser(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Println("Guardando usuario")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(request.Body)
	var body userInput
	err := decoder.Decode(&body)
	handleError(err)

	u, err2  := json.Marshal(body)
	handleError(err2)
	uuid := genXid()
	fmt.Println(string(u))
	redisClient.Do("SET", uuid, u)
	responseBody:= insertUserResponse{Response: "OK", User: userResponse{UUID: uuid}}
	sal, _ := json.Marshal(responseBody)
	_, _ = fmt.Fprintf(writer, string(sal))

}

func createAccount(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Println("Guardando cuenta")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(request.Body)
	var body accountInput
	err := decoder.Decode(&body)
	handleError(err)

	uuid := genXid()

	dt := time.Now().String()

	tx, err := database.Begin()
	handleError(err)

	fil := filaAv {}

	filas, err34 := database.Query("SELECT clabe, referencia, ntarjeta FROM available WHERE flag = 0")
	handleError(err34)
	for filas.Next() {
		_ = filas.Scan(&fil.clabe, &fil.reference, &fil.number)
		break
	}

	_, err3 := database.Exec("UPDATE available set flag = 1 WHERE clabe = '" + fil.clabe + "';")
	fmt.Println(fil.clabe)
	if err3 != nil {
		_ = tx.Rollback()
		log.Fatal(err3)
	}
	handleError(tx.Commit())

	sqlq := "INSERT INTO account set account_uuid= '" + uuid +
		"', balance= '" + "10000.00" + "', user_uuid = '" +
		body.User.UUID + "', clabe ='" + fil.clabe +
		"', facade_id ='cacao', reference = '" + fil.reference + "'"

	q2, err2 := database.Exec(sqlq)
	handleError(err2)
	fmt.Println(q2)
	responseBody:= insertAccountResponse{
		Response: "OK",
		User:userOutputAccount{
			UUID:body.User.UUID,
			Account:accountOutputAccount{
				Balance:10000.0, Clabe: fil.clabe, Reference:fil.reference, FacadeID:"cacao",
				Card:cardOutputAccount{
					Current:currentOutputAccount{
						Adate: dt, Expmyy:"05/2023",Number:fil.number, Status:"activa"},
				}}}}

	a, err2  := json.Marshal(responseBody.User.Account)
	handleError(err2)
	redisClient.Do("SET", uuid, a)
	sal, _ := json.Marshal(responseBody)
	_, _ = fmt.Fprintf(writer, string(sal))

}

func getAccount(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Obteniendo cuentas")
	urlParams := mux.Vars(request)
	id := urlParams["uuid"]
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	queryAccount := "SELECT * FROM account WHERE user_uuid='" + id +"';"
	fmt.Println(queryAccount)
	filas := database.QueryRow(queryAccount)
	account := account{}
	var x string
	var y string
	_ = filas.Scan(&x, &account.Balance, &y,
		&account.Clabe, &account.Reference, &account.FacadeID)


	writer.Header().Set("Content-Type", "application/json")

	responseBody := getAccountResponse{Response:"OK", User:getAccountUser{Uuid:id, Account:account}}

	sal, _ := json.Marshal(responseBody)
	_, _ = fmt.Fprintf(writer, string(sal))
}

type getAccountResponse struct {
	Response string  `json:"response"`
	User getAccountUser `json:"user"`
}

type getAccountUser struct {
	Uuid string `json:"uuid"`
	Account account `json:"account"`
}

type userResponse struct {
	UUID string `json:"uuid"`
}

type insertUserResponse struct {
	Response string `json:"response"`
	User userResponse  `json:"user"`
}

type insertAccountResponse struct {
	Response string `json:"response"`
	User userOutputAccount `json:"user"`
}

type userOutputAccount struct {
	UUID string `json:"uuid"`
	Account accountOutputAccount `json:"account"`
}

type accountOutputAccount struct {
	Clabe string `json:"clabe"`
	Reference string `json:"reference"`
	FacadeID string `json:"facade_id"`
	Balance float32 `json:"balance"`
	Card cardOutputAccount `json:"card"`
}

type cardOutputAccount struct {
	Current currentOutputAccount `json:"current"`
}

type currentOutputAccount struct {
	Number string `json:"number"`
	Expmyy string `json:"exp_mmyy"`
	Status string `json:"status"`
	Adate string `json:"adate"`
}

type accounts struct {
	Accounts []dbaccount `json:"accounts"`
}

type dbaccount struct {
	AccountUUID string  `json:"account_uuid"`
	Balance     float32 `json:"balance"`
	UserUUID    string  `json:"user_uuid"`
	Clabe       string  `json:"clabe"`
	Reference   string  `json:"reference"`
	FacadeID    string  `json:"facade_id"`
}

type account struct {
	Clabe     string  `json:"clabe"`
	Reference string  `json:"reference"`
	FacadeID  string  `json:"facade_id"`
	Balance   float32 `json:"balance"`
}

type card struct {
	Current    current   `json:"current"`
	Historical []current `json:"historical"`
}

type current struct {
	Number  string `json:"number"`
	Expmmyy string `json:"exp_mmyy"`
	Code    string `json:"code"`
	Status  string `json:"status"`
	Adate   string `json:"adate"`
	Cdate   string `json:"cdate"`
}

type accountInput struct {
	User userInputAccount `json:"user"`
}

type userInputAccount struct {
	UUID string `json:"uuid"`
	Account accountInputUser `json:"account"`
}

type accountInputUser struct {
	Card cardInputUser `json:"card"`
}

type cardInputUser struct {
	Current currentInputUser `json:"current"`
}

type currentInputUser struct {
	Code string `json:"code"`
}

type userInput struct {
	User user `json:"user"`
}

type user struct {
	Fname   string  `json:"fname"`
	Gname   string  `json:"gname"`
	Gender  string  `json:"gender"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
	Bdate   string  `json:"bdate"`
	Address address `json:"address"`
}

type address struct {
	Shipping  adr `json:"shipping"`
}

type adr struct {
	Country        string `json:"country"`
	Locality       string `json:"locality"`
	Street         string `json:"street"`
	Municipality   string `json:"municipality"`
	Town           string `json:"town"`
	OutsideNumber  string `json:"outside_number"`
	InteriorNumber string `json:"interior_number"`
	Zip            string `json:"zip"`
}

type userInsideTransaction struct {
	UUID              string  `json:"uuid"`
	Account AccountInputTransaction `json:"account"`
}

type AccountInputTransaction struct {
	Transaction transactionInputTransaction `json:"transaction"`
}
type transactionInputTransaction struct {

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

type transaction struct {
	User userInsideTransaction `json:"user"`
}


func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type filaAv struct {
	 number string
	 clabe string
	 reference string
}

type badResponse struct {
	Response    string    `json:"response"`
	Detail    string `json:"detail"`
}
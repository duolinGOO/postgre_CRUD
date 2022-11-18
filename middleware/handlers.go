package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"

	"postgres-go/models"
	"strconv"
)

type response struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv(".env"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request. %v", err)
	}
	insertId := insertStock(stock)
	res := response{
		ID:      insertId,
		Message: "Succuessfuly created the stock",
	}
	json.NewEncoder(w).Encode(res)
}
func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	var id int64
	sqlStatement := `INSERT INTO stocks VALUES($1,$2,$3) RETURNING stockid`
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to insert the stock. %v", err)
	}
	return id
}
func GetStock(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}
func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	var stock models.Stock
	defer db.Close()
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned ")
	case nil:
		return stock, nil
	default:
		fmt.Println("Unable to scan the row")
	}
	return stock, err
}
func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStock()
	if err != nil {
		log.Fatalf("Unable to get all stocks. %v", err)
	}
	json.NewEncoder(w).Encode(stocks)
}
func getAllStock() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stocks []models.Stock
	sqlStatement := `INSERT * FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Printf("Unable to get all stocks. %v", err)
	}
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable scan the row. %v", err)
		}
		stocks = append(stocks, stock)

	}
	return stocks, err
}
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}
	deletedRows := deleteStock(int64(id))
	res := response{
		ID:      int64(id),
		Message: fmt.Sprintf("Stock deleted successfully. Total rows affected %v", deletedRows),
	}
	json.NewEncoder(w).Encode(res)
}
func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE * FROM stocks WHERE stockid=$1`
	row, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to delete the row. %v", err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows affected : %v", rowsAffected)
	return rowsAffected

}
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request . %v", err)
	}
	updatedRow := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully, Total rows affected %v", updatedRow)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}
func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE stocks SET name=$2,price=$3,company=$4 WHERE stockid=$1`
	row, err := db.Exec(sqlStatement, stock.StockId, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	rowAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows affected : %v", rowAffected)
	return rowAffected

}

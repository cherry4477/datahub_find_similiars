package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
   "net"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Ds struct {
	db *sql.DB
}

var (
	ds = new(Ds)
)

func init() {
	DB_ADDR := os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
	DB_PORT := os.Getenv("MYSQL_PORT_3306_TCP_PORT")
	DB_DATABASE := os.Getenv("MYSQL_ENV_MYSQL_DATABASE")
	DB_USER := os.Getenv("MYSQL_ENV_MYSQL_USER")
	DB_PASSWORD := os.Getenv("MYSQL_ENV_MYSQL_PASSWORD")
	DB_URL := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`, DB_USER, DB_PASSWORD, DB_ADDR, DB_PORT, DB_DATABASE)

	log.Println("connect to", DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Printf("error: %s\n", err)
	} else {
		ds.db = db
	}
}

//======================================================
//
//======================================================

type FindSimiliarService struct {
	
}

//======================================================
//
//======================================================

func onBuildSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(getBuildSimiliarDataItemsResult(r))
}

func onSearchSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(getSearchSimiliarDataItemsResult(r))
}

func onServiceError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonResult("Unsupported URL", nil))
}

//======================================================
//
//======================================================

func runBuildThread () {
   
}

//======================================================
//
//======================================================

var port = flag.Int("port", 6666, "server port")

func main() {
   go runBuildThread ()
   
	http.HandleFunc("/similiars/search", onSearchSimiliarDataItems) // ?data_item_id=n
	http.HandleFunc("/similiars/build", onBuildSimiliarDataItems)   // ?data_item_id=n
	http.HandleFunc("/", onServiceError)

	address := fmt.Sprintf(":%d", *port)
	log.Printf("Listening at %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil)) // will block here
}

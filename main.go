package main

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
)

type Result struct {
	
}

func buildSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	r.ParseForm ()
	itemIdString := r.FormValue ("data_item_id")
	itemId, err := strconv.Atoi (itemIdString)
	if err != nil {
		
		return
	}
	
	items, err := buildSimiliarDataItems (itemId)
	if err != nil {
		
		return
	}
	
	
}


func searchSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	r.ParseForm ()
	itemIdString := r.FormValue ("data_item_id")
	itemId, err := strconv.Atoi (itemIdString)
	if err != nil {
		
		return
	}
	
	items, err := searchSimiliarDataItems (itemId)
	if err != nil {
		
		return
	}
	
	
}

func main() {
	http.HandleFunc("/similiars/search", searchSimiliarDataItems)
	http.HandleFunc("/similiars/build", buildSimiliarDataItems)
	
	address := fmt.Sprintf (":%d", port)
	
	log.Fatal (http.ListenAndServe(address, nil)) // will block here
}

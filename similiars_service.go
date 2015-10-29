package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//======================================================
//
//======================================================

type Result struct {
	Similiars []*SimiliarDataItem `json:"similiars"`
}

//======================================================
//
//======================================================

func jsonResult(errorMessage string, items []*SimiliarDataItem) []byte {
	if errorMessage != "" {
		items = nil
	}

	if items == nil {
		items = []*SimiliarDataItem{}
	}

	result := Result{Similiars: items}

	data, err := json.Marshal(&result)
	if err != nil {
		return []byte(`{"result":false,"error":"json error"}`)
	} else {
		return data
	}
}

//======================================================
//
//======================================================

func getBuildSimiliarDataItemsResult(r *http.Request) []byte {
	r.ParseForm()
	itemIdString := r.FormValue("data_item_id")
	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		return jsonResult("Invalid data_item_id", nil)
	}

	err = buildSimiliarDataItems(itemId)
	if err != nil {
		return jsonResult(fmt.Sprintf("Build similiars error: %s", err), nil)
	}

	return jsonResult("", nil)
}

func getSearchSimiliarDataItemsResult(r *http.Request) []byte {
	r.ParseForm()
	itemIdString := r.FormValue("data_item_id")
	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		return jsonResult("Invalid data_item_id", nil)
	}

	items, err := searchSimiliarDataItems(itemId)
	if err != nil {
		return jsonResult(fmt.Sprintf("Search similiars error: %s", err), nil)
	}

	return jsonResult("", items)
}

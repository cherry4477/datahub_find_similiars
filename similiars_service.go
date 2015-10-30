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
	Ok        bool                `json:"ok"`
	Error     string              `json:"error,omitempty"`
	Similiars []*SimiliarDataItem `json:"similiars,omitempty"`
}

const JsonErrorMessage string = `{"ok":false,"error":"json error"}`

//======================================================
//
//======================================================

func jsonResult(errorMessage string, items []*SimiliarDataItem) []byte {
	result := &Result{}
	
	if errorMessage != "" {
		result.Ok = false
		result.Error = errorMessage
	} else {
		result.Ok = true
		result.Similiars = items
	}

	data, err := json.Marshal(&result)
	if err != nil {
		return []byte(JsonErrorMessage)
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

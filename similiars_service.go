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
	Ok        bool                `json:"result"` // true | false
	Err       string              `json:"error,omitempty"`
	Similiars []*SimiliarDataItem `json:"similiars,omitempty"`
}

//======================================================
//
//======================================================

func jsonResult(w http.ResponseWriter, errorMessage string, items []*SimiliarDataItem) {
	result := Result{}
	if errorMessage == "" {
		result.Ok = true
		result.Similiars = items
	} else {
		result.Ok = false
		result.Err = errorMessage
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data, err := json.Marshal(&result)
	if err != nil {
		w.Write([]byte(`{"result":false,"error":"json error"}`))
	} else {
		w.Write(data)
	}
}

//======================================================
//
//======================================================

func onBuildSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	itemIdString := r.FormValue("data_item_id")
	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		jsonResult(w, "Invalid data_item_id", nil)
		return
	}

	err = buildSimiliarDataItems(itemId)
	if err != nil {
		jsonResult(w, fmt.Sprintf("Build similiars error: %s", err), nil)
		return
	}

	jsonResult(w, "", nil)
}

func onSearchSimiliarDataItems(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	itemIdString := r.FormValue("data_item_id")
	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		jsonResult(w, "Invalid data_item_id", nil)
		return
	}

	items, err := searchSimiliarDataItems(itemId)
	if err != nil {
		jsonResult(w, fmt.Sprintf("Search similiars error: %s", err), nil)
		return
	}

	jsonResult(w, "", items)
}

func onServiceError(w http.ResponseWriter, r *http.Request) {
	jsonResult(w, "Unsupported url", nil)
}

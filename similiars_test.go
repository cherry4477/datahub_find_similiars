package main

import (
	"fmt"
	"testing"
)

//=====================================================
// temp
//=====================================================

var (
	SimpleDataItems = []*SimpleDataItem{
		&SimpleDataItem{
			id:            1,
			keywords:      "创业;大赛;网易;堵车",
			name:          "网易创业大赛",
			respositoryId: 100,
		},
		&SimpleDataItem{
			id:            2,
			keywords:      "跑马;大赛",
			name:          "世界跑马大赛",
			respositoryId: 101,
		},
		&SimpleDataItem{
			id:            3,
			keywords:      "西二旗;堵车",
			name:          "西二旗堵车",
			respositoryId: 101,
		},
		&SimpleDataItem{
			id:            4,
			keywords:      "创业;大赛;网易;堵车",
			name:          "网易创业大赛",
			respositoryId: 100,
		},
	}

	SimpleDataItemsMap = make(map[int32]*SimpleDataItem)
)

func init() {
	for _, item := range SimpleDataItems {
		SimpleDataItemsMap[int32(item.id)] = item
	}
}

func _retrieveSimpleDataItem(id int32) (*SimpleDataItem, error) {
	item := SimpleDataItemsMap[id]
	if item == nil {
		return nil, fmt.Errorf("Item (id=%s), is not found", id)
	}

	return item, nil
}

func _retrieveAllSimpleDataItems() ([]*SimpleDataItem, error) {
	//sql := `select DATAITEM_ID, REPOSITORY_ID, DATAITEM_NAME, KEY_WORDS from DH_DATAITEM

	return SimpleDataItems, nil
}

//=====================================================
//
//=====================================================

// set DB env before testing
// set DB=datahub:datahub@tcp(10.1.235.96:3306)/datahub?charset=utf8
// export DB=datahub:datahub@tcp(10.1.235.96:3306)/datahub?charset=utf8
func TestSearch(t *testing.T) {
	t.Logf("TestSearch, DB_URL = %s\n", DB_URL)
	items, err := searchSimiliarDataItems(1011)
	if err != nil {
		t.Errorf("TestSearch error: %s", err)
		return
	}

	t.Logf("TestSearch: number of result: %d\n", len(items))
	for _, item := range items {
		t.Logf("item.Dataitem_name = %s\n", item.Dataitem_name)
	}
}

// set DB env before testing
// set DB=datahub:datahub@tcp(10.1.235.96:3306)/datahub?charset=utf8
//func TestBuild (t *testing.T) {
//	err := buildSimiliarDataItems (1011)
//	if err != nil {
//		t.Errorf ("TestBuild error: %s", err)
//		return
//	}
//}

func TestGetSimiliars(t *testing.T) {
	item, err := _retrieveSimpleDataItem(1)
	if err != nil {
		t.Errorf("_retrieveSimpleDataItem err = %s\n", err)
		return
	}
	t.Logf("item.name = %s\n", item.name)

	allSimpleDataItems, err := _retrieveAllSimpleDataItems()
	if err != nil {
		t.Errorf("_retrieveAllSimpleDataItems err = %s\n", err)
		return
	}
	t.Logf("len (allSimpleDataItems) = %d\n", len(allSimpleDataItems))

	similiarSimpleDataItems := retrieveSimiliarSimpleDataItems(item, allSimpleDataItems, 30.0)
	if len(similiarSimpleDataItems) < 1 || similiarSimpleDataItems[0].id != 4 || similiarSimpleDataItems[0].score != 100.0 {
		t.Errorf("the first item id should be 4\n")
	}

	t.Logf("len (similiarSimpleDataItems) = %d\n", len(similiarSimpleDataItems))
	for i, item := range similiarSimpleDataItems {
		t.Logf("i=%d, item.name = %s\n", i, item.name)
	}
}

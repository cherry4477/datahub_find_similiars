package main

import (
	"sort"
	"strings"
	//"encoding/json"
	"fmt"
	"database/sql"
)

//=====================================================
// to rewrite
//=====================================================

func getDB () *sql.DB {
	return ds.db
}

//type DataItem struct ...

//=====================================================
// search
//=====================================================

func searchSimiliarDataItems (forItemId int) ([]*DataItem, error) {
	maxRows := 20
	sql := fmt.Sprintf (`select 
				a.DATAITEM_ID, a.SIMILAR_ID, a.SCORE
				from DH_SIMILAR a 
				WHERE a.DATAITEM_ID=%d 
				LIMIT %d 
				ORDER BY a.SCORE DESC
				`, forItemId, maxRows)
	rows, err := getDB ().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	ids := make ([]int, maxRows)
	numIds := 0
	for rows.Next() && numIds < maxRows {
		dummy := 0
		if err := rows.Scan(&dummy, &ids[numIds], &dummy); err != nil {
			return nil, err
		}
		numIds ++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	items := make ([]*DataItem, numIds)
	numItems := 0
	for i := 0; i < numIds; i ++ {
		item, err := retrieveDataItem (ids [i])
		if err != nil {
			items [numItems] = item
			numItems ++
		}
	}
	
	return items [:numItems], nil
}

func retrieveDataItem(id int) (*DataItem, error) {
	var item DataItem
	sql := fmt.Sprintf (`select 
				a.DATAITEM_ID, a.DATAITEM_NAME, a.KEY_WORDS, a.REPOSITORY_ID, 
				a.PRICE, a.ICO_NAME, a.USER_ID, a.PERMIT_TYPE 
				from DH_DATAITEM a 
				where a.DATAITEM_ID=%d
				`, id)
	err := getDB ().QueryRow(sql).Scan(
				&item.Dataitem_id, &item.Dataitem_name, &item.Key_words, &item.Repository_id,
				&item.Price, &item.Ico_name, &item.User_id, &item.Permit_type)
	if err != nil {
		return nil, err
	}
	
	return &item, nil
}

//=====================================================
// build
//=====================================================

func buildSimiliarDataItems (forItemId int) (error) {
	err := deleteSimiliarDataItems (forItemId)
	if err != nil {
		return err
	}
	
	forItem, err := retrieveSimpleDataItem (forItemId)
	if err != nil {
		return err
	}
	
	allItems, err := retrieveAllSimpleDataItems ()
	if err != nil {
		return err
	}
	
	similiarItems := retrieveSimiliarSimpleDataItems (forItem, allItems, 70.0)
	for _, item := range similiarItems {
		insertSimiliarDataItem (forItem.id, item.id, int(item.score))
	}
	
	return nil
}

func deleteSimiliarDataItems (forItemId int) error {
	sql := fmt.Sprintf ("DELETE FROM DH_SIMILAR a where a.DATAITEM_ID=%d", forItemId)
	result, err := getDB ().Exec(sql)
	if err != nil {
		return err
	}
	
	_, err = result.RowsAffected ()
	
	return nil
}

func insertSimiliarDataItem (forItemId int, similiarItemId int, score int) {
	sql := fmt.Sprintf ("INSERT INTO DH_SIMILAR (DATAITEM_ID, SIMILAR_ID, SCORE) VALUES (%d, %d, %d)", forItemId, similiarItemId, score)
	result, err := getDB ().Exec(sql)
	if err != nil {
		// ...
		return
	}
	
	_, err = result.LastInsertId ()
}

//=====================================================
// ... 
//=====================================================

type SimpleDataItem struct {
	id            int
	name          string
	keywords      string //
	respositoryId int
	
	score float64
}

func retrieveSimpleDataItem(id int) (*SimpleDataItem, error) {
	var item SimpleDataItem
	sql := fmt.Sprintf ("select a.DATAITEM_ID, a.DATAITEM_NAME, a.KEY_WORDS, a.REPOSITORY_ID from DH_DATAITEM a where a.DATAITEM_ID=%d", id)
	err := getDB ().QueryRow(sql).Scan(&item.id, &item.name, &item.keywords, &item.respositoryId)
	if err != nil {
		return nil, err
	}
	
	return &item, nil
}

func retrieveAllSimpleDataItems() ([]*SimpleDataItem, error) {
	maxRows := 100
	sql := fmt.Sprintf ("select a.DATAITEM_ID, a.DATAITEM_NAME, a.KEY_WORDS, a.REPOSITORY_ID from DH_DATAITEM a LIMIT %d", maxRows)
	rows, err := getDB ().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	items := make ([]*SimpleDataItem, maxRows)
	num := 0
	for rows.Next() && num < maxRows {
		item := &SimpleDataItem{}
		if err := rows.Scan(&item.id, &item.name, &item.keywords, &item.respositoryId); err != nil {
			return nil, err
		}
		
		items [num] = item
		num ++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items [:num], nil
}

//=====================================================
// ...
//=====================================================

type ItemStats struct {
	items  []*SimpleDataItem
	num    int
}
func (stat ItemStats) Len() int {
	return stat.num
}
func (stat ItemStats) Swap(i, j int) {
	stat.items[i], stat.items[j] = stat.items[j], stat.items[i]
}
func (stat ItemStats) Less(i, j int) bool {
	return stat.items[i].score > stat.items[j].score
}

//=====================================================
// ...
//=====================================================

type ParsedSimpleDataItem struct {
	dataItem            *SimpleDataItem
	splitedKeywords     []string
	splitedNameSegments []string
}

func retrieveSimiliarSimpleDataItems(dataItem *SimpleDataItem, allDataItems []*SimpleDataItem, minimumScore float64) []*SimpleDataItem {
	pdi := &ParsedSimpleDataItem{
		dataItem:            dataItem,
		splitedKeywords:     strings.Split(dataItem.keywords, ";"),
		splitedNameSegments: splitNameIntoSegments(dataItem.name),
	}

	numItems := len(allDataItems)

	itemStats := ItemStats{}
	itemStats.items = make([]*SimpleDataItem, numItems)
	itemStats.num = 0
	for i := 0; i < numItems; i++ {
		item := allDataItems[i]
		if item.id == dataItem.id {
			continue
		}

		item.score = compareDataItemSimilarityScore(pdi, item)
		//fmt.Printf ("data (%d): %f\n", item.id, score)
		if minimumScore >= 0.0 && item.score < minimumScore {
			continue
		}

		itemStats.items[itemStats.num] = item
		itemStats.num++
	}

	sort.Sort(itemStats)
	if itemStats.num <= 20 {
		return itemStats.items[:itemStats.num]
	} else {
		return itemStats.items[:20]
	}
}

func splitNameIntoSegments(name string) []string {
	runes := []rune(name)
	num := len(runes)
	segments := make([]string, num-1)
	for i := 0; i < num-1; i++ {
		segments[i] = string(runes[i : i+2])
	}
	return segments
}

func calculateKeywordsScore(weight float64, words1 []string, words2 []string) float64 {
	k := 0

	n := len(words1)
	m := len(words2)
	if n == 0 || m == 0 {
		return 0.0
	}

	for i := 0; i < n; i++ {
		word := words1[i]
		for j := 0; j < m; j++ {
			if strings.Compare(word, words2[j]) == 0 {
				k++
			}
		}
	}

	return weight * 2.0 * float64(k) / float64(m+n)
}

func compareDataItemSimilarityScore(di1 *ParsedSimpleDataItem, di2 *SimpleDataItem) float64 {
	// key_words的相似分数 = 70 * k * 2 /(m+n)
	// dataitem_name的相似分数 = 20 * k * 2 /(m+n)
	// repository_id相同，加20 
	score1 := calculateKeywordsScore(70.0, di1.splitedKeywords, strings.Split(di2.keywords, ";"))
	score2 := calculateKeywordsScore(20.0, di1.splitedNameSegments, splitNameIntoSegments(di2.name))
	score := score1 + score2
	if di1.dataItem.respositoryId == di2.respositoryId {
		score += 10.0
	}

	return score
}

package main

import (
	"errors"
	"sort"
	"strings"
   "math"
	//"encoding/json"
	"database/sql"
	"fmt"
)

//=====================================================
//
//=====================================================

const (
	DEBUG = true
   MinimumSimiliariry = 14.0
)

func dbError(err error, secureErrMessage string) error {
	if DEBUG { //
		return err
	} else {
		return errors.New(secureErrMessage)
	}
}

//=====================================================
// search
//=====================================================

type SimiliarDataItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Views        int    `json:"views"`
	Follows      int    `json:"follows"`
	Downloads    int    `json:"downloads"`
	Stars        int    `json:"stars"`
	Refresh_date string `json:"refresh_date"`
	Usability    int    `json:"usability"`
}

func searchSimiliarDataItems(forItemId int) ([]*SimiliarDataItem, error) {
	db := getDB()
	if db == nil {
		return nil, errors.New("db is not initilized yet")
	}

	maxRows := 20
	sql := fmt.Sprintf(`select 
				a.SIMILAR_ID
				from DH_SIMILAR a 
				WHERE a.DATAITEM_ID=%d 
				ORDER BY a.SCORE DESC
				LIMIT %d 
				`, forItemId, maxRows)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, dbError(err, "db query error")
	}
	defer rows.Close()

	ids := make([]int, maxRows)
	numIds := 0
	for rows.Next() && numIds < maxRows {
		if err := rows.Scan(&ids[numIds]); err != nil {
			return nil, dbError(err, "db query scan error")
		}
		numIds++
	}
	if err := rows.Err(); err != nil {
		return nil, dbError(err, "db query rows error")
	}

	items := make([]*SimiliarDataItem, numIds)
	numItems := 0
	for i := 0; i < numIds; i++ {
		item, err := retrieveSimiliarDataItem(db, ids[i])
		if err == nil {
			items[numItems] = item
			numItems++
		}
	}

	return items[:numItems], nil
}

func retrieveSimiliarDataItem(db *sql.DB, id int) (*SimiliarDataItem, error) {
	var item SimiliarDataItem
	sql := fmt.Sprintf(`select 
				a.DATAITEM_ID, a.DATAITEM_NAME, 
				a.VIEWS, a.FOLLOWS, a.DOWNLOADS, a.STARS, 
				a.REFRESH_DATE, a.USABILITY 
				from DH_DATAITEMUSAGE a 
				where a.DATAITEM_ID=%d
				`, id)
	err := db.QueryRow(sql).Scan(
		&item.ID, &item.Name,
		&item.Views, &item.Follows, &item.Downloads, &item.Stars,
		&item.Refresh_date, &item.Usability)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

//=====================================================
// build
//=====================================================

func buildSimiliarDataItems(forItemId int) error {
	db := getDB()
	if db == nil {
		return errors.New("db is not initilized yet")
	}

	err := deleteSimiliarDataItems(db, forItemId)
	if err != nil {
		return dbError(err, "db delete error")
	}

	forItem, err := retrieveSimpleDataItem(db, forItemId)
	if err != nil {
		return dbError(err, "db retrieve simple data item error")
	}

	allItems, err := retrieveAllSimpleDataItems(db)
	if err != nil {
		return dbError(err, "db retrieve items error")
	}

	similiarItems := findSimiliarSimpleDataItems(forItem, allItems, MinimumSimiliariry)
	for _, item := range similiarItems {
		insertSimiliarDataItem(db, forItem.id, item.id, int(math.Ceil (item.score)))
	}

	return nil
}

func deleteSimiliarDataItems(db *sql.DB, forItemId int) error {
	sql := fmt.Sprintf("DELETE FROM DH_SIMILAR where DATAITEM_ID=%d", forItemId)
	result, err := db.Exec(sql)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return nil
}

func insertSimiliarDataItem(db *sql.DB, forItemId int, similiarItemId int, score int) {
	sql := fmt.Sprintf("INSERT INTO DH_SIMILAR (DATAITEM_ID, SIMILAR_ID, SCORE) VALUES (%d, %d, %d)", forItemId, similiarItemId, score)
	result, err := db.Exec(sql)
	if err != nil {
		// ...
		return
	}

	_, err = result.LastInsertId()
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

func retrieveSimpleDataItem(db *sql.DB, id int) (*SimpleDataItem, error) {
	var item SimpleDataItem
	sql := fmt.Sprintf("select a.DATAITEM_ID, a.DATAITEM_NAME, a.KEY_WORDS, a.REPOSITORY_ID from DH_DATAITEM a where a.DATAITEM_ID=%d", id)
	err := db.QueryRow(sql).Scan(&item.id, &item.name, &item.keywords, &item.respositoryId)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func retrieveAllSimpleDataItems(db *sql.DB) ([]*SimpleDataItem, error) {
	maxRows := 100
	sql := fmt.Sprintf("select a.DATAITEM_ID, a.DATAITEM_NAME, a.KEY_WORDS, a.REPOSITORY_ID from DH_DATAITEM a LIMIT %d", maxRows)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*SimpleDataItem, maxRows)
	num := 0
	for rows.Next() && num < maxRows {
		item := &SimpleDataItem{}
		if err := rows.Scan(&item.id, &item.name, &item.keywords, &item.respositoryId); err != nil {
			return nil, err
		}

		items[num] = item
		num++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items[:num], nil
}

//=====================================================
// ...
//=====================================================

type ItemStats struct {
	items []*SimpleDataItem
	num   int
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

func findSimiliarSimpleDataItems(dataItem *SimpleDataItem, allDataItems []*SimpleDataItem, minimumScore float64) []*SimpleDataItem {
	pdi := &ParsedSimpleDataItem{
		dataItem:            dataItem,
		splitedKeywords:     splitKeywords(dataItem.keywords),
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

func splitKeywords(keywords string) []string {
	words := strings.Split(keywords, ";")
   num := len (words)
   index := 0
   for i := 0; i < num; i++ {
      if len (strings.TrimSpace (words[0])) > 0 {
         words [index] = words [i]
         index ++
      }
   }
   return words [:index]
}

func calculateKeywordsScore(weight float64, words1 []string, words2 []string) float64 {
	k := 0

	fmt.Printf("words1 = %s\n", strings.Join(words1, ","))
	fmt.Printf("words2 = %s\n", strings.Join(words1, ","))

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
   
   fmt.Printf("k = %d, m = %d, n = %d\n", k, m, n)

	return weight * 2.0 * float64(k) / float64(m+n)
}

func compareDataItemSimilarityScore(di1 *ParsedSimpleDataItem, di2 *SimpleDataItem) float64 {
	// key_words的相似分数 = 70 * k * 2 /(m+n)
	// dataitem_name的相似分数 = 20 * k * 2 /(m+n)
	// repository_id相同，加20
   
   fmt.Printf ("Words1=%s\n", di1.dataItem.keywords)
   fmt.Printf ("Words2=%s\n", di2.keywords)
   fmt.Printf ("name1=%s\n", di1.dataItem.name)
   fmt.Printf ("name2=%s\n", di2.name)
   
	score1 := calculateKeywordsScore(70.0, di1.splitedKeywords, splitKeywords(di2.keywords))
	score2 := calculateKeywordsScore(20.0, di1.splitedNameSegments, splitNameIntoSegments(di2.name))
	score := score1 + score2
	if di1.dataItem.respositoryId == di2.respositoryId {
		score += 10.0
	}

	fmt.Printf("score1 = %f, score2 = %f, score = %f \n", score1, score2, score)

	return score
}

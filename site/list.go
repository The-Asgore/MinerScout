package site

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

var recordFileTime int64

func init() {
	recordFileTime = time.Now().Unix()
}

var RWMutexChan = make(chan int, 1)

func GetSite() (<-chan []string, int) {
	siteChan := make(chan []string, 1)

	siteList, err := os.Open("top-1m.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(siteList *os.File) {
		err := siteList.Close()
		if err != nil {
		}
	}(siteList)

	siteListCSVReader := csv.NewReader(siteList)
	allContent, err := siteListCSVReader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer close(siteChan)
		for _, line := range allContent {
			siteChan <- line
		}
	}()

	return siteChan, len(allContent)
}

func RecordSiteResult(record []string) (errorMsg error) {
	RWMutexChan <- 1
	resultFile, errorMsg := os.OpenFile(fmt.Sprintf("site_result_%v.csv", recordFileTime), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer func(resultFile *os.File) {
		err := resultFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resultFile)
	resultWriter := csv.NewWriter(resultFile)
	errorMsg = resultWriter.Write(record)
	resultWriter.Flush()
	<-RWMutexChan
	return
}

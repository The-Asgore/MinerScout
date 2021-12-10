package main

import (
	"MinerScout/browser"
	"MinerScout/site"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"log"
	"sort"
)

var errorChan = make(chan error, 100)
var reminderChan = make(chan interface{}, 100)

func main() {
	var miningPoolCount = map[string]float64{
		"CoinHive":     0,
		"CrypotoNoter": 0,
		"NFWebMiner":   0,
		"JSECoin":      0,
		"Webmine":      0,
		"CryptoLoot":   0,
		"CoinImp":      0,
		"DeepMiner":    0,
		"Monerise":     0,
		"Coinhave":     0,
		"Cpufun":       0,
		"Minr":         0,
		"Mineralt":     0,
	}
	var siteScannedCounter = 0
	var siteRefreshChan = make(chan string, 20)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	text := textList()
	bar := barChart()
	cCurve := cpuCurve()
	mCurve := memCurve()
	progress := progressBar()
	pMessageBox := processedMessageBox()
	eMessageBox := errorMessageBox()
	uiRender := func() {
		ui.Render(text, bar, cCurve, mCurve, progress, pMessageBox, eMessageBox)
	}
	uiEvent := ui.PollEvents()

	siteChan, sitesNum := site.GetSite()
	progress.Label = fmt.Sprintf("%v | %v", 0, sitesNum)

	go func() {
	EndScan:
		for {
			select {
			case err, ok := <-errorChan:
				if ok {
					eMessageBox.Text = fmt.Sprintf("[%v](fg:red)", err.Error())
					uiRender()
				} else {
					break EndScan
				}
			case message := <-reminderChan:
				switch v := message.(type) {
				case browser.SiteUrl:
					tempSlice := append([]string{}, string(v))
					text.Rows = append(tempSlice, text.Rows[0:13]...)
					uiRender()
				case string:
					miningPoolCount[v] += 1
					bar.Data, bar.Labels = sortMapByValue(miningPoolCount)
					uiRender()
				case map[string]float64:
					if len(cCurve.Data[0]) < 56 {
						cCurve.Data[0] = append(cCurve.Data[0], v["cpu"])
						mCurve.Data[0] = append(mCurve.Data[0], v["mem"])
					} else {
						cCurve.Data[0] = append(cCurve.Data[0][1:55], v["cpu"])
						mCurve.Data[0] = append(mCurve.Data[0][1:55], v["mem"])
					}
					uiRender()
				}
			case site2BeRefresh := <-siteRefreshChan:
				progress.Percent = 100 * siteScannedCounter / sitesNum
				progress.Label = fmt.Sprintf("%v | %v", siteScannedCounter, sitesNum)
				pMessageBox.Text = site2BeRefresh
				uiRender()
			}
		}
	}()

	for topSite := range siteChan {
		select {
		case e := <-uiEvent:
			switch e.ID {
			case "<C-c>":
				return
			}
		default:
			siteScannedCounter += 1
			siteRefreshChan <- topSite[1]
			browser.StartScan(topSite[0], topSite[1], reminderChan, errorChan)
		}
	}

	close(errorChan)
	close(reminderChan)
	close(siteRefreshChan)
	return
}

type Pair struct {
	Key   string
	Value float64
}
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func sortMapByValue(m map[string]float64) (data []float64, lab []string) {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	for _, v := range p {
		data = append(data, v.Value)
		lab = append(lab, v.Key)
	}
	return
}

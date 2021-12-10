package browser

import (
	"MinerScout/collecter"
	"MinerScout/site"
	"fmt"
	"log"
	"time"
)

type BeforeChromeInit interface {
	BeforeChromeInit([]string, chan<- error) error
}

type AfterChromeInit interface {
	AfterChromeInit([]string, chan<- interface{}, chan<- error) error
}

type BeforePatternMatch interface {
	BeforePatternMatch([]string, chan<- error) error
}

type AfterPatternMatch interface {
	AfterPatternMatch(string, []string, chan<- interface{}, chan<- error) error
}

type FinalTask interface {
	FinalTask([]string, chan<- error)
}

type hookerFunc struct {
}

func (hf *hookerFunc) AfterChromeInit(resultRecord []string, reminderChan chan<- interface{}, errorChan chan<- error) (errorMsg error) {
	time.Sleep(1 * time.Second)
	if usage, errorMsg := collecter.GetCPUAndMemChrome(); errorMsg != nil {
		errorChan <- errorMsg
	} else {
		resultRecord[5] = fmt.Sprintf("%v", usage["cpu"])
		resultRecord[6] = fmt.Sprintf("%v", usage["mem"])
		reminderChan <- usage
	}

	return
}

func (hf *hookerFunc) AfterPatternMatch(matchResult string, resultRecord []string, reminderChan chan<- interface{}, errorChan chan<- error) (errorMsg error) {
	if matchResult != "" {
		resultRecord[2] = "true"
		resultRecord[3] = "true"
		resultRecord[4] = matchResult
		time.Sleep(4 * time.Second)
		if usage, errorMsg := collecter.GetCPUAndMemChrome(); errorMsg != nil {
			errorChan <- errorMsg
		} else {
			resultRecord[7] = fmt.Sprintf("%v", usage["cpu"])
			resultRecord[8] = fmt.Sprintf("%v", usage["mem"])
			resultRecord[9] = fmt.Sprintf("%v", usage["num"])
			reminderChan <- usage
			reminderChan <- matchResult
		}

		return
	} else {
		resultRecord[2] = "true"
		resultRecord[3] = "false"
		return
	}
}

func (hf *hookerFunc) FinalTask(resultRecord []string, errorChan chan<- error) {
	go func() {
		if err := site.RecordSiteResult(resultRecord); err != nil {
			errorChan <- err
			log.Fatal()
		}
	}()
}

package browser

import (
	"MinerScout/collecter"
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"regexp"
	"time"
)

func CreateChromeInitTask() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
	}
}

func CreateChromeTask(targetSite string, htmlContent *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(targetSite),
		chromedp.WaitReady("html"),
		chromedp.OuterHTML("document.querySelector(\"html\")", htmlContent, chromedp.ByJSPath),
	}
}

func MatchMiningKeywords(htmlContent2BeMatch *string) <-chan string {
	matchResultChan := make(chan string)
	miningKeywordsMap := make(map[string]*regexp.Regexp)
	miningKeywordsMap["CoinHive"] = regexp.MustCompile(`new CoinHive\.Anonymous|coinhive.com/lib/coinhive.min.js|authedmine.com/lib/`)
	miningKeywordsMap["CrypotoNoter"] = regexp.MustCompile(`minercry.pt/processor.js|\.User\(addr`)
	miningKeywordsMap["NFWebMiner"] = regexp.MustCompile(`new NFMiner|nfwebminer.com/lib/`)
	miningKeywordsMap["JSECoin"] = regexp.MustCompile(`load.jsecoin.com/load`)
	miningKeywordsMap["Webmine"] = regexp.MustCompile(`webmine.cz/miner`)
	miningKeywordsMap["CryptoLoot"] = regexp.MustCompile(`CRLT\.anonymous|webmine.pro/lib/crlt.js`)
	miningKeywordsMap["CoinImp"] = regexp.MustCompile(`www.coinimp.com/scripts|new CoinImp.Anonymous|new Client.Anonymous|freecontent.stream|freecontent.data|freecontent.date`)
	miningKeywordsMap["DeepMiner"] = regexp.MustCompile(`new deepMiner.Anonymous | deepMiner.js`)
	miningKeywordsMap["Monerise"] = regexp.MustCompile(`apin.monerise.com | monerise_builder`)
	miningKeywordsMap["Coinhave"] = regexp.MustCompile(`minescripts\.info’`)
	miningKeywordsMap["Cpufun"] = regexp.MustCompile(`snipli.com/[A-Za-z]+" data-id=’`)
	miningKeywordsMap["Minr"] = regexp.MustCompile(`abc\.pema\.cl|metrika\.ron\.si|cdn\.rove\.cl|host\.dns\.ga|static\.hk\.rs|hallaert\.online|st\.kjli\.fi|minr\.pw|cnt\.statistic\.date|cdn\.static-cnt\.bid|ad\.g-content\.bid|cdn\.jquery-uim\.download’`)
	miningKeywordsMap["Mineralt"] = regexp.MustCompile(`ecart\.html\?bdata=|/amo\.js">|mepirtedic\.com’`)

	go func() {
		defer close(matchResultChan)
		for pattern, miningRegexp := range miningKeywordsMap {
			if miningRegexp.MatchString(*htmlContent2BeMatch) {
				matchResultChan <- pattern
				return
			}
		}
	}()

	return matchResultChan
}

func SetChromeCTX() (ctx context.Context, cancel context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.Flag("headless", false),
		//chromedp.ExecPath("C:\\Program Files (x86)\\Microsoft\\Edge Beta\\Application\\msedge.exe"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.20 Safari/537.36 Edg/97.0.1072.21"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel = chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	return
}

type SiteUrl string

func RunChromeTasks(ctx context.Context, siteNo string, targetSite string, hookFunc interface{}, reminderChan chan<- interface{}, errorChan chan<- error) <-chan int {
	var htmlContent string
	var statusChan = make(chan int)

	resultRecord := make([]string, 10)
	resultRecord[0] = siteNo
	resultRecord[1] = targetSite
	resultRecord[2] = "false"

	go func() {
		if t, ok := hookFunc.(FinalTask); ok {
			defer func() {
				close(statusChan)
				t.FinalTask(resultRecord, errorChan)
			}()
		} else {
			defer close(statusChan)
		}

		if t, ok := hookFunc.(BeforeChromeInit); ok {
			if err := t.BeforeChromeInit(resultRecord, errorChan); err != nil {
				errorChan <- err
				return
			}
		}

		if err := chromedp.Run(ctx, CreateChromeInitTask()); err != nil {
			errorChan <- err
			return
		} else if t, ok := hookFunc.(AfterChromeInit); ok {
			if err := t.AfterChromeInit(resultRecord, reminderChan, errorChan); err != nil {
				errorChan <- err
				return
			}
		}

		schemedUrl := targetSite
		parsedUrl, err := url.Parse(targetSite)
		if err != nil {
			errorChan <- err
			return
		} else if parsedUrl.Scheme == "" {
			schemedUrl = "https://" + targetSite
		}
		if err := chromedp.Run(ctx, CreateChromeTask(schemedUrl, &htmlContent)); err != nil {
			schemedUrl = "http://" + targetSite
			if err := chromedp.Run(ctx, CreateChromeTask(schemedUrl, &htmlContent)); err != nil {
				errorChan <- err
				return
			} else if t, ok := hookFunc.(BeforePatternMatch); ok {
				if err := t.BeforePatternMatch(resultRecord, errorChan); err != nil {
					errorChan <- err
					return
				}
			}
		} else if t, ok := hookFunc.(BeforePatternMatch); ok {
			if err := t.BeforePatternMatch(resultRecord, errorChan); err != nil {
				errorChan <- err
				return
			}
		}

		matchResult := <-MatchMiningKeywords(&htmlContent)
		if matchResult != "" {
			reminderChan <- SiteUrl(schemedUrl)
			if t, ok := hookFunc.(AfterPatternMatch); ok {
				if err := t.AfterPatternMatch(matchResult, resultRecord, reminderChan, errorChan); err != nil {
					errorChan <- err
					return
				}
			}
		} else {
			return
		}
	}()
	return statusChan
}

func StartScan(siteNo string, targetSite string, reminderChan chan<- interface{}, errorChan chan<- error) {
	defer func() {
		for ii := 0; ii < 3; ii++ {
			if collecter.IsChromeClosed() {
				return
			} else {
				time.Sleep(1 * time.Second)
			}
		}
		log.Fatalln("cannot kill chrome process")
	}()
	ctx, cancel := SetChromeCTX()
	defer func(ctx context.Context) {
		err := chromedp.Cancel(ctx)
		if err != nil {
		}
	}(ctx)
	defer cancel()

	statusChan := RunChromeTasks(ctx, siteNo, targetSite, &hookerFunc{}, reminderChan, errorChan)
	select {
	case <-statusChan:
	case <-time.After(30 * time.Second):
	}
}

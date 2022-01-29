package PageActions

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"log"
	"strconv"
	"time"
)

func scrollPage(ctx context.Context) error {
	log.Println("Scrolling page to the bottom...")
	err := chromedp.Run(ctx,
		chromedp.SendKeys(`a`, kb.PageDown+kb.PageDown, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`a`, kb.PageDown+kb.PageDown, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`a`, kb.End, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		return errors.New("error while scrolling the page")
	}

	return nil
}

func GetPageData(
	ctx context.Context,
	page string,
	totalPages *int,
) ([]*cdp.Node, []*cdp.Node, []*cdp.Node, error) {
	var names []*cdp.Node
	var prices []*cdp.Node
	var links []*cdp.Node
	var pageControllerTotal string
	var pageControllerCurrent string

	BaseUrl := "https://www.shopee.com.br/search?keyword="

	log.Println("Going to the requested page...")
	// Go to requested page
	errVisit := chromedp.Run(ctx,
		chromedp.Navigate(BaseUrl+page),
		chromedp.Sleep(5*time.Second),
	)
	if errVisit != nil {
		return []*cdp.Node{}, []*cdp.Node{}, []*cdp.Node{}, errVisit
	}

	// Scroll page
	errScroll := scrollPage(ctx)
	if errScroll != nil {
		return []*cdp.Node{}, []*cdp.Node{}, []*cdp.Node{}, errScroll
	}

	log.Println("Getting data from the page...")
	// Get page data
	errGet := chromedp.Run(ctx,
		chromedp.Nodes(`._10Wbs-`, &names, chromedp.ByQueryAll),
		chromedp.Nodes(".zp9xm9 > span:nth-child(2)", &prices, chromedp.ByQueryAll),
		chromedp.Nodes(".col-xs-2-4 > a", &links, chromedp.ByQueryAll),
		chromedp.Text(`.shopee-mini-page-controller__total`, &pageControllerTotal, chromedp.NodeVisible),
		chromedp.Text(`.shopee-mini-page-controller__current`, &pageControllerCurrent, chromedp.NodeVisible),
	)

	if errGet != nil {
		return []*cdp.Node{}, []*cdp.Node{}, []*cdp.Node{}, errGet
	}

	// Define total pages with page controller
	if *totalPages == 1 {
		*totalPages, _ = strconv.Atoi(pageControllerTotal)
	}

	log.Printf("Got %v registries from page %v of %v total pages", len(names), pageControllerCurrent, *totalPages)

	return names, prices, links, nil
}

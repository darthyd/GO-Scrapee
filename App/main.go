package App

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/darthyd/go-webscrapper-shopee/App/PageActions"
	"github.com/darthyd/go-webscrapper-shopee/App/Structs"
	"github.com/darthyd/go-webscrapper-shopee/App/Utilities"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Scrap struct {
	Query         string
	RequiredQuery []string
	MaxPrice      float64
	Pages         int
	Offset        int
}

func createContext() (context.Context, context.CancelFunc, context.CancelFunc) {
	// opts
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("window-size", fmt.Sprintf("%d,%d", 1920, 1080)),
	)

	// create allocator with options
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)

	// create context with allocator
	ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return ctx, cancelAlloc, cancelCtx
}

func requestPageData(
	ctx *context.Context,
	s Scrap,
	i *int,
	totalPages *int,
) ([]*cdp.Node, []*cdp.Node, []*cdp.Node) {
	fmt.Println("")
	log.Println("Requesting a page...")
	names, prices, links, err := PageActions.GetPageData(
		*ctx, s.Query+"&order=asc&sortBy=price&page="+fmt.Sprintf("%d", *i), totalPages)
	if err != nil {
		log.Fatal(err.Error())
	}
	return names, prices, links
}

func filterByRequires(
	reqs *[]string,
	max *float64,
	name *cdp.Node,
	price *cdp.Node,
	link *cdp.Node,
) (Structs.Product, error) {

	nameStr := name.Children[0].NodeValue

	priceStr := price.Children[0].NodeValue
	dotPrice := strings.Replace(priceStr, ".", "", 1)
	dotPrice = strings.Replace(dotPrice, ",", ".", -1)

	parsedPrice, err := strconv.ParseFloat(dotPrice, 64)
	if err != nil {
		return Structs.Product{}, err
	}

	for _, req := range *reqs {
		if !strings.Contains(strings.ToLower(nameStr), strings.ToLower(req)) {
			return Structs.Product{}, nil
		}
	}

	if parsedPrice >= *max && *max != 0 {
		return Structs.Product{}, nil
	}

	return Structs.Product{
		Name:  nameStr,
		Price: parsedPrice,
		Link:  "https://www.shopee.com.br" + link.AttributeValue("href"),
	}, nil
}

func (s Scrap) Scrapee() ([]byte, error) {

	if s.Pages == 0 {
		s.Pages = 1
	}

	ctx, cancelAlloc, cancelContext := createContext()
	defer cancelAlloc()
	defer cancelContext()

	ProdList := Structs.ProductList{}
	totalPages := s.Pages + s.Offset

	for i := s.Offset; i < totalPages; i++ {
		names, prices, links := requestPageData(&ctx, s, &i, &totalPages)

		var m int
		m, err := Utilities.SliceMin([]int{len(names), len(prices)})
		if err != nil {
			m = len(names)
		}

		for i := 0; i < m; i++ {
			product, err := filterByRequires(
				&s.RequiredQuery,
				&s.MaxPrice,
				names[i],
				prices[i],
				links[i],
			)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			if product.Name != "" {
				ProdList.Products = append(ProdList.Products, product)
			}
		}
	}

	fmt.Println("")
	log.Printf("Finished: Got a total of %v registries\n", len(ProdList.Products))

	sort.Slice(ProdList.Products, func(i, j int) bool {
		return ProdList.Products[i].Price < ProdList.Products[j].Price
	})

	return json.Marshal(ProdList.Products)
}

func RequestScrap(q Scrap) ([]byte, error) {
	return q.Scrapee()
}

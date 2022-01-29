package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/darthyd/go-webscrapper-shopee/PageActions"
	"github.com/darthyd/go-webscrapper-shopee/Structs"
	"github.com/darthyd/go-webscrapper-shopee/Utilities"
	"log"
)

func main() {
	// opts
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("window-size", fmt.Sprintf("%d,%d", 1920, 1080)),
	)

	// create allocator with options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context with allocator
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var keySearch string = "patriots tom brady"
	ProdList := Structs.ProductList{}

	totalPages := 1
	for i := 0; i < totalPages; i++ {
		fmt.Println("")
		log.Println("Requesting a page...")
		names, prices, links, err := PageActions.GetPageData(
			ctx, keySearch+"&page="+fmt.Sprintf("%d", i), &totalPages)
		if err != nil {
			log.Fatal(err.Error())
		}

		var m int
		m, err = Utilities.SliceMin([]int{len(names), len(prices)})
		if err != nil {
			m = len(names)
		}

		for i := 0; i < m; i++ {
			p := Structs.Product{
				Name:  names[i].Children[0].NodeValue,
				Price: prices[i].Children[0].NodeValue,
				Link:  "https://www.shopee.com.br" + links[i].AttributeValue("href"),
			}
			ProdList.Products = append(ProdList.Products, p)
		}
	}

	//j, err := json.Marshal(ProdList)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Println(string(j))
	fmt.Println("")
	log.Printf("Finished: Got a total of %v registries\n", len(ProdList.Products))
}

package Structs

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

type ProductList struct {
	Products []Product
}

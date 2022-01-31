package Server

import (
	"github.com/darthyd/go-webscrapper-shopee/App"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	Port string
}

func (s Server) rootHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	pagesStr := r.URL.Query().Get("p")
	pages, err := strconv.Atoi(pagesStr)
	if err != nil {
		pages = 1
	}
	offsetStr := r.URL.Query().Get("o")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	maxPriceStr := r.URL.Query().Get("m")
	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		maxPrice = 0
	}
	requiredStr := r.URL.Query().Get("r")
	required := strings.Split(requiredStr, " ")

	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := App.RequestScrap(App.Scrap{
		Query:         query,
		RequiredQuery: required,
		Pages:         pages,
		Offset:        offset,
		MaxPrice:      maxPrice,
	})
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) Serve() {
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/search", s.searchHandler)
	log.Fatal(http.ListenAndServe(s.Port, nil))
}

func UpServer(port string) {
	s := Server{Port: ":" + port}
	s.Serve()
}

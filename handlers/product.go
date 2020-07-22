package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/subhamsriv/microservices-in-go/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Product")

	lp := data.GetProducts()

	err := lp.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Product) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Unable to convert to string", http.StatusBadRequest)
	}

	p.l.Println("Handle Put Product")

	prod := &data.Product{}

	err = prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(prod, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

}

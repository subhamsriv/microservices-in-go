package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/subhamsriv/microservices-in-go/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProduct(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		//pass id in uri
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Printf("id is %v", id)
		p.updateProduct(w, r, id)
		return
	}
	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Product")

	lp := data.GetProducts()

	err := lp.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Product) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Product) updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	p.l.Println("Handle Put Product")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(prod, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

}

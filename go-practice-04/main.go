package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "PS5 Slim", Price: 1000, Stock: 5},
	{ID: 2, Name: "Keyboard", Price: 120, Stock: 10},
	{ID: 3, Name: "Mouse", Price: 50, Stock: 20},
}

func main() {

	http.HandleFunc("/products", getProductsHandler)
	http.HandleFunc("/product", getProductByIDHandler)
	http.HandleFunc("/health", getHealthHandler)
	http.HandleFunc("/products/count", getProductsCountHandler)
	http.HandleFunc("/products/create", createProductHandler)

	fmt.Println("Server runinng on :8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("server error:", err)
	}

}

func getProductsCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count := len(products)
	response := map[string]int{"count": count}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, "failed to encode products", http.StatusInternalServerError)
		return
	}

}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	var newProd Product
	err := json.NewDecoder(r.Body).Decode(&newProd)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	
	if newProd.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	for _, product := range products{
		if product.Name == newProd.Name {
			http.Error(w, "Product with this name already exists", http.StatusConflict)
			return
		}
	}

	if newProd.Price <= 0 {
		http.Error(w, "Price must be greater than zero", http.StatusBadRequest)
		return
	}

	
	newProd.ID = len(products) + 1

	
	products = append(products, newProd)

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201
	json.NewEncoder(w).Encode(newProd)
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	product, err := findProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, "failed to encode product", http.StatusInternalServerError)
		return
	}
}

func findProductByID(id int) (*Product, error) {

	for i := range products {
		if products[i].ID == id {
			return &products[i], nil
		}
	}
	return nil, fmt.Errorf("product with id %d not found", id)
}

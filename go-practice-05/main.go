package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var db *sql.DB

func main() {
	connStr := "host=localhost port=5435 user=postgres password=postgres dbname=go_products sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open database:", err)
	} 

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
	log.Println("Connected to PostgreSQL successfully!")

	http.HandleFunc("/health", getHealthHandler)
	http.HandleFunc("/products", getProductsHandler)
	http.HandleFunc("/product", getProductByIDHandler)

	log.Println("Server running on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("server error:", err)
	}

}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query the database for products
	rows, err := db.Query("SELECT id, name, price, stock FROM products ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to query products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "error while reading products", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "failed to encode products", http.StatusInternalServerError)
		return
	}

}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	var product Product

	err = db.QueryRow(
		"SELECT id, name, price, stock FROM products WHERE id = $1",
		id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "failed to encode product", http.StatusInternalServerError)
		return
	}
}

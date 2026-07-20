package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	connStr := "host=localhost port=5435 user=postgres password=admin123 dbname=go-postgres sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open database:", err)
	} 

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
	log.Println("Connected to PostgreSQL successfully!")

	http.HandleFunc("GET /health", getHealthHandler)
	http.HandleFunc("GET /products", getProductsHandler)
	http.HandleFunc("GET /product", getProductByIDHandler)
	http.HandleFunc("POST /product", createProductHandler)
	http.HandleFunc("DELETE /product", deleteProductHandler)
	http.HandleFunc("PUT /product", updateProductHandler)

	log.Println("Server running on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("server error:", err)
	}

}

func createProductHandler(w http.ResponseWriter,r *http.Request){
	// 1. Just Post method is valid
	if r.Method != http.MethodPost {
		w.Header().Set("Allow",http.MethodPost)
		http.Error(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		return
	}

	// 2 . Decode Json input
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil{
		http.Error(w,"Invalid JSON input",http.StatusBadRequest)
		return
	}

	// 3. Validation
    if newProduct.Name == "" || newProduct.Price <= 0 || newProduct.Stock < 0 {
        http.Error(w, "Name, Price > 0 and Stock >= 0 are required", http.StatusBadRequest)
        return
    }

	// 4 . Insert into db and retrieve Id
	err := db.QueryRow(
        "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
        newProduct.Name, newProduct.Price, newProduct.Stock,
    ).Scan(&newProduct.ID)


	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
            http.Error(w, "Product name already exists", http.StatusConflict) // 409
            return
        }
        // اینجا بعداً خطای Duplicate Name را هندل می‌کنیم
        http.Error(w, "Failed to create product"+err.Error(), http.StatusInternalServerError)
        return
    }

	// 5 Return product created Id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)

}

func updateProductHandler(w http.ResponseWriter,r *http.Request){
	// 1 : Check http method
	if r.Method != http.MethodPut {
		w.Header().Set("Allow",http.MethodPut)
		http.Error(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		return
	}
	
	// 2 : Check product with id exist
	idParam := r.URL.Query().Get("id")
	id ,err := strconv.Atoi(idParam)
	if err != nil{
		http.Error(w,"Invalid product ID",http.StatusBadRequest)
		return
	}

	// 3 : Read product detail from request body
	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct) ; err != nil{
		http.Error(w,"Invalid JSON body",http.StatusBadRequest)
		return
	}

	// 4 : Validation input
	if updatedProduct.Name == "" || updatedProduct.Price <= 0 || updatedProduct.Stock < 0 {
		http.Error(w,"Invalid product data : Name,Price > 0 and stock >= 0 are required",http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4",
		updatedProduct.Name, updatedProduct.Price, updatedProduct.Stock, id,
	)

	if err != nil{
		if strings.Contains(err.Error(),"unique constraint"){
			http.Error(w,"Product name allready exists",http.StatusConflict)// 409
			return
		}

		http.Error(w,"Failed to update product: "+err.Error(),http.StatusInternalServerError)
		return
	}


	// 6 : Check product with id exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound) // 404
		return
	}

	// 7 : Success response with updated info
	updatedProduct.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)

}

func deleteProductHandler(w http.ResponseWriter,r *http.Request){
	// 1 : Check Http method
	if r.Method != http.MethodDelete{
		w.Header().Set("Allow",http.MethodDelete)
		http.Error(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		return
	}


	// 2 : Check id exists in query param
	idParam := r.URL.Query().Get("id")
	if idParam == ""{
		http.Error(w,"Missing product ID",http.StatusBadRequest)
		return
	}

	// 3 : Check id from query param format
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w,"Invalid product ID format",http.StatusBadRequest)
		return
	}

	// 4 : Execute delete product query
	result,err := db.Exec("DELETE FROM products WHERE id = $1",id)
	if err != nil{
		http.Error(w,"Failed to delete product: "+err.Error(),http.StatusInternalServerError)
		return
	}
	
	// 5 : Check record with this id exists
	rowAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w,"Failed to check affected rows",http.StatusInternalServerError)
		return
	}

	if rowAffected == 0 {
		http.Error(w,"Product not found",http.StatusNotFound)
		return
	}

	// 6 : response answer
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
	"message": "Product deleted successfully",
	})
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

package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	products := []Product{
		{ID: 1, Name: "PS5 Slim", Price: 1000, Stock: 5},
		{ID: 2, Name: "Keyboard", Price: 120, Stock: 10},
		{ID: 3, Name: "Mouse", Price: 50, Stock: 20},
	}

	for _, value := range products {
		fmt.Println(value.GetInfo())
	}

	foundProduct, err := findProductByID(products, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Found:", foundProduct.GetInfo())
	}

	productMap := map[int]Product{
		1: {ID: 1, Name: "PS5 Slim", Price: 1000, Stock: 5},
		2: {ID: 2, Name: "Keyboard", Price: 120, Stock: 10},
		3: {ID: 3, Name: "Mouse", Price: 50, Stock: 20},
	}

	product := productMap[2]
	fmt.Println("Product info via map is :", product.GetInfo())

	
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	fmt.Println(string(data))
}

func findProductByID(products []Product, id int) (*Product, error) {
	for i := range products {
		if products[i].ID == id {
			return &products[i], nil
		}
	}
	return nil, fmt.Errorf("product with Id %d not found", id)
}

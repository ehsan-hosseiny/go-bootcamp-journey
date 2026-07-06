package main

import "fmt"

type Product struct {
	ID    int `json:"id"`
	Name  string `json:"name"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}

func (p *Product) GetInfo() string {
	return fmt.Sprintf("Name: %s, Price: %d, Stock: %d", p.Name, p.Price, p.Stock)
}



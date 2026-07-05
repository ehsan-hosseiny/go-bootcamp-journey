package main

import (
	"errors"
	"fmt"
)

type Product struct {
	ID    int
	Name  string
	Price int
	Stock int
}

type Displayable interface {
	GetInfo() string
}

func (p *Product) Sell(quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity should more than zero")
	}
	if p.Stock < quantity {
		return fmt.Errorf("Insufficient Stock: requested %d,but only %d available",quantity,p.Stock)
	}

	p.Stock -= quantity
	return nil
}

func (p *Product) UpdatePrice(newPrice int) error {
	if newPrice <= 0 {
		return errors.New("The new price must be valid.")
	}
	p.Price = newPrice
	return nil
}

func (p *Product) Restock(stock int) error {
	if stock <= 0 {
		return errors.New("stock must be greater than zero")
	}
	p.Stock += stock
	return nil
}

func (p *Product) GetInfo() string {
	return fmt.Sprintf("Name is: %s, Price is: %d, Stock is: %d", p.Name, p.Price, p.Stock)
}

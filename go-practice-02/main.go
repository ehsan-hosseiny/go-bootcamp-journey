package main

import "fmt"

func main() {
	// Create new product
	p := Product{
		ID:    1,
		Name:  "PS5 Slim",
		Price: 1000,
		Stock: 5,
	}

	fmt.Printf("Initial Product: %v\n", p)

	err := p.Sell(2)
	if err != nil {
		fmt.Println("Error selling:", err)
	} else {
		fmt.Println("Sold 2 units successfully!")
	}

	err = p.Sell(10)
	if err != nil {
		fmt.Printf("Expected Error: %s\n", err)
	}

	err = p.UpdatePrice(-10)
	if err != nil {
		fmt.Printf("Price update Error: %s\n", err)
		
	}

	productInfo(&p)

	fmt.Printf("Final Product State: %+v\n", p)
}

func productInfo(d Displayable) {
	fmt.Printf("Product Info: %s\n", d.GetInfo())
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func processOrder(orderId int, wg *sync.WaitGroup, resultChan chan string) {
	defer wg.Done()
	fmt.Printf("Processing order %d\n", orderId)
	time.Sleep(1 * time.Second)

	resultChan <- fmt.Sprintf("Order %d processed successfully", orderId)
}

func main() {
	var wg sync.WaitGroup
	resultChan := make(chan string)

	orders := []int{101,102,103,104}

	for _,orderID := range orders{
		wg.Add(1)
		go processOrder(orderID,&wg,resultChan)

	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println(result)
	}

	fmt.Println("All orders processed")

}

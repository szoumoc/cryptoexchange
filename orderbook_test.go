package main

import (
	"fmt"
	"testing"
)

// func TestOrderBook(t *testing.T) {
// 	// Create an order book
// 	orderBook := OrderBook{
// 		Bids: &[]Limit{},
// 		Asks: &[]Limit{},
// 	}

// func TestLimit(t *testing.T) {
// 	l := NewLimit(10_000)
// 	buyOrderA := NewOrder(100, true)
// 	buyOrderB := NewOrder(200, true)
// 	buyOrderC := NewOrder(300, true)
// 	l.AddOrder(buyOrderA)
// 	l.AddOrder(buyOrderB)
// 	l.AddOrder(buyOrderC)

// 	l.DeleteOrder(buyOrderB)
// 	fmt.Printf("Limit: %v\n", l)

// }
func TestOrderBook(t *testing.T) {
	ob := NewOrderBook()
	buyOrderA := NewOrder(100, true)
	buyOrderB := NewOrder(200, true)

	ob.PlaceOrder(10_000, buyOrderA)
	ob.PlaceOrder(18_000, buyOrderB)

	fmt.Printf("%+v", ob)
}

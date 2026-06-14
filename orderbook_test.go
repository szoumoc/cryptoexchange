package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %v, got %v", b, a)
	}
}
func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(100, true)
	buyOrderB := NewOrder(200, true)
	buyOrderC := NewOrder(300, true)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)
	fmt.Printf("Limit: %v\n", l)

}
func TestOrderBook(t *testing.T) {
	ob := NewOrderBook()
	sellOrder := NewOrder(100, false)

	ob.PlaceOrderLimit(18_000, sellOrder)
	ob.PlaceOrderLimit(18_000, sellOrder)
	assert(t, len(ob.asks), 1)
}
func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()
	sellOrderA := NewOrder(100, false)
	ob.PlaceOrderLimit(20_000, sellOrderA)

	buyOrder := NewOrder(200, true)
	matches := ob.PlaceMarketOrder(buyOrder)
	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10)
	fmt.Printf("%+v", matches)
}

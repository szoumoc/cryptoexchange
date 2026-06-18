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
	sellOrderA := NewOrder(300, false)
	ob.PlaceOrderLimit(20_000, sellOrderA)

	buyOrder := NewOrder(200, true)
	matches := ob.PlaceMarketOrder(buyOrder)
	fmt.Printf("%+v", matches)
	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 100.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 200.0)
	assert(t, matches[0].Price, 20_000.0)
	assert(t, buyOrder.isFilled(), true)
	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultipleMatches(t *testing.T) {
	ob := NewOrderBook()
	sellOrderA := NewOrder(5, true)
	sellOrderB := NewOrder(8, true)
	sellOrderC := NewOrder(10, true)
	ob.PlaceOrderLimit(10_000, sellOrderA)
	ob.PlaceOrderLimit(5_000, sellOrderC)
	ob.PlaceOrderLimit(20_000, sellOrderB)

	buyOrder := NewOrder(20, false)
	matches := ob.PlaceMarketOrder(buyOrder)
	fmt.Printf("%+v", matches)
	assert(t, len(matches), 3)
	// assert(t, len(ob.asks), 0)
	assert(t, len(ob.bids), 1)
	assert(t, ob.BidTotalVolume(), 3.0)
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrderA := NewOrder(5, true)
	buyOrderB := NewOrder(8, true)
	buyOrderC := NewOrder(10, true)
	ob.PlaceOrderLimit(10_000, buyOrderA)
	ob.PlaceOrderLimit(5_000, buyOrderC)
	ob.PlaceOrderLimit(20_000, buyOrderB)
	fmt.Printf("%+v", ob.bids)
	assert(t, len(ob.bids), 3)
	ob.cancelOrder(buyOrderB)
	assert(t, len(ob.bids), 2)
	assert(t, ob.BidTotalVolume(), 15.0)
}

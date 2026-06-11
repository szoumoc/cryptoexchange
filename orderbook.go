package main

import (
	"fmt"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}
type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

func NewOrder(size float64, bid bool) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("Size: %f", o.Size)
}

type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

type OrderBook struct {
	Bids       []*Limit
	Asks       []*Limit
	BidLimits  map[float64]*Limit
	AsksLimits map[float64]*Limit
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size

}

func (l *Limit) String() string {
	return fmt.Sprintf("Price: %f, TotalVolume: %f", l.Price, l.TotalVolume)
}

func (l *Limit) DeleteOrder(o *Order) {
	for i, order := range l.Orders {
		if order == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			break
		}
	}
	l.TotalVolume -= o.Size
	o.Limit = nil
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:       []*Limit{},
		Asks:       []*Limit{},
		BidLimits:  make(map[float64]*Limit),
		AsksLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceOrder(price float64, o *Order) []Match {
	//1. try to match the order
	// matching logic here

	//2. if not matched, add to the order book
	if o.Size > 0.0 {
		ob.add(price, o)
	}
	return []Match{} // Placeholder return, replace with actual matches
}

func (ob *OrderBook) add(price float64, o *Order) {
	var limit *Limit
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AsksLimits[price]
	}
	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(o)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AsksLimits[price] = limit
		}
	}
}

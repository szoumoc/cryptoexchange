package main

import (
	"fmt"
	"sort"
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

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

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
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestBid struct {
	Limits
}

func (a ByBestBid) Len() int           { return len(a.Limits) }
func (a ByBestBid) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestBid) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestAsk struct {
	Limits
}

func (b ByBestAsk) Len() int           { return len(b.Limits) }
func (b ByBestAsk) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestAsk) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

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
	sort.Sort(l.Orders)
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

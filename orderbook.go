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

func (l *Limit) Fill(o *Order) []Match {
	var (
		matches        []Match
		ordersToDelete []*Order
	)

	for _, order := range l.Orders {
		match := l.filledOrder(order, o)
		matches = append(matches, match)

		l.TotalVolume -= match.SizeFilled
		if order.isFilled() {
			ordersToDelete = append(ordersToDelete, order)
		}
		if o.isFilled() {
			break
		}
	}
	for _, order := range ordersToDelete {
		l.DeleteOrder(order)
	}
	return matches
}

func (o *Order) isFilled() bool {
	return o.Size == 0.0
}

func (l *Limit) filledOrder(a, b *Order) Match {
	var ask, bid *Order
	var sizeFilled float64
	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}
	if bid.Size > ask.Size {
		sizeFilled = ask.Size
		bid.Size -= sizeFilled
		ask.Size = 0
	} else {
		sizeFilled = bid.Size
		ask.Size -= sizeFilled
		bid.Size = 0
	}
	return Match{
		Bid:        bid,
		Ask:        ask,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}
}

type OrderBook struct {
	bids       []*Limit
	asks       []*Limit
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
		bids:       []*Limit{},
		asks:       []*Limit{},
		BidLimits:  make(map[float64]*Limit),
		AsksLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceMarketOrder(o *Order) []Match {
	matches := []Match{}
	if o.Bid {
		if o.Size > ob.AskTotalVolume() {
			panic("Not enough liquidity to fill the order")
		}
		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)
		}
	} else {
		if o.Size > ob.BidTotalVolume() {
			panic("Not enough liquidity to fill the order")
		}
		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)
		}
	}
	return matches
}

func (ob *OrderBook) PlaceOrderLimit(price float64, o *Order) {
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
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AsksLimits[price] = limit
		}
	}
}

func (ob *OrderBook) BidTotalVolume() float64 {
	total := 0.0
	for i := 0; i < len(ob.bids); i++ {
		total += ob.bids[i].TotalVolume
	}
	return total
}
func (ob *OrderBook) AskTotalVolume() float64 {
	total := 0.0
	for i := 0; i < len(ob.asks); i++ {
		total += ob.asks[i].TotalVolume
	}
	return total
}
func (ob *OrderBook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *OrderBook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}

package main

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/szoumoc/matchingengine/orderbook"
)

func main() {
	e := echo.New()
	ex := NewExchange()
	e.POST("/order", ex.handlerFunc)
	e.Start(":8080")
}

type Market string

const (
	BTCUSD Market = "BTCUSD"
	ETHUSD Market = "ETHUSD"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.OrderBook
}

func NewExchange() *Exchange {
	return &Exchange{
		orderbooks: make(map[Market]*orderbook.OrderBook),
	}
}

type OrderType string

const (
	LimitOrder  OrderType = "LIMIT"
	MarketOrder OrderType = "MARKET"
)

type PlaceOrderRequest struct {
	Type  OrderType
	Bid   bool
	Size  float64
	Price float64
}

func (ex *Exchange) handlerFunc(c echo.Context) error {
	var req PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return err
	}
	return c.String(200, "Hello, World!")
}

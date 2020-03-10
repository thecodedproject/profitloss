package profitloss

import (
	"github.com/shopspring/decimal"
)

type CalcType int

const (
	CalcTypeAverage = 0
	CalcTypeUnknown = 1
	calcTypeSentinal = 2
)

type OrderType string

const (
	OrderTypeBid = "BID"
	OrderTypeAsk = "ASK"
)

type CompletedOrder struct {
	Price decimal.Decimal
	Volume decimal.Decimal
	Type OrderType
}

type Report struct {
	Type CalcType

	RealisedGain decimal.Decimal

	AverageBuyPrice decimal.Decimal
	AverageSellPrice decimal.Decimal

	BaseBought decimal.Decimal
	BaseSold decimal.Decimal
	CounterBought decimal.Decimal
	CounterSold decimal.Decimal

	BaseBalance decimal.Decimal
	CounterBalance decimal.Decimal

	TotalVolume decimal.Decimal
	OrderCount int64

}

func Add(r Report, orders ...CompletedOrder) Report {

	for _, o := range orders {
		orderCost := o.Volume.Mul(o.Price)

		if o.Type == OrderTypeBid {
			r.BaseBought = r.BaseBought.Add(o.Volume)
			r.CounterSold = r.CounterSold.Add(orderCost)
			r.AverageBuyPrice = r.CounterSold.Div(r.BaseBought)

			r.BaseBalance = r.BaseBalance.Add(o.Volume)
			r.CounterBalance = r.CounterBalance.Sub(orderCost)
		}	else {
			r.CounterBought = r.CounterBought.Add(orderCost)
			r.BaseSold = r.BaseSold.Add(o.Volume)
			r.AverageSellPrice = r.CounterBought.Div(r.BaseSold)

			r.BaseBalance = r.BaseBalance.Sub(o.Volume)
			r.CounterBalance = r.CounterBalance.Add(orderCost)
		}

		volumeForRealisedGain := decimal.Min(r.BaseBought, r.BaseSold)
		r.RealisedGain = r.AverageSellPrice.Sub(r.AverageBuyPrice).Mul(volumeForRealisedGain)

		r.TotalVolume = r.TotalVolume.Add(o.Volume)
		r.OrderCount++
	}

	return r
}


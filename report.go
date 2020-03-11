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
	BaseFee decimal.Decimal
	CounterFee decimal.Decimal
	Type OrderType
}

type Report struct {
	Type CalcType `json:"type"`

	RealisedGain decimal.Decimal `json:"realised_gain"`

	AverageBuyPrice decimal.Decimal `json:"average_buy_price"`
	AverageSellPrice decimal.Decimal `json:"average_sell_price"`

	BaseBought decimal.Decimal `json:"base_bought"`
	BaseSold decimal.Decimal `json:"base_sold"`
	BaseFees decimal.Decimal `json:"base_fees"`
	CounterBought decimal.Decimal `json:"counter_bought"`
	CounterSold decimal.Decimal `json:"counter_sold"`
	CounterFees decimal.Decimal `json:"counter_fees"`

	BaseBalance decimal.Decimal `json:"base_balance"`
	CounterBalance decimal.Decimal `json:"counter_balance"`

	TotalVolume decimal.Decimal `json:"total_volume"`
	OrderCount int64 `json:"order_count"`

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

		r.CounterFees = r.CounterFees.Add(o.CounterFee)
		r.CounterBalance = r.CounterBalance.Sub(o.CounterFee)

		r.RealisedGain = r.AverageSellPrice.Sub(r.AverageBuyPrice).Mul(volumeForRealisedGain).Sub(r.CounterFees)

		r.BaseFees = r.BaseFees.Add(o.BaseFee)
		r.TotalVolume = r.TotalVolume.Add(o.Volume)
		r.OrderCount++
	}

	return r
}


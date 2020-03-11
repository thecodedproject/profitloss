package profitloss_test

import (
	"fmt"
	"github.com/thecodedproject/profitloss"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	"testing"
)

const (
	Bid = profitloss.OrderTypeBid
	Ask = profitloss.OrderTypeAsk
)

func assertDecimalsEqual(t *testing.T, expected, actual decimal.Decimal, i ...interface{}) {

	initialMessage := fmt.Sprint(i...)
	if len(initialMessage) != 0 {
		initialMessage = initialMessage + " "
	}
	assert.Equalf(t, 0, expected.Cmp(actual), "%sExpected: %s  Actual: %s", initialMessage, expected, actual)
}

func assertReportsEqual(t *testing.T, e, a profitloss.Report) {

	assertDecimalsEqual(t, e.RealisedGain, a.RealisedGain, "RealisedGain")

	assertDecimalsEqual(t, e.AverageBuyPrice, a.AverageBuyPrice, "AverageBuyPrice")
	assertDecimalsEqual(t, e.AverageSellPrice, a.AverageSellPrice, "AverageSellPrice")

	assertDecimalsEqual(t, e.BaseBought, a.BaseBought, "BaseBought")
	assertDecimalsEqual(t, e.BaseSold, a.BaseSold, "BaseSold")
	assertDecimalsEqual(t, e.BaseFees, a.BaseFees, "BaseFees")
	assertDecimalsEqual(t, e.CounterBought, a.CounterBought, "CounterBought")
	assertDecimalsEqual(t, e.CounterSold, a.CounterSold, "CounterSold")
	assertDecimalsEqual(t, e.CounterFees, a.CounterFees, "CounterFees")

	assertDecimalsEqual(t, e.BaseBalance, a.BaseBalance, "BaseBalance")
	assertDecimalsEqual(t, e.CounterBalance, a.CounterBalance, "CounterBalance")

	assertDecimalsEqual(t, e.TotalVolume, a.TotalVolume, "TotalVolume")
	assert.Equal(t, e.OrderCount, a.OrderCount)
}

func D(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func TestAveragePriceReport(t *testing.T) {

	testCases := []struct {
		Name string
		Inital profitloss.Report
		Orders []profitloss.CompletedOrder
		Expected profitloss.Report
	}{
		{
			Name: "No orders gives zero gain",
		},
		{
			Name: "Multiple buy and sell orders with more buy volume than sell volume uses avterage prices and total volume sold for realised gain",
			Orders: []profitloss.CompletedOrder{
				{
					Price: D(100.0),
					Volume: D(15.0),
					Type: Bid,
				},
				{
					Price: D(150.0),
					Volume: D(5.0),
					Type: Ask,
				},
				{
					Price: D(250.0),
					Volume: D(7.5),
					Type: Bid,
				},
				{
					Price: D(170.0),
					Volume: D(11.0),
					Type: Ask,
				},
			},
			Expected: profitloss.Report{
				RealisedGain: D(220.0),
				AverageBuyPrice: D(150.0),
				AverageSellPrice: D(163.75),
				BaseBought: D(22.5),
				BaseSold: D(16.0),
				CounterBought: D(2620.0),
				CounterSold: D(3375.0),
				BaseBalance: D(6.5),
				CounterBalance: D(-755.0),
				TotalVolume: D(38.5),
				OrderCount: 4,
			},
		},
		{
			Name: "Multiple buy and sell orders with more sell volume than buy volume uses average prices and total buy volume",
			Orders: []profitloss.CompletedOrder{
				{
					Price: D(150.0),
					Volume: D(15.0),
					Type: Ask,
				},
				{
					Price: D(200.0),
					Volume: D(25.0),
					Type: Ask,
				},
				{
					Price: D(175.0),
					Volume: D(20.0),
					Type: Bid,
				},
				{
					Price: D(160.0),
					Volume: D(10.0),
					Type: Bid,
				},
			},
			Expected: profitloss.Report{
				RealisedGain: D(337.5),
				AverageBuyPrice: D(170.0),
				AverageSellPrice: D(181.25),
				BaseBought: D(30.0),
				BaseSold: D(40.0),
				CounterBought: D(7250.0),
				CounterSold: D(5100.0),
				BaseBalance: D(-10.0),
				CounterBalance: D(2150.0),
				TotalVolume: D(70.0),
				OrderCount: 4,
			},
		},
		{
			Name: "Multiple buy and sell orders with some base fees remove base fees from base balance",
			Orders: []profitloss.CompletedOrder{
				{
					Price: D(100.0),
					Volume: D(5.0),
					BaseFee: D(1.1),
					Type: Bid,
				},
				{
					Price: D(200.0),
					Volume: D(15.0),
					BaseFee: D(2.2),
					Type: Bid,
				},
				{
					Price: D(120.0),
					Volume: D(6.0),
					BaseFee: D(3.3),
					Type: Ask,
				},
				{
					Price: D(140.0),
					Volume: D(9.0),
					BaseFee: D(4.4),
					Type: Ask,
				},
			},
			Expected: profitloss.Report{
				RealisedGain: D(-645),
				AverageBuyPrice: D(175.0),
				AverageSellPrice: D(132.0),
				BaseBought: D(20.0),
				BaseSold: D(15.0),
				BaseFees: D(11.0),
				CounterBought: D(1980.0),
				CounterSold: D(3500.0),
				BaseBalance: D(5.0),
				CounterBalance: D(-1520.0),
				TotalVolume: D(35.0),
				OrderCount: 4,
			},
		},
		{
			Name: "Multiple buy and sell orders with some counter fees remove base fees from counter balance and realised gain",
			Orders: []profitloss.CompletedOrder{
				{
					Price: D(200.26),
					Volume: D(5.0),
					CounterFee: D(1.1),
					Type: Bid,
				},
				{
					Price: D(140.0),
					Volume: D(5.0),
					CounterFee: D(2.2),
					Type: Bid,
				},
				{
					Price: D(120.0),
					Volume: D(8.5),
					CounterFee: D(3.3),
					Type: Ask,
				},
				{
					Price: D(160.0),
					Volume: D(1.5),
					CounterFee: D(4.4),
					Type: Ask,
				},
			},
			Expected: profitloss.Report{
				RealisedGain: D(-452.3),
				AverageBuyPrice: D(170.13),
				AverageSellPrice: D(126.0),
				BaseBought: D(10.0),
				BaseSold: D(10.0),
				CounterBought: D(1260.0),
				CounterSold: D(1701.3),
				CounterFees: D(11.0),
				BaseBalance: D(0.0),
				CounterBalance: D(-452.3),
				TotalVolume: D(20.0),
				OrderCount: 4,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			report := profitloss.Add(test.Inital, test.Orders...)
			assertReportsEqual(t, test.Expected, report)
		})
	}
}


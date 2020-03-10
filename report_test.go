package profitloss_test

import (
	"github.com/thecodedproject/profitloss"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	Bid = profitloss.OrderTypeBid
	Ask = profitloss.OrderTypeAsk
)

func assertDecimalsEqual(t *testing.T, expected, actual decimal.Decimal, i ...interface{}) {

	expectedF, ok := expected.Float64()
	require.True(t, ok)
	actualF, ok := actual.Float64()
	require.True(t, ok)
	assert.Equal(t, expectedF, actualF, i...)
}

func assertReportsEqual(t *testing.T, e, a profitloss.Report) {

	assertDecimalsEqual(t, e.RealisedGain, a.RealisedGain)

	assertDecimalsEqual(t, e.AverageBuyPrice, a.AverageBuyPrice)
	assertDecimalsEqual(t, e.AverageSellPrice, a.AverageSellPrice)

	assertDecimalsEqual(t, e.BaseBought, a.BaseBought)
	assertDecimalsEqual(t, e.BaseSold, a.BaseSold)
	assertDecimalsEqual(t, e.CounterBought, a.CounterBought)
	assertDecimalsEqual(t, e.CounterSold, a.CounterSold)

	assertDecimalsEqual(t, e.BaseBalance, a.BaseBalance)
	assertDecimalsEqual(t, e.CounterBalance, a.CounterBalance)

	assertDecimalsEqual(t, e.TotalVolume, a.TotalVolume)
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
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			report := profitloss.Add(test.Inital, test.Orders...)
			assertReportsEqual(t, test.Expected, report)
		})
	}
}


package tcsi_exporter

import (
	"context"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func NewClient(token string) (*sdk.RestClient, error) {
	log.Println("Setup TCS Investments client")
	client := sdk.NewRestClient(token)
	ctx := context.TODO()
	if _, err := client.Accounts(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

type PortfolioCollector struct {
	Client       *sdk.RestClient
	Descriptions PortfolioDesciptions
}

type PortfolioDesciptions struct {
	AverageSpend  *prometheus.Desc
	ItemsTotal    *prometheus.Desc
	ExceptedYield *prometheus.Desc
	MarketPrice   *prometheus.Desc
}

func NewPortfolioCollector(client *sdk.RestClient) (*PortfolioCollector, error) {
	labels := []string{"account_type", "instrument_type", "figi", "isin", "ticker", "human_name", "currency"}

	tcsiPositionAvgSpend := prometheus.NewDesc(
		"tcsi_position_spend_avg",
		"Average spend on position in instrument currency",
		labels,
		nil)

	tcsiPositionExpectedYield := prometheus.NewDesc(
		"tcsi_position_excepted_yield_total",
		"Execepted total yield of position at portfolio in instrument curency",
		labels,
		nil)

	tcsiPositionItemsTotal := prometheus.NewDesc(
		"tcsi_position_items_count",
		"Total items in position",
		labels,
		nil)

	tcsiPositionMarketPrice := prometheus.NewDesc(
		"tcsi_position_market_price_avg",
		"Market price of instrument",
		labels,
		nil)

	log.Println("Create PortfolioCollector")
	return &PortfolioCollector{
		Client: client,
		Descriptions: PortfolioDesciptions{
			AverageSpend:  tcsiPositionAvgSpend,
			ItemsTotal:    tcsiPositionItemsTotal,
			ExceptedYield: tcsiPositionExpectedYield,
			MarketPrice:   tcsiPositionMarketPrice,
		},
	}, nil
}

func (c PortfolioCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Descriptions.AverageSpend
	ch <- c.Descriptions.ExceptedYield
	ch <- c.Descriptions.ItemsTotal
	ch <- c.Descriptions.MarketPrice
}

func (c PortfolioCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.TODO()
	accounts, err := c.Client.Accounts(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, account := range accounts {
		portfolio, err := c.Client.Portfolio(ctx, account.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, position := range portfolio.Positions {
			labels := []string{
				string(account.Type),            // account_type
				string(position.InstrumentType), // instrument_type
				position.FIGI,                   //figi
				position.ISIN,                   //isin
				position.Ticker,                 //ticker
				position.Name,                   //human_name
				string(position.AveragePositionPrice.Currency), //currency
			}

			ch <- prometheus.MustNewConstMetric(
				c.Descriptions.AverageSpend,
				prometheus.GaugeValue,
				position.AveragePositionPrice.Value,
				labels...)

			ch <- prometheus.MustNewConstMetric(
				c.Descriptions.ExceptedYield,
				prometheus.GaugeValue,
				position.ExpectedYield.Value,
				labels...)

			ch <- prometheus.MustNewConstMetric(
				c.Descriptions.ItemsTotal,
				prometheus.GaugeValue,
				position.Balance,
				labels...)

			ch <- prometheus.MustNewConstMetric(
				c.Descriptions.MarketPrice,
				prometheus.GaugeValue,
				position.AveragePositionPrice.Value+position.ExpectedYield.Value/position.Balance,
				labels...)
		}
	}
}

type CashCollector struct {
	Client      *sdk.RestClient
	Description *prometheus.Desc
}

func NewCashCollector(client *sdk.RestClient) (*CashCollector, error) {
	labels := []string{"account_type", "currency"}

	tcsiCurrencyBalance := prometheus.NewDesc(
		"tcsi_currency_balance",
		"Free currency at portfolio",
		labels,
		nil)

	log.Println("Create CashCollector")
	return &CashCollector{
		Client:      client,
		Description: tcsiCurrencyBalance,
	}, nil
}

func (c CashCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Description
}

func (c CashCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.TODO()
	accounts, err := c.Client.Accounts(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, account := range accounts {
		portfolio, err := c.Client.Portfolio(ctx, account.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, cur := range portfolio.Currencies {
			labels := []string{
				string(account.Type), // account_type
				string(cur.Currency), // currency
			}

			ch <- prometheus.MustNewConstMetric(
				c.Description,
				prometheus.GaugeValue,
				cur.Balance,
				labels...)
		}
	}
}

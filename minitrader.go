package forexbot

import (
	"errors"
	"log"
	"time"
)

type Minitrader struct {
	Epic                 string
	InvestmentPercentage float64
	Timeframe            Timeframe
	Status               MinitraderStatus
	MarketStatus         MinitraderMarketStatus
	Strategy             Strategy

	capitalClient       *CapitalClientAPI
	candlesChannel      chan Candles // TODO: Implement "Pipeline" Pattern To Handle Larger Data Efficiently
	activeDealReference string
}

type MinitraderStatus string

const (
	// healthy statuses
	INITIALIZING      MinitraderStatus = "INITIALIZING"
	RUNNING           MinitraderStatus = "RUNNING"
	HOLDING           MinitraderStatus = "HOLDING"
	SELL_ORDER_ACTIVE MinitraderStatus = "SELL_ORDER_ACTIVE"
	BUY_ORDER_ACTIVE  MinitraderStatus = "BUY_ORDER_ACTIVE"

	// error statuses
	ERROR_ON_UPDATE_CANDLES_DATA MinitraderStatus = "ERROR_ON_UPDATE_CANDLES_DATA"
	ERROR_ON_MAKING_ORDER        MinitraderStatus = "ERROR_ON_MAKING_ORDER"
)

type MinitraderMarketStatus string

const (
	TRADEABLE MinitraderMarketStatus = "TRADEABLE"
	CLOSED    MinitraderMarketStatus = "CLOSED"
)

type Timeframe string

const (
	MINUTE    Timeframe = "MINUTE"
	MINUTE_5  Timeframe = "MINUTE_5"
	MINUTE_15 Timeframe = "MINUTE_15"
	MINUTE_30 Timeframe = "MINUTE_30"
	HOUR      Timeframe = "HOUR"
	HOUR_4    Timeframe = "HOUR_4"
	DAY       Timeframe = "DAY"
	WEEK      Timeframe = "WEEK"
)

func NewMinitrader(epic string, investmentPercentage float64, timeframe Timeframe, strategy Strategy) *Minitrader {
	return &Minitrader{
		Epic:                 epic,
		InvestmentPercentage: investmentPercentage,
		Timeframe:            timeframe,
		Strategy:             strategy,
		Status:               INITIALIZING,
		candlesChannel:       make(chan Candles),
	}
}

func (minitrader *Minitrader) Start() {
	for candles := range minitrader.candlesChannel {
		signal, price := minitrader.Strategy(candles)
		err := minitrader.Effect(signal, price)
		if err != nil {
			minitrader.Status = ERROR_ON_MAKING_ORDER
			close(minitrader.candlesChannel)
			log.Printf("Error On Minitrader Effect: %v", err)
			return
		}
	}
}

func (minitrader *Minitrader) Effect(signal Signal, price float64) error {
	if minitrader.MarketStatus == CLOSED {
		return errors.New("Unable To Do Trading; Market Closed")
	}

	// make a buy/sell order and wait 3:30 minutes or less if order has been completed before wait time.
	if minitrader.Status == RUNNING && signal == BUY {
		minitrader.Status = BUY_ORDER_ACTIVE
		err := minitrader.makeOrderAndWaitUntilComplete(minitrader.Epic, BUY, LIMIT, price)
		if err != nil {
			return err
		}
	} else if minitrader.Status == HOLDING && signal == SELL {
		minitrader.Status = SELL_ORDER_ACTIVE
		err := minitrader.makeOrderAndWaitUntilComplete(minitrader.Epic, SELL, LIMIT, price)
		if err != nil {
			return err
		}
	}

	return nil
}

func (minitrader *Minitrader) makeOrderAndWaitUntilComplete(epic string, signal Signal, orderType OrderType, targetPrice float64) error {
	var amount float64
	var err error
	var dealReference string

	// get amount available from preferred account or get amount from position/order confirmation
	if signal == BUY {
		amount, err = minitrader.getAmountFromPreferredAccount() // TODO: Need to return a good quantity base on price and available usd
	} else {
		amount, err = minitrader.getAmountFromPositionOrderConfirmation()
	}
	if err != nil {
		return err
	}

	// create a working order and retry if it fails
	orderResponse, err := minitrader.createWorkingOrderWithRetries(epic, signal, orderType, targetPrice, amount)
	if err != nil {
		return err
	}
	dealReference = orderResponse.DealReference

	// check if the working order status it was successfully completed
	confirmationStatus, err := minitrader.waitUntilConfirmationWithRetries(dealReference)
	if err != nil {
		return err
	}
	if confirmationStatus == string(DELETED) {
		minitrader.Status = RUNNING
		return nil
	}

	// update minitrader status
	if signal == BUY {
		minitrader.Status = HOLDING
		minitrader.activeDealReference = dealReference
	} else {
		minitrader.Status = RUNNING
	}

	return nil
}

func (minitrader *Minitrader) getAmountFromPreferredAccount() (float64, error) {
	var amount float64
	tryCounter := 0
	for tryCounter < 3 {
		accountsResponse, err := minitrader.capitalClient.GetAllAccounts()
		if err != nil {
			tryCounter++
			time.Sleep(time.Second * 5)
			continue
		}

		for _, account := range accountsResponse.Accounts {
			if account.Preferred {
				amount = account.Balance.Available
				break
			}
		}
		break
	}
	if tryCounter == 3 {
		return 0, errors.New("Failed to get amount from preferred account")
	}
	return amount, nil
}

func (minitrader *Minitrader) getAmountFromPositionOrderConfirmation() (amount float64, err error) {
	tryCounter := 0
	for tryCounter < 3 {
		positionOrderResponse, err := minitrader.capitalClient.GetPositionOrderConfirmation(minitrader.activeDealReference)
		if err != nil {
			tryCounter++
			time.Sleep(time.Second * 5)
			continue
		}

		amount = positionOrderResponse.Size
		break
	}
	if tryCounter == 3 {
		return 0, err
	}
	return amount, nil
}

func (minitrader *Minitrader) createWorkingOrderWithRetries(epic string, signal Signal, orderType OrderType, targetPrice float64, amount float64) (workingOrderResponse WorkingOrderResponse, err error) {
	tryCounter := 0
	for tryCounter < 3 {
		workingOrderResponse, err = minitrader.capitalClient.CreateWorkingOrder(epic, signal, orderType, targetPrice, amount)
		if err != nil {
			tryCounter++
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	if tryCounter == 3 {
		return workingOrderResponse, err
	}
	return workingOrderResponse, nil
}

func (minitrader *Minitrader) waitUntilConfirmationWithRetries(dealReference string) (orderPositionStatus string, err error) {
	tryCounter := 0
	for tryCounter < 3 {
		confirmationResponse, err := minitrader.capitalClient.GetPositionOrderConfirmation(dealReference)
		if err != nil {
			tryCounter++
			time.Sleep(time.Second * 5)
			continue
		}

		orderPositionStatus = confirmationResponse.Status
		break
	}
	if tryCounter == 3 {
		return "", err
	}
	return orderPositionStatus, nil
}

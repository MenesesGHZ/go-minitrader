package gominitrader

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Minitrader struct {
	Epic                 string
	Timeframe            Timeframe
	Status               MinitraderStatus
	MarketStatus         MinitraderMarketStatus
	Strategy             Strategy
	InvestmentPercentage float64
	StopLossPercentage   float64

	capitalClient       *CapitalClientAPI
	candlesChannel      chan Candles // TODO: Implement "Pipeline" Pattern To Handle Larger Data Efficiently
	activeDealReference string

	payedPrice                   float64
	volatileAmountAvailable      float64
	volatileInvestmentPercentage float64
}

type MinitraderStatus string

const (
	// healthy statuses
	NEW               MinitraderStatus = "NEW"
	RUNNING           MinitraderStatus = "RUNNING"
	HOLDING           MinitraderStatus = "HOLDING"
	SELL_ORDER_ACTIVE MinitraderStatus = "SELL_ORDER_ACTIVE"
	BUY_ORDER_ACTIVE  MinitraderStatus = "BUY_ORDER_ACTIVE"

	// error statuses
	ERROR_ON_UPDATE_CANDLES_DATA MinitraderStatus = "ERROR_ON_UPDATE_CANDLES_DATA"
	ERROR_ON_MAKING_ORDER        MinitraderStatus = "ERROR_ON_MAKING_ORDER"
	ERROR_ON_DELETING_ORDER      MinitraderStatus = "ERROR_ON_DELETING_ORDER"
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

var TimeframeMinuteMap = map[Timeframe]int{
	MINUTE:    1,
	MINUTE_5:  5,
	MINUTE_15: 15,
	MINUTE_30: 30,
	HOUR:      60,
	HOUR_4:    240,
	DAY:       1440,
	WEEK:      10080,
}

func NewMinitrader(epic string, investmentPercentage float64, stopLossPercentage float64, timeframe Timeframe, strategy Strategy) *Minitrader {
	return &Minitrader{
		Epic:                         epic,
		InvestmentPercentage:         investmentPercentage,
		Timeframe:                    timeframe,
		Strategy:                     strategy,
		Status:                       NEW,
		StopLossPercentage:           stopLossPercentage,
		candlesChannel:               make(chan Candles),
		volatileInvestmentPercentage: investmentPercentage,
	}
}

func (minitrader *Minitrader) Start(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	defer close(minitrader.candlesChannel)
	minitrader.Status = RUNNING
	for candles := range minitrader.candlesChannel {
		signal, price := minitrader.Strategy(candles)
		err := minitrader.Effect(signal, price)
		if err != nil {
			log.Printf("Error On Minitrader Effect: %v", err)
			return
		}
	}
}

func (minitrader *Minitrader) Effect(signal Signal, price float64) error {
	log.Printf("Epic: %s - Signal: %v - Price: %v", minitrader.Epic, signal, price)

	if minitrader.MarketStatus == CLOSED {
		return errors.New("Unable To Do Trading; Market Closed")
	}

	// quick sell out
	if minitrader.payedPrice*(100-minitrader.StopLossPercentage) > price {
		err := minitrader.deleteOrder(minitrader.activeDealReference)
		if err != nil {
			minitrader.Status = ERROR_ON_DELETING_ORDER
			return err
		}
	}

	// make a buy/sell order and wait 3:30 minutes or less if order has been completed before wait time.
	if (minitrader.Status == RUNNING && signal == BUY) || (minitrader.Status == HOLDING && signal == SELL) {
		err := minitrader.makeOrderAndWaitUntilComplete(minitrader.Epic, signal, LIMIT, price)
		if err != nil {
			minitrader.Status = ERROR_ON_MAKING_ORDER
			return err
		}
	}

	return nil
}

func (minitrader *Minitrader) makeOrderAndWaitUntilComplete(epic string, signal Signal, orderType OrderType, targetPrice float64) error {
	var amount float64
	var err error
	var dealReference string

	// update minitrader status and get amount available
	// from preferred account or get amount from position/order confirmation
	if signal == BUY {
		minitrader.Status = BUY_ORDER_ACTIVE
		amount = minitrader.volatileAmountAvailable // TODO; Double check if quantity good (amount)
	} else {
		minitrader.Status = SELL_ORDER_ACTIVE
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

	// update minitrader status, active deal reference and payed price
	if signal == BUY {
		minitrader.Status = HOLDING
		minitrader.activeDealReference = dealReference
		minitrader.payedPrice = targetPrice
	} else {
		minitrader.Status = RUNNING
		minitrader.activeDealReference = ""
		minitrader.payedPrice = 0.0
	}

	return nil
}

func (minitrader *Minitrader) getAmountFromPositionOrderConfirmation() (amount float64, err error) { // TODO: Unused
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

func (minitrader *Minitrader) deleteOrder(dealReference string) error {
	_, err := minitrader.capitalClient.DeleteWorkingOrder(dealReference)
	if err != nil {
		return err
	}
	minitrader.activeDealReference = ""
	minitrader.Status = RUNNING
	minitrader.payedPrice = 0.0

	return nil
}

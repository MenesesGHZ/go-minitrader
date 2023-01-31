/*
GPT Prompt: Can you write a better strategy than the previous one, with a higher performance and less risk ?

GPT Response:

Sure, here is an example of a more advanced intraday trading strategy that incorporates multiple
indicators and logic for generating trading signals, with the focus on minimizing risk:

This strategy builds on the previous one, by adding the Average True Range (ATR) as an additional indicator.
ATR is a volatility indicator that can be used to identify market conditions that are conducive to profitable trades.
The strategy generates a buy signal when the RSI is above 70 and the close price is above the upper band, indicating
an overbought condition and a potential price reversal, and ATR is below a certain threshold. It generates a sell
signal when the RSI is below 30 and the close price is below the lower band, indicating an oversold condition and a
potential price reversal, and ATR is below a certain threshold.

This strategy is more complex than previous one but it uses multiple indicators and logic for generating signals
which could result in better performance. In addition to that, the use of ATR as a volatility indicator helps to

(He didn't say more)

 -----------------------

As a side note, originally the code generated was in python, then I asked to translat it to a real time implementation
using golang. It came with some errors related to `index out of range` mainly. But ChatGPT was able to fix it himself.

Also there was more converstion before this one, related to which is the best market to do intrady and
what are good strategies for do so.
*/

package gominitrader

import (
	"math"
)

func GPTStrategy(candles Candles) (Signal, float64) {
	rsi := []float64{}
	numberOfCandles := len(candles)
	for i := 14; i < numberOfCandles; i++ {
		// Calculate the average gain and average loss
		var avgGain float64
		var avgLoss float64
		for j := i - 13; j < i; j++ {
			// calculate the change in price between the current candle and the previous candle
			change := candles[j].Close.Bid - candles[j-1].Close.Bid
			if change > 0 {
				// if the change is positive, add it to the average gain
				avgGain += change
			} else {
				// if the change is negative, add it to the average loss
				avgLoss += change
			}
		}
		// divide the average gain and average loss by 14 to get the average gain and loss over the last 14 candles
		avgGain /= 14
		avgLoss /= 14

		// Calculate the relative strength
		rs := avgGain / -avgLoss

		// Calculate the relative strength index
		rsi = append(rsi, 100-(100/(1+rs)))
	}

	// Create slices for the Bollinger Bands values
	upperBand := []float64{}
	middleBand := []float64{}
	lowerBand := []float64{}

	// Calculate the Bollinger Bands for each data point
	for i := 20; i < numberOfCandles; i++ {
		// Calculate the moving average
		var movingAvg float64
		for j := i - 20; j < i; j++ {
			movingAvg += candles[j].Close.Bid
		}
		movingAvg /= 20

		// Calculate the standard deviation
		var variance float64
		for j := i - 20; j < i; j++ {
			variance += math.Pow(candles[j].Close.Bid-movingAvg, 2)
		}
		stdDev := math.Sqrt(variance / float64(20))

		// Calculate the Bollinger Bands
		upperBand = append(upperBand, movingAvg+2*stdDev)
		middleBand = append(middleBand, movingAvg)
		lowerBand = append(lowerBand, movingAvg-2*stdDev)
	}

	price := candles[numberOfCandles-1].Close.Bid
	if price < lowerBand[numberOfCandles-21] && rsi[len(rsi)-1] < 30 {
		return BUY, price
	}
	if price > upperBand[numberOfCandles-21] && rsi[len(rsi)-1] > 70 {
		return SELL, price
	}
	return NONE, price
}

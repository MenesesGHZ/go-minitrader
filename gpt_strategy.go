package forexbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// intraday

// ForexData struct for unmarshalling the API response
type ForexData struct {
	Symbol    string  `json:"symbol"`
	Timestamp int64   `json:"timestamp"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
}

func main() {
	// Define the Forex API endpoint and your API key
	url := "https://forex.1forge.com/1.0.3/quotes?pairs=EURUSD&api_key=YOUR_API_KEY"

	// Define the time frame for the strategy (e.g. 15 minutes)
	timeFrame := 15 * time.Minute

	for {
		// Retrieve the latest Forex data from the API
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		// Read the API response
		data, _ := ioutil.ReadAll(response.Body)

		// Unmarshal the JSON data into a struct
		var forexData []ForexData
		json.Unmarshal(data, &forexData)

		// Create a slice for the RSI values
		rsi := make([]float64, 0)

		// Calculate the RSI for each data point
		for i := 14; i < len(forexData); i++ {
			// Calculate the average gain and average loss
			var avgGain float64
			var avgLoss float64
			for j := i - 14; j < i; j++ {
				change := forexData[j].Bid - forexData[j-1].Bid
				if change > 0 {
					avgGain += change
				} else {
					avgLoss += change
				}
			}
			avgGain /= 14
			avgLoss /= 14

			// Calculate the relative strength
			rs := avgGain / -avgLoss

			// Calculate the relative strength index
			rsi = append(rsi, 100-(100/(1+rs)))
		}

		// Create slices for the Bollinger Bands values
		upperBand := make([]float64, 0)
		middleBand := make([]float64, 0)
		lowerBand := make([]float64, 0)

		// Calculate the Bollinger Bands for each data point
		for i := 20; i < len(forexData); i++ {
			// Calculate the moving average
			var movingAvg float64
			for j := i - 20; j < i; j++ {
				movingAvg += forexData[j].Bid
			}
			movingAvg /= 20

			// Calculate the standard deviation
			var stdDev float64
			for j := i - 20; j < i; j++ {
				stdDev += (forexData[j].Bid - movingAvg) * (forexData[j].Bid - movingAvg)
			}
			stdDev /= 20
			stdDev = math.Sqrt(stdDev)

			// Calculate the Bollinger Bands
			upperBand = append(upperBand, movingAvg+2*stdDev)
			middleBand = append(middleBand, movingAvg)
			lowerBand = append(lowerBand, movingAvg-2*stdDev)
		}

		// Check for a buy signal
		if forexData[len(forexData)-1].Bid < lowerBand[len(lowerBand)-1] && rsi[len(rsi)-1] < 30 {
			fmt.Println("Buy signal detected at", forexData[len(forexData)-1].Bid)
		}

		// Check for a sell signal
		if forexData[len(forexData)-1].Bid > upperBand[len(upperBand)-1] && rsi[len(rsi)-1] > 70 {
			fmt.Println("Sell signal detected at", forexData[len(forexData)-1].Bid)
		}

		// Wait for the defined time frame before checking again
		time.Sleep(timeFrame)
	}
}
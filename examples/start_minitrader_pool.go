package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	gominitrader "github.com/menesesghz/go-minitrader"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	capitalEmail := os.Getenv("CAPITAL_EMAIL")
	capitalApiKey := os.Getenv("CAPITAL_API_KEY")
	capitalApiKeyPassword := os.Getenv("CAPITAL_API_KEY_PASSWORD")

	capitalClient, err := gominitrader.NewCapitalClient(capitalEmail, capitalApiKey, capitalApiKeyPassword, true)
	if err != nil {
		log.Fatal(err)
	}

	// Trade the USD/JPY currency pair with a 2% stop loss and a 0.35% upperbound profit target using 15-minute intervals and the GPTStrategy.
	minitraderUSDJPY := gominitrader.NewMinitrader(
		"USDJPY", 25, 2, 0.35, gominitrader.MINUTE_15, gominitrader.GPTStrategy,
	)

	// Trade the USD/CAD currency pair with a 2% stop loss and a 0.35% upperbound profit target using 15-minute intervals and the GPTStrategy.
	minitraderUSDCAD := gominitrader.NewMinitrader(
		"USDCAD", 25, 2, 0.35, gominitrader.MINUTE_15, gominitrader.GPTShortTermStrategy,
	)

	// Trade the USD/MXN currency pair with a 2% stop loss and a 0.35% upperbound profit target using 15-minute intervals and the GPTStrategy.
	minitraderUSDMXN := gominitrader.NewMinitrader(
		"USDMXN", 50, 2, 0.35, gominitrader.MINUTE_15, gominitrader.GPTShortTermStrategy,
	)

	minitraderPool, _ := gominitrader.NewMinitraderPool(capitalClient, minitraderUSDJPY, minitraderUSDCAD, minitraderUSDMXN)
	minitraderPool.Start()
}

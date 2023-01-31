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
	minitraderUSDJPY := gominitrader.NewMinitrader(
		"USDJPY", 50, 5, gominitrader.MINUTE_15, gominitrader.GPTStrategy,
	)
	minitraderUSDCAD := gominitrader.NewMinitrader(
		"USDCAD", 50, 5, gominitrader.MINUTE_15, gominitrader.GPTStrategy,
	)

	minitraderPool, _ := gominitrader.NewMinitraderPool(capitalClient, minitraderUSDJPY, minitraderUSDCAD)
	minitraderPool.Start()
}

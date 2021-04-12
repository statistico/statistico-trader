package main

import (
	"context"
	"fmt"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	q := app.Queue()
	h := app.MarketHandler()
	ctx := context.Background()

	for {
		fmt.Println("Polling queue for messages...")

		markets := q.ReceiveMarkets()

		for m := range markets {
			h.HandleEventMarket(ctx, m)
		}
	}
}

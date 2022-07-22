package api

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	pfsync "github.com/umee-network/umee/price-feeder/pkg/sync"
)

const (
	tickerSleep = 500 * time.Millisecond
)

type API struct {
	logger zerolog.Logger
	closer *pfsync.Closer

	candles map[string]Candle
	tickers map[string]Ticker
}

func NewAPI() API {
	return API{}
}

// Start starts the api process in a blocking fashion.
func (a *API) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			a.closer.Close()

		default:
			a.logger.Debug().Msg("starting api tick")

			if err := a.tick(ctx); err != nil {
				a.logger.Err(err).Msg("api tick failed")
			}

			time.Sleep(tickerSleep)
		}
	}
}

// tick is where we'll query the osmosis node,
// update current pricing data, and then update the
// websocket subscriptions.
func (a *API) tick(ctx context.Context) error {
	return nil
}

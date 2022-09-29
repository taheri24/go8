package main

import (
	"context"
	"github.com/gmhafiz/go8/ent/gen"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return &gen.Trade{}
		},
	}

	dataChan := make(chan *gen.Trade, 1_000)
	workerCount := 10
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i := 0; i < workerCount; i++ {

		go func() {
			client := gen.NewClient()
			defer client.Close()
			for trade := range dataChan {
				client.
					Trade.
					Create().
					SetInstrumentID(trade.InstrumentID).
					SetOpen(trade.Open).
					SetClose(trade.Close).
					SetDateEn(trade.DateEn).
					SetHigh(trade.High).
					SetLow(trade.Low).Save(ctx)
				pool.Put(trade)
			}
		}()
	}
	randCount := 1_000_000
	for i := 0; i < randCount; i++ {
		freshTrade := pool.Get().(*gen.Trade)
		*freshTrade = gen.Trade{
			ID:           0,
			InstrumentID: uint(math.Floor(float64(i) / 1000)),
			DateEn:       time.Now(),
			Open:         float64(rand.Int()),
			Close:        float64(rand.Int()),
			High:         float64(rand.Int()),
			Low:          float64(rand.Int()),
		}
		dataChan <- freshTrade
	}
}

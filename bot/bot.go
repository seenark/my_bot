package bot

import (
	"fmt"
	"mybot/binance"
	"mybot/binance/timeframe"
	"mybot/candle"
	"mybot/mixins"
	"time"

	"github.com/markcheno/go-talib"
)

var (
	slowEmaList []candle.Ema = []candle.Ema{
		{
			Time:  0,
			Value: 0,
		},
		{
			Time:  0,
			Value: 0,
		},
	}
	fastEmaList []candle.Ema = []candle.Ema{
		{
			Time:  0,
			Value: 0,
		},
		{
			Time:  0,
			Value: 0,
		},
	}
)

func Process(b *binance.BinanceClient) {
	pair := "SOLUSDT"
	// canLong := false
	// canShort := false
	// getSlowEma(b, pair)
	// getFastEma(b, pair)
	for {
		now := time.Now()
		min := now.Minute()
		second := now.Second()
		fetchEma(b, pair, min, second)
		fmt.Print("Fast Ema ", fastEmaList[1].Value)
		mixins.PrintTime(fastEmaList[1].Time)
		fmt.Print("Slow Ema ", slowEmaList[1].Value)
		mixins.PrintTime(slowEmaList[1].Time)
		// pause
		time.Sleep(55 * time.Second)
	}

}

func fetchEma(b *binance.BinanceClient, pair string, min int, second int) {
	if slowEmaList[0].Time == 0 || slowEmaList[0].Value == 0 {
		setSlowEma(b, pair)
	}
	if fastEmaList[0].Time == 0 || fastEmaList[0].Value == 0 {
		setFastEma(b, pair)
	}
	// slow ema every 1 min
	if second == 00 {
		setSlowEma(b, pair)

		// fast ema every 30 min
		if min == 00 || min == 30 {
			setFastEma(b, pair)
		}
	}
}

func setFastEma(b *binance.BinanceClient, pair string) {
	emas := getFastEma(b, pair)
	fastEmaList[0] = emas[len(emas)-2]
	fastEmaList[1] = emas[len(emas)-1]
	fmt.Printf("fastEmaList: %v\n", fastEmaList)
}

func setSlowEma(b *binance.BinanceClient, pair string) {
	emas := getSlowEma(b, pair)
	slowEmaList[0] = emas[len(emas)-2]
	slowEmaList[1] = emas[len(emas)-1]
	fmt.Printf("slowEmaList: %v\n", slowEmaList)
}

func getSlowEma(b *binance.BinanceClient, pair string) []candle.Ema {
	var shift time.Duration = 100
	shift2 := 100

	end := time.Now()
	fmt.Printf("end: %v\n", end)
	// start := end.AddDate(0, 0, -1)
	start := end.Add((-1 * shift) * 1 * time.Minute)
	// fmt.Printf("start: %v\n", start)

	slowKlines := b.GetKLine(pair, timeframe.M1, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), shift2)
	slowCandles := binance.KLinesToCandles(slowKlines)
	slowHaCadles := candle.CandlesToHeikinAshies(*slowCandles)

	slowEma := talib.Ema(candle.GetHACloseSeries(*slowHaCadles), 16)
	slowEmaSeries, err := candle.MakeEmaSeries(*slowHaCadles, slowEma)
	if err != nil {
		fmt.Println(err)
	}
	// mixins.PrintPretty(slowEmaSeries)
	return slowEmaSeries
}

func getFastEma(b *binance.BinanceClient, pair string) []candle.Ema {
	var shift time.Duration = 5
	shift2 := 5
	end := time.Now()
	fmt.Printf("end: %v\n", end)
	// start := end.AddDate(0, 0, -1)
	start := end.Add((-1 * shift) * 30 * time.Minute)
	// fmt.Printf("start: %v\n", start)

	fastKlines := b.GetKLine(pair, timeframe.M30, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), shift2)
	fastCandles := binance.KLinesToCandles(fastKlines)
	fastHaCadles := candle.CandlesToHeikinAshies(*fastCandles)

	fastEma := talib.Ema(candle.GetHACloseSeries(*fastHaCadles), 1)
	fastEmaSeries, err := candle.MakeEmaSeries(*fastHaCadles, fastEma)
	if err != nil {
		fmt.Println(err)
	}
	// mixins.PrintPretty(fastEmaSeries)
	return fastEmaSeries
}

package lotus

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/modules/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

var (
	marketStorageAsk []*storagemarket.StorageAsk
	bestMinerAsk     []*storagemarket.StorageAsk
	trustedMiners    []*storagemarket.StorageAsk
	usedMiner        []string
	usedBest         []string
)

func AskPrice(apiCtx *apicontext.APIContext) {
	// ctx, cancel := context.WithCancel(apiCtx.Context)
	// defer cancel()

	// nodeApi := apiCtx.FullNode
	logrus.Infof("init price")
	for _, t := range apiCtx.Config.Rule.Trusted {
		price, err := types.ParseFIL(t.Price)
		if err != nil {
			logrus.Warnf("err price of trusted miner %s\n", t.Miner)
			continue
		}
		taddr, err := address.NewFromString(t.Miner)
		if err != nil {
			logrus.Warnf("err miner of trusted miner %s\n", t.Miner)
			continue
		}
		trustedMiners = append(trustedMiners, &storagemarket.StorageAsk{
			Price:        types.BigInt(price),
			MinPieceSize: 0,
			MaxPieceSize: 0,
			Miner:        taddr,
			Timestamp:    0,
		})
	}
	// timer := util.Ticker(ctx, time.Minute*30)

	// doAsk(apiCtx, nodeApi)

	// for {
	// 	select {
	// 	case <-timer:
	// 		logrus.Info("start ask price")
	// 		doAsk(apiCtx, nodeApi)
	// 	case <-ctx.Done():
	// 		logrus.Info("exist askPrice")
	// 		return
	// 	}
	// }
}

func doAsk(apiCtx *apicontext.APIContext, nodeapi api.FullNode) {
	res, err := getAsks(apiCtx, nodeapi)
	if err != nil {
		logrus.Errorf("get ask: %s\n", err.Error())
		return
	}
	best := make([]*storagemarket.StorageAsk, 0)
	market := make([]*storagemarket.StorageAsk, 0)
	fmt.Printf("miner Price_per_GiB Verified_Price_per_GiB Min_size Max_size Timestamp Expiry SeqNo\n")
	for _, ask := range res {
		if util.ArrInclude(apiCtx.Config.Rule.Best, ask.Miner.String()) {
			best = append(best, ask)
		} else {
			market = append(market, ask)
		}
		fmt.Printf("%s %s %s %s %s %s %s %s\n", ask.Miner, types.FIL(ask.Price), types.FIL(ask.VerifiedPrice), types.SizeStr(types.NewInt(uint64(ask.MinPieceSize))),
			types.SizeStr(types.NewInt(uint64(ask.MaxPieceSize))), types.NewInt(uint64(ask.Timestamp)), types.NewInt(uint64(ask.Expiry)), types.NewInt(ask.SeqNo))
	}
	bestMinerAsk = best
	marketStorageAsk = market
}
func getAsks(apictx *apicontext.APIContext, api api.FullNode) ([]*storagemarket.StorageAsk, error) {
	ctx := apictx.Context
	color.Blue(".. getting miner list")
	miners, err := api.StateListMiners(ctx, types.EmptyTSK)
	if err != nil {
		return nil, xerrors.Errorf("getting miner list: %w", err)
	}
	logrus.Printf("miner len %d", len(miners))

	var lk sync.Mutex
	var found int64
	var withMinPower []address.Address
	done := make(chan struct{})

	go func() {
		defer close(done)

		var wg sync.WaitGroup
		wg.Add(len(miners))

		throttle := make(chan struct{}, 100)
		for _, miner := range miners {
			throttle <- struct{}{}
			go func(miner address.Address) {
				defer wg.Done()
				defer func() {
					<-throttle
				}()

				power, err := api.StateMinerPower(ctx, miner, types.EmptyTSK)
				if err != nil {
					return
				}

				if power.HasMinPower { // TODO: Lower threshold
					atomic.AddInt64(&found, 1)
					lk.Lock()
					withMinPower = append(withMinPower, miner)
					lk.Unlock()
				}
			}(miner)
		}
	}()

loop:
	for {
		select {
		case <-time.After(150 * time.Millisecond):
			logrus.Printf("\r* Found %d miners with power", atomic.LoadInt64(&found))
		case <-done:
			break loop
		}
	}
	logrus.Printf("\r* Found %d miners with power\n", atomic.LoadInt64(&found))

	color.Blue(".. querying asks")

	var asks []*storagemarket.StorageAsk
	var queried, got int64

	done = make(chan struct{})
	go func() {
		defer close(done)

		var wg sync.WaitGroup
		wg.Add(len(withMinPower))

		throttle := make(chan struct{}, 50)
		for _, miner := range withMinPower {
			throttle <- struct{}{}
			go func(miner address.Address) {
				defer wg.Done()
				defer func() {
					<-throttle
					atomic.AddInt64(&queried, 1)
				}()

				ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
				defer cancel()

				mi, err := api.StateMinerInfo(ctx, miner, types.EmptyTSK)
				if err != nil {
					return
				}
				if mi.PeerId == nil {
					return
				}

				ask, err := api.ClientQueryAsk(ctx, *mi.PeerId, miner)
				if err != nil {
					return
				}

				atomic.AddInt64(&got, 1)
				if util.ArrInclude(apictx.Config.Rule.Disable, ask.Miner.String()) { // 黑名单,跳过
					return
				}
				lk.Lock()
				asks = append(asks, ask)
				lk.Unlock()
			}(miner)
		}
	}()

loop2:
	for {
		select {
		case <-time.After(150 * time.Millisecond):
			logrus.Printf("\r* Queried %d asks, got %d responses", atomic.LoadInt64(&queried), atomic.LoadInt64(&got))
		case <-done:
			break loop2
		}
	}
	logrus.Printf("\r* Queried %d asks, got %d responses\n", atomic.LoadInt64(&queried), atomic.LoadInt64(&got))

	sort.Slice(asks, func(i, j int) bool {
		return asks[i].Price.LessThan(asks[j].Price)
	})

	return asks, nil
}

func GetMarketStorageAsk() []interface{} {
	ret := make([]interface{}, 0)
	for _, ask := range marketStorageAsk {
		ret = append(ret, map[string]interface{}{
			"miner":                  ask.Miner,
			"price_per_gib":          types.FIL(ask.Price).String(),
			"verified_price_per_gib": types.FIL(ask.VerifiedPrice).String(),
			"min_size":               types.SizeStr(types.NewInt(uint64(ask.MinPieceSize))),
			"max_size":               types.SizeStr(types.NewInt(uint64(ask.MaxPieceSize))),
			"timestamp":              types.NewInt(uint64(ask.Timestamp)),
			"expiry":                 types.NewInt(uint64(ask.Expiry)),
			"seqno":                  types.NewInt(ask.SeqNo),
		})
	}

	for _, ask := range bestMinerAsk {
		ret = append(ret, map[string]interface{}{
			"miner":                  ask.Miner,
			"price_per_gib":          types.FIL(ask.Price).String(),
			"verified_price_per_gib": types.FIL(ask.VerifiedPrice).String(),
			"min_size":               types.SizeStr(types.NewInt(uint64(ask.MinPieceSize))),
			"max_size":               types.SizeStr(types.NewInt(uint64(ask.MaxPieceSize))),
			"timestamp":              types.NewInt(uint64(ask.Timestamp)),
			"expiry":                 types.NewInt(uint64(ask.Expiry)),
			"seqno":                  types.NewInt(ask.SeqNo),
		})
	}

	for _, ask := range trustedMiners {
		ret = append(ret, map[string]interface{}{
			"miner":                  ask.Miner,
			"price_per_gib":          types.FIL(ask.Price).String(),
			"verified_price_per_gib": types.FIL(ask.VerifiedPrice).String(),
			"min_size":               types.SizeStr(types.NewInt(uint64(ask.MinPieceSize))),
			"max_size":               types.SizeStr(types.NewInt(uint64(ask.MaxPieceSize))),
			"timestamp":              types.NewInt(uint64(ask.Timestamp)),
			"expiry":                 types.NewInt(uint64(ask.Expiry)),
			"seqno":                  types.NewInt(ask.SeqNo),
		})
	}

	return ret
}

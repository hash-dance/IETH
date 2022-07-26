package lotus

import (
	"context"
	"fmt"
	"sort"

	"github.com/docker/go-units"
	"github.com/filecoin-project/go-address"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
	db "github.com/guowenshuai/ieth/db/mongo"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/modules/util"
	itypes "github.com/guowenshuai/ieth/types"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-cidutil/cidenc"
	"github.com/multiformats/go-multibase"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	bestDealCounts = 0
)

// StartDeal 发起交易
// @param from string 钱包地址
// @param fCID string 文件cid
// @param miner string 目标矿工
// @param price big.Int 交易的订单金额
// @param duration int 订单存储时间,天
func startDeal(ctx context.Context, nodeApi api.FullNode, from, miner address.Address, fcid cid.Cid, price big.Int, epoch abi.ChainEpoch) (string, error) {
	// return "xxxxxxxxxxxxxxxxxxxxxxxx", nil
	if abi.ChainEpoch(epoch) < build.MinDealDuration {
		epoch = build.MinDealDuration
	}
	fref := &storagemarket.DataRef{
		TransferType: storagemarket.TTGraphsync,
		Root:         fcid,
	}

	dcap, err := nodeApi.StateVerifiedClientStatus(ctx, from, types.EmptyTSK)
	if err != nil {
		return "", err
	}
	isVerified := dcap != nil
	proposal, err := nodeApi.ClientStartDeal(ctx, &api.StartDealParams{
		Data:               fref,
		Wallet:             from,
		Miner:              miner,
		EpochPrice:         types.BigInt(price),
		MinBlocksDuration:  uint64(epoch),
		DealStartEpoch:     abi.ChainEpoch(-1),
		FastRetrieval:      true,
		VerifiedDeal:       isVerified,
		ProviderCollateral: big.NewInt(0),
	})
	if err != nil {
		return "", err
	}
	e := cidenc.Encoder{Base: multibase.MustNewEncoder(multibase.Base32)}
	str := e.Encode(*proposal)
	fmt.Println(str)
	return str, nil
}

// 构造一个订单,并发起交易
func buildDeal(apiCtx *apicontext.APIContext, from string, miner address.Address, fcid string, price big.Int, epoch abi.ChainEpoch) (string, error) {
	ctx, cancel := context.WithCancel(apiCtx.Context)
	defer cancel()
	nodeApi := apiCtx.FullNode
	faddr, err := address.NewFromString(from)
	if err != nil {
		return "", fmt.Errorf("failed to parse 'from' address: %w", err)
	}

	data, err := cid.Parse(fcid)
	if err != nil {
		return "", err
	}
	dealcid, err := startDeal(ctx, nodeApi, faddr, miner, data, price, epoch)
	if err != nil {
		return "", err
	}
	logrus.Infof("dealcid [%s] make deal [%s] to miner [%s], price is [%s], epoch %s\n", dealcid, fcid, miner.String(), types.FIL(price), epoch)
	return dealcid, nil
}

// 获取当前交易池中正在交易的订单
// 如果交易传输量达到限制, 不再添加新的交易
func getPendingDeal(apiCtx *apicontext.APIContext, nodeApi api.FullNode) (int, error) {
	deals, _, err := listDeals(apiCtx.Context, nodeApi)
	if err != nil {
		return 0, err
	}
	num := 0
	for _, deal := range deals {
		if isPendingDeal(deal.LocalDeal.State) {
			num++
		}
	}
	return num, nil
}

func isPendingDeal(status storagemarket.StorageDealStatus) bool {
	return status == storagemarket.StorageDealStartDataTransfer || status == storagemarket.StorageDealTransferring
}

// 从demo获取所有交易
func listDeals(ctx context.Context, full api.FullNode) (deals []*itypes.Deal, dealsCids []string, err error) {
	localDeals, err := full.ClientListDeals(ctx)
	if err != nil {
		return
	}
	sort.Slice(localDeals, func(i, j int) bool {
		return localDeals[i].CreationTime.Before(localDeals[j].CreationTime)
	})

	for _, localDeal := range localDeals {
		deals = append(deals, dealFromDealInfo(ctx, full, localDeal))
		dealsCids = append(dealsCids, localDeal.ProposalCid.String())
	}
	return
}

func getDeal(ctx context.Context, full api.FullNode, dealcid string) (*itypes.Deal, error) {
	propcid, err := cid.Decode(dealcid)
	if err != nil {
		return nil, err
	}

	dealinfo, err := full.ClientGetDealInfo(ctx, propcid)
	if err != nil {
		logrus.Errorf("clientGetDealInfo %s", err.Error())
		return nil, err
	}
	return dealFromDealInfo(ctx, full, *dealinfo), nil
}

func listTransfers(ctx context.Context, full api.FullNode) (int, error) {
	channels, err := full.ClientListDataTransfers(ctx)
	if err != nil {
		return 0, err
	}
	ongoingCount := 0
	for _, channel := range channels {
		if channel.Status == datatransfer.Ongoing {
			logrus.Infof("deal transfer going %d\t%d", channel.TransferID, units.BytesSize(float64(channel.Transferred)))
			ongoingCount += 1
		}
	}
	return ongoingCount, nil
}

func dealFromDealInfo(ctx context.Context, full api.FullNode, v api.DealInfo) *itypes.Deal {
	if v.DealID == 0 {
		return &itypes.Deal{
			LocalDeal: v,
			OnChain:   nil,
		}
	}

	onChain, err := full.StateMarketStorageDeal(ctx, v.DealID, types.EmptyTSK)
	if err != nil {
		return &itypes.Deal{LocalDeal: v}
	}

	return &itypes.Deal{
		LocalDeal: v,
		OnChain:   onChain,
	}
}

// FetchCanDealCids 计算一次可发单量
func FetchCanDealCids(apiCtx *apicontext.APIContext, askNum int, maxPrice string, toMarker, tosort bool) ([]*itypes.DealWithMiner, error) {
	maxDealOne := apiCtx.Config.Setting.MaxDealOne
	canDeals := make([]*itypes.DealWithMiner, 0)
	allDealCids := make([]string, 0)
	maxPriceFil, err := types.ParseFIL(maxPrice)
	if err != nil {
		return nil, err
	}
	logrus.Infof("maxPriceFil %s", maxPriceFil)
	// 订单比例调0
	bestDealCounts = 0

	dealsCol := apiCtx.MongoClient.Collection(db.DealsCollectionName)
	matchPip := bson.D{{"$group", bson.D{{"_id", "$payloadcid"}, {"total", bson.D{{"$sum", 1}}}, {"miners", bson.D{{"$push", "$minerid"}}}}}}
	showInfoCursor, err := dealsCol.Aggregate(context.Background(), mongo.Pipeline{matchPip})
	if err != nil {
		return canDeals, err
	}
	defer showInfoCursor.Close(context.Background())
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.Background(), &showsWithInfo); err != nil {
		return canDeals, err
	}
	for _, v := range showsWithInfo {
		ccid := fmt.Sprintf("%s", v["_id"])
		allDealCids = append(allDealCids, ccid)
		total := int(v["total"].(int32))
		if total > maxDealOne {
			continue
		}
		miners := make([]string, 0)
		for _, m := range v["miners"].(bson.A) {
			miners = append(miners, m.(string))
		}
		for i := 0; i < maxDealOne-total; i++ {
			nextMiner := selectMiners(apiCtx.Config.Rule, types.BigInt(maxPriceFil), miners, 0, askNum, 100)
			if nextMiner == nil {
				continue
			}
			todo := &itypes.DealWithMiner{
				PayloadCid: ccid,
				Miner:      nextMiner,
			}
			canDeals = append(canDeals, todo)
			miners = append(miners, nextMiner.Miner.String()) // 同一时间同样的交易不能是一个矿工
			if len(canDeals) >= askNum {
				goto RET
			}
		}
	}
RET:
	// 如果满足int(maxDeal), 则直接返回
	if len(canDeals) >= askNum {
		return canDeals[:askNum], nil
	}
	// 从ipfscids表中找出新的cid
	findops := options.Find()
	findops.SetLimit(1000)
	if tosort { // 如果指定该参数, 则升序排序, 默认降序
		findops.SetSort(bson.D{{"filesizebytes", 1}})
	} else {
		findops.SetSort(bson.D{{"filesizebytes", -1}})
	}
	fileCol := apiCtx.MongoClient.Collection(db.IpfsCollectionName)
	showInfoCursor, err = fileCol.Find(context.Background(), bson.D{{"payloadcid", bson.D{{"$nin", allDealCids}}}}, findops)
	if err != nil {
		return canDeals, err
	}
	var results []*itypes.IpfsData
	if err := showInfoCursor.All(context.TODO(), &results); err != nil {
		return canDeals, err
	}
	diff := askNum - len(canDeals)
	for _, v := range results {
		miners := make([]string, 0)
		for i := 0; i < maxDealOne; i++ {
			if diff == 0 {
				return canDeals, nil
			}
			nextMiner := selectMiners(apiCtx.Config.Rule, types.BigInt(maxPriceFil), miners, uint64(v.FileSizeBytes), askNum, 100)
			if nextMiner == nil {
				continue
			}
			todo := &itypes.DealWithMiner{
				PayloadCid: v.PayloadCid,
				Miner:      nextMiner,
			}
			canDeals = append(canDeals, todo)
			miners = append(miners, nextMiner.Miner.String()) // 同一时间同样的交易不能是一个矿工
			diff--
		}
	}
	return canDeals, nil
}

func selectMiners(rule *itypes.Rule, maxPrice types.BigInt, miners []string, size uint64, askNum, maxPercent int) *storagemarket.StorageAsk {
	priorityMiners := make([]*storagemarket.StorageAsk, 0)
	priorityMiners = append(priorityMiners, trustedMiners...)
	// priorityMiners = append(priorityMiners, bestMinerAsk...)

	if (bestDealCounts * 100 / askNum) < maxPercent {
		// 优先矿工
		if len(usedBest) == len(priorityMiners) {
			usedBest = make([]string, 0)
		}

		for _, m := range priorityMiners {
			if util.ArrInclude(miners, m.Miner.String()) {
				continue
			}
			if m.Price.GreaterThan(maxPrice) {
				continue
			}
			if util.ArrInclude(usedBest, m.Miner.String()) {
				continue
			}
			if !checkSizeOk(size, m) {
				continue
			}
			bestDealCounts += 1
			usedBest = append(usedBest, m.Miner.String())
			return m
		}
	}

	var selected *storagemarket.StorageAsk
Loop:
	for _, m := range marketStorageAsk {
		if util.ArrInclude(rule.Disable, m.Miner.String()) { // 黑名单,跳过
			continue
		}
		if util.ArrInclude(miners, m.Miner.String()) {
			continue
		}
		if util.ArrInclude(usedMiner, m.Miner.String()) {
			continue
		}
		if m.Price.GreaterThan(maxPrice) {
			continue
		}
		if !checkSizeOk(size, m) {
			continue
		}
		// return m
		selected = m
		usedMiner = append(usedMiner, m.Miner.String())
		break
	}
	if selected == nil {
		usedMiner = make([]string, 0)
		goto Loop
	}
	return selected
	// return nil
}

func checkSizeOk(size uint64, ask *storagemarket.StorageAsk) bool {
	if size == 0 {
		return true
	}
	if ask.MaxPieceSize == 0 && ask.MinPieceSize == 0 {
		return true
	}

	if size > uint64(ask.MinPieceSize) && size < uint64(ask.MaxPieceSize) {
		return true
	}
	return false
}

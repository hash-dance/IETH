package lotus

import (
	"errors"
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/build"
	apicontext "github.com/guowenshuai/ieth/modules/context"
)

type DealPool struct {
	ctx      *apicontext.APIContext
	wallet   string
	duration abi.ChainEpoch
	limit    int // 限制队列数
}

var dealPoolIInstance *DealPool

func StartDealPool(apiCtx *apicontext.APIContext) error {
	_, err := GetDealPool(apiCtx)
	if err != nil {
		return err
	}
	// go dealpool.watch()
	return nil
}

// GetDealPool 新建一个实例
func GetDealPool(apiCtx *apicontext.APIContext) (*DealPool, error) {
	if dealPoolIInstance == nil {
		duration := apiCtx.Config.Setting.Duration
		if duration > 366*3 {
			return nil, errors.New("duration is day, too large")
		}
		dur := 24 * time.Hour * time.Duration(duration)
		epochs := abi.ChainEpoch(dur / (time.Duration(build.BlockDelaySecs) * time.Second))
		dealPoolIInstance = &DealPool{
			ctx:      apiCtx,
			limit:    apiCtx.Config.Setting.MaxDealTransfers,
			wallet:   apiCtx.Config.Setting.Wallet,
			duration: epochs,
		}
	}
	return dealPoolIInstance, nil
}

// func (p *DealPool) watch() {
// 	nodeApi := p.ctx.FullNode
// 	waittime := time.Minute * 1
// 	dealsCol := p.ctx.MongoClient.Collection(db.DealsCollectionName)

// 	for {
// 		select {
// 		case <-p.ctx.Context.Done():
// 			logrus.Info("manually exist make deals")
// 			return
// 		default:
// 			logrus.Info("start watch deals")

// 			var dealtodo *types.DealTodo
// 			// 检查交易池中的正在执行的任务数是否最大
// 			num, err := listTransfers(p.ctx.Context, nodeApi)
// 			if err != nil {
// 				logrus.Errorf("get pending deal err: %s\n", err)
// 				goto NEXT
// 			}
// 			if num >= p.limit { // 传输达到上限, 不添加交易
// 				logrus.Warnf("%d deals transfers now, wait", num)
// 				goto NEXT
// 			}
// 			logrus.Infof("%d deals transfers now, start add next", num)

// 			// todo 发起一笔订单交易
// 			dealtodo, err = p.findOneDealTodo()
// 			if err != nil {
// 				logrus.Error(err.Error())
// 				goto NEXT
// 			}
// 			logrus.Infof("准备发单 %+v", dealtodo)
// 			if dealcid, err := buildDeal(p.ctx, dealtodo.From, dealtodo.Miner, dealtodo.Fcid, dealtodo.Price, dealtodo.Epochs); err != nil {
// 				logrus.Errorf("build deal err: %s\n", err.Error())
// 			} else {
// 				logrus.Infof("send deal success: %s", dealcid)
// 				deal, err := getDeal(p.ctx.Context, p.ctx.FullNode, dealcid)
// 				if err != nil {
// 					logrus.Errorf("get deal from daemon %s", err.Error())
// 					continue
// 				}
// 				insertDeal(dealtodo.Fcid, deal, dealsCol)
// 			}
// 		NEXT:
// 			time.Sleep(waittime)
// 		}
// 	}
// }

// func insertDeal(fcid string, d *itypes.Deal, dealCol *mongo.Collection) {
// 	price := ltypes.FIL(ltypes.BigMul(d.LocalDeal.PricePerEpoch,
// 		ltypes.NewInt(d.LocalDeal.Duration)))

// 	dat := itypes.DealInfo{
// 		PayloadCid:    fcid,
// 		DealId:        uint64(d.LocalDeal.DealID),
// 		DealCid:       d.LocalDeal.ProposalCid.String(),
// 		PieceCID:      d.LocalDeal.PieceCID.String(),
// 		MinerId:       d.LocalDeal.Provider.String(),
// 		Message:       d.LocalDeal.Message,
// 		Duration:      d.LocalDeal.Duration,
// 		PricePerEpoch: d.LocalDeal.PricePerEpoch.String(),
// 		TotalPrice:    price.String(),
// 		DealSizeBytes: d.LocalDeal.Size,
// 		IsActive:      false,
// 		DealState:     d.LocalDeal.State,
// 		DealStartAt:   d.LocalDeal.CreationTime,
// 		DealUpdateAt:  time.Now(),
// 		ODealInfo:     d.LocalDeal,
// 		OOnChain:      d.OnChain,
// 	}
// 	updateOrInsert(dealCol, dat)
// }

// func (p *DealPool) findOneDealTodo() (*types.DealTodo, error) {
// 	/*
// 	 * 1. 获数据库中的文件列表
// 	 * 2. 遍历文件
// 	 * 3. 文件和矿工号组合, 查看是否有订单
// 	 * 	3.1 如果没有, 向该矿工发单, 结束
// 	 *	3.2 如果有订单
// 	 *		3.2.1 判断是否总订单达到发单限制
// 	 *		3.2.2 判断是否对该矿工达到发单数量限制和时间限制(3d), 满足发单条件则发单
// 	 */
// 	maxDealOne := p.ctx.Config.Setting.MaxDealOne       // 单文件最大订单数目, 10
// 	minerMaxDeals := p.ctx.Config.Setting.MinerMaxDeals // 单文件对矿池最大订单数

// 	/******** 1.0 获取已经达到最大订单数的条目
// 	db.dealsinfo.aggregate([
// 		{"$group": {"_id": "$payloadcid", "count": {"$sum": 1}}},
// 		{"$match": {"count": {"$lt": 3}}},
// 		{"$group": {"_id": 0, "payloadcids": {"$push": "$_id"}}}
// 		])
// 	*/
// 	dealsCol := p.ctx.MongoClient.Collection(db.DealsCollectionName)
// 	grouppip := bson.D{{"$group", bson.D{{"_id", "$payloadcid"}, {"count", bson.D{{"$sum", 1}}}}}}
// 	matchpip := bson.D{{"$match", bson.D{{"count", bson.D{{"$gte", maxDealOne}}}}}}
// 	grouppip2 := bson.D{{"$group", bson.D{{"_id", 0}, {"payloadcids", bson.D{{"$push", "$_id"}}}}}}
// 	c1, err := dealsCol.Aggregate(context.Background(),
// 		mongo.Pipeline{grouppip, matchpip, grouppip2})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer c1.Close(context.Background())

// 	var dealWorking []bson.M
// 	if err = c1.All(context.Background(), &dealWorking); err != nil {
// 		return nil, err
// 	}
// 	// 交易表中满单的cid
// 	var payloadcids interface{}
// 	if len(dealWorking) == 0 {
// 		logrus.Info("empty dealWorking")
// 		payloadcids = []string{}
// 	} else {
// 		payloadcids = dealWorking[0]["payloadcids"]
// 	}
// 	logrus.Debugf("满单payloadcids: %s %s", payloadcids, reflect.TypeOf(payloadcids))

// 	/********** 1.1 获取文件列表中的其他数据
// 	 */
// 	// 从ipfscids表中找出可发单的cid
// 	findops := options.Find()
// 	findops.SetProjection(bson.D{{"_id", 0}, {"payloadcid", 1}})
// 	fileCol := p.ctx.MongoClient.Collection(db.IpfsCollectionName)
// 	// 查询去掉交易满单的文件
// 	c2, err := fileCol.Find(context.Background(), bson.D{{"payloadcid", bson.D{{"$nin", payloadcids}}}}, findops)
// 	if err != nil {
// 		logrus.Errorf("find  %s", err.Error())
// 		return nil, err
// 	}
// 	defer c2.Close(context.Background())
// 	// 读取所有的数据
// 	payloadcids_todo := make([]string, 0)
// 	for c2.Next(context.Background()) {
// 		var rec map[string]string
// 		if err := c2.Decode(&rec); err != nil {
// 			logrus.Errorf("decode  %s", err.Error())
// 			return nil, err
// 		} else {
// 			payloadcids_todo = append(payloadcids_todo, rec["payloadcid"])
// 		}
// 	}
// 	logrus.Debugf("cid could to send %+v", payloadcids_todo)

// 	/***	2. 遍历文件
// 	 * 		3. 文件和矿工号组合, 查看是否有订单
// 	 */
// 	for _, fcid := range payloadcids_todo {
// 		checktime := time.Now().Add(-time.Hour * time.Duration(p.ctx.Config.Setting.DealTimeout))
// 		c3, err := dealsCol.Find(context.Background(), bson.D{
// 			{"payloadcid", fcid},
// 			// {"$or", bson.A{
// 			// 	bson.D{
// 			// 		{"dealstartat", bson.D{{"$lte", checktime}}},
// 			// 		{"dealstate", storagemarket.StorageDealActive},
// 			// 	},
// 			// 	bson.D{
// 			// 		{"dealstartat", bson.D{{"$gt", checktime}}},
// 			// 	},
// 			// }},
// 		})
// 		if err != nil {
// 			logrus.Errorf("find playloadcid %s:  %s", fcid, err.Error())
// 			return nil, err
// 		}
// 		defer c3.Close(context.Background())
// 		var dealinfos []*itypes.DealInfo
// 		if err := c3.All(context.TODO(), &dealinfos); err != nil {
// 			logrus.Errorf("decode dealinfos %s:  %s", fcid, err.Error())
// 			return nil, err
// 		}
// 		// 3.2 对单个矿工订单上限判断
// 		for _, m := range trustedMiners {
// 			var sumOfDeal int
// 			// 查看矿工单量
// 			for _, d := range dealinfos {
// 				// 近期的单子,一个顶十个,为了不在短期发送发多订单
// 				if d.DealStartAt.After(checktime) {
// 					sumOfDeal += 10
// 				} else {
// 					if d.MinerId == m.Miner.String() {
// 						sumOfDeal += 1
// 					}
// 				}
// 			}
// 			if sumOfDeal < minerMaxDeals {
// 				// 满足对该矿工进行发单
// 				dealtodo := itypes.DealTodo{
// 					From:   p.wallet,
// 					Miner:  m.Miner,
// 					Fcid:   fcid,
// 					Price:  m.Price,
// 					Epochs: p.duration,
// 				}
// 				return &dealtodo, nil
// 			}
// 		}
// 	}
// 	return nil, errors.New("no fcid to done")
// }

// /*
// // 通过外部命令添加任务
// func (p *DealPool) AddTask(ops *types.MakeDealOptions) error {
// 	if p.Queue.Length() > 0 {
// 		return fmt.Errorf("can't add task, queue has %d task now", p.Queue.Length())
// 	}
// 	if ops.Duration > 366*3 {
// 		return errors.New("duration is day, too large")
// 	}
// 	dur := 24 * time.Hour * time.Duration(ops.Duration)
// 	epochs := abi.ChainEpoch(dur / (time.Duration(build.BlockDelaySecs) * time.Second))
// 	// 获取可交易的数据cid列表
// 	logrus.Infof("add %d task, max price %s", ops.Nums, ops.Price)
// 	canDeals, err := FetchCanDealCids(p.ctx, ops.Nums, ops.Price, ops.Market, ops.Sort)
// 	if err != nil {
// 		return err
// 	}
// 	// 将待执行交易放入执行队列中
// 	for i, d := range canDeals {
// 		deal := &types.DealTodo{
// 			From:   ops.Wallet,
// 			Miner:  d.Miner.Miner,
// 			Fcid:   d.PayloadCid,
// 			Price:  d.Miner.Price,
// 			Epochs: epochs,
// 		}
// 		logrus.Printf("add to queue [%d] %+v\n", i, deal)
// 		p.addToQueue(deal)
// 		logrus.Printf("queue length %d\n", p.Queue.Length())

// 	}
// 	return nil
// }

// func (p *DealPool) addToQueue(d *types.DealTodo) {
// 	p.mux.Lock()
// 	defer p.mux.Unlock()
// 	p.Queue.Add(d)
// }

// func (p *DealPool) getFromQueue() *types.DealTodo {
// 	p.mux.Lock()
// 	defer p.mux.Unlock()
// 	if p.Queue.Length() > 0 {
// 		return (p.Queue.Remove()).(*types.DealTodo)
// 	}
// 	return nil
// }

// func (p *DealPool) CleanQueue() {
// 	p.mux.Lock()
// 	defer p.mux.Unlock()
// 	p.Queue = queue.New()
// 	usedMiner = make([]string, 0)
// 	usedBest = make([]string, 0)
// }
// */

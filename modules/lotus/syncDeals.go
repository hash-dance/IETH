package lotus

import (
	"context"
	"time"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/api"
	db "github.com/guowenshuai/ieth/db/mongo"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	itypes "github.com/guowenshuai/ieth/types"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SyncDeals 列出订单
func SyncDeals(apiCtx *apicontext.APIContext) {
	ctx := apiCtx.Context
	nodeApi := apiCtx.FullNode
	waittime := time.Minute * 30

	// timer := util.Ticker(ctx, time.Minute*20)
	dealCol := apiCtx.MongoClient.Collection(db.DealsCollectionName)
	for {
		select {
		case <-ctx.Done():
			logrus.Info("exist sync deals")
			return
		default:
			logrus.Info("start sync deals")
			err := doSync(apiCtx, nodeApi, dealCol)
			if err != nil {
				logrus.Errorf("%s\n", err.Error())
			}
			time.Sleep(waittime)

		}
	}
}

func doSync(ctx *apicontext.APIContext, full api.FullNode, dealCol *mongo.Collection) error {
	// 获取dealtimeout小时以前的单子进行状态更新
	// 短期的单子不管状态更新
	// isdeal: 1, 已经离线导入的单子
	filter := bson.D{{"createdtime", bson.D{{"$lt",
		time.Now().Add(-time.Hour * time.Duration(ctx.Config.Setting.DealTimeout))}}},
		{"isdeal", 1}}
	cur, err := dealCol.Find(ctx.Context, filter)
	if err != nil {
		logrus.Errorf("dealinfo find: %s", err.Error())
		return err
	}
	defer cur.Close(context.TODO()) // 完成后关闭游标

	var results []*itypes.DealInfo
	cur.All(context.TODO(), &results)
	if err := cur.Err(); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("find deals in db %+v", results)
	// 找出所有的, dealcid
	for _, d := range results {
		dealonline, err := getDeal(ctx.Context, ctx.FullNode, d.Dealcid)
		if err != nil {
			logrus.Errorf("get deal to sync: %s", err.Error())
			continue
		}

		if err := updateOrInsert(dealCol, itypes.DealInfo{
			Dealcid:   dealonline.LocalDeal.ProposalCid.String(),
			Status:    dealonline.LocalDeal.State,
			Statusmsg: storagemarket.DealStates[dealonline.LocalDeal.State],
			Dealid:    uint64(dealonline.LocalDeal.DealID),
		}); err != nil {
			return err
		}
	}
	return nil
}

func updateOrInsert(col *mongo.Collection, d itypes.DealInfo) error {
	// return nil
	filter := bson.D{{"dealcid", d.Dealcid}}
	logrus.Infof("update deals %s %s", d.Dealcid, d.Statusmsg)
	var updatedDocument bson.M
	if err := col.FindOneAndUpdate(context.Background(), filter, bson.D{{
		"$set", bson.D{
			{"status", d.Status},        // 订单状态
			{"statusmsg", d.Statusmsg},  // 订单状态信息
			{"updatedtime", time.Now()}, // 更新时间
			{"dealid", d.Dealid},
		},
	}}).Decode(&updatedDocument); err != nil {
		logrus.Errorf("updateOrInsert.FindOneAndUpdate: %s", err.Error())
		return err
	}
	return nil
}

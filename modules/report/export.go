package report

// func Export(apiCtx *apicontext.APIContext, checkAll bool) ([]string, error) {
// 	ctx := apiCtx.Context
// 	fileCol := apiCtx.MongoClient.Collection(db.IpfsCollectionName)
// 	dealCol := apiCtx.MongoClient.Collection(db.DealsCollectionName)
// 	// 查找出在数据库中已经完成状态的订单
// 	var filter bson.D
// 	filter = bson.D{{"dealstate", storagemarket.StorageDealActive}}
// 	if checkAll {
// 		filter = bson.D{{"dealstate", bson.D{{"$in", []storagemarket.StorageDealStatus{
// 			storagemarket.StorageDealActive,
// 			storagemarket.StorageDealSealing,
// 		}}}}}
// 	}
// 	cur, err := dealCol.Find(ctx, filter)
// 	if err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	defer cur.Close(context.TODO()) // 完成后关闭游标

// 	var deals []*itypes.DealInfo
// 	cur.All(context.TODO(), &deals)
// 	if err := cur.Err(); err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	// 找出已经完成dealcid
// 	dealWithCid := make(map[string][]*itypes.DealInfo, 0)
// 	finishedDealCid := make([]string, 0)
// 	for idx, d := range deals {
// 		logrus.Infof("dealCid %s, fayloadcid %s", d.DealCid, d.PayloadCid)
// 		finishedDealCid = append(finishedDealCid, d.PayloadCid)
// 		if _, ok := dealWithCid[d.PayloadCid]; !ok {
// 			dealWithCid[d.PayloadCid] = make([]*itypes.DealInfo, 0)
// 		}
// 		dealWithCid[d.PayloadCid] = append(dealWithCid[d.PayloadCid], deals[idx])
// 	}

// 	filter = bson.D{{"payloadcid", bson.D{{"$in", finishedDealCid}}}}
// 	cur, err = fileCol.Find(ctx, filter)
// 	if err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	defer cur.Close(context.TODO()) // 完成后关闭游标

// 	var files []*itypes.IpfsData
// 	cur.All(context.TODO(), &files)
// 	if err := cur.Err(); err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	lines := make([]string, 0)
// 	lines = append(lines, fmt.Sprintf("deal_id\tdeal_cid\tminer_id\t"+
// 		"payload_cid\tfilename\tdir\tfile_format\tdeal_state\tepoch\t" +
// 		"file_size_bytes\tdeal_size_in_bytes\tdate"))
// 	for _, fl := range files {
// 		logrus.Infof("payloadcid %s", fl.PayloadCid)
// 		currentDeal := make([]*itypes.DealInfo, 0)
// 		if v, ok := dealWithCid[fl.PayloadCid]; ok {
// 			currentDeal = v
// 		}
// 		for _, d := range currentDeal {
// 			lines = append(lines, fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s", d.DealId, d.DealCid, d.MinerId,
// 				fl.PayloadCid, fl.FileName, fl.Dir, fl.FileFormat, storagemarket.DealStates[d.DealState], d.Duration,
// 				fl.FileSizeBytes, d.DealSizeBytes, d.DealStartAt))
// 		}
// 	}
// 	return lines, nil
// }

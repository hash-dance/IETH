package types

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
)

type MakeDealOptions struct {
	Wallet    string `json:"wallet"`    // 发单钱包地址
	Duplicate int64  `json:"duplicate"` // 存储几份
	Duration  int64  `json:"duration"`  // 存储时长, 天
	Price     string `json:"price"`     // 存储单价
	Market    bool   `json:"market"`    // 是否存储到市场
	Nums      int    `json:"nums"`      // 添加多少交易
	Sort      bool   `json:"sort"`      // 排序, 如果指定该参数,则升序排序
}

type Deal struct {
	LocalDeal api.DealInfo
	OnChain   *api.MarketDeal
}

type DealTodo struct {
	From   string
	Miner  address.Address
	Fcid   string
	Price  big.Int
	Epochs abi.ChainEpoch
}

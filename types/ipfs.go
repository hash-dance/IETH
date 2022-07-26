package types

import (
	"time"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
)

type IPFSPushOptions struct {
	Recursive bool   `json:"recursive"`
	Path      string `json:"path"`
}

type IPFSListOptions struct {
	Cid       string `json:"cid"`
	Recursive bool   `json:"recursive"`
}

type IpfsData struct {
	PayloadCid     string    `json:"payload_cid"`
	FileName       string    `json:"file_name"`
	Dir            string    `json:"dir"`
	IsDir          bool      `json:"is_dir"`
	FileFormat     string    `json:"file_format"`
	FileSizeBytes  int64     `json:"file_size_bytes"`
	Group          string    `json:"group"`
	CuratedDataset string    `json:"curated_dataset"`
	CreatedTime    time.Time `json:"created_time"`
	UpdatedTime    time.Time `json:"updated_time"`
}

// type DealInfo struct {
// 	PayloadCid    string                          `json:"payload_cid"`     // 文件cid
// 	DealId        uint64                          `json:"deal_id"`         // 发单机交易id
// 	DealCid       string                          `json:"deal_cid"`        // 交易cid
// 	PieceCID      string                          `json:"piece_cid"`       // PieceCID
// 	MinerId       string                          `json:"miner_id"`        // 存储矿机矿工号
// 	Message       string                          `json:"message"`         // 信息
// 	Duration      uint64                          `json:"duration"`        // 存储时间
// 	PricePerEpoch string                          `json:"price_per_epoch"` // 存储单价
// 	TotalPrice    string                          `json:"total_price"`     // 存储消耗总价
// 	DealSizeBytes uint64                          `json:"deal_size_bytes"` // 交易切片size
// 	IsActive      bool                            `json:"is_active"`       // 是否激活
// 	DealState     storagemarket.StorageDealStatus `json:"deal_state"`      // 交易状态
// 	DealStartAt   time.Time                       `json:"deal_start_at"`   // 交易发起时间
// 	DealUpdateAt  time.Time                       `json:"deal_update_at"`  // 交易更新时间
// 	ODealInfo     interface{}                     `json:"o_deal_info"`     // 原始dealinfo
// 	OOnChain      interface{}                     `json:"o_on_chain"`      // 原始onchain
// }

type DealInfo struct {
	Filecid     string                          `json:"filecid"`
	Dealcid     string                          `json:"dealcid"`
	Dealid      uint64                          `json:"dealid"`
	Miner       string                          `json:"miner"`
	Price       int                             `json:"price"`
	Duration    int                             `json:"duration"`
	Wallet      string                          `json:"wallet"`
	Isdeal      int                             `json:"isdeal"`
	Status      storagemarket.StorageDealStatus `json:"status"`
	Statusmsg   string                          `json:"statusmsg"`
	CreatedTime time.Time                       `json:"created_time"`
	UpdatedTime time.Time                       `json:"updated_time"`
}

type DealWithMiner struct {
	PayloadCid string                    `json:"payload_cid"`
	Miner      *storagemarket.StorageAsk `json:"miner"`
}

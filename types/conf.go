package types

const CONFIG = "config"

// Config Config struct
type Config struct {
	Mongodb  *Mongodb  `yaml:"mongodb" json:"mongodb"`
	Lotus    *Lotus    `yaml:"lotus" json:"lotus"`
	Ipfs     *Ipfs     `yaml:"ipfs" json:"ipfs"`
	BaseConf *BaseConf `yaml:"baseConf" json:"baseConf"`
	Setting  *Setting  `yaml:"setting" json:"setting"`
	Rule     *Rule     `yaml:"rule" json:"rule"`
}

type Rule struct {
	Disable []string   `yaml:"disable" json:"disable"`
	Best    []string   `yaml:"best" json:"best"`
	Trusted []*Trusted `yaml:"trusted" json:"trusted"`
}

type Setting struct {
	DealTimeout      int    `yaml:"dealTimeout" json:"dealTimeout"`           // 订单的有效时间
	MaxDealTransfers int    `yaml:"maxDealTransfers" json:"maxDealTransfers"` // 交易池一次的发单的上限
	MinerMaxDeals    int    `yaml:"minerMaxDeals" json:"minerMaxDeals"`
	MaxDealOne       int    `yaml:"maxDealOne" json:"maxDealOne"` // 一个文件可存储的最大数量
	Wallet           string `yaml:"wallet" json:"wallet"`
	Duration         int    `yaml:"duration" json:"duration"`
}

type Lotus struct {
	Token   string `yaml:"token" json:"token"`
	Address string `yaml:"address" json:"address"`
}
type Ipfs struct {
	Token   string `yaml:"token" json:"token"`
	Address string `yaml:"address" json:"address"`
}
type BaseConf struct {
	Repo  string `yaml:"repo" json:"repo"`   // 仓库根目录
	Debug bool   `yaml:"debug" json:"debug"` // enable debug
	Cors  bool   `yaml:"cors" json:"cors"`   // enable cors http

	LogFormat   string `yaml:"logFormat" json:"logFormat"`     // logsFormat
	LogPath     string `yaml:"logPath" json:"logPath"`         // log output path
	LogDispatch bool   `yaml:"logDispatch" json:"logDispatch"` // dispatch to different file

	HTTPListenPort int  `yaml:"httpListenPort" json:"httpListenPort"` // server port
	Monitor        bool `yaml:"monitor" json:"monitor"`               // 是否开启监控

	SessionTimeOut int64 `yaml:"sessionTimeOut" json:"sessionTimeOut"` // session timeout seconds

	SSL        bool   `yaml:"ssl" json:"ssl"`
	SSLCrtFile string `yaml:"sslCrtFile" json:"sslCrtFile"`
	SSLKeyFile string `yaml:"sslKeyFile" json:"sslKeyFile"`
}

type Mongodb struct {
	Server   string `yaml:"server" json:"server"`
	NoAuth   bool   `yaml:"noAuth" json:"noAuth"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
}

type Trusted struct {
	Miner string `yaml:"miner" json:"miner"`
	Price string `yaml:"price" json:"price"`
}

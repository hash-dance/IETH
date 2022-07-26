# IETH
发单工具
## 功能说明
1. 对接ipfs和lotus-daemon API接口
2. 将指定的文件上传到ipfs
3. 指定价格和矿工,开始发单

## 编译和依赖
环境: go1.14.6
依赖: 见go.mod

### 构建
```bash
make clean all
```
构建完成会生成`ieth`和`ieth-cmd`两个可执行程序 

## 架构描述
ieth是一个c/s架构的程序

服务器端`ieth`程序,负责维护数据库,上传文件和发单

客户端`ieth-cmd`程序,只需向服务端发出请求,服务端将会挖陈工作并返回结果

##　配置和运行
依赖: mongo:3.4

### 服务器端
`ieth`启动需要提供配置文件, 在同一级目录下的`conf.yaml`文件

启动前请先创建mongo数据库

执行如下命令启动server端, 建议挂起
```bash
./ieth
```

注意: `repo`配置参数,不同应用发单需求的服务端,请指定不同的仓库路径, 不能重复, 该配置项用于客户端识别服务端

```yaml
mongodb: # 配置mongo
  server: 127.0.0.1:27017
  noAuth: false #是否认证
  username: admin
  password: admin
  database: ieth
lotus: # lotus配置
  token: "xxxxxx"
  address: "/ip4/172.18.x.x/tcp/1234/http"
ipfs: # ipfs配置
  token: "xxxxxx"
  address: "172.18.x.x:5001"
baseConf:
  repo: ~/.ieth # 不同应用发单需求的服务端,请指定不同的仓库路径, 不能重复
  debug: true
  cors: true
  logFormat: text
  logPath: ./log
  logDispatch: false
  httpListenPort: 8088 # 服务端监听端口
  monitor: true        # 是否程序指标监控
  sessionTimeout: 30
  ssl: false
  sslCrtFile: ""
  sslKeyFile: ""
setting:
  maxDealTransfers: 20    # 指定订单传输最大值
  minerMaxDeals: 2      # 传输给每个矿工的订单数
  maxDealOne: 5         # 每个文件传输最大份数
rule:
  disable:  # 拒绝的矿工列表, 不会发单到这些矿工
    - f066596
    - f062933
  best:     # 优质矿工列表, 会先发单的这些矿工, 但是矿工必须能在存储市场查询
    - f062931
    - f022072
  trusted: # 信任矿工列表, 用来配置我们自己的矿工, 优先级高于优质矿工, 不会向信任矿工询价, 直接发单
    - miner: f0116287
      price: 0
    - miner: f010035
      price: 0
    - miner: f021255
      price: 0
```

### 客户端
`ieth-cmd`启动需要指定服务端的`repo`目录,从而访问指定的服务端程序

通过配置环境变量`IETH_PATH`指定`repo`目录,表明要访问的服务器

```bash
IETH_PATH=~/.ieth
```

`ieth-cmd`通过`--help`查看帮助信息

```bash
 ./ieth-cmd 
NAME:
   ieth-cli - everything to ipfs and filecoin!

USAGE:
   ieth-cmd [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

COMMANDS:
   ipfs     Tools for deal with ipfs node
   lotus    Tools for deal with lotus miner
   export   Tools for report
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```
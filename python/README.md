## 数据库

### cars
存储准备好的离线数据

程序参数
inputpath outputpath

```js
{
    "_id": ObjectId("60306b0ff28f4dfbb68ceb59"),
    "inputpath": "/root/data/a.tar",,  
    "outputpath": "/root/data/a.tar.car",
    "piececid": "bafykbzaceaj5kyurab6bohchm3ivyj27ixhpfg5nvrlrrwzlmvytkbitlbf3w",
    "filecid": "bafykbzaceaj5kyurab6bohchm3ivyj27ixhpfg5nvrlrrwzlmvytkbitlbf3w",
    "group": "Fedora",
    "filesizebytes": NumberLong("722425856"), 
    "dealsizebytes": NumberLong("722425856"), 
    "createdtime": ISODate("2021-02-20T01:51:11.245Z"),
    "updatedtime": ISODate("2021-02-20T01:51:11.245Z")
}
```

### offline-deals

存储已经发送的离线订单

```sh
lotus client deal --manual-piece-cid=CID --manual-piece-size=datasize  --from addr <Data CID> <miner> <price> <duration>
```

读取cars表内容, 发别向不同矿机发起离线交易

```yml
miners:
    - miner: f0148143
        price: 0
    - miner: f0155467
        price: 0
    - miner: f0392734
        price: 0
setting:
  dealTimeout: 20 ## 订单查询更新
  minerMaxDeals: 2
  maxDealNums: 10
  wallet: f3tgyeflr6s5aubvjilx4kvhqhbba5couxadsgl3xrsyjb2ay2zp3gurtmd3c22nqsztdoqg7qrkbnfh7tttlq
  duration: 366
```

```js
{
    "_id": ObjectId("60306b0ff28f4dfbb68ceb59"),
    "filecid": "bafykbzaceaj5kyurab6bohchm3ivyj27ixhpfg5nvrlrrwzlmvytkbitlbf3w",  ## 和cars中的对应关系
    "dealcid": "bafyreifb6bdckjoyblin6e4urbgnfvl5yykxoeo5odb5xmk53eo2xes4ta",  ## 发布离线交易后的订单号
    "miner": "f010035",
    "price": 0,
    "duration": 111111,
    "wallet": "f3tgyeflr6s5aubvjilx4kvhqhbba5couxadsgl3xrsyjb2ay2zp3gurtmd3c22nqsztdoqg7qrkbnfh7tttlq",
    "isdeal": 1,                                                                ## 0: 未成交; 1: 已成交
    "status": 7,
    "statusmsg": "Active",
    "createdtime": ISODate("2021-02-20T01:51:11.245Z"),
    "updatedtime": ISODate("2021-02-20T01:51:11.245Z")
}
```

### 矿机自动接单程序

监控`offlinedeals`表中, 发给本矿池的订单, `isdeal === 0`, 然后从cars表中获取原始数据, 执行import订单操作,
更新`offlinedeals.isdeal`


### offline-deals守护进程

更新`offlinedeals`中的订单状态



# 脚本运行

## makecars

制作car包

```
python3 makecars.py --input /home/xjyt/iput --output /home/xjyt/output
```


## makedeals

制作离线订单

```
python3 makedeals.py -c conf.yaml
```

## 

导入离线订单

```
python3 makedealimport.py -c import.yaml
```
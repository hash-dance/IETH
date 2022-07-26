
'''
python3 makedeals.py -c conf.yaml -n 64
'''

import subprocess
import pymongo

import logging
import base64
import sys, getopt, os, re
from datetime import datetime, timedelta
import time
import math
import yaml
 
logging.getLogger().setLevel(logging.DEBUG)
sh = logging.StreamHandler(stream=sys.stdout)    # output to standard output
sh.setFormatter(logging.Formatter("%(asctime)s [%(levelname)s] %(message)s") )
logging.getLogger().addHandler(sh)
logger = logging.getLogger()



config = {}

def Usage():
    print(''' python makecars [options]
    -h, --help
    -n, --num num of deals
    -c, --config [conf.yaml]
    ''')

def runcmd(cmds):
    logger.info("start run [%s]" % cmds)
    try:
        subp = subprocess.Popen(cmds, stderr=subprocess.PIPE, stdout=subprocess.PIPE, shell=True)
        output, err = subp.communicate()
        # 判断命令是否执行成功
        status = 1 if err else 0
        ret = None
        if status == 0:
            ret = str(output, encoding = "utf-8")
        else:
            ret = str(err, encoding = "utf-8")
        return ret.strip()
            # print("失败")
    except Exception as e:
        logger.error("cmds [%s] exec err %s" % (cmds, str(e)))
        return None

def parseretlines(retlines):
    try:
        resjson = {}
        for ln in retlines.split("\n"):
            lns = list(map(lambda x: x.strip(), ln.split(":")))
            resjson[lns[0]] = lns[1]
        return resjson
    except Exception as e:
        logger.warning("%s" % str(e))
        return retlines.split("\n")

def countWordInArr(word, arr):
    count = 0
    for item in arr:
        if word == item:
            count += 1
    return count

def getHeight():
    return runcmd("lotus chain list --count 1 | awk  -F\":\" '{print $1}'")



def makeadeal(car, miner):
    wallet = config['setting']['wallet']
    duration = config['setting']['duration']
    dealcid = ""
    logger.info("start make deal %s to miner %s %f" % (car['filecid'], miner['miner'], miner['price']))

    current = int(getHeight())
    logger.info("current is %d" % current)
    current += 2880*6
    ret = runcmd("lotus client deal --start-epoch %d --manual-piece-cid=%s --manual-piece-size=%d --from %s %s %s %f %d" % 
        (current, car['piececid'], car['dealsizebytes'], wallet, car['filecid'],  miner['miner'], miner['price'], duration))
    if ret == None:
        logger.error("make deal err")
        return None
    else:
        if "failed" in ret:
            logger.error("make deal err %s" % (ret))
            return None
        else:
            dealcid = ret
    logger.info("make deal success %s" % dealcid)
    return {
        "filecid": car['filecid'],  ## 和cars中的对应关系
        "dealcid": dealcid,  ## 发布离线交易后的订单号
        "miner": miner['miner'],
        "price": miner['price'],
        "duration": config['setting']['duration'],
        "wallet": config['setting']['wallet'],
        "isdeal": 0,               ## 0: 未成交; 1: 已成交
        "status": 0,
        "statusmsg": "",
        "createdtime": datetime.now(),
        "updatedtime": datetime.now(),
    }


def makeofflinedeal(db):
    maxDealNums = config['setting']["maxDealNums"]
    minerMaxDeals = config['setting']["minerMaxDeals"]
    miners = config['miners']
    finishedFilecids = []
    ## 获取发满的单子
    for item in db.offlinedeals.aggregate([{"$group": {"_id": "$filecid", "miners": {"$push": "$miner"}}}]):
        if len(item['miners']) >= maxDealNums:
            finishedFilecids.append(item['_id'])

    for car in db.cars.find({"filecid": {"$nin": finishedFilecids}}):
        for miner in miners:
            car2miner = list(db.offlinedeals.find({"filecid": car['filecid'], "miner": miner['miner']}))
            if len(car2miner) < minerMaxDeals:
                dealinfo = makeadeal(car, miner)
                if dealinfo != None:
                    db.offlinedeals.insert_one(dealinfo) ## 插入到数据库
                    return

if __name__ == '__main__':

    count = 64
    try:
        opts, args = getopt.getopt(sys.argv[1:],"hc:n:",["help","config=","num="])
    except getopt.GetoptError:
        Usage()
        sys.exit(2)
    for opt, arg in opts:
        if opt in('-h', '--help'):
            Usage()
            sys.exit(2)
        elif opt in ("-n", "--num"):
            count = int(arg)
        elif opt in ("-c", "--config"):
            f = "conf.yaml"
            if arg != "":
                f = arg
            if not os.path.exists(f):
                logger.error("config file not exist %s" % f)
            file = open(f, 'r', encoding="utf-8")
            #使用文件对象作为参数
            config = yaml.full_load(file)      

    logger.info("config => %s\n" % config)

    # 客户端连接
    myclient = pymongo.MongoClient("mongodb://%s:%s@%s" % (config["mongodb"]["username"], config["mongodb"]["password"], config["mongodb"]["server"]))
    dbcli = myclient[config["mongodb"]["database"]]
    dbcli.cars.create_index([("inputpath", 1), ("outputpath",1), ("filecid", 1)], unique=True)
    dbcli.offlinedeals.create_index([("dealcid", 1)], unique=True)

    
    while count > 0:
        logger.info("loops")
        makeofflinedeal(dbcli)
        time.sleep(1)
        count -= 1

    
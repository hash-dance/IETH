
'''
python3 makedealimport.py -c import.yaml
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

def makedealimport(db):
    # 离线订单 > 25个再做交易
    dealstodo = list(db.offlinedeals.find({"miner": config['miner'], "isdeal": 0}))
    if len(dealstodo) < 25:
        logger.warning("dealstodo less than 25")
        return

    # 执行所有待做的交易
    # lotus-miner storage-deals import-data <dealCid> <carFilePath>
    for deal in dealstodo:
        car = db.cars.find_one({"filecid": deal['filecid']})
        cmd = "lotus-miner storage-deals import-data %s %s" % (deal['dealcid'], car['outputpath'])
        logger.info("start import deal [%s]" % cmd)
        ret = runcmd(cmd)
        db.offlinedeals.find_one_and_update({"dealcid": deal['dealcid']}, 
            {"$set": {"isdeal": 1, "updatedtime": datetime.now()}})
        break

if __name__ == '__main__':
    f = "import.yaml"
    try:
        opts, args = getopt.getopt(sys.argv[1:],"hc:",["help","config="])
    except getopt.GetoptError:
        Usage()
        sys.exit(2)
    for opt, arg in opts:
        if opt in('-h', '--help'):
            Usage()
            sys.exit(2)
        elif opt in ("-c", "--config"):
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


    while True:
        logger.info("loops")
        makedealimport(dbcli)
        time.sleep(1)
        
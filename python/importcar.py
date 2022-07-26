
'''
python3 makecars.py --input /home/xjyt/iput --output /home/xjyt/output
'''

import subprocess
import pymongo

import logging
import base64
import sys, getopt, os, re
from datetime import datetime, timedelta
import time
import math

logging.getLogger().setLevel(logging.DEBUG)
sh = logging.StreamHandler(stream=sys.stdout)    # output to standard output
sh.setFormatter(logging.Formatter("%(asctime)s [%(levelname)s] %(message)s") )
logging.getLogger().addHandler(sh)
logger = logging.getLogger()



inputpath = ""
outputpath = ""
mongaddr = "192.168.1.157:27037"
mongouname = "admin"
mongopass = "admin"
mongodb = "iso-offline-deals"


def Usage():
    print(''' python makecars [options]
    -i, --input input path
    -o, --output output path
    -a, --addr mongo addr
    -u, --uname mongo username
    -p, --pass mongo password
    -d, --db mongo db name
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


def calDealsize(x):
    return 127*( math.pow(2, math.ceil( math.log( math.ceil ( x /127 ), 2 ) ) ) )


# 1. gen car
# 2. gen piececid
# 3. cal dealsize
# 4. gen filecid
def importcar(inputfile, outputfile):
    logger.info("start import car gen cid %s" % (outputfile))
    ret = runcmd("lotus client import %s" % outputfile)
    if ret == None:
        logger.error("gen filecid err")
        return None
    else:
        formatret = parseretlines(ret)
        logger.info("gen filecid ret %s", formatret)
        for ln in formatret:
            print("========%s" % ln)
            matchObj = matchObj = re.search(r'Root\s+(\w+)', ln, re.I)
            print(matchObj)
            if matchObj:
                filecid = matchObj.group(1)
                break
        if filecid == "":
            logger.error( "parse filecid err")
            return None
        logger.info("import car success %s %s" % (outputfile, filecid))
    return "success"


if __name__ == '__main__':

    try:
        opts, args = getopt.getopt(sys.argv[1:],"hi:o:a:u:p:d:",["help","input=","output=","addr=","uname=","pass=","db="])
    except getopt.GetoptError:
        Usage()
        sys.exit(2)
    for opt, arg in opts:
        if opt in('-h', '--help'):
            Usage()
            sys.exit(2)
        elif opt in ("-i", "--input"):
            inputpath = arg
        elif opt in ("-o", "--output"):
            outputpath = arg
        elif opt in ("-a", "--addr"):
            mongaddr = arg
        elif opt in ("-u", "--uname"):
            mongouname = arg
        elif opt in ("-P", "--pass"):
            mongopass = arg
        elif opt in ("-d", "--db"):
            mongodb = arg

    # for subp in [inputpath, outputpath]:
    #     if not os.path.exists(subp):
    #         print("input or output %s not exists" % subp)
    #         Usage()
    #         sys.exit(2)            
    # 客户端连接
    myclient = pymongo.MongoClient("mongodb://%s:%s@%s" % (mongouname, mongopass, mongaddr))
    dbcli = myclient[mongodb]
    dbcli.cars.create_index([("inputpath", 1), ("outputpath",1), ("filecid", 1)], unique=True)

    for car in list(dbcli.cars.find()):
        print(car['outputpath'])
        importcar("", car['outputpath'])


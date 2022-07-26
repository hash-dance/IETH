
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
def gencarinfo(inputfile, outputfile):
    piececid = ""
    filecid = ""
    filesizebytes = 0
    dealsizebytes = 0
    logger.info("===============================================")
    logger.info("start gen car %s => %s" % (inputfile, outputfile))
    ret = runcmd("lotus client generate-car %s %s" % (inputfile, outputfile))
    if ret == None:
        logger.error("gen car err")
        return None
    else:
        logger.info("gen car success %s => %s" % (inputfile, outputfile))
    
    logger.info("start gen piececid %s" % (outputfile))
    ret = runcmd("lotus client commP %s" % outputfile)
    if ret == None:
        logger.error("gen piececid err")
        return None
    else:
        formatret = parseretlines(ret)
        logger.info("gen piececid ret %s", formatret)
        logger.info("gen piececid success %s" % (outputfile))
        piececid = formatret['CID']
    filesizebytes = os.path.getsize(outputfile)
    dealsizebytes = calDealsize(filesizebytes)

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
    return {
        "inputpath": inputfile,
        "outputpath": outputfile,
        "piececid": piececid,
        "filecid": filecid,
        "filesizebytes": filesizebytes,
        "dealsizebytes": dealsizebytes,
    }



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

    for subp in [inputpath, outputpath]:
        if not os.path.exists(subp):
            print("input or output %s not exists" % subp)
            Usage()
            sys.exit(2)            
    # 客户端连接
    myclient = pymongo.MongoClient("mongodb://%s:%s@%s" % (mongouname, mongopass, mongaddr))
    dbcli = myclient[mongodb]
    dbcli.cars.create_index([("inputpath", 1), ("outputpath",1), ("filecid", 1)], unique=True)



    fileList = os.listdir(inputpath)
    for subfile in fileList:
        # 构造输入输出文件全路径
        fullinputfile = os.path.join(inputpath, subfile)
        fulloutputfile = os.path.join(outputpath,  "%s.car"%(subfile))
        logger.info("start handle file %s ", fullinputfile)
        # 检查已经存在
        if dbcli.cars.find_one({"inputpath": fullinputfile}):
            logger.info("%s already handled, next" % fullinputfile)
            continue


        # 解析组信息
        group = ""
        matchObj = re.match( r'(.*)_\d+.\w+$', subfile, re.M|re.I)
        if matchObj:
            group = matchObj.group(1)
        else:
            logger.warning("file %s invalid" % fullinputfile)
            continue
        fileinfo = gencarinfo(fullinputfile,fulloutputfile)
        if fileinfo == None:
            logger.error("file %s gencarinfo err, check log" % fullinputfile)
            continue
        fileinfo['group'] = group
        fileinfo['createdtime'] = datetime.now()
        fileinfo['updatedtime'] = datetime.now()
        # 5. insert to mongo
        dbcli.cars.insert_one(fileinfo)
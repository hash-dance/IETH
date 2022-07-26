
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

mongaddr = "172.18.6.72:27077"
mongouname = "xjyt"
mongopass = "xjyt"
mongodb = "iso-offline-deals"

def Usage():
    print(''' python makecars [options]
    -a, --addr mongo addr
    -u, --uname mongo username
    -p, --pass mongo password
    -d, --db mongo db name
    ''')

if __name__ == '__main__':
    try:
        opts, args = getopt.getopt(sys.argv[1:],"ha:u:p:d:",["help","addr=","uname=","pass=","db="])
    except getopt.GetoptError:
        Usage()
        sys.exit(2)
    for opt, arg in opts:
        if opt in('-h', '--help'):
            Usage()
            sys.exit(2)
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
    # 客户端连接
    myclient = pymongo.MongoClient("mongodb://%s:%s@%s" % (mongouname, mongopass, mongaddr))
    dbcli = myclient[mongodb]
    with open("/home/gws/WORK/星际/竞赛2/竞赛3.2/iso_info.txt", "r") as fh:
        for ln in fh.readlines():
            lns = ln.split()
            info = {
                "group": lns[0].upper(),
                "logo": lns[1],
                "description": " ".join(lns[2:]),
            }

            print(info)
            dbcli.groups.insert_one(info)
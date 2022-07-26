import math


def calDealsize(x):
    return 127*( math.pow(2, math.ceil( math.log( math.ceil ( x /127 ), 2 ) ) ) )

x = 3261412864
pieceSize = calDealsize(x)
askPrice = 0.000000000400000000
print("piece size %s => %s" % (x, pieceSize))
needAtto = pieceSize*askPrice/math.pow(2,30)
needAtto2= pieceSize*askPrice/math.pow(10, 9)
print(needAtto2*math.pow(10,18))


'''

filesize: 文件大小(字节数)
askPrice: 查询矿工的单价
duration: 存储时长(天数)

# 计算实际交易文件大小(pieceSize)
function dealPieceSize(filesize) {
    return 127*( Math.pow(2, Math.ceil( Math.log2( Math.ceil ( filesize /127 ) ) ) ) )    
}

pieceSize = dealPieceSize(filesize)

# 计算单笔交易的预估单价
estimatedUnit = pieceSize*askPrice/Math.pow(10, 9)

# 计算单笔交易的预估总价
estimatedCost = estimatedUnit * duration * 2880


'''
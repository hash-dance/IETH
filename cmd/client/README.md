# Go API client for swagger

ieth 工具

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: last
- Package version: 1.0.0
- Build package: io.swagger.codegen.languages.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```golang
import "./swagger"
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ExportApi* | [**ExportReportGet**](docs/ExportApi.md#exportreportget) | **Get** /export/report | 导出报告
*IpfsApi* | [**IpfsListPost**](docs/IpfsApi.md#ipfslistpost) | **Post** /ipfs/list | 从ipfs获取cid信息
*IpfsApi* | [**IpfsPushPost**](docs/IpfsApi.md#ipfspushpost) | **Post** /ipfs/push | 存储到ipfs
*LotusApi* | [**LotusDealCleanGet**](docs/LotusApi.md#lotusdealcleanget) | **Get** /lotus/deal/clean | 清空发单
*LotusApi* | [**LotusDealPushPost**](docs/LotusApi.md#lotusdealpushpost) | **Post** /lotus/deal/push | 发单到lotus
*LotusApi* | [**LotusDealStatusGet**](docs/LotusApi.md#lotusdealstatusget) | **Get** /lotus/deal/status | 发单状态
*PriceApi* | [**PriceGet**](docs/PriceApi.md#priceget) | **Get** /price | 获取矿工价格


## Documentation For Models

 - [EmptyObject](docs/EmptyObject.md)
 - [EmptyObject1](docs/EmptyObject1.md)
 - [EmptyObject2](docs/EmptyObject2.md)
 - [EmptyObject3](docs/EmptyObject3.md)
 - [EmptyObject4](docs/EmptyObject4.md)
 - [EmptyObject5](docs/EmptyObject5.md)
 - [EmptyObject6](docs/EmptyObject6.md)
 - [InlineResponse200](docs/InlineResponse200.md)
 - [InlineResponse2001](docs/InlineResponse2001.md)
 - [InlineResponse2002](docs/InlineResponse2002.md)
 - [InlineResponse2002Data](docs/InlineResponse2002Data.md)
 - [InlineResponse2003](docs/InlineResponse2003.md)
 - [InlineResponse200Data](docs/InlineResponse200Data.md)


## Documentation For Authorization
 Endpoints do not require authorization.


## Author




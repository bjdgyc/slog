# slog

## Introduction

simple log
简单日志是一个用golang编写的简单日志类

该类库提供了`FATAL、ERROR、WARN、INFO、DEBUG` 五个日志等级

同时提供了一个记录访问日志的接口 `SetAccessFile(name, accessFile string)`
该接口可根据不同的业务写入不同的文件
该日志文件内部实现了按天切割

## Installation

`go get gopkg.in/bjdgyc/slog.v1`


## example

```go

package main

import "gopkg.in/bjdgyc/slog.v1"

func main() {
	//设置日志等级
	slog.SetLogLevel("DEBUG")
	//设置普通日志文件
	slog.SetLogfile("/var/log/info.log")
	//设置access日志文件 (自动按天切割,不包含日志等级,可以设置多个)
	slog.SetAccessFile("access", "/var/log/access.log")
	//设置订单日志文件
	slog.SetAccessFile("order", "/var/log/order.log")

	//LogRecord: 2016/07/05 18:06:00 main.go:104: [ERROR] 错误信息
	slog.Error("错误信息")
	//LogRecord: 2016/07/05 01:28:34 订单信息
	slog.Access("order", "订单信息")
}


```




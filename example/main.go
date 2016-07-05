package main

import "github.com/bjdgyc/slog"

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

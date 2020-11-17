package app

import (
	"micode.be.xiaomi.com/systech/soa/xlog"
	"micode.be.xiaomi.com/systech/soa/xrpc"
)

var (
	g_rpc *xrpc.XRpc
)

func initRpc(environ string) (err error) {

	if Config().RpcEnable == 0 {
		return
	}

	rpc, err := xrpc.NewXRpc(
		Config().GroupName,
		Config().ServiceName,
		environ,
		Config().RootNode,
		Config().LocalRegister,
	)

	if err != nil {
		xlog.Warn("init rpc failed,err:%v", err)
		return
	}

	g_rpc = rpc
	return
}

func Rpc() *xrpc.XRpc {
	return g_rpc
}

package app

import (
	"micode.be.xiaomi.com/systech/asset/xredis"
	"micode.be.xiaomi.com/systech/base/xutil"
	"reflect"
)

type RedisMgr struct {
}

func NewXRedisIns(i interface{}, name string) (redis *xredis.XRedis, err error) {
	obj := reflect.ValueOf(i)
	host := obj.FieldByName("Host")
	if host.IsValid() == false {
		return nil, xutil.NewError("host field not found")
	}
	port := obj.FieldByName("Port")
	if port.IsValid() == false {
		return nil, xutil.NewError("port field not found")
	}
	auth := obj.FieldByName("Auth")
	if auth.IsValid() == false {
		return nil, xutil.NewError("auth field not found")
	}
	connTimeout := obj.FieldByName("ConnTimeout")
	if connTimeout.IsValid() == false {
		return nil, xutil.NewError("connTimeout not found")
	}
	readTimeout := obj.FieldByName("ReadTimeout")
	if readTimeout.IsValid() == false {
		return nil, xutil.NewError("readTimeout not found")
	}
	writeTimeout := obj.FieldByName("WriteTimeout")
	if writeTimeout.IsValid() == false {
		return nil, xutil.NewError("writeTimeout not found")
	}
	maxOpenConn := obj.FieldByName("MaxOpenConn")
	if maxOpenConn.IsValid() == false {
		return nil, xutil.NewError("maxOpenConn not found")
	}

	redis = xredis.NewXRedis(name)
	if Config().KerDirect == 1 {
		xlog.Notice("connect to ker_proxy directly")
		err = redis.OpenKer(auth.String(),
			Config().LocalRegister,
			Config().GroupName,
			Config().ServiceName,
			Config().KerService,
			uint(connTimeout.Int()),
			uint(readTimeout.Int()),
			uint(writeTimeout.Int()),
			uint(maxOpenConn.Int()),
		)
	} else {
		xlog.Notice("connect to ker_proxy through lvs")
		err = redis.Open(host.String(),
			int(port.Int()),
			auth.String(),
			uint(connTimeout.Int()),
			uint(readTimeout.Int()),
			uint(writeTimeout.Int()),
			uint(maxOpenConn.Int()),
		)
	}

	return
}

func (p *RedisMgr) Close() {

}

func initRedis() (mgr *RedisMgr, err error) {
	mgr = &RedisMgr{}

	return mgr, nil
}

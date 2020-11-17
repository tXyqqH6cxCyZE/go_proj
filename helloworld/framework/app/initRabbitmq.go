package app

import (
	"micode.be.xiaomi.com/systech/asset/xmq"
	"micode.be.xiaomi.com/systech/base/xutil"
	"reflect"
)

type RabbitmqMgr struct {
}

func NewXRabbitmq(i interface{}, name string) (mq *xmq.XAmqpProduce, err error) {
	obj := reflect.ValueOf(i)
	host := obj.FieldByName("Host")
	if host.IsValid() == false {
		return nil, xutil.NewError("host field not found")
	}
	port := obj.FieldByName("Port")
	if port.IsValid() == false {
		return nil, xutil.NewError("port field not found")
	}
	user := obj.FieldByName("User")
	if user.IsValid() == false {
		return nil, xutil.NewError("user field not found")
	}
	passwd := obj.FieldByName("Passwd")
	if passwd.IsValid() == false {
		return nil, xutil.NewError("passwd field not found")
	}
	exchange := obj.FieldByName("Exchange")
	if exchange.IsValid() == false {
		return nil, xutil.NewError("exchange field not found")
	}
	exchange_type := obj.FieldByName("ExchangeType")
	if exchange.IsValid() == false {
		return nil, xutil.NewError("exchange field not found")
	}
	vhost := obj.FieldByName("Vhost")
	if vhost.IsValid() == false {
		return nil, xutil.NewError("vhost field not found")
	}
	connTimeout := obj.FieldByName("ConnTimeout")
	if connTimeout.IsValid() == false {
		return nil, xutil.NewError("conntimeout not found")
	}
	readTimeout := obj.FieldByName("ReadTimeout")
	if readTimeout.IsValid() == false {
		return nil, xutil.NewError("read timeout not found")
	}
	maxOpenConn := obj.FieldByName("MaxOpenConn")
	if maxOpenConn.IsValid() == false {
		return nil, xutil.NewError("maxOpenConn not found")
	}

	mq = xmq.NewXAmqpProduce(name, int(connTimeout.Int()), int(readTimeout.Int()))
	err = mq.Open(
		host.String(),
		int(port.Int()),
		user.String(),
		passwd.String(),
		vhost.String(),
		exchange.String(),
		exchange_type.String(),
		int(maxOpenConn.Int()),
	)

	return
}

func (p *RabbitmqMgr) Close() {

}

func initRabbitmq() (*RabbitmqMgr, error) {
	mgr := &RabbitmqMgr{}

	return mgr, nil
}

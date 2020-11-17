package app

import (
	"reflect"

	"micode.be.xiaomi.com/systech/asset/xnsq"
	"micode.be.xiaomi.com/systech/base/xutil"
)

type NsqConsumerMgr struct {
}

func NewXNsqConsumerIns(i interface{}) (nsqConsumer *xnsq.Consumer, err error) {
	obj := reflect.ValueOf(i)
	topic := obj.FieldByName("Topic")
	if topic.IsValid() == false {
		return nil, xutil.NewError("topic field not found")
	}
	channel := obj.FieldByName("Channel")
	if channel.IsValid() == false {
		return nil, xutil.NewError("channel field not found")
	}
	auth := obj.FieldByName("Auth")
	if auth.IsValid() == false {
		return nil, xutil.NewError("auth field not found")
	}

	ordered := obj.FieldByName("Ordered")
	if ordered.IsValid() == false {
		return nil, xutil.NewError("ordered field not found")
	}

	extend := obj.FieldByName("Extend")
	if extend.IsValid() == false {
		return nil, xutil.NewError("extend field not found")
	}

	lookups := obj.FieldByName("Lookups")
	if lookups.IsValid() == false {
		return nil, xutil.NewError("lookups not found")
	}
	lookupsStr, ok := lookups.Interface().([]string)
	if !ok {
		return nil, xutil.NewError("lookups convert err")
	}

	readTimeout := obj.FieldByName("ReadTimeout")
	if readTimeout.IsValid() == false {
		return nil, xutil.NewError("readTimeout not found")
	}
	maxInFlight := obj.FieldByName("MaxInFlight")
	if maxInFlight.IsValid() == false {
		return nil, xutil.NewError("maxInFlight not found")
	}

	cfg := xnsq.NewConfig()

	// 秘钥名，生产活消费时用，从mis平台获取
	cfg.AuthSecret = auth.String()

	// 有序设置，不设置true无法消费有序mq
	if ordered.String() == "true" || ordered.String() == "t" {
		cfg.EnableOrdered = true
	}

	nsqConsumer, err = xnsq.NewConsumer(topic.String(), channel.String(), cfg)
	if err != nil {
		panic(err)
	}
	nsqConsumer.LookupAddr = lookupsStr

	//设置是否是消费带extend消息的msg，即带扩展消息的topic的msg
	ext := extend.String()
	if ext == "true" || ext == "t" || ext == "1" || ext == "TRUE" || ext == "T" {
		nsqConsumer.SetConsumeExt(true)
	}

	return
}

func (p *NsqConsumerMgr) Close() {

}

func (p *NsqConsumerMgr) Start() (err error) {

	return nil
}

func initNsqConsumer() (mgr *NsqConsumerMgr, err error) {
	mgr = &NsqConsumerMgr{}

	return mgr, nil
}

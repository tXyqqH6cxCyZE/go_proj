package app

import (
	"reflect"

	"github.com/spaolacci/murmur3"

	"micode.be.xiaomi.com/systech/asset/xnsq"
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/soa/xlog"
)

type NsqProducerMgr struct {
}

func NewXNsqProducerIns(i interface{}) (nsqPro *xnsq.ProducerMgr, err error) {
	obj := reflect.ValueOf(i)
	topic := obj.FieldByName("Topic")
	if topic.IsValid() == false {
		return nil, xutil.NewError("topic field not found")
	}
	topicStr, ok := topic.Interface().([]string)
	if !ok{
		return nil, xutil.NewError("topic convert err")
	}

	auth := obj.FieldByName("Auth")
	if auth.IsValid() == false {
		return nil, xutil.NewError("auth field not found")
	}
	poolSize := obj.FieldByName("PoolSize")
	if poolSize.IsValid() == false {
		return nil, xutil.NewError("poolSize field not found")
	}

	lookups := obj.FieldByName("Lookups")
	if lookups.IsValid() == false {
		return nil, xutil.NewError("lookups not found")
	}
	lookupsStr, ok := lookups.Interface().([]string)
	if !ok{
		return nil, xutil.NewError("lookups convert err")
	}

	pubTimeout := obj.FieldByName("PubTimeout")
	if pubTimeout.IsValid() == false {
		return nil, xutil.NewError("pubTimeout not found")
	}

	cfg := xnsq.NewConfig()

	//有序生产时，以partitionKey hash 到不同分区时用
	cfg.Hasher = murmur3.New32()

	// 秘钥名，生产活消费时用，从mis平台获取
	cfg.AuthSecret = auth.String()

	// 默认2， 最大值 100， 最小值 1； 建议参考并发量 设置
	cfg.ProducerPoolSize = int(poolSize.Int())

	nsqPro, err = xnsq.NewTopicProducerMgr(topicStr, cfg)
	if err != nil {
		return nil, xutil.NewError("create nsq producerMgr err : %v", err)
	}

	nsqPro.SetLogger(xlog.GetLogger(), xnsq.LogLevelInfo)

	nsqPro.AddLookupdNodes(lookupsStr)

	return
}

func (p *NsqProducerMgr) Close() {

}

func initNsqProducer() (mgr *NsqProducerMgr, err error) {
	mgr = &NsqProducerMgr{}

	return mgr, nil
}

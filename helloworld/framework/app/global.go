package app

import (
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/soa/xlog"
	"reflect"
	"sync/atomic"
	"time"
)

type GlobalDataMgr struct {
	globalAtomic atomic.Value
}

type GlobalData struct {
	*RedisMgr
	*DbMgr
	*RabbitmqMgr
	*KylinMgr
	*NsqConsumerMgr
	*NsqProducerMgr
}

type Closer interface {
	Close()
}

var (
	globalDataMgr GlobalDataMgr
)

func Global() *GlobalData {
	return globalDataMgr.getGlobal()
}

func (p *GlobalDataMgr) getGlobal() *GlobalData {
	value := p.globalAtomic.Load()
	if value == nil {
		return nil
	}

	global, ok := value.(*GlobalData)
	if !ok {
		return nil
	}

	return global
}

func (p *GlobalDataMgr) Init() (err error) {

	globalData := &GlobalData{}
	err = p.init(globalData)
	if err != nil {
		err = xutil.NewError("initApp failed, err:%v", err)
		xlog.Warn("InitApp failed, err:%v", err)
		return
	}

	//配置加载成功，进行切换
	old := p.getGlobal()
	p.globalAtomic.Store(globalData)

	//TODO 关闭老的资源, 资源必须实现 Close接口
	go func(old *GlobalData) {

		if old == nil {
			return
		}

		time.Sleep(time.Second * 3600)
		elems := reflect.ValueOf(old).Elem()
		for i := 0; i < elems.NumField(); i++ {
			field := elems.Field(i)
			if !field.CanSet() {
				continue
			}

			fieldValue := field.Interface()
			closer, ok := fieldValue.(Closer)
			if ok {
				closer.Close()
			}
		}
	}(old)
	return
}
func GlobalInit() {
	globalDataMgr = GlobalDataMgr{}
	globalData := &GlobalData{}
	globalData.DbMgr = &DbMgr{}
	globalData.RedisMgr = &RedisMgr{}
	globalData.KylinMgr = &KylinMgr{}
	globalData.NsqConsumerMgr = &NsqConsumerMgr{}
	globalData.NsqProducerMgr = &NsqProducerMgr{}
	globalDataMgr.globalAtomic.Store(globalData)
}
//在此初始化基础模块
//比如redis、db等全局对象
//在package app外，用app.Global().RedisXXXXPool的方式引用全局对象
//天然支持热加载，也就是修改redis或db的配置，无需重启程序
func (p *GlobalDataMgr) init(global *GlobalData) (err error) {

	global.RedisMgr, err = initRedis()
	if err != nil {
		return err
	}

	global.DbMgr, err = initDb()
	if err != nil {
		return err
	}

	global.RabbitmqMgr, err = initRabbitmq()
	if err != nil {
		return err
	}

	global.KylinMgr, err = initKylin()
	if err != nil {
		return err
	}

	global.NsqProducerMgr, err = initNsqProducer()
	if err != nil {
		return err
	}

	global.NsqConsumerMgr, err = initNsqConsumer()
	if err != nil {
		return err
	}

	return nil
}

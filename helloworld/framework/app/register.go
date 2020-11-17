package app

import (
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/soa/xlog"
	"micode.be.xiaomi.com/systech/soa/xregister"
)

var (
	registerHelper *RegisterHelper
)

type RegisterHelper struct {
	registerArray []*xregister.XEtcdRegister
	localIP       string
	environ       string
	cluster       string
}

func RegisterHelperInstance() *RegisterHelper {
	return registerHelper
}

func InitRegister(environ, cluster string) error {

	registerHelper = &RegisterHelper{}

	localIP, err := xutil.GetLocalIP()
	if err != nil {
		return xutil.NewError("xutil.GetLocalIP failed, err:%v", err)
	}

	registerHelper.localIP = localIP
	registerHelper.environ = environ
	registerHelper.cluster = cluster
	return registerHelper.Register()
}

func (p *RegisterHelper) Close(registerArray []*xregister.XEtcdRegister) {
	if len(registerArray) == 0 {
		return
	}

	for _, v := range registerArray {
		v.Close()
	}

	registerArray = registerArray[0:0]
}

func (p *RegisterHelper) register(etcdHost string) (r *xregister.XEtcdRegister, err error) {

	r = xregister.NewXEtcdRegister()
	err = r.Open(etcdHost, Config().RootNode, p.environ)
	if err != nil {

		xlog.Fatal(
			"register open failed, addr:%v group:%s servic:%s err:%v",
			etcdHost,
			Config().GroupName,
			Config().ServiceName,
			err,
		)
		return
	}

	err = r.SetAuth(Config().EtcdAuth)
	if err != nil {
		xlog.Fatal("set auth failed, auth:%s", Config().EtcdAuth)
		return
	}

	err = r.Register(Config().GroupName,
		Config().ServiceName, p.localIP,
		Config().ListenPort,
	)
	if err != nil {
		xlog.Fatal(
			"register failed, addr:%v group:%s servic:%s err:%v",
			etcdHost,
			Config().GroupName,
			Config().ServiceName,
			err,
		)

		return
	}

	r.Run()
	return
}

func (p *RegisterHelper) SetOffline(offline bool) {

	for _, v := range p.registerArray {
		v.SetOffline(offline)
	}
}

func (p *RegisterHelper) Register() error {

	if Config().ConfigType == "local" {
		return nil
	}

	var registerArray []*xregister.XEtcdRegister

	registerHosts := Config().CrossRegister
	for _, v := range registerHosts {

		register, err := p.register(v)
		if err != nil {
			continue
		}

		register.SetTag(p.cluster)
		registerArray = append(registerArray, register)
	}

	if len(registerArray) > 0 {
		p.Close(p.registerArray)
		p.registerArray = registerArray
	}

	return nil
}

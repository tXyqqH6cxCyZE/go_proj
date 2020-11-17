package app

import (
	"errors"
	"io/ioutil"
	"path"
	"sync/atomic"
	"time"

	"micode.be.xiaomi.com/systech/base/xconfig"
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/soa/ratelimit"
	"micode.be.xiaomi.com/systech/soa/thrift"
	"micode.be.xiaomi.com/systech/soa/xlog"
)

const (
	DefRootNode     = "/b2c_service"
	DefEnvironNode  = "/environ"
	DefClusterNode  = "/cluster"
	DefAppConfNode  = "/appconf"
	DefFlowConfNode = "/flowconf"
	DefAdminPort    = 50999
)

// ker proxy
const (
	DefKerDirect              = 0
	DefKerService             = "ker_proxy"
	DefAfterLoadConfigTimeout = 30
)

var (
	RootPath   string
	ConfigPath string
)

type AppConfig struct {
	ListenPort  int
	AdminPort   int
	ServiceName string
	GroupName   string
	LogLevel    string
	LogPath     string
	Output      string
	DoSplit     string

	ConfigAddr     string   //配置中心地址
	ConfigType     string   //配置类型
	LocalRegister  string   //跨机房注册中心配置,支持多个
	CrossRegister  []string //跨机房注册中心配置,支持多个
	RegisterCenter []string
	RootNode       string //etcd 中的根节点
	EnvironNode    string //etcd环境节点
	ClusterNode    string //etcd cluster节点
	AppConfNode    string //etcd应用配置节点
	FlowConfNode   string //etcd流量配置节点

	StatInterval int    //统计间隔
	StatAddr     string //统计上报的地址

	RpcEnable        int    //是否开启rpc模块
	KerDirect        int    //是否直连ker代理
	KerService       string // ker_proxy服务名称，默认为ker_proxy，不用特殊配置
	EtcdAuth         string
	EnableIPAuth     bool //是否开启ip白名单
	EnableMethodAuth bool
	ApiIOLog         map[string]int
	LogGlobalFlag    int
	UserConfig

	TracingEnabled      bool
	TracingSamplerType  string
	TracingSamplerParam string

	RateLimitEnabled bool
	RateLimitConf    string

	AfterLoadConfigTimeout int // hook中OnAfterLoadConfig方法超时时间 单位s
}

type XField struct {
	section         string
	sectionOrigin   string
	fieldName       string
	fieldNameOrigin string
	fieldType       string
	itemType        string
}

type ConfigMgr struct {
	configAtomic atomic.Value
	config       xconfig.XConfig
	initialize   bool
	fields       []XField
	plugin       XPlugin
}

var (
	appConfigMgr ConfigMgr
)

func Config() *AppConfig {
	return appConfigMgr.getConfig()
}

func (p *AppConfig) GetString(section, key string) (result string, err error) {

	if appConfigMgr.config == nil {
		err = errors.New("invalid config instance")
		return
	}

	result, err = appConfigMgr.config.GetString(section, key)
	return
}

func (p *AppConfig) GetInt(section, key string) (result int, err error) {

	if appConfigMgr.config == nil {
		err = errors.New("invalid config instance")
		return
	}

	result, err = appConfigMgr.config.GetInt(section, key)
	return
}

func (p *AppConfig) GetKeyList(section string) (result []string, err error) {
	if appConfigMgr.config == nil {
		err = errors.New("invalid config instance")
		return
	}
	return appConfigMgr.config.GetKeyList(section)
}

func (p *AppConfig) initConfig(xconfig xconfig.XConfig) {
	appConfigMgr.config = xconfig
}

func InitBaseConf(confPath, filename string) (err error) {

	RootPath = path.Join(confPath, "/../") + "/"
	ConfigPath = confPath
	return appConfigMgr.InitBaseConf(confPath, filename)
}

func InitAppConf(confPath, filename, environ string) (err error) {
	return appConfigMgr.InitAppConf(confPath, filename, environ)
}

func (p *ConfigMgr) readBaseConf(config *AppConfig) (err error) {

	serviceName, err := p.config.GetString("xbase", "service")
	if err != nil {
		if !p.initialize {
			return
		}

		serviceName = Config().ServiceName
	}

	config.ServiceName = serviceName

	config.GroupName, err = p.config.GetString("xbase", "group")
	if err != nil {
		if !p.initialize {
			return
		}

		config.GroupName = Config().GroupName
	}

	rootNode, err := p.config.GetString("xbase", "root_node")
	if err != nil {
		rootNode = DefRootNode
		err = nil
	}
	config.RootNode = rootNode

	environNode, err := p.config.GetString("xbase", "environ_node")
	if err != nil {
		environNode = DefEnvironNode
		err = nil
	}
	config.EnvironNode = path.Join(rootNode, environNode)

	clusterNode, err := p.config.GetString("xbase", "cluster_node")
	if err != nil {
		clusterNode = DefClusterNode
		err = nil
	}
	config.ClusterNode = path.Join(rootNode, clusterNode)

	appConfNode, err := p.config.GetString("xbase", "appconf_node")
	if err != nil {
		appConfNode = DefAppConfNode
		err = nil
	}
	config.AppConfNode = appConfNode

	flowConfNode, err := p.config.GetString("xbase", "flowconf_node")
	if err != nil {
		flowConfNode = DefFlowConfNode
		err = nil
	}
	config.FlowConfNode = flowConfNode

	configAddr, err := p.config.GetString("xbase", "config_addr")
	if err != nil {
		if !p.initialize {
			return
		}

		configAddr = Config().ConfigAddr
	}
	config.ConfigAddr = configAddr

	configType, err := p.config.GetString("xbase", "config_type")
	if err != nil {
		if !p.initialize {
			return
		}

		configType = Config().ConfigType
	}
	config.ConfigType = configType

	auth, err := p.config.GetString("xbase", "auth")
	if err != nil {
		if p.initialize {
			auth = Config().EtcdAuth
		}

		err = nil
	}

	config.EtcdAuth = auth
	return nil
}

func (p *ConfigMgr) appendRegister(conf *AppConfig, register string) {

	found := false
	for _, v1 := range conf.CrossRegister {

		if register == v1 {
			found = true
			break
		}
	}

	if !found {
		conf.CrossRegister = append(conf.CrossRegister, register)
	}

	return
}

func (p *ConfigMgr) readConf() (err error) {

	backupConfig := &AppConfig{}
	err = p.readBaseConf(backupConfig)
	if err != nil {
		return
	}

	backupConfig.ListenPort, err = p.config.GetInt("server", "port")
	if err != nil {
		return
	}

	backupConfig.AdminPort, err = p.config.GetInt("server", "admin_port")
	if err != nil {
		err = nil
		backupConfig.AdminPort = DefAdminPort
	}

	etcdAuth, _ := p.config.GetString("server", "auth")
	if len(etcdAuth) > 0 {
		backupConfig.EtcdAuth = etcdAuth
	}

	//默认开启权限验证
	backupConfig.EnableIPAuth = true
	enableIPAuth, err := p.config.GetInt("server", "ip_auth")
	if err == nil && enableIPAuth == 0 {
		backupConfig.EnableIPAuth = false
	}

	backupConfig.EnableMethodAuth = true
	enableMethodAuth, err := p.config.GetInt("server", "method_auth")
	if err == nil && enableMethodAuth == 0 {
		backupConfig.EnableMethodAuth = false
	}

	// TODO: 现在如果不写该配置, 默认是不开启tracing的, 待tracing稳定后可调整为默认开启
	tracingEnabled, err := p.config.GetString("tracing", "enabled")
	if err == nil && tracingEnabled == "true" {
		backupConfig.TracingEnabled = true
	}

	tracingSamplerType, err := p.config.GetString("tracing", "sampler_type")
	if err == nil {
		backupConfig.TracingSamplerType = tracingSamplerType
	}

	tracingSamplerParam, err := p.config.GetString("tracing", "sampler_param")
	if err == nil {
		backupConfig.TracingSamplerParam = tracingSamplerParam
	}

	// 限流配置，没有该配置，默认不开启
	backupConfig.RateLimitEnabled = ratelimit.SwitchOff
	rateLimitEnabled, err := p.config.GetString("ratelimit", "enabled")
	if err == nil && rateLimitEnabled == "true" {
		backupConfig.RateLimitEnabled = ratelimit.SwitchOn
	}

	rateLimitConf, err := p.config.GetString("ratelimit", "conf")
	if err == nil {
		backupConfig.RateLimitConf = rateLimitConf
	}

	backupConfig.RpcEnable, err = p.config.GetInt("rpc", "enable")
	if err != nil {
		backupConfig.RpcEnable = 0
		err = nil
	}

	backupConfig.KerDirect, err = p.config.GetInt("ker", "direct")
	if err != nil {
		backupConfig.KerDirect = DefKerDirect
		err = nil
	}

	backupConfig.KerService, err = p.config.GetString("ker", "ker_service")
	if err != nil || len(backupConfig.KerService) == 0 {
		backupConfig.KerService = DefKerService
	}

	var apiList []string
	var ioflag int
	if apiList, err = p.config.GetKeyList("apiiolog"); err == nil {
		flag := make(map[string]int)
		for _, api := range apiList {
			ioflag, err = p.config.GetInt("apiiolog", api)
			if err != nil {
				continue
			}
			flag[api] = ioflag
		}
		backupConfig.ApiIOLog = flag
	}

	if ioflag, err = p.config.GetInt("log", "IO_data"); err == nil {
		backupConfig.LogGlobalFlag = ioflag
	}

	backupConfig.LogLevel, err = p.config.GetString("log", "log_level")
	if err != nil {
		return
	}

	backupConfig.LogPath, err = p.config.GetString("log", "log_path")
	if err != nil {
		return
	}

	backupConfig.Output, err = p.config.GetString("log", "output")
	if err != nil {
		backupConfig.Output = "file"
		err = nil
	}

	backupConfig.DoSplit, err = p.config.GetString("log", "dosplit")
	if err != nil {
		backupConfig.DoSplit = "true"
		err = nil
	}

	localRegister, err := p.config.GetString("register", "local_register")
	if err != nil {
		return
	}

	backupConfig.LocalRegister = localRegister
	p.appendRegister(backupConfig, localRegister)

	crossRegister, err := p.config.GetArray("register", "cross_register")
	if err == nil {
		for _, v := range crossRegister {
			if inst, ok := v.(string); ok {
				p.appendRegister(backupConfig, inst)
			}
		}
	}

	centerRegister, err := p.config.GetArray("register", "register_center")
	if err == nil {
		for _, v := range centerRegister {
			if inst, ok := v.(string); ok {
				p.appendRegister(backupConfig, inst)
			}
		}
	}

	backupConfig.StatInterval, err = p.config.GetInt("stat", "stat_interval")
	if err != nil {
		return
	}
	backupConfig.StatAddr, err = p.config.GetString("stat", "stat_addr")
	if err != nil {
		return
	}

	backupConfig.AfterLoadConfigTimeout, err = p.config.GetInt("hook", "after_load_config_timeout[int]")
	if err != nil {
		err = nil
		backupConfig.AfterLoadConfigTimeout = DefAfterLoadConfigTimeout
	}

	err = p.readUserConfig(&backupConfig.UserConfig)
	if err != nil {
		return
	}

	//配置加载成功，进行切换
	p.configAtomic.Store(backupConfig)
	return nil
}

func (p *ConfigMgr) OnReload() (errRet error) {

	if p.plugin != nil {
		p.plugin.OnBeforeReloadConfig()
	}

	err := p.readConf()
	if err != nil {
		return err
	}

	g_status.loadConfTime = time.Now()

	thrift.SetIPAuthEnable(Config().EnableIPAuth)
	thrift.SetMethodAuthEnable(Config().EnableMethodAuth)
	thrift.SetStatInterval(Config().StatInterval)
	thrift.SetStatUrl(Config().StatAddr)
	thrift.SetPrintGlobalFlag(Config().LogGlobalFlag)
	thrift.SetPrintApiFlag(Config().ApiIOLog)
	xlog.SetLevel(Config().LogLevel)

	err = globalDataMgr.Init()
	if err != nil {
		errRet = err
		xlog.Warn("on reload of InitApp failed, err:%v", err)
	}

	err = ratelimit.Init(Config().RateLimitEnabled, Config().RateLimitConf)
	if err != nil {
		errRet = err
		xlog.Warn("on reload of InitApp failed, err:%v", err)
	} else {
		xlog.Notice("Reload rateLimiter success")
	}

	if p.plugin != nil {
		p.plugin.OnAfterReloadConfig()
	}

	err = RegisterHelperInstance().Register()
	if err != nil {
		errRet = xutil.NewError("initApp err:%v register:%v", errRet, err)
		xlog.Warn("on reload of register failed, err:%v", err)
	}
	return nil
}

func (p *ConfigMgr) getConfig() *AppConfig {
	value := p.configAtomic.Load()
	if value == nil {
		return nil
	}

	config, ok := value.(*AppConfig)
	if !ok {
		return nil
	}

	return config
}

func (p *ConfigMgr) InitBaseConf(confPath, filename string) (err error) {

	fullpath := path.Join(confPath, filename)
	value, err := ioutil.ReadFile(fullpath)
	if err != nil {
		err = xutil.NewError("ReadFile:%s failed, err:%v", fullpath, err)
		return
	}

	iniConfig := xconfig.NewXIniConfig()
	err = iniConfig.Unmarshal(value)
	if err != nil {
		err = xutil.NewError("ini unmarshal failed, err:%v", err)
		return
	}

	p.config = iniConfig

	configData := &AppConfig{}
	err = p.readBaseConf(configData)
	if err != nil {
		err = xutil.NewError("read base conf failed, err:%v", err)
		return
	}

	if len(configData.GroupName) == 0 {
		panic("please input group name in section[xbase] in scm_config.ini")
	}

	if len(configData.ServiceName) == 0 {
		panic("please input service name in section[xbase] in scm_config.ini")
	}
	//配置加载成功，进行切换
	p.configAtomic.Store(configData)
	p.config.Close()
	p.initialize = true

	return
}

func (p *ConfigMgr) InitAppConf(confPath, filename, environ string) (err error) {

	appConfNode := path.Join(p.getConfig().RootNode, p.getConfig().GroupName, p.getConfig().ServiceName,
		environ, p.getConfig().AppConfNode)
	flowConfNode := path.Join(p.getConfig().RootNode, p.getConfig().GroupName, p.getConfig().ServiceName,
		environ, p.getConfig().FlowConfNode)

	p.config = xconfig.NewXIniConfigWithNode([]string{flowConfNode, appConfNode})

	err = p.config.Open(confPath, filename)
	if err != nil {
		return
	}

	p.config.AddObserver(p)
	return p.readConf()
}

func ConfigInit(xconfig xconfig.XConfig) {
	appConfigMgr = ConfigMgr{}
	appConfig := &AppConfig{}
	appConfig.UserConfig = UserConfig{}
	appConfigMgr.configAtomic.Store(appConfig)
	appConfig.initConfig(xconfig)
	appConfigMgr.readConf()
}

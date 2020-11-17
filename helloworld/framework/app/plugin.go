package app

import (
	"errors"
	"micode.be.xiaomi.com/systech/soa/xlog"
	"sync"
)

type XPlugin interface {
	//程序启动后调用
	OnStartup() (err error)
	//加载配置文件之前调用
	OnBeforeLoadConfig() (err error)
	//加载配置文件之后，服务正式运行之前调用
	OnAfterLoadConfig() (err error)
	//热加载配置之前调用
	OnBeforeReloadConfig() (err error)
	//热加载配置之后调用
	OnAfterReloadConfig() (err error)
	//程序退出时调用
	OnShutdown() (err error)
}

type XPluginMgr struct {
	plugins []XPlugin
	lock    *sync.RWMutex
}

var (
	g_pluginMgr = &XPluginMgr{
		lock: new(sync.RWMutex),
	}
)

func RegisterPlugin(plugin XPlugin) (err error) {

	if plugin == nil {
		err = errors.New("invalid plugin nil")
		return
	}

	g_pluginMgr.registerPlugin(plugin)
	return
}

func (p *XPluginMgr) registerPlugin(plugin XPlugin) {

	p.lock.Lock()
	g_pluginMgr.plugins = append(g_pluginMgr.plugins, plugin)
	p.lock.Unlock()
}

func (p *XPluginMgr) OnStartup() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnStartup()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

func (p *XPluginMgr) OnBeforeLoadConfig() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnBeforeLoadConfig()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

func (p *XPluginMgr) OnAfterLoadConfig() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnAfterLoadConfig()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

func (p *XPluginMgr) OnBeforeReloadConfig() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnBeforeReloadConfig()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

func (p *XPluginMgr) OnAfterReloadConfig() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnAfterReloadConfig()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

func (p *XPluginMgr) OnShutdown() (err error) {

	p.lock.RLock()
	for i := 0; i < len(g_pluginMgr.plugins); i++ {
		err = g_pluginMgr.plugins[i].OnShutdown()
		if err != nil {
			xlog.Fatal("run plugin failed, err:%v", err)
			continue
		}
	}

	p.lock.RUnlock()
	return
}

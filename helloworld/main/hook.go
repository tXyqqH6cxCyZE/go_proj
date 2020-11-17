package main

type XMainHook struct {
}

//程序启动后调用
func (p *XMainHook) OnStartup() (err error) {
	return
}

//加载配置文件之前调用
func (p *XMainHook) OnBeforeLoadConfig() (err error) {
	return
}

//加载配置文件之后，服务正式运行之前调用
func (p *XMainHook) OnAfterLoadConfig() (err error) {
	return
}

//热加载配置之前调用
func (p *XMainHook) OnBeforeReloadConfig() (err error) {
	return
}

//热加载配置之后调用
func (p *XMainHook) OnAfterReloadConfig() (err error) {
	return
}

//程序退出时调用
func (p *XMainHook) OnShutdown() (err error) {
	return
}

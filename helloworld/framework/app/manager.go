package app

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"micode.be.xiaomi.com/systech/base/xfile"
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/crypto/xaes"
	"micode.be.xiaomi.com/systech/soa/thrift"
	"micode.be.xiaomi.com/systech/soa/xlog"
)

const (
	ErrDecodeRequest = 1001
	ErrUnsupportCmd  = 1002
)

const (
	DefAdminSecretKey = "abcdefghij123456"
	DefValidSecond    = 60
	DefPrefixNode     = "b2c_service"
	MaxWaitExitTime   = 8 * time.Second
)

var (
	serviceOfflineFile = "service_status"
	serviceOffline     = false
)

type XStatus struct {
	startTime    time.Time
	loadConfTime time.Time
}

var (
	g_status = &XStatus{}
)

func getCurPath() (curPath string, err error) {

	curPath, err = xutil.GetExecPath()
	if err != nil {
		fmt.Println("GetExecPath failed, err:", err)
		return
	}

	curPath = path.Join(curPath, "..")
	return
}

func init() {

	g_status.startTime = time.Now()
	g_status.loadConfTime = time.Now()

	curPath, err := getCurPath()
	if err != nil {
		return
	}

	fileName := path.Join(curPath, serviceOfflineFile)
	data, err := ioutil.ReadFile(fileName)
	if err == nil && len(data) > 0 {
		offlineStr := strings.TrimSpace(string(data))
		if offlineStr == "true" {
			serviceOffline = true
		}
	}
}

type AdminCommandRequest struct {
	Username string `json:"username"`
	Cmd      string `json:"cmd"`
	Time     string `json:"time"`
	Interval int    `json:"interval"`
	Filename string `json:"filename"`
}

func forceStop() {

	pid := strconv.Itoa(os.Getpid())
	//等待超时，直接kill -9 退出
	time.Sleep(MaxWaitExitTime)
	xlog.Warn("wait exit timeout[%v], kill -9 it", pid)
	exCmd := exec.Command("kill", "-SIGKILL", pid)
	exCmd.Run()
}

func restartHandle(w http.ResponseWriter) (err error) {

	pid := strconv.Itoa(os.Getpid())
	exCmd := exec.Command("kill", "-SIGHUP", pid)
	respJson(w, 0, nil)

	go forceStop()
	err = exCmd.Run()
	return
}

func offlineHandle(w http.ResponseWriter, offline bool) (err error) {

	serviceOffline = offline
	registerHelper.SetOffline(offline)
	respJson(w, 0, nil)

	status := "false"
	curPath, err := getCurPath()
	if err != nil {
		return
	}

	fileName := path.Join(curPath, serviceOfflineFile)
	if serviceOffline == true {
		status = "true"
	}

	ioutil.WriteFile(fileName, []byte(status), 0755)
	return
}

func statusHandle(w http.ResponseWriter) (err error) {

	//同步上下线状态
	registerHelper.SetOffline(serviceOffline)
	respStatusJson(w, 200, nil)
	return
}

func getBinaryMd5() (result string) {
	return Version
}

func getBuildInfo() (version, date, md5 string) {

	md5 = getBinaryMd5()
	versionPath := path.Join(RootPath, "version")
	isFile, _ := xfile.IsFile(versionPath)
	if !isFile {
		return
	}

	value, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return
	}

	data := string(value)
	lines := strings.Split(data, "\n")
	for _, v := range lines {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}

		segs := strings.Split(v, "=")
		if len(segs) < 2 {
			continue
		}

		switch v[0] {
		case 'v':
			version = segs[1]
		case 'b':
			date = segs[1]
		}
	}

	return
}

func respStatusJson(w http.ResponseWriter, code int, err error) {

	data := make(map[string]interface{})

	start_interval := time.Now().Unix() - g_status.startTime.Unix()
	data["start_time"] = g_status.startTime.Format("2006-01-02 15:04:05")
	data["start_interval"] = start_interval

	data["load_conf_time"] = g_status.loadConfTime.Format("2006-01-02 15:04:05")
	data["load_conf_interval"] = time.Now().Unix() - g_status.loadConfTime.Unix()

	lastUpdateAuth := thrift.GetLastUpdateAuthTime()
	data["load_auth_time"] = lastUpdateAuth.Format("2006-01-02 03:04:05 PM")
	data["load_auth_interval"] = time.Now().Unix() - lastUpdateAuth.Unix()

	data["ip_auth"] = Config().EnableIPAuth
	data["method_auth"] = Config().EnableMethodAuth
	data["log_level"] = Config().LogLevel
	data["offline"] = serviceOffline

	version, date, md5 := getBuildInfo()
	data["version"] = version
	data["build"] = date
	data["md5"] = md5

	//内存相关统计
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	data["alloc"] = stat.Alloc
	data["sys"] = stat.Sys
	data["heap_alloc"] = stat.HeapAlloc
	data["heap_sys"] = stat.HeapSys

	pack := make(map[string]interface{})
	pack["code"] = code
	pack["msg"] = "success"
	pack["data"] = data

	result, err := json.Marshal(pack)
	if err != nil {
		xlog.Fatal("[admin]json.Marshal failed, err:%v", err)
		return
	}

	w.Write(result)
	xlog.Debug("offline:[%v] start_interval[%v]", serviceOffline, start_interval)
}

func respJson(w http.ResponseWriter, code int, err error) {

	data := make(map[string]interface{})
	data["code"] = code
	data["msg"] = "success"

	if err != nil {
		data["msg"] = err.Error()
	}

	result, err := json.Marshal(data)
	if err != nil {
		xlog.Fatal("[admin]json.Marshal failed, err:%v", err)
		return
	}

	w.Write(result)
}

func adminHandle(w http.ResponseWriter, r *http.Request) {

	cmd, err := decodeRequest(r)
	if err != nil {
		xlog.Warn("deocode request failed, err:%v", err)
		respJson(w, ErrDecodeRequest, err)
		return
	}

	switch cmd.Cmd {
	case "offline":
		err = offlineHandle(w, true)
	case "online":
		err = offlineHandle(w, false)
	case "restart":
		err = restartHandle(w)
	case "status":
		err = statusHandle(w)
	default:
		err = fmt.Errorf("unsupported cmd[%s]", cmd.Cmd)
		respJson(w, ErrUnsupportCmd, err)
	}

	if err != nil {
		xlog.Warn("exec cmd[%s] failed, user[%s], err:%v", cmd.Cmd, cmd.Username, err)
		return
	}

	if err == nil && (cmd.Cmd == "offline" || cmd.Cmd == "online" || cmd.Cmd == "restart") {
		xlog.Notice("exec cmd[%s], user[%s]", cmd.Cmd, cmd.Username)
	}

	return
}

func RunManager(adminPort int) {

	go func(port int) {

		http.HandleFunc("/admin/cmd", adminHandle)
		portStr := fmt.Sprintf(":%d", port)

		for {
			err := http.ListenAndServe(portStr, nil)
			xlog.Fatal(err)
			time.Sleep(1 * time.Second)
		}
	}(adminPort)
}

func decodeRequest(r *http.Request) (cmd *AdminCommandRequest, err error) {

	data := r.FormValue("Data")
	if len(data) == 0 {
		err = fmt.Errorf("request data empty")
		return
	}

	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		xlog.Warn("[admin]decode string failed, err:%v", err)
		return
	}

	data, err = decrypt(dataByte, DefAdminSecretKey)
	if err != nil {
		xlog.Warn("[admin]decrypt data failed, err:%v", err)
		return
	}

	cmd = &AdminCommandRequest{}
	err = json.Unmarshal([]byte(data), cmd)
	if err != nil {
		xlog.Warn("[admin]decrypt data failed, err:%v", err)
		return
	}

	err = validateCommand(cmd)
	if err != nil {
		xlog.Warn("[admin] cmd is not valid, err:%v data[%+v] raw[%s]", err, cmd, data)
		return
	}

	return
}

func decrypt(cryptedData []byte, secretKey string) (string, error) {
	nodes := []string{
		DefPrefixNode,
		Config().GroupName,
		Config().ServiceName,
	}
	key := getEncryptKey(nodes, secretKey)
	return xaes.DecryptCBC(key, key, cryptedData)
}

func getEncryptKey(nodes []string, encryptKey string) (keyStr string) {
	node1Md5 := []rune(fmt.Sprintf("%x", md5.Sum([]byte(nodes[0]))))
	node2Md5 := []rune(fmt.Sprintf("%x", md5.Sum([]byte(nodes[1]))))
	node3Md5 := []rune(fmt.Sprintf("%x", md5.Sum([]byte(nodes[2]))))
	node4Md5 := []rune(fmt.Sprintf("%x", md5.Sum([]byte(encryptKey))))

	part1 := node1Md5[0:8]
	part2 := node2Md5[0:8]
	part3 := node3Md5[0:8]
	part4 := node4Md5[0:8]

	var key []rune
	key = append(key, part1...)
	key = append(key, part2...)
	key = append(key, part3...)
	key = append(key, part4...)

	keyStr = strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(string(key))))))
	return
}

func validateCommand(cmd *AdminCommandRequest) (err error) {

	t, err := time.Parse("2006-01-02 15:04:05", cmd.Time)
	if err != nil {
		return
	}

	now := time.Now()
	validTime := t.Add(time.Second * DefValidSecond)
	if !validTime.After(now) {
		err = fmt.Errorf("cmd[%+v] is expired, now[%v]", cmd, now)
	}

	return
}

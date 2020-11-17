package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"micode.be.xiaomi.com/systech/base/xconfig"
	"micode.be.xiaomi.com/systech/base/xutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

const (
	XDefaultHintsName  = "hints"
	XDefStructPrefix   = "auto"
	XDefRedisPrefix    = "redis_"
	XDefMySqlPrefix    = "mysql_"
	XDefPGPrefix       = "pg_"
	XDefRabbitmqPrefix = "rabbitmq_"
)

func Generate(path, filename string) (err error) {

	InitAppConf(path, filename, "test")
	return appConfigMgr.Generate()
}

var (
	redisIns    []string
	mysqlIns    []string
	pgIns       []string
	rabbitmqIns []string
)

type RedisTpl struct {
	Names []string
}
type DbTpl struct {
	Names   []string
	PGNames []string
}

type RabbitmqTpl struct {
	Names []string
}

func (p *ConfigMgr) GenerateDb() (err error) {
	t := template.New("initMysql")
	t, _ = t.Parse(`
	package app
	import (
		"micode.be.xiaomi.com/systech/soa/xdb"
		"micode.be.xiaomi.com/systech/base/xutil"
        "fmt"
		"reflect"
	)

	type DbMgr struct{
		{{range .Names}}
		{{.}} *xdb.XDb
		{{end}}

		{{range .PGNames}}
		{{.}} *xdb.XDb
		{{end}}
	}

	func NewXDb(i interface{}, name, dbType string)(db *xdb.XDb, err error){
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
		if passwd.IsValid() == false{
			return nil, xutil.NewError("passwd field not found")
		}
		database := obj.FieldByName("Database")
		if database.IsValid() == false{
			return nil, xutil.NewError("database field not found")
		}
		connTimeout := obj.FieldByName("ConnTimeout")
		if connTimeout.IsValid() == false {
			return nil, xutil.NewError("conntimeout not found")
		}
        charset := obj.FieldByName("Charset")
		if charset.IsValid() == false {
			return nil, xutil.NewError("charset not found")
		}
		maxIdleConn := obj.FieldByName("MaxIdleConn")
		if maxIdleConn.IsValid() == false {
			return nil, xutil.NewError("maxIdleConn not found")
		}
		maxOpenConn := obj.FieldByName("MaxOpenConn")
		if maxOpenConn.IsValid() == false {
			return nil, xutil.NewError("maxOpenConn not found")
		}

        var dsn string
        switch (dbType) {
        case "mysql":
            dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=%dms&parseTime=true&loc=Local",
                user.String(),
                passwd.String(),
                host.String(),
                int(port.Int()),
                database.String(),
                charset.String(),
                int(connTimeout.Int()),
            )

            db, err = xdb.NewXDbWithMysql(database.String(), dsn)
            if err != nil {
                return
            }
        case "pg":
            dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
                user.String(),
                passwd.String(),
                host.String(),
                int(port.Int()),
                database.String(),
                int(connTimeout.Int()),
            )

            db, err = xdb.NewXDbWithPG(database.String(), dsn)
            if err != nil {
                return
            }
        default:
            err = fmt.Errorf("invalid database type[%s]", dbType)
            return
        }

        db.SetMaxOpenConns(int(maxOpenConn.Int()))
        db.SetMaxIdleConns(int(maxIdleConn.Int()))

        return
	}

    func (p *DbMgr) Close() {
		{{range .Names}}
		if (p.{{.}} != nil) {
            p.{{.}}.Close()
        }
		{{end}}

		{{range .PGNames}}
		if (p.{{.}} != nil) {
            p.{{.}}.Close()
        }
		{{end}}
    }

	func initDb() (*DbMgr, error){
		mgr := &DbMgr{}
		{{range .Names}}
		{{.}}, err := NewXDb(Config().{{.}}, "{{.}}", "mysql")
		if err != nil{
			return nil, err
		}

        mgr.{{.}} = {{.}}
		{{end}}

        {{range .PGNames}}
		{{.}}, err := NewXDb(Config().{{.}}, "{{.}}", "pg")
		if err != nil{
			return nil, err
		}

        mgr.{{.}} = {{.}}
		{{end}}

		return mgr, nil
	}`)

	exePath, err := xutil.GetExecPath()
	if err != nil {
		return
	}
	filename := path.Join(exePath, "../framework/app/initDb.go")
	tpl, _ := os.Create(filename)
	v := DbTpl{Names: mysqlIns, PGNames: pgIns}
	t.Execute(tpl, v)
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func (p *ConfigMgr) GenerateRabbitmq() (err error) {
	t := template.New("initRabbitMq")
	t, _ = t.Parse(`
	package app
	import (
		"micode.be.xiaomi.com/systech/asset/xmq"
		"micode.be.xiaomi.com/systech/base/xutil"
		"reflect"
	)


	type RabbitmqMgr struct{
		{{range .Names}}
		{{.}} *xmq.XAmqpProduce
		{{end}}
	}

	func NewXRabbitmq(i interface{}, name string)(mq *xmq.XAmqpProduce, err error){
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
		if passwd.IsValid() == false{
			return nil, xutil.NewError("passwd field not found")
		}
		exchange := obj.FieldByName("Exchange")
		if exchange.IsValid() == false{
			return nil, xutil.NewError("exchange field not found")
		}
		exchange_type := obj.FieldByName("ExchangeType")
		if exchange.IsValid() == false{
			return nil, xutil.NewError("exchange field not found")
		}
		vhost := obj.FieldByName("Vhost")
		if vhost.IsValid() == false{
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
		{{range .Names}}
		if (p.{{.}} != nil) {
            p.{{.}}.Close()
        }
		{{end}}
    }

	func initRabbitmq() (*RabbitmqMgr, error){
		mgr := &RabbitmqMgr{}
		{{range .Names}}
		{{.}}, err := NewXRabbitmq(Config().{{.}}, "{{.}}")
		if err != nil{
			return nil, err
		}

        mgr.{{.}} = {{.}}
		{{end}}
		return mgr, nil
	}`)
	exePath, err := xutil.GetExecPath()
	if err != nil {
		return
	}
	filename := path.Join(exePath, "../framework/app/initRabbitmq.go")
	tpl, _ := os.Create(filename)
	v := RabbitmqTpl{Names: rabbitmqIns}
	t.Execute(tpl, v)
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func (p *ConfigMgr) GenerateRedis() (err error) {
	t := template.New("initRedis")

	t, _ = t.Parse(`
	package app
	import (
		"reflect"

		"micode.be.xiaomi.com/systech/asset/xredis"
		"micode.be.xiaomi.com/systech/base/xutil"
	)

	type RedisMgr struct{
		{{range .Names}}
		{{.}} xredis.XRedisInterface
		{{end}}
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
		{{range .Names}}
		if (p.{{.}} != nil) {
            p.{{.}}.Close()
        }
		{{end}}
    }

	func initRedis() (mgr *RedisMgr, err error){
		mgr = &RedisMgr{}
		{{range .Names}}
		mgr.{{.}}, err = NewXRedisIns(Config().{{.}}, "{{.}}")
		if err != nil{
			return nil, err
		}
		{{end}}
		return mgr, nil
	}`)
	exePath, err := xutil.GetExecPath()
	if err != nil {
		return
	}
	filename := path.Join(exePath, "../framework/app/initRedis.go")
	tpl, _ := os.Create(filename)
	v := RedisTpl{Names: redisIns}
	t.Execute(tpl, v)
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}
func (p *ConfigMgr) GenerateRedisMock() (err error) {
	t := template.New("initRedisMock")
	t, _ = t.Parse(`
	package app
	import (
	"github.com/golang/mock/gomock"
	"testing"
	"micode.be.xiaomi.com/systech/asset/xredis"
	)
	func InitGlobalRedisMockObj(t *testing.T)*xredis.MockXRedisInterface{
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		var redisMockObj *xredis.MockXRedisInterface
		redisMockObj = xredis.NewMockXRedisInterface(mockCtrl)
		GlobalInit()
		{{range .Names}}
		Global().{{.}} = redisMockObj
		{{end}}
		return redisMockObj
	}`)
	exePath, err := xutil.GetExecPath()
	if err != nil {
		return err
	}
	filename := path.Join(exePath, "../framework/app/initRedisMock.go")
	tpl, _ := os.Create(filename)
	v := RedisTpl{Names: redisIns}
	t.Execute(tpl, v)
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func (p *ConfigMgr) GenerateSqlMock() (err error) {
	t := template.New("initSqlMock")
	t, _ = t.Parse(`
	package app
	import (
	"github.com/golang/mock/gomock"
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"fmt"
	"micode.be.xiaomi.com/systech/base/xsql"
	"micode.be.xiaomi.com/systech/soa/xdb"
	)
	func InitGlobalSqlMockObj(t *testing.T) sqlmock.Sqlmock{
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB,sqlmockObj,err := sqlmock.New()
		if err != nil {
			fmt.Println("sqlmock error : ",err)
		}
		//defer mockDB.Close()
		xsqlDB := xsql.NewDb(mockDB,"sqlmock")
		db := &xdb.XDb{DB:xsqlDB}
		GlobalInit()
		{{range .Names}}
		Global().{{.}} = db
		{{end}}
		{{range .PGNames}}
		Global().{{.}} = db
		{{end}}
		return sqlmockObj
	}`)
	exePath, err := xutil.GetExecPath()
	if err != nil {
		return err
	}
	filename := path.Join(exePath, "../framework/app/initSqlMock.go")
	tpl, _ := os.Create(filename)
	v := DbTpl{Names: mysqlIns, PGNames: pgIns}
	t.Execute(tpl, v)
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}
func (p *ConfigMgr) GenerateUserConf() (err error) {
	allConfig := p.config.GetAllConfig()
	keys := make([]string, len(allConfig))
	i := 0
	for k, _ := range allConfig {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		p.generate(k, allConfig[k])
	}

	code := "package app\n\n"
	code += "import(\n"
	code += "\"micode.be.xiaomi.com/systech/base/xutil\"\n"
	code += ")\n\n"

	varCode, err := p.generateVarDefine()
	if err != nil {
		return
	}

	funcCode, err := p.generateFuncDefine()
	if err != nil {
		return
	}

	code += varCode + funcCode

	exePath, err := xutil.GetExecPath()
	if err != nil {
		return
	}
	filename := path.Join(exePath, "../framework/app/user_config.go")
	err = ioutil.WriteFile(filename, []byte(code), 0755)
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	return

}
func (p *ConfigMgr) Generate() (err error) {
	if err = p.GenerateUserConf(); err != nil {
		return fmt.Errorf("generate user conf error:%v", err)
	}
	if err = p.GenerateRedis(); err != nil {
		return fmt.Errorf("generate redis error:%v", err)
	}
	if err = p.GenerateDb(); err != nil {
		return fmt.Errorf("generate db error:%v", err)
	}
	if err = p.GenerateRabbitmq(); err != nil {
		return fmt.Errorf("generate rabbitmq error:%v", err)
	}
	if err = p.GenerateRedisMock(); err != nil {
		return fmt.Errorf("generate redis mock error:%v", err)
	}
	if err = p.GenerateSqlMock(); err != nil {
		return fmt.Errorf("generate sql mock error:%v", err)
	}
	return nil
}

func (p *ConfigMgr) generateVarDefine() (code string, err error) {

	for _, item := range p.fields {
		switch item.fieldType {
		case "struct":
			code += fmt.Sprintf("type %s%s struct {\n", XDefStructPrefix, item.fieldName)
		case "stop":
			code += "}\n\n"
		case "array":
			code += fmt.Sprintf("%s []%s\n", item.fieldName, item.itemType)
		default:
			code += fmt.Sprintf("%s %s\n", item.fieldName, item.fieldType)
		}
	}

	code += "type UserConfig struct {\n"
	for _, item := range p.fields {
		switch item.fieldType {
		case "struct":
			code += fmt.Sprintf("%s %s%s\n", item.fieldName, XDefStructPrefix, item.fieldName)
		}
	}

	code += "}\n"
	return
}

func (p *ConfigMgr) generateFuncDefine() (code string, err error) {

	for _, item := range p.fields {
		switch item.fieldType {
		case "struct":
			format := "func (p *ConfigMgr)read%sConf(conf *UserConfig)(err error) {\n\n"
			code += fmt.Sprintf(format, item.fieldName)
		case "stop":
			code += "return\n"
			code += "}\n\n"
		case "array":
			confCode, errConf := p.generateArrConf(item)
			if errConf != nil {
				err = errConf
				return
			}
			code += confCode
		default:
			confCode, errConf := p.generateConf(item)
			if errConf != nil {
				err = errConf
				return
			}
			code += confCode
		}
	}

	code += "func (p *ConfigMgr) readUserConfig(conf *UserConfig) (err error) {\n"
	code += "reserverSection := [...]string{\"xbase\", \"log\", \"server\", \"rpc\", \"stat\", \"register\"}\n\n"
	for _, item := range p.fields {
		switch item.fieldType {
		case "struct":
			code += fmt.Sprintf("err = p.read%sConf(conf)\n", item.fieldName)
			code += "if err != nil {\n"
			code += "isReserve := false\n"
			code += "for _, v := range reserverSection {\n"
			code += "if v == \"" + item.sectionOrigin + "\"{\n"
			code += "isReserve = true\n"
			code += "break\n}\n}\n\n"
			code += "if isReserve == false {\n"
			code += "return \n}\n"
			code += "}\n\n"
		}
	}

	code += "return\n"
	code += "}\n"

	return
}

func (p *ConfigMgr) generateConf(item XField) (code string, err error) {
	switch item.itemType {
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int", "uint", "int64", "uint64":
		code += fmt.Sprintf("%s, err := p.config.GetInt64(\"%s\", \"%s\")\n",
			item.fieldName, item.sectionOrigin, item.fieldNameOrigin)
	case "string":
		code += fmt.Sprintf("%s, err := p.config.GetString(\"%s\", \"%s\")\n",
			item.fieldName, item.sectionOrigin, item.fieldNameOrigin)
	case "float", "float32", "float64":
		code += fmt.Sprintf("%s, err := p.config.GetFloat64(\"%s\", \"%s\")\n",
			item.fieldName, item.sectionOrigin, item.fieldNameOrigin)
	default:
		err = xutil.NewError("unsupported type:%s", item.itemType)
		return
	}

	code += "if err != nil {\n"
	code += fmt.Sprintf("err = xutil.NewError(\"read %s.%s failed, err:%s\", err)\n", item.sectionOrigin, item.fieldNameOrigin, "%v")
	code += "return\n"
	code += "}\n"

	if item.itemType == "string" {
		code += fmt.Sprintf("conf.%s.%s = %s\n\n", item.section,
			item.fieldName, item.fieldName)
	} else {
		code += fmt.Sprintf("conf.%s.%s = %s(%s)\n\n", item.section,
			item.fieldName, item.itemType, item.fieldName)
	}

	return
}

func (p *ConfigMgr) transform(key string) string {
	var str []rune
	need := true
	for _, v := range key {
		if need {
			if v >= 'a' && v <= 'z' {
				str = append(str, v+'A'-'a')
			} else {
				str = append(str, v)
			}
			need = false
			continue
		}

		if v == '_' {
			need = true
			continue
		}

		str = append(str, v)
	}

	return string(str)
}

func (p *ConfigMgr) getHints(valueMap map[string]interface{}) (hintsMap map[string]string) {
	hintsMap = make(map[string]string)
	hints, ok := valueMap[XDefaultHintsName]
	if !ok {
		return
	}

	strHints, ok := hints.(string)
	if !ok {
		return
	}

	hintsArr := strings.Split(strHints, ",")
	for _, v := range hintsArr {
		item := strings.Split(v, ":")
		if len(item) != 2 {
			continue
		}

		key := strings.TrimSpace(item[0])
		value := strings.TrimSpace(item[1])
		hintsMap[key] = value
	}

	return
}

func (p *ConfigMgr) generateArrConf(item XField) (code string, err error) {
	switch item.itemType {
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int", "uint", "int64", "uint64":
		code += fmt.Sprintf("%s, err := p.config.GetArrayInt(\"%s\", \"%s\")\n",
			item.fieldName, item.sectionOrigin, item.fieldNameOrigin)
	case "string":
		code += fmt.Sprintf("%s, err := p.config.GetArrayString(\"%s\", \"%s\")\n",
			item.fieldName, item.sectionOrigin, item.fieldNameOrigin)
	default:
		err = xutil.NewError("unsupported type:%s", item.itemType)
		return
	}

	code += "if err != nil {\n"
	code += fmt.Sprintf("err = xutil.NewError(\"read %s.%s failed, err:%s\", err)\n", item.sectionOrigin, item.fieldNameOrigin, "%v")
	code += "return\n"
	code += "}\n"

	if item.itemType == "string" {
		code += fmt.Sprintf("conf.%s.%s = %s\n\n", item.section,
			item.fieldName, item.fieldName)
	} else {
		code += fmt.Sprintf("for _, v := range %s {\n\n", item.fieldName)
		code += fmt.Sprintf("conf.%s.%s = append(conf.%s.%s, %s(v))\n",
			item.section, item.fieldName, item.section, item.fieldName, item.itemType)
		code += "}\n"
	}

	return
}

func (p *ConfigMgr) generate(section string, value interface{}) {

	//忽略内部保留的节
	if section == xconfig.XConfigSection {
		return
	}

	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return
	}

	sectionAlias := p.transform(section)
	field := XField{
		fieldName:     sectionAlias,
		sectionOrigin: section,
		fieldType:     "struct",
	}
	if strings.HasPrefix(strings.ToLower(section), XDefRedisPrefix) {
		redisIns = append(redisIns, sectionAlias)
	}
	if strings.HasPrefix(strings.ToLower(section), XDefMySqlPrefix) {
		mysqlIns = append(mysqlIns, sectionAlias)
	}

	if strings.HasPrefix(strings.ToLower(section), XDefPGPrefix) {
		pgIns = append(pgIns, sectionAlias)
	}

	if strings.HasPrefix(strings.ToLower(section), XDefRabbitmqPrefix) {
		rabbitmqIns = append(rabbitmqIns, sectionAlias)
	}

	p.fields = append(p.fields, field)
	hintsMap := p.getHints(valueMap)

	keys := make([]string, len(valueMap))
	i := 0
	for k, _ := range valueMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := valueMap[k]
		startIndex := strings.LastIndex(k, "[")
		endIndex := strings.LastIndex(k, "]")

		var itemType string
		realKey := k
		if startIndex > 0 && endIndex > startIndex {
			realKey = k[0:startIndex]
			itemType = k[startIndex+1 : endIndex]
			itemType = strings.TrimSpace(itemType)
		}

		key := p.transform(realKey)

		if len(itemType) == 0 {
			itemType = hintsMap[k]
			if len(itemType) == 0 {
				itemType = "string"
			}
		}

		if itemType == "float" {
			itemType = "float32"
		}

		fieldType := itemType
		switch v.(type) {
		case []string:
			fieldType = "array"
		}
		field = XField{
			section:         sectionAlias,
			sectionOrigin:   section, //原始的section
			fieldName:       key,
			fieldNameOrigin: k, //原始的key
			fieldType:       fieldType,
			itemType:        itemType,
		}

		p.fields = append(p.fields, field)
	}

	field = XField{
		fieldName: section,
		fieldType: "stop",
	}

	p.fields = append(p.fields, field)
}

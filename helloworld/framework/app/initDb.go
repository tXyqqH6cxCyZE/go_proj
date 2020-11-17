package app

import (
	"fmt"
	"micode.be.xiaomi.com/systech/base/xutil"
	"micode.be.xiaomi.com/systech/soa/xdb"
	"reflect"
)

type DbMgr struct {
}

func NewXDb(i interface{}, name, dbType string) (db *xdb.XDb, err error) {
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
	database := obj.FieldByName("Database")
	if database.IsValid() == false {
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
	switch dbType {
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
}

func initDb() (*DbMgr, error) {
	mgr := &DbMgr{}
	return mgr, nil
}

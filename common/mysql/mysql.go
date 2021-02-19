package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	cfg "tx-interview/common/configparse"
	se "tx-interview/common/error"
	log "tx-interview/common/formatlog"
)

type MySQLUtil struct {
	db          *sql.DB
	initialized bool
}

var DB = MySQLUtil{db: nil, initialized: false}

func (m *MySQLUtil) InitConn() {
	m.CloseConn()

	db, err := sql.Open("mysql", cfg.GlobalConf.GetStr("mysql", "datasource"))
	if err != nil {
		log.Errorf("[mysql] 创建MySQL链接串失败")
		panic(err)
	}
	db.SetMaxOpenConns(cfg.GlobalConf.GetInt("mysql", "maxconns"))
	db.SetMaxIdleConns(cfg.GlobalConf.GetInt("mysql", "idelconns"))

	err = db.Ping()
	if err != nil {
		log.Errorf("[mysql] MySQL 链接失败")
		panic(err)
	}

	m.db = db
	m.initialized = true
	log.Infoln("[mysql] MySQL初始化成功")
}

func (m *MySQLUtil) CloseConn() {
	if m.initialized {
		m.db.Close()
		m.db = nil
		m.initialized = false
	}
}

func (m *MySQLUtil) GetConn() *sql.DB {
	if m.initialized == false {
		log.Errorln("[mysql] MySQL 还未初始化， 获取conn 失败")
		return nil
	}
	return m.db
}

func (m *MySQLUtil) GetTx() *sql.Tx {
	if m.initialized == false {
		log.Errorln("[mysql] MySQL 还未初始化， 获取tx 失败")
		return nil
	}

	tx, err := m.db.Begin()
	if err != nil {
		log.Errorf("[mysql] MySQL %s，获取tx 失败", err.Error())
		return nil
	}
	return tx
}

// 查询返回0行或者1行的数据
func (m *MySQLUtil) SingleRowQuery(sql string, args []interface{}, result ...interface{}) (int64, error) {
	if m.initialized == false {
		log.Errorln("[mysql] MySQL 还未初始化， 查询失败")
		return -1, se.DBError()
	}

	tx := m.GetTx()
	if tx == nil {
		return -1, se.New("tx is nil")
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Errorf("[mysql] MySQL SingleRowQuery prepare错误, %v", err.Error())
		tx.Rollback()
		return -1, se.DBError()
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Errorf("[mysql] MySQL SingleRowQuery query错误, %v", err.Error())
		stmt.Close()
		tx.Rollback()
		return -1, se.DBError()
	}
	var cnt int64 = 0
	for rows.Next() {
		err := rows.Scan(result...)
		if err != nil {
			log.Errorf("[mysql] MySQL SingleRowQuery rows.Next错误, %v", err.Error())
			rows.Close()
			stmt.Close()
			tx.Rollback()
			return -1, se.DBError()
		} else {
			cnt += 1
			break
		}
	}
	err = rows.Err()
	if err != nil {
		log.Errorf("[mysql] MySQL SingleRowQuery rows.Err错误, %v", err.Error())
		rows.Close()
		stmt.Close()
		tx.Rollback()
		return -1, se.DBError()
	}
	rows.Close()
	stmt.Close()
	tx.Commit()
	return cnt, nil

}

// 插入单条数据
func (m *MySQLUtil) SimpleInsert(sql string, args ...interface{}) error {
	if m.initialized == false {
		log.Errorln("[mysql] MySQL 还未初始化， 查询失败")
		return se.DBError()
	}

	tx := m.GetTx()
	if tx == nil {
		return se.New("tx is nil")
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Errorf("[mysql] MySQL SimpleInsert prepare错误, %v", err.Error())
		tx.Rollback()
		return se.DBError()
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Errorf("[mysql] MySQL SimpleInsert exec错误, %v", err.Error())
		stmt.Close()
		tx.Rollback()
		return se.DBError()
	}
	stmt.Close()
	err = tx.Commit()
	if err != nil {
		log.Errorf("[mysql] MySQL SimpleInsert commit错误, %v", err.Error())
		return se.DBError()
	}
	return nil
}

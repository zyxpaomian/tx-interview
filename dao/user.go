package dao

import (
	se "tx-interview/common/error"
	log "tx-interview/common/formatlog"
	"tx-interview/common/mysql"
	"tx-interview/structs"
)

type UserDao struct {
}

// 获取所有用户信息
func (u *UserDao) ListAllUsers() ([]*structs.User, error) {
	result := []*structs.User{}

	tx := mysql.DB.GetTx()
	if tx == nil {
		return nil, se.New("tx is nil")
	}

	sql := `select USERNAME, PASSWORD from USER;`
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Errorf("获取用户信息错误, sql: %s ,错误信息: %s", sql, err.Error())
		tx.Rollback()
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("获取用户信息错误, sql: %s ,错误信息: %s", sql, err.Error())
		stmt.Close()
		tx.Rollback()
		return nil, err
	}
	for rows.Next() {
		user := &structs.User{}
		err := rows.Scan(&user.UserName, &user.PassWord)
		if err != nil {
			log.Errorf("获取用户信息错误, sql: %s ,错误信息: %s", sql, err.Error())
			rows.Close()
			stmt.Close()
			tx.Rollback()
			return nil, err
		} else {
			result = append(result, user)
		}
	}
	rows.Close()
	stmt.Close()
	tx.Commit()

	return result, nil
}

// 认证用户密码是否正确
func (u *UserDao) UserAuth(username, password string) (*structs.User, error) {
	result := &structs.User{}
	cnt, err := mysql.DB.SingleRowQuery("SELECT USERNAME, PASSWORD FROM USER where USERNAME=? AND PASSWORD=password(?);", []interface{}{username, password}, &result.UserName, &result.PassWord)

	if err != nil {
		log.Errorln("用户认证失败")
		return nil, se.DBError()
	}
	if cnt == 0 {
		log.Errorln("用户认证失败")
		return nil, se.AuthError()
	}
	return result, nil
}

// 创建token, 默认超时时间一天
func (u *UserDao) TokenSave(token, username string) error {
	err := mysql.DB.SimpleInsert("insert into TOKEN(TOKEN, USERNAME, EXPIRE_TIME) values (?,?,date_sub(NOW(),interval -1 day))", token, username)
	if err != nil {
		log.Errorf("生成token失败, 失败原因: %v", err.Error())
		return se.DBError()
	}
	return nil
}

// 认证token是否正确
func (u *UserDao) TokenAuth(token string) (string, error) {
	result := &structs.User{}
	cnt, err := mysql.DB.SingleRowQuery("select USERNAME from TOKEN where TOKEN=? and EXPIRE_TIME >= NOW()", []interface{}{token}, &result.UserName)

	if err != nil {
		log.Errorln("token认证失败")
		return "", se.DBError()
	}
	if cnt == 0 {
		log.Errorln("token认证失败")
		return "", se.AuthError()
	}
	return result.UserName, nil
}

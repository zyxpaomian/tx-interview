package common

import (
	"encoding/json"
	"net/http"
	log "tx-interview/common/formatlog"
)

type ResMsgS struct {
	Msg string `json:"msg"`
}

func ResMsg(res http.ResponseWriter, code int, msg string) {
	resMsg_ := ResMsgS{Msg: msg}
	result := ""
	if code != 200 {
		b, err := json.Marshal(resMsg_)
		if err != nil {
			log.Errorf("[http] 生成json串失败, %v", err.Error())
			result = `{"msg": "内部错误"}`
			code = 500
		} else {
			result = string(b)
		}
	} else {
		result = msg
	}

	res.WriteHeader(code)
	res.Write([]byte(result))
}

func ParseJsonStr(str string, body interface{}) error {
	data := []byte(str)
	err := json.Unmarshal(data, body)
	if err != nil {
		log.Errorf("[http] 报文解析失败, %v", err.Error())
		return err
	} else {
		return nil
	}
}

func ReqBodyInvalid(res http.ResponseWriter) {
	ResMsg(res, 400, "请求报文格式错误")
}

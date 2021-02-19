package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"runtime"
	"tx-interview/common"
	log "tx-interview/common/formatlog"
	"tx-interview/controller"
)

type WWWMux struct {
	r *mux.Router
}

func New() *WWWMux {
	return &WWWMux{r: mux.NewRouter()}
}

func (m *WWWMux) GetRouter() *mux.Router {
	return m.r
}

// 记录日志
func AccessLogHandler(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[http] %s - %s", r.Method, r.RequestURI)
		h(w, r)
	}
}

// 用户认证
func AccessAuthHandler(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if value, ok := r.Header["Auth-Token"]; ok {
			username, err := controller.UserController.GetTokenUser(value[0])
			if err != nil {
				common.ResMsg(w, 400, "用户认证失败")
				return
			} else {
				log.Infof("[http] %s - %s - %s", r.Method, r.RequestURI, username)
			}
		} else {
			common.ResMsg(w, 400, "用户认证失败")
			return
		}
		h(w, r)
	}
}

// 注册URL映射
func (m *WWWMux) RegistURLMapping(path string, method string, needAuth bool, handle func(http.ResponseWriter, *http.Request)) {
	log.Infof("[http] URL注册映射, path: %v, method: %v, handle: %v", path, method, runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name())
	if needAuth == true {
		handle = AccessAuthHandler(handle)
	} else {
		handle = AccessLogHandler(handle)
	}
	m.r.HandleFunc(path, handle).Methods(method)
}

package handle

import (
	"tx-interview/http"
	//go_http "net/http"
)

func InitHandle(r *http.WWWMux) {
	// api相关的接口
	initAPIMapping(r)
}

func initAPIMapping(r *http.WWWMux) {
	// 用户认证
	r.RegistURLMapping("/v1/api/user/userauth", "POST", false, apiUserAuth)
	// 获取所有用户
	r.RegistURLMapping("/v1/api/user/getalluser", "GET", true, apiGetAllUser)
	// 获取所有镜像
	r.RegistURLMapping("/v1/api/user/getallimage", "GET", true, apiGetAllImage)
	// 启动容器
	r.RegistURLMapping("/v1/api/user/containercreate", "POST", false, apiRunContainer)
}

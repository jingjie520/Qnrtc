package main

import (
	"core/utils/config"
	Log "core/utils/logutil"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/rtc"
)

var (
	manager *rtc.Manager
	appID   string
)

func init() {

	//加载配置文件
	Conf := config.Conf
	if Conf == nil {
		fmt.Println("加载配置文件conf.ini失败。请检查当前目录下是否存在该文件。")
	}

	//初始化七牛云
	appID = Conf.AppID

	mac := qbox.NewMac(Conf.AccessKey, Conf.SecretKey)
	manager = rtc.NewManager(mac)
}

func main() {
	//注册处理函数，用户连接，自动调用指定的处理函数

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/api/v1.qnrtc/getToken/roomName/{roomName}/userId/{userId}", GetToken)

	//监听绑定
	http.ListenAndServe(":80", router)
}

//Index 默认首页
func Index(writer http.ResponseWriter, request *http.Request) {
	Log.RequestInfo(request)
	fmt.Fprintln(writer, "Welcome Go! ")
}

//GetToken 获取Token
func GetToken(writer http.ResponseWriter, request *http.Request) {
	Log.RequestInfo(request)

	vars := mux.Vars(request)
	roomName := vars["roomName"]
	userID := vars["userId"]

	token, err := manager.GetRoomToken(rtc.RoomAccess{AppID: appID, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 36000, Permission: "user"})

	if err == nil {
		writer.Write([]byte(token))
	} else {
		fmt.Fprintln(writer, "获取TOKEN失败。", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Unknwon/goconfig"
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
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		fmt.Println("加载配置文件conf.ini失败。请检查当前目录下是否存在该文件。")
		return
	}

	sec, err := cfg.GetSection("super")
	if err != nil {
		fmt.Println("加载配置super配置项失败。请检查配置文件。")
		return
	}

	accessKey := getConfigByKey(sec, "accessKey")
	secretKey := getConfigByKey(sec, "secretKey")
	appID = getConfigByKey(sec, "appID")

	mac := qbox.NewMac(accessKey, secretKey)
	manager = rtc.NewManager(mac)
}

func getConfigByKey(cfg map[string]string, key string) (value string) {
	if val, ok := cfg[key]; !ok {
		fmt.Println("加载配置super.", key, "失败。请检查配置文件。")
		os.Exit(0)
	} else {
		value = val
	}
	return
}

func logHandle(request *http.Request) {
	fmt.Println("Clinet:", request.RemoteAddr, "\nURL:", request.URL)
}

func main() {
	//注册处理函数，用户连接，自动调用指定的处理函数

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/api/v1.qnrtc/getToken/roomName/{roomName}/userId/{userId}", GetToken)

	//监听绑定
	http.ListenAndServe(":80", router)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	logHandle(request)
	fmt.Fprintln(writer, "Welcome Go! ")
}

func GetToken(writer http.ResponseWriter, request *http.Request) {
	logHandle(request)

	vars := mux.Vars(request)
	roomName := vars["roomName"]
	userID := vars["userId"]

	token, err := getRoomToken(appID, roomName, userID)

	if err == nil {
		writer.Write([]byte(token))
	} else {
		fmt.Fprintln(writer, "获取TOKEN失败。", err)
	}

}

func getRoomToken(appId, roomName, userID string) (token string, err error) {
	token, err = manager.GetRoomToken(rtc.RoomAccess{AppID: appId, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 36000, Permission: "user"})
	return
}

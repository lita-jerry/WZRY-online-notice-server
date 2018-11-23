package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"io/ioutil"
	"net/http"
	"strings"
)

type Data struct {
	RoleName	string 	// 昵称
	GameOnline	int		// 在线状态 0:离线 1:在线 2:正在游戏
	RoleBigIcon	string 	// 头像
	JobName		string 	// 段位
	AllStar 	int 	// 段位升级需要星星的数量
	RankingStar string 	// 当前星星
	TotalCount 	int	    // 总场数
	WinRate		string 	// 胜率
	MvpNum 		int 	// MVP数量
	RoleUrl 	string 	// 过往赛季概况的页面url
}

type ResultData struct {
	Result 	   int
	ReturnCode int
	ReturnMsg  string
	Time       string
	Data       Data
}

var r ResultData

func init() {
}

func main() {
	// 读取配置文件, 获取游戏助手的Session及其配置信息

	// 激活监听函数, 轮询角色状态
	stateChanged := make(chan bool)
	getUserStateError := make(chan string) // 防止出现错误后继续访问接口
	go lestenEventStart(stateChanged, getUserStateError)

	// 启动服务

	for {
		select {
		case <- stateChanged:
			// push
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "stateChanged: ", r.Data.GameOnline)

		case <- getUserStateError:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "getUserStateError: ", r.ReturnMsg)
			return
		}
	}
}

func lestenEventStart(c chan bool, e chan string) {

	currentState := 0

	for {
		getUserState()
		// fmt.Println(r)

		if r.Result != 0 {
			e <- r.ReturnMsg
			return
		}

		if r.Data.GameOnline != currentState {
			c <- true 
			currentState = r.Data.GameOnline
		}

		time.Sleep(60 * time.Second)
	}
	
}

func getUserState() {
	url := "https://ssl.kohsocialapp.qq.com:10001/game/rolecard"

	timestamp := strconv.FormatInt(time.Now().UnixNano() / 1e6, 10)
	fmt.Println(timestamp)
	payload := strings.NewReader("apiVersion=4&cChannelId=0&cClientVersionCode=2018092102&cClientVersionName=2.36.102&cCurrentGameId=20001&cDeviceCPU=ARM64&cDeviceId=9c46d1c18fea063ce7ca478d8691e9ca8218914b&cDeviceMem=3134406656&cDeviceModel=iPhone&cDeviceNet=WiFi&cDeviceSP=%E4%B8%AD%E5%9B%BD%E8%81%94%E9%80%9A&cDeviceScreenHeight=736&cDeviceScreenWidth=414&cGameId=20001&cGzip=1&cRand="+timestamp+"&cSystem=ios&cSystemVersionCode=12.1&cSystemVersionName=iOS&friendUserId=&gameId=20001&isMI=0&myRoleId=892934648&platType=ios&roleId=886872615&token=UUwgnsPM&userId=460403972&versioncode=2018092102")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cookie", "accessToken=15_JYoFMFjsUFK5lwjorQYbyCVstof6kBtkmjfu6xsUX5Di4--gHEcRkzDSz2JTFN9p9TcScZF8prPpqd-mT0AaE_cDRsHqv8WTnQeMuGRFINc; appOpenId=oFhrws5IUYTYRF7hnKV_9SYOgbNY")
	req.Header.Add("host", "ssl.kohsocialapp.qq.com:10001")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	// req.Header.Add("accept-encoding", "br, gzip, deflate") // 没有解压, 有乱码
	req.Header.Add("accept-encoding", "utf-8")
	req.Header.Add("connection", "keep-alive")
	req.Header.Add("accept", "*/*")
	req.Header.Add("user-agent", "smoba/2.36.102 (iPhone; iOS 12.1; Scale/3.00)")
	req.Header.Add("accept-language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	req.Header.Add("noencrypt", "1")
	req.Header.Add("content-length", "563")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	// var r ResultData
	json.Unmarshal([]byte(string(body)), &r)
	// fmt.Println(r)
	// r.Data.GameOnline = 3
}

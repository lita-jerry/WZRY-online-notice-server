package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// 读取配置文件, 获取游戏助手的Session及其配置信息
	// 启动服务
	// 激活监听函数, 轮询角色状态
	getUserState()
}

func getUserState() {
	url := "https://ssl.kohsocialapp.qq.com:10001/game/rolecard"

	payload := strings.NewReader("apiVersion=4&cChannelId=0&cClientVersionCode=2018092102&cClientVersionName=2.36.102&cCurrentGameId=20001&cDeviceCPU=ARM64&cDeviceId=9c46d1c18fea063ce7ca478d8691e9ca8218914b&cDeviceMem=3134406656&cDeviceModel=iPhone&cDeviceNet=WiFi&cDeviceSP=%E4%B8%AD%E5%9B%BD%E8%81%94%E9%80%9A&cDeviceScreenHeight=736&cDeviceScreenWidth=414&cGameId=20001&cGzip=1&cRand=1542890464800&cSystem=ios&cSystemVersionCode=12.1&cSystemVersionName=iOS&friendUserId=&gameId=20001&isMI=0&myRoleId=892934648&platType=ios&roleId=886872615&token=UUwgnsPM&userId=460403972&versioncode=2018092102")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cookie", "accessToken=15_JYoFMFjsUFK5lwjorQYbyCVstof6kBtkmjfu6xsUX5Di4--gHEcRkzDSz2JTFN9p9TcScZF8prPpqd-mT0AaE_cDRsHqv8WTnQeMuGRFINc; appOpenId=oFhrws5IUYTYRF7hnKV_9SYOgbNY")
	req.Header.Add("host", "ssl.kohsocialapp.qq.com:10001")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept-encoding", "br, gzip, deflate")
	req.Header.Add("connection", "keep-alive")
	req.Header.Add("accept", "*/*")
	req.Header.Add("user-agent", "smoba/2.36.102 (iPhone; iOS 12.1; Scale/3.00)")
	req.Header.Add("accept-language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	req.Header.Add("noencrypt", "1")
	req.Header.Add("content-length", "563")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "9e8bdaf5-b624-4a21-1c90-31c22523008a")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

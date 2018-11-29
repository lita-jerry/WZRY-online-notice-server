package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var startTime = 11 // 每天开始时间, 单位:小时
var duration = 12  // 运行持续时间, 单位:小时

var roleId = "886872615" // 用户id

var token = "SydKGR8N"
var cookie = "accessToken=16_NxZhosEUeFmVHTOt0h3_lMcKytG5SU4CPZaUdqi1uHpqoh867Jud8djQeapTaxjaf13KDZHqq6Z0ko6o7kxQ8QV0AxsCnI0648K69B2GEeo; appOpenId=oFhrws5IUYTYRF7hnKV_9SYOgbNY"

var eventloopCount = 0

var returnData ResultData
var lastState = "0"

type Data struct {
	RoleName    string      // 昵称
	GameOnline  interface{} // 在线状态 0:离线 1:在线 2:正在游戏
	RoleBigIcon string      // 头像
	JobName     string      // 段位
	AllStar     int         // 段位升级需要星星的数量
	RankingStar string      // 当前星星
	TotalCount  int         // 总场数
	WinRate     string      // 胜率
	MvpNum      int         // MVP数量
	RoleUrl     string      // 过往赛季概况的页面url
}

type ResultData struct {
	Result     int
	ReturnCode int
	ReturnMsg  string
	Time       string
	Data       Data
}

func main() {

	// 激活监听函数, 轮询角色状态
	stateChangedChan := make(chan string)
	getUserStateErrorChan := make(chan string) // 防止出现错误后继续访问接口
	stopLestenEventChan := make(chan bool)     // 停止监听用户状态

	go lestenEventStart(stateChangedChan, getUserStateErrorChan, stopLestenEventChan)

	// 启动服务
	go func() {
		http.HandleFunc("/", getStateServer)
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			log.Fatal("ListenAndServer: ", err)
		}
	}()

	// 监听goroutine消息
	for {
		select {
		case changed := <-stateChangedChan:
			// push
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "stateChanged: ", changed)
			if lastState == "0" && changed == "1" {
				go sendOnlineStateMSG()
				lastState = changed
			}

		case err := <-getUserStateErrorChan:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "getUserStateError: ", err)
			stopLestenEventChan <- true
			return
		}
	}
}

func getStateServer(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	b, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println("json err:", err)
		return
	}
	fmt.Fprint(w, string(b))
}

func lestenEventStart(changedChan chan string, errorChan chan string, stopChan chan bool) {

	// var state interface{} = "0"
	requestFinishedChan := make(chan bool)

	for {

		// 判断是否在可执行时间范围内
		currentTimeHour := time.Now().Hour()
		if !(currentTimeHour >= startTime && currentTimeHour <= (startTime+duration)) {
			time.Sleep(60 * time.Second)
			continue
		}

		var resultData ResultData

		go getUserState(&resultData, requestFinishedChan, 3)

		select {
		case <-requestFinishedChan:
			// 读数据
			// fmt.Println(resultData)

			returnData = resultData // 传给服务端返回的值

			if resultData.Result != 0 || resultData.ReturnCode != 0 {
				errorChan <- resultData.ReturnMsg
				return
			}
			// 判断是否有变化
			if currentState, ok := resultData.Data.GameOnline.(string); ok {
				if currentState != lastState {
					changedChan <- currentState
				}
			}
			// 随机间隔
			rand.Seed(time.Now().UTC().UnixNano())
			interval := rand.Intn(33) + 20
			eventloopCount += 1
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), ": ", interval, "秒后再次请求 当前循环次数:", eventloopCount, " 当前状态:", lastState)
			time.Sleep(time.Duration(interval) * time.Second)
			continue

		case <-time.After(3 * time.Second):
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), ": Timeout!")
			continue
		case <-stopChan:
			return
		}
	}

}

func getUserState(response *ResultData, requestFinishedChan chan bool, timeout int) {

	var startTime = time.Now()

	url := "https://ssl.kohsocialapp.qq.com:10001/game/rolecard"

	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	// fmt.Println(timestamp)

	payload := strings.NewReader("apiVersion=4&cChannelId=0&cClientVersionCode=2018092102&cClientVersionName=2.36.102&cCurrentGameId=20001&cDeviceCPU=ARM64&cDeviceId=9c46d1c18fea063ce7ca478d8691e9ca8218914b&cDeviceMem=3134406656&cDeviceModel=iPhone&cDeviceNet=WiFi&cDeviceSP=%E4%B8%AD%E5%9B%BD%E8%81%94%E9%80%9A&cDeviceScreenHeight=736&cDeviceScreenWidth=414&cGameId=20001&cGzip=1" +
		"&cRand=" + timestamp +
		"&cSystem=ios&cSystemVersionCode=12.1&cSystemVersionName=iOS&friendUserId=&gameId=20001&isMI=0&myRoleId=892934648&platType=ios" +
		"&roleId=" + roleId +
		"&token=" + token +
		"&userId=460403972&versioncode=2018092102")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cookie", cookie)
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

	// var response ResultData
	json.Unmarshal([]byte(string(body)), response)

	// 解析gameOnline, 并转成string类型
	r_state := "0"
	if _state, ok := response.Data.GameOnline.(int); ok {
		r_state = strconv.Itoa(_state)
	} else if _state, ok := response.Data.GameOnline.(float64); ok {
		r_state = strconv.FormatFloat(_state, 'G', -1, 64)
	} else if _state, ok := response.Data.GameOnline.(string); ok {
		r_state = _state
	}
	response.Data.GameOnline = r_state

	// 判断是否超时
	var endTime = time.Now()
	if endTime.Sub(startTime).Seconds() < float64(timeout) {
		requestFinishedChan <- true
	}
}

func sendOnlineStateMSG() {
	url := "https://u.ifeige.cn/api/message/send"

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	payload := strings.NewReader("{\"secret\":\"31a6732b116424c43e35e371e079459d\",\"app_key\":\"4ad89f97be527df43992c53c27390edf\",\"template_id\":\"odbtM60xeWsGXJEnOn6XL7MzvOJ3mMkkvP0lruBW4Og\",\"url\":\"\",\"data\":{\"first\":{\"value\":\"楚慷慷上线啦！\",\"color\":\"#173177\"},\"keyword1\":{\"value\":\"+我yi个\",\"color\":\"#173177\"},\"keyword2\":{\"value\":\"" + currentTime + "\",\"color\":\"#173177\"},\"keyword3\":{\"value\":\"唐山\",\"color\":\"#173177\"},\"keyword4\":{\"value\":\"\",\"color\":\"#173177\"},\"remark\":{\"value\":\"\",\"color\":\"#173177\"}}}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "873a72e7-2d34-4b2c-a1ab-c5ea543e6a8c")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

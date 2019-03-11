package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"coding.net/cherrysd/wxserver/message"
	"coding.net/cherrysd/wxserver/util"
	"encoding/json"
	"time"
)

// Server 微信后台实例
type Server struct {
	checkToken string
	appid      string
	appsecret  string
	tokenInfo  AccessTokenInfo
	handleMap  map[HandleType]interface{}
}

// HandleType 消息处理器类型
type HandleType string

// 消息处理器类型枚举
const (
	RawHandle        = "RawHandle"
	TextHandle       = "TextHandle"
	ImageHandle      = "ImageHandle"
	VoiceHandle      = "VoiceHandle"
	VideoHandle      = "VideoHandle"
	ShortVideoHandle = "ShortVideoHandle"
	EventHandle      = "EventHandle"
)

// NewServer 创建底层服务实例
func NewServer(checkToken string) *Server {
	newServer := new(Server)
	newServer.checkToken = checkToken
	newServer.handleMap = make(map[HandleType]interface{})
	return newServer
}

func (svr *Server) init() {
	log.Println("Server Init")
}

func (svr *Server) GetAccessToken() string {
	if svr.tokenInfo.ErrCode != 0 || time.Now().Sub(svr.tokenInfo.lastTime).Seconds() > float64(svr.tokenInfo.ExpiresIn) {
		//token超时了
		if svr.appid == "" || svr.appsecret == "" {
			return ""
		}
		svr.updateAccessToken(svr.appid, svr.appsecret)
	}
	return svr.tokenInfo.AccessToken
}

func (svr *Server) updateAccessToken(appid string, appsecret string) {
	svr.tokenInfo = AccessTokenInfo{}
	url := fmt.Sprintf(AccessTokenURL, appid, appsecret)
	response, err := http.Get(url)

	if err != nil {
		return
	}

	result, responseError := ioutil.ReadAll(response.Body)
	if responseError != nil {
		return
	}

	resultJSON := AccessTokenInfo{}
	jsonError := json.Unmarshal(result, &resultJSON)
	svr.tokenInfo = resultJSON
	svr.tokenInfo.lastTime = time.Now()
	if jsonError != nil {
		return
	}

	if resultJSON.ErrCode == 0 && resultJSON.AccessToken != "" {
		return
	}
	return
}

// SetAppInfo 设置appid与appsecret
func (svr *Server) SetAppInfo(appid string, appsecret string) {
	svr.appid = appid
	svr.appsecret = appsecret
}

// AppHandle 微信公众号消息入口
func (svr *Server) serverHandle(w http.ResponseWriter, r *http.Request) {
	contentBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Content Error")
		return
	}
	var requestMsg message.RawMessage
	requestMsg, err = message.ParseMsg(contentBytes)
	if err != nil {
		log.Println("Parse Request Message Error")
		return
	}

	if svr.handleMap[RawHandle] == nil {
		svr.defaultMessageHandle(requestMsg, w)
	} else {
		handleFunc := svr.handleMap[RawHandle].(func(message.RawMessage, http.ResponseWriter))
		handleFunc(requestMsg, w)
	}
}

// ConnectServer 与wx公众号第一次连接，给腾讯做校验使用
func (svr *Server) ConnectServer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signatureIn := strings.Join(r.Form["signature"], "")
	if util.CheckSignature(svr.checkToken, timestamp, nonce, signatureIn) {
		echostr := strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
	} else {
		fmt.Fprintf(w, "")
	}
}

func (svr *Server) defaultMessageHandle(msg message.RawMessage, w http.ResponseWriter) {
	switch msg.MsgType {
	case message.TextMsg:
		if svr.handleMap[TextHandle] != nil {
			text := new(message.Text)
			text.FromUserName = msg.FromUserName
			text.ToUserName = msg.ToUserName
			text.Content = msg.Content
			text.CreateTime = msg.CreateTime
			text.MsgID = msg.MsgID
			handleFunc := svr.handleMap[TextHandle].(func(*message.Text, http.ResponseWriter))
			handleFunc(text, w)
		}
	case message.ImageMsg:
		if svr.handleMap[ImageHandle] != nil {
			image := new(message.Image)
			image.FromUserName = msg.FromUserName
			image.ToUserName = msg.ToUserName
			image.MediaID = msg.MediaID
			image.CreateTime = msg.CreateTime
			image.MsgID = msg.MsgID
			handleFunc := svr.handleMap[ImageHandle].(func(*message.Image, http.ResponseWriter))
			handleFunc(image, w)
		}
	case message.VoiceMsg:
		if svr.handleMap[VoiceHandle] != nil {
			voice := new(message.Voice)
			voice.FromUserName = msg.FromUserName
			voice.ToUserName = msg.ToUserName
			voice.MediaID = msg.MediaID
			voice.CreateTime = msg.CreateTime
			voice.MsgID = msg.MsgID
			handleFunc := svr.handleMap[VoiceHandle].(func(*message.Voice, http.ResponseWriter))
			handleFunc(voice, w)
		}
	case message.ShortvideoMsg:
		if svr.handleMap[ShortVideoHandle] != nil {
			//todo
		}
	case message.VideoMsg:
		if svr.handleMap[VideoHandle] != nil {
			video := new(message.Video)
			video.FromUserName = msg.FromUserName
			video.ToUserName = msg.ToUserName
			video.MediaID = msg.MediaID
			video.CreateTime = msg.CreateTime
			video.MsgID = msg.MsgID
			video.Title = msg.Title
			video.Description = msg.Description
			handleFunc := svr.handleMap[VideoHandle].(func(*message.Video, http.ResponseWriter))
			handleFunc(video, w)
		}
	case message.EventMsg:
		if svr.handleMap[EventHandle] != nil {
			event := new(message.Event)
			event.FromUserName = msg.FromUserName
			event.ToUserName = msg.ToUserName
			event.CreateTime = msg.CreateTime
			event.Event = msg.Event
			event.EventKey = msg.Event
			event.Longitude = msg.Longitude
			event.Latitude = msg.Latitude
			event.Ticket = msg.Ticket
			event.Precision = msg.Precision
			handleFunc := svr.handleMap[EventHandle].(func(*message.Event, http.ResponseWriter))
			handleFunc(event, w)
		}
	}
}

// RegisterHandle 注册消息处理器
func (svr *Server) RegisterHandle(handleType HandleType, handle interface{}) {
	// 加入HandleMap之前做个类型检查
	switch handleType {
	case RawHandle:
		handleFunc := handle.(func(message.RawMessage, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case TextHandle:
		handleFunc := handle.(func(*message.Text, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case ImageHandle:
		handleFunc := handle.(func(*message.Image, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case VoiceHandle:
		handleFunc := handle.(func(*message.Voice, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case ShortVideoHandle:
		handleFunc := handle.(func(*message.ShortVideo, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case VideoHandle:
		handleFunc := handle.(func(*message.Video, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	case EventHandle:
		handleFunc := handle.(func(*message.Event, http.ResponseWriter))
		if handleFunc != nil {
			svr.handleMap[handleType] = handle
		}
	}
}

// Start 启动服务，监听80端口
func (svr *Server) Start() {
	http.HandleFunc("/", svr.serverHandle)
	http.HandleFunc("/check", svr.ConnectServer)
	err := http.ListenAndServe(":80", nil)
	if err == nil {
		fmt.Println("server error", err)
	}
}

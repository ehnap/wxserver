package message

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"coding.net/cherrysd/wxserver/util"
)

type responseMusicMessage struct {
	PublicMessage
	Title        string `xml:"Music>Title"`
	Description  string `xml:"Music>Description"`
	MusicURL     string `xml:"Music>MusicUrl"`
	HQMusicURL   string `xml:"Music>HQMusicUrl"`
	ThumbMediaID string `xml:"Music>ThumbMediaId"`
}

// Music 回复音乐逻辑消息体(供外部使用)
type Music struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MediaID      string
	Title        string
	Description  string
}

func (rtmsg *Music) formatLogicMsg() (string, error) {
	destMsg := responseVideoMessage{}
	destMsg.MediaID = rtmsg.MediaID
	destMsg.ToUserName = rtmsg.ToUserName
	destMsg.FromUserName = rtmsg.FromUserName
	destMsg.Title = rtmsg.Title
	destMsg.Description = rtmsg.Description
	if rtmsg.CreateTime == 0 {
		destMsg.CreateTime = util.GetCurrTimeStamp()
	} else {
		destMsg.CreateTime = rtmsg.CreateTime
	}
	destMsg.MsgType = MusicMsg

	responseRawXMLMsg, err := xml.Marshal(destMsg)
	if err != nil {
		log.Println("Build Response Text Message Error")
		return "", err
	}
	result := string(responseRawXMLMsg)
	return result, err
}

// Send 向服务器发送文字消息
func (rtmsg *Music) Send(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
	w.WriteHeader(200)
	strResponseMsg, err := rtmsg.formatLogicMsg()
	fmt.Fprint(w, strResponseMsg)
	return err
}

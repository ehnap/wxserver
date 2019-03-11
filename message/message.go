package message

import (
	"encoding/xml"
)

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Type 消息类型
type Type string

// 消息类型枚举
const (
	TextMsg       = "text"
	ImageMsg      = "image"
	LocationMsg   = "location"
	LinkMsg       = "link"
	VoiceMsg      = "voice"
	VideoMsg      = "video"
	MusicMsg      = "music"
	ArticleMsg    = "news"
	ShortvideoMsg = "shortvideo"
	EventMsg      = "event"
)

// Message 微信消息体
type requestMessage struct {
	PublicMessage
	//基本消息
	MsgID        int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	MediaID      string  `xml:"MediaId"`
	PicURL       string  `xml:"PicUrl"`
	Format       string  `xml:"Format"`
	Recognition  string  `xml:"Recognition"`
	ThumbMediaID string  `xml:"ThumbMediaId"`
	Title        string  `xml:"Title"`
	Description  string  `xml:"Description"`
	URL          string  `xml:"Url"`
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64 `xml:"Scale"`
	Label        string  `xml:"Label"`
	Event        string  `xml:"Event"`
	EventKey     string  `xml:"EventKey"`
	Ticket       string  `xml:"Ticket"`
	Latitude     float64 `xml:"Latitude"`
	Longitude    float64 `xml:"Longitude"`
	Precision    float64 `xml:"Precision"`
}

// PublicMessage 公共微信消息头数据
type PublicMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      Type     `xml:"MsgType"`
}

// RawMessage 服务器传来的微信逻辑消息体
type RawMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      Type
	MsgID        int64
	Content      string
	PicURL       string
	MediaID      string
	Format       string
	Recognition  string
	ThumbMediaID string
	Title        string
	Description  string
	URL          string
	LocationX    float64
	LocationY    float64
	Scale        float64
	Label        string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float64
	Longitude    float64
	Precision    float64
}

// ParseMsg 解析服务器发来的消息
func ParseMsg(contentBytes []byte) (RawMessage, error) {
	msg := requestMessage{}
	requestMsg := RawMessage{}
	err := xml.Unmarshal(contentBytes, &msg)

	requestMsg.Content = msg.Content
	requestMsg.MsgType = msg.MsgType
	requestMsg.FromUserName = msg.FromUserName
	requestMsg.ToUserName = msg.ToUserName
	requestMsg.MsgID = msg.MsgID
	requestMsg.CreateTime = msg.CreateTime
	requestMsg.MediaID = msg.MediaID
	requestMsg.PicURL = msg.PicURL
	requestMsg.Format = msg.Format
	requestMsg.Recognition = msg.Recognition
	requestMsg.ThumbMediaID = msg.ThumbMediaID
	requestMsg.Title = msg.Title
	requestMsg.URL = msg.URL
	requestMsg.LocationX = msg.LocationX
	requestMsg.LocationY = msg.LocationY
	requestMsg.Scale = msg.Scale
	requestMsg.Label = msg.Label
	requestMsg.Event = msg.Event
	requestMsg.EventKey = msg.EventKey
	requestMsg.Ticket = msg.Ticket
	requestMsg.Latitude = msg.Latitude
	requestMsg.Longitude = msg.Longitude
	requestMsg.Precision = msg.Precision
	return requestMsg, err
}

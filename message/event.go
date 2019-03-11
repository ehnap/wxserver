package message

// Event 事件消息体
type Event struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float64
	Longitude    float64
	Precision    float64
}

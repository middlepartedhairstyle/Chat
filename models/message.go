package models

import (
	"context"
	"fmt"
	Kafka "github.com/middlepartedhairstyle/HiWe/kafka"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/segmentio/kafka-go"
	"strconv"
)

// 消息类型
const (
	MessageTypeText  uint8 = iota + 1 //文本
	MessageTypeImage                  //图片
	MessageTypeVoice                  //音频
)

const (
	MediaFriend uint8 = iota + 1 //好友消息
	MediaGroup                   //群消息
)

const (
	ChatWithFriend = "f" //订阅消息者管道，后期更改为用户具体管道
	ChatWithGroup  = "g"
	ChatGroupUser  = "gu" //群用户具体通道
)

const maxUserNum = 100
const maxGroupNum = 50

// UserChatMessage 用户消息
type UserChatMessage struct {
	FromID      uint   `json:"from_id"`      //用户ID(好友),群用户ID(群)
	ToID        uint   `json:"to_id"`        //好友ID,群ID,
	MessageType uint8  `json:"message_type"` //消息类型,如图片,文字等
	Media       uint8  `json:"media"`        //消息种类,如群消息和好友消息
	Message     string `json:"message"`      //消息主体
}

type Information interface {
	SetTopic(topic uint) string
	MessageDispose(producers map[string]*Kafka.Producer, fromID uint, message []byte, args ...map[string]uint)
	SetConsumerID() uint
	GetInformation(opts ...string) map[string]interface{}
}

// GetInformation 获取消息信息
func (userMessage *UserChatMessage) GetInformation(opts ...string) map[string]interface{} {
	var info map[string]interface{} = make(map[string]interface{})
	for _, opt := range opts {
		switch opt {
		case "fromID":
			info["fromID"] = userMessage.FromID
		case "toID":
			info["toID"] = userMessage.ToID
		case "messageType":
			info["messageType"] = userMessage.MessageType
		case "media":
			info["media"] = userMessage.Media
		case "message":
			info["message"] = userMessage.Message
		}
	}
	return info
}

func NewUserChatMessage() *UserChatMessage {
	return &UserChatMessage{}
}

// SetTopic 设置消息话题发送点或获取点
func (userMessage *UserChatMessage) SetTopic(topic uint) string {
	switch userMessage.Media {
	case MediaFriend:
		var tpId = topic/maxUserNum + 1
		return fmt.Sprintf("%s%s%d", ChatWithFriend, "tp", tpId) //例如ftp1,ftp2
		//return "test"
	case MediaGroup:
		var tpId = topic/maxGroupNum + 1
		return fmt.Sprintf("%s%s%d", ChatWithGroup, "tp", tpId) //例如gtp1,gtp2
	default:
		return ""
	}
}

// SetConsumerID 设置消费者id
func (userMessage *UserChatMessage) SetConsumerID() uint {
	switch userMessage.Media {
	//消息种类为好友消息
	case MediaFriend:
		var friend mySQL.Friends
		friend.ID = userMessage.ToID
		b := friend.FindTwoUserID()
		if b {
			if friend.UserOneID == userMessage.FromID {
				return friend.UserTwoID
			} else {
				return friend.UserOneID
			}
		} else {
			return 0
		}
	//消息种类为群消息
	case MediaGroup:
		return userMessage.FromID
	default:
		return 0
	}
}

// MessageDispose 消息处理
func (userMessage *UserChatMessage) MessageDispose(producers map[string]*Kafka.Producer, fromID uint, message []byte, args ...map[string]uint) {
	var infoVerify map[string]uint
	if len(args) > 0 {
		infoVerify = args[0]
	}
	switch userMessage.Media {
	case MediaFriend:
		userMessage.IsFriendMessage(producers, fromID, message, infoVerify)
	case MediaGroup:
		userMessage.IsGroupMessage(producers, fromID, message, infoVerify)
	default:
		panic("unhandled default case")

	}
}

// IsFriendMessage 消息类型为好友消息的处理方式
func (userMessage *UserChatMessage) IsFriendMessage(producers map[string]*Kafka.Producer, fromID uint, message []byte, infoVerify map[string]uint) {
	//消息正确性验证
	if userMessage.FromID == fromID && userMessage.JudgeFriend(infoVerify) {
		var producer *Kafka.Producer
		toFriendID := infoVerify[ChatWithFriend+strconv.Itoa(int(userMessage.ToID))]
		//查看是否有该topic,有就使用没有就新建
		if producers[userMessage.SetTopic(toFriendID)] != nil {
			producer = producers[userMessage.SetTopic(toFriendID)]
		} else {
			topic := userMessage.SetTopic(toFriendID)
			producer = Kafka.NewProducer(Kafka.SetProducerTopic(topic))
			if !producer.GetTopic(topic) {
				if !producer.CreateTopicWithRetention(topic) {
					return
				}
			}
			producers[userMessage.SetTopic(toFriendID)] = producer
		}

		//fmt.Println(toFriendID)
		key := []byte(ChatWithFriend + strconv.Itoa(int(toFriendID)))
		err := producer.WriteData(&key, &message)
		if err != nil {
			fmt.Println(err)
		}

		//自己接收自己的信息,测试部分
		if producers[userMessage.SetTopic(fromID)] != nil {
			producer = producers[userMessage.SetTopic(fromID)]
		} else {
			topic := userMessage.SetTopic(fromID)
			producer = Kafka.NewProducer(Kafka.SetProducerTopic(topic))
			if !producer.GetTopic(topic) {
				if !producer.CreateTopicWithRetention(topic) {
					return
				}
			}
			producers[userMessage.SetTopic(fromID)] = producer
		}
		key = []byte(ChatWithFriend + strconv.Itoa(int(fromID)))
		err = producer.WriteData(&key, &message)
		if err != nil {
			return
		}

		//将数据持续化存入服务器
		go userMessage.FriendMessageTypeDispose()
	} else {
		return
	}

}

// IsGroupMessage 消息类型为群消息的处理方式
func (userMessage *UserChatMessage) IsGroupMessage(producers map[string]*Kafka.Producer, fromID uint, message []byte, infoVerify map[string]uint) {
	//消息来源正确性验证
	if userMessage.FromID != fromID {
		return
	}
	var producer *Kafka.Producer
	if producers[userMessage.SetTopic(userMessage.ToID)] != nil {
		producer = producers[userMessage.SetTopic(userMessage.ToID)]
	} else {
		topic := userMessage.SetTopic(userMessage.ToID)
		producer = Kafka.NewProducer(Kafka.SetProducerTopic(topic))
		if !producer.GetTopic(topic) {
			if !producer.CreateTopicWithRetention(topic) {
				return
			}
		}
		producers[userMessage.SetTopic(userMessage.ToID)] = producer
	}

	key := []byte(ChatGroupUser + strconv.Itoa(int(userMessage.ToID)))
	err := producer.WriteData(&key, &message)
	if err != nil {
		fmt.Println(err)
	}

	go userMessage.GroupMessageTypeDispose()
}

// FriendMessageTypeDispose 好友消息处理
func (userMessage *UserChatMessage) FriendMessageTypeDispose() bool {

	switch userMessage.MessageType {
	//消息类型为文本
	case MessageTypeText:
		friendMessage := mySQL.NewFriendMessage(userMessage.FromID, userMessage.ToID, userMessage.MessageType, &(userMessage.Message))
		return friendMessage.CreateFriendMessage()
	//消息类型为图片
	case MessageTypeImage:
		return true
	//消息类型为语音
	case MessageTypeVoice:
		return true
	default:
		return false
	}
}

// GroupMessageTypeDispose 群消息处理
func (userMessage *UserChatMessage) GroupMessageTypeDispose() bool {
	switch userMessage.MessageType {
	//消息类型为文本
	case MessageTypeText:
		groupMessage := mySQL.NewGroupMessage(userMessage.FromID, userMessage.ToID, userMessage.MessageType, &(userMessage.Message))
		return groupMessage.CreateGroupMessage()
	//消息类型为图片
	case MessageTypeImage:
		return true
	//消息类型为语音
	case MessageTypeVoice:
		return true
	default:
		return false
	}
}

// JudgeFriend 判断是否为好友
func (userMessage *UserChatMessage) JudgeFriend(infoVerify map[string]uint) bool {

	if infoVerify[ChatWithFriend+strconv.Itoa(int(userMessage.ToID))] != 0 {
		return true
	} else {
		friend := mySQL.NewFriend(mySQL.SetFriendID(userMessage.ToID), mySQL.SetUserOneID(userMessage.FromID))
		toFriendID, b := friend.IsFriendUseFriendID()
		if b {
			infoVerify[ChatWithFriend+strconv.Itoa(int(userMessage.ToID))] = toFriendID
			return true
		}
		return false
	}
}

// JudgeGroupUser 判断是否是群用户成员
func (userMessage *UserChatMessage) JudgeGroupUser() {

}

// GetFriendMessage 获取好友发送的消息
func (userMessage *UserChatMessage) GetFriendMessage(id uint, ws *WebSocketClient) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovering from panic in GetMessage:", r)
		}
	}()
	var message []byte // 消息
	var tpId = id/maxUserNum + 1
	topic := fmt.Sprintf("%s%s%d", ChatWithFriend, "tp", tpId) // 例如 ftp1, ftp2
	consumer := Kafka.NewConsumer(Kafka.SetConsumerTopic(topic), Kafka.SetConsumerGroupID(ChatWithFriend+strconv.Itoa(int(id))))
	defer func(consumer *kafka.Reader) {
		err := consumer.Close()
		if err != nil {

		}
	}(consumer) // 确保消费者关闭

	for {
		m, err := consumer.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}
		if string(m.Key) == ChatWithFriend+strconv.Itoa(int(id)) {
			message = m.Value
			ws.messageList <- message
		} else {
			if err = consumer.CommitMessages(context.Background(), m); err != nil {
				fmt.Printf("提交偏移量失败: %v\n", err)
			}
		}
	}
}

// GetGroupMessage 获取群发送的消息
func (userMessage *UserChatMessage) GetGroupMessage(id uint, ws *WebSocketClient) {}

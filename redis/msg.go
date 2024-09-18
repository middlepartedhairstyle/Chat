package redis

import (
	"context"
	"fmt"
)

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, message string) error {
	err := Rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		fmt.Println("redis publish err:", err)
		return err
	} else {
		fmt.Println("redis publish success")
		return nil
	}
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Rdb.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("redis subscribe err:", err)
		return "", err
	} else {
		return msg.Payload, nil
	}
}

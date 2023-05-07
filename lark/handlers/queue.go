package handlers

import (
	"encoding/json"
	"errors"
	"lark/db"
	"time"
)

type IDiscordQueue struct {
	LarkChatId string `json:"larkChatId"`
	Time       int64  `json:"time"`
}

func discordQueueCheck(larkChatId string) error {
	for {
		isLock := db.GetCache().Get(DiscordLockKey)
		if isLock == "" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	db.GetCache().SetCustom(DiscordLockKey, "lock", time.Duration(2)*time.Second)
	defer db.GetCache().Clear(DiscordLockKey)

	// 下面是正常逻辑
	queue := db.GetCache().GetInterface(DiscordQueueKey)

	if queue != nil {
		var queueList []IDiscordQueue
		if err := json.Unmarshal(queue, &queueList); err != nil {
			return nil
		}
		queueList = discordQueueAutoOutDie(queueList)
		db.GetCache().SetInterfaceNotTimeLimit(DiscordQueueKey, queueList)

		for _, item := range queueList {
			if item.LarkChatId == larkChatId {
				return errors.New("🤖️：您存在任务正在运行中，请稍后再试～")
			}
		}

		if len(queueList) > 3 {
			return errors.New("🤖️：当前任务过多，请稍后再试～")
		}
	}
	return nil
}

func discordQueueSet(larkChatId string) {
	queue := db.GetCache().GetInterface(DiscordQueueKey)

	if queue != nil {
		var queueList []IDiscordQueue
		if err := json.Unmarshal(queue, &queueList); err != nil {
			return
		}
		queueList = discordQueueAutoOutDie(queueList)
		queueList = append(queueList, IDiscordQueue{
			LarkChatId: larkChatId,
			Time:       time.Now().Unix(),
		})
		db.GetCache().SetInterfaceNotTimeLimit(DiscordQueueKey, queueList)
	} else {
		db.GetCache().SetInterface(DiscordQueueKey, []IDiscordQueue{{
			LarkChatId: larkChatId,
			Time:       time.Now().Unix(),
		}})
	}
}

func discordQueueDel(larkChatId string) {
	queue := db.GetCache().GetInterface(DiscordQueueKey)

	if queue != nil {
		var queueList []IDiscordQueue
		if err := json.Unmarshal(queue, &queueList); err != nil {
			return
		}
		queueList = discordQueueAutoOutDie(queueList)
		newQueueList := make([]IDiscordQueue, 0)
		for _, item := range queueList {
			if item.LarkChatId != larkChatId {
				newQueueList = append(newQueueList, item)
			}
		}
		db.GetCache().SetInterfaceNotTimeLimit(DiscordQueueKey, newQueueList)
	}
}

func discordQueueAutoOutDie(queueList []IDiscordQueue) []IDiscordQueue {
	currentTime := time.Now().Unix()
	newQueueList := make([]IDiscordQueue, 0)
	for _, item := range queueList {
		if item.Time+30*60 > currentTime {
			newQueueList = append(newQueueList, item)
		}
	}

	return newQueueList
}

package cron_msg_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"message-nest/models"
	"message-nest/pkg/constant"
	"message-nest/service/cron_service"
	"message-nest/service/send_message_service"
)

type MsgCronTask struct {
}

func (s MsgCronTask) Register() {
	// 获取所有的定时消息任务
	limit := 10000
	filter := make(map[string]interface{})
	filter["enable"] = 1
	data, err := models.GetCronMessages(0, limit, "", filter)
	if err != nil {
		logrus.Errorf("获取定时消息任务失败！原因：%s", err.Error())
		return
	}
	if len(data) == 0 {
		logrus.Infof("没有定时消息任务需要注册")
		return
	}
	//注册定时任务
	for _, msg := range data {
		AddCronMsgToCronServer(msg)
	}
	length := len(data)
	if length > 0 {
		logrus.Infof("完成用户自定义的定时消息注册，注册个数：%d", length)
	}
}

// AddCronMsgToCronServer 注册定时任务到定时服务
func AddCronMsgToCronServer(msg models.CronMessages) {
	if msg.Enable != 1 {
		return
	}
	taskId := cron_service.AddTask(cron_service.ScheduledTask{
		Schedule: msg.Cron,
		Job: func() {
			CronMsgSendF(msg)
		},
	})
	constant.CronMsgIdMapMemoryCache[msg.ID] = taskId
	logrus.Infof("新增定时消息成功，消息id: %s，消息名: %s，当前任务总数：%d", msg.ID, msg.Name, len(constant.CronMsgIdMapMemoryCache))
}

// 执行任务的构造函数
func CronMsgSendF(msg models.CronMessages) {
	logrus.Infof("开始只能执行定时消息发送任务: %s，消息名: %s", msg.ID, msg.Name)
	task, err := models.GetTaskByID(msg.TaskID)
	if err != nil {
		logrus.Infof("消息任务不存在: %s ", msg.TaskID)
		return
	}
	sender := send_message_service.SendMessageService{
		TaskID:   task.ID,
		Title:    msg.Title,
		Text:     msg.Content,
		URL:      msg.Url,
		CallerIp: fmt.Sprintf("[CrondTask] [%s] ID: %s", task.Name, task.ID),
		DefaultLogger: logrus.WithFields(logrus.Fields{
			"prefix": "[Cron Message]",
		}),
	}
	taskData, _ := sender.SendPreCheck()
	_, err = sender.Send(taskData)
	if err != nil {
		logrus.Error("执行定时消息失败：%s", err.Error())
		return
	}
}

// UpdateCronMsgToCronServer 更新定时服务的任务
func UpdateCronMsgToCronServer(msg models.CronMessages) {
	if entryId, ok := constant.CronMsgIdMapMemoryCache[msg.ID]; ok {
		// 先删除之前的定时任务
		delete(constant.CronMsgIdMapMemoryCache, msg.ID)
		cron_service.RemoveTask(entryId)
		// 再注册新的定时任务
		AddCronMsgToCronServer(msg)
	} else {
		// 注册新的定时任务
		AddCronMsgToCronServer(msg)
	}
	logrus.Infof("完成定时消息的定时更新，消息id: %s，当前任务总数：%d", msg.ID, len(constant.CronMsgIdMapMemoryCache))
}

// RemoveCronMsgToCronServer 删除定时任务中心的任务
func RemoveCronMsgToCronServer(msg models.CronMessages) {
	if entryId, ok := constant.CronMsgIdMapMemoryCache[msg.ID]; ok {
		// 先删除之前的定时任务
		delete(constant.CronMsgIdMapMemoryCache, msg.ID)
		cron_service.RemoveTask(entryId)
	}
	logrus.Infof("删除定时消息完成，消息id: %s，剩余任务总数：%d", msg.ID, len(constant.CronMsgIdMapMemoryCache))
}

// StartUpMsgCronTask 启动注册定时任务
func StartUpMsgCronTask() {
	MsgCronTask{}.Register()
}

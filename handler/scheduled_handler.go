package handler

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	loader "sf-bot/handler/load"
	"sf-bot/handler/service"
)

// StartGroupMorningPaperTimer 2023/6/6 lamkeizyi - 工作日9点半运行
func StartGroupMorningPaperTimer() {
	if loader.LoadBool("group.morning_paper") {
		fmt.Println("initialize 9:30 timer")
		c := cron.New()
		c.AddFunc("1 30 9 ? * 2,3,4,5,6", func() {
			if groupHandler.morningPaperMode == TEXT {
				if dayTextService, err := service.GetPicDayTextService(); err == nil {
					log.Println("cuz has error that send pic . convert to text to sending... ")
					groupHandler.sendText(dayTextService.ToString())
				}
			} else {
				service.StartPicDayService(func(name string) {
					if err := groupHandler.sendPic(name); err != nil {
						//convert to text
						if dayTextService, err := service.GetPicDayTextService(); err == nil {
							log.Println("cuz has error that send pic . convert to text to sending... ")
							groupHandler.sendText(dayTextService.ToString())
						}

					}
				})
			}
		})
		c.Start()
		select {}
	}
}

// StartGroupMoyuTimer 2023/6/6 lamkeizyi - 工作日10点
func StartGroupMoyuTimer() {
	if loader.LoadBool("group.moyu") {
		fmt.Println("initialize 10:00 timer")
		c := cron.New()
		c.AddFunc("1 0 10 ? * 2,3,4,5,6", func() {
			service.StartMoyuPicDayService(func(name string) {
				groupHandler.sendPic(name)
			})
		})
		c.Start()
		select {}
	}
}

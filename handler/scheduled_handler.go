package handler

import "github.com/eatmoreapple/openwechat"

//todo 2023/6/4 lamkeizyi -

type ScheduledHandler struct {
	everyDay    bool
	groupNames  []string
	friendNames []string
	timeFormat  string
	source      map[string]string
}

func CreateScheduledHandler(weixin *openwechat.User) {

}

func (s *ScheduledHandler) createTask(startTime int64, everyDay bool) {

}

func (s *ScheduledHandler) startTask() {

}

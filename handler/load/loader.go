package loader

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

// var groupNames map[string]string
var file *ini.File

func init() {
	load, err := ini.Load("kkbot.ini")
	if err != nil {
		log.Println("can't find kkbot.ini")
		panic(err)
	}
	file = load
	log.Println("load kkbot.ini ... ")
}

func Exist(str string) bool {
	name := GroupName()
	split := strings.Split(name, ",")
	for _, key := range split {
		if strings.Contains(str, key) {
			return true
		}
	}
	return false
}

func GroupName() string {
	key := file.Section("group").Key("group_name").String()
	return key
}

func Load(str string) string {
	split := strings.Split(str, ".")
	key := file.Section(split[0]).Key(split[1]).String()
	return key
}

func LoadBool(str string) bool {
	//index := strings.Index(str, ".")
	split := strings.Split(str, ".")
	key, _ := file.Section(split[0]).Key(split[1]).Bool()
	return key
}

func LoadInt(str string) int {
	split := strings.Split(str, ".")
	key := file.Section(split[0]).Key(split[1]).String()
	atoi, _ := strconv.Atoi(key)
	return atoi
}

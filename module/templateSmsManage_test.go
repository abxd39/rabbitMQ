package module_test

import (
	"log"
	"sctek.com/typhoon/th-platform-gateway/service"
	"testing"
	"time"
)

func TestTemplateSmsManage_AboutIdInfo(t *testing.T) {
	log.Println("测试manageMq能否可以正常压入数据")
	temp:=new(service.LogicService)
	for i:=0;i<100;i++{
		time.Sleep(1*time.Second)
		log.Printf("第%d次推送消息到manageMq",i)
		temp.AboutIdInfo(12)
	}
}


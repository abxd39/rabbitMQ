package module

import (
	"log"
	"testing"
	"time"
)

func TestTemplateSmsManage_AboutIdInfo(t *testing.T) {
	log.Println("测试manageMq能否可以正常压入数据")
	temp:=new(TemplateSmsManage)
	temp.AboutIdInfo(12)
	for i:=0;i<100;i++{
		time.Sleep(1*time.Second)
		log.Printf("第%d次推送消息到manageMq",i)
		new(TemplateSmsManage).AboutIdInfo(12)
	}
}


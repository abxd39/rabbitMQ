package module

import (
	"fmt"
	Log "github.com/sirupsen/logrus"
	"sctek.com/typhoon/th-platform-gateway/common"
	"time"
)

type TemplateSmsManage struct {
	Id             int       `xorm:"not null pk autoincr INT(10)"`
	CorpId         int       `xorm:"default 0 INT(10)"`
	MallId         int       `xorm:"default 0 INT(10)"`
	UserId         int       `xorm:"default 0 comment('操作人id') INT(10)"`
	Username       string    `xorm:"default '' comment('用户名') VARCHAR(20)"`
	Mobile         string    `xorm:"default '' comment('指定手机号') VARCHAR(20)"`
	TemplateId     int       `xorm:"default 0 comment('template_sms表的主键id') INT(10)"`
	AcceptUserType int       `xorm:"default 1 comment('1-指定会员；2-全部会员;') TINYINT(1)"`
	SendType       int       `xorm:"default 1 comment('1-定时发送;2-即时发送') TINYINT(1)"`
	SendTime       time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('发送的时间') DATETIME"`
	SendStatus     int       `xorm:"default 1 comment('1-待发送；2-已发送；3-已取消') TINYINT(1)"`
	SendCount      int       `xorm:"default 0 comment('发送数量') INT(10)"`
	Created        time.Time `xorm:"not null DATETIME"`
	Updated        time.Time `xorm:"not null DATETIME"`
	Delete         int       `xorm:"default 0 comment('0-正常；1-删除') TINYINT(1)"`
}

//获取定时任务
func(t*TemplateSmsManage)GetManageCron()([]TemplateSmsManage,error){
	list:=make([]TemplateSmsManage,0)
	err:=common.DB.Where("send_type=1").Where("send_status=1").Find(&list)
	if err!=nil{
		return nil,err
	}
	return list,nil
}

func (t *TemplateSmsManage) GetManageOfId(id int) (error) {
	Log.Infoln("从mq 中获取消息id")
	engine := common.DB
	has, err := engine.Where("id=?", id).Get(t)
	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("给定的模板id=%d的条目不存在！！！", id)
	}
	return nil

}


func (t*TemplateSmsManage) UpdateCount(id,count int)error{
	Log.Infof("修改数据库状态id=%v,发送的数量为=%v\r\n",id,count)
	engine:=common.DB
	if count ==0{
		return t.UpdateSendStatus(id)
	}
	has,err:=engine.Where("id=?",id).Get(t)
	if err!=nil{
		return err
	}
	if !has{
		return fmt.Errorf("数据库状态更新失败")
	}
	_,err=engine.Cols("send_count","updated","send_status").Where("id=?",id).Update(&TemplateSmsManage{
		SendCount:count+t.SendCount,
		SendStatus:2,
		Updated:time.Now(),
	})
	if err!=nil{
		return err
	}
	return nil
}

func (t*TemplateSmsManage)UpdateSendStatus(id int)error{
	Log.Infoln("修改数据库发送状态")
	has,err:=common.DB.Cols("send_status").Where("id=?",id).Update(&TemplateSmsManage{
		SendStatus:1,
	})
	if err!=nil{
		return err
	}
	if has<=0{
		return fmt.Errorf("数据库状态更新失败")
	}
	return nil
}
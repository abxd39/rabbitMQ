package module

import (
	"fmt"
	"log"
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

func (t *TemplateSmsManage) AboutIdInfo(id int) error {
	engine := common.DB
	has, err := engine.Where("id=?", id).Where("send_status=1").Get(t)
	if err != nil {
		common.Log.Errorln(err)
		return err
	}
	//判断怎么发送 发送那些人
	if !has {
		err = fmt.Errorf("给定的模板id=%d的条目不存在！！！", id)
		common.Log.Errorf(err.Error())
		return err

	}
	if t.AcceptUserType == 1 { //指定会员
		if t.SendType == 1 { //即时发
			log.Println("指定会员——即时发送")
			new(TemplateSmsType).SearchOfManageId(t.Id, t.TemplateId)
		} else if t.SendType == 2 { //定时发
			//启动定时器
			str := "2006-01-02 15:04:05"
			str = t.SendTime.Format(str)
			log.Println("定时发送短息时间是" + str)
		ConditionUser:
			for {
				time.Sleep(5 * time.Second)
				tim := time.Now()
				//再次查询数据库是否已经被取消掉了
				has, err = engine.Where("id=?", id).Where("send_status=1").Get(t)
				if err != nil {
					return err
				}
				if !has {
					//发送已经取消
					return nil
				}
				//判断是否到发送时间
				if t.SendTime.Unix() <= tim.Unix() {
					log.Println("指定会员——定时发送")
					new(TemplateSmsType).SearchOfManageId(t.Id, t.TemplateId)
					break ConditionUser
				}
			}

		}

	} else if t.AcceptUserType == 2 { //全部会员
		//获取短息模板
		message, err := new(TemplateSms).GetText(t.TemplateId)
		if err != nil {
			common.Log.Infoln(err)
			return err
		}
		if t.SendType == 1 { //即时发
			log.Println("全员会员——即时发送")
			new(MemberInfo).SendMessageEveryOne(message)
		} else if t.SendType == 2 { //定时发
			//启动定时器
			log.Println("全会员——定时发送")
		EveryMark:
			for {
				time.Sleep(5 * time.Second)
				tim := time.Now()
				//再次查询数据库是否已经被取消掉了
				has, err = engine.Where("id=?", id).Where("send_status=1").Get(t)
				if err != nil {
					return err
				}
				if !has {
					//发送已经取消
					return nil
				}
				//判断是否到发送时间
				if t.SendTime.Unix() <= tim.Unix() {
					log.Println("指定会员——定时发送")
					new(MemberInfo).SendMessageEveryOne(message)
					break EveryMark
				}
			}
		}

	}else if t.AcceptUserType ==3{//指定手机号发送
		message, err := new(TemplateSms).GetText(t.TemplateId)
		if err != nil {
			common.Log.Infoln(err)
			return err
		}
		if t.SendType == 1 { //即时发
			log.Println("全员会员——即时发送")
			new(MemberInfo).SendMessageOfPhone(t.Mobile,message)
		} else if t.SendType == 2 { //定时发
			PhoneMar:
				for{
					time.Sleep(5 * time.Second)
					tim := time.Now()
					//再次查询数据库是否已经被取消掉了
					has, err = engine.Where("id=?", id).Where("send_status=1").Get(t)
					if err != nil {
						return err
					}
					if !has {
						//发送已经取消
						return nil
					}
					//判断是否到发送时间
					if t.SendTime.Unix() <= tim.Unix() {
						log.Println("指定手机——定时发送")
						new(MemberInfo).SendMessageOfPhone(t.Mobile,message)
						break PhoneMar
					}
				}

		}
	}
	//t.SendStatus = 2
	//count, err := engine.Cols("send_status").Where("id=?", id).Update(t)
	//if err != nil {
	//	common.Log.Errorf("短息发送的状态更新数据库失败！！！%q", err)
	//	return err
	//}
	//if count == 0 {
	//	common.Log.Infoln("短息发送的状态更新数据库失败！！！")
	//}
	return nil
}

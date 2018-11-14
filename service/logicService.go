package service

import (
	"fmt"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/module"
	"strings"
	"time"
)

type LogicService struct {
}

func (l *LogicService) SendMessage(id, templateId int) {
	list, err := new(module.TemplateSmsType).GetSmsTypeOfId(id)
	if err != nil {
		common.Log.Errorln(err)
		return
	}
	message, err := new(module.TemplateSms).GetText(templateId)
	if err != nil {
		return
	}
	for _, v := range list {
		fmt.Println("指定发送给谁", v.Type)
		if v.Type == 1 { //姓别
			err := l.SendMessageForSex(id, v.TypeData, message)
			if err != nil {
				common.Log.Errorln(err)
				continue
			}
		} else if v.Type == 2 { //会员等级
			err:=l.SendMessageForGrade(id, v.TypeData, message)
			if err!=nil{
				common.Log.Errorln(err)
				continue
			}

		} else if v.Type == 3 { //会员生日
			err:=l.SendMessageOfBirthDay(id, v.TypeData, message)
			if err!=nil{
				common.Log.Errorln(err)
				continue
			}
		}
	}
	return
}
//根据会员生日
func (l*LogicService)SendMessageOfBirthDay(id int,typeDate,msg string)error  {
	list,err:=new(module.MemberInfo).GetAllMember()
	if err!=nil{
		return err
	}
	subList := strings.Split(typeDate, ",")
	for _, v := range list {
		if len(v.Birthday) != len("1991-06-14") {
			continue
		}
		month := v.Birthday[5:7]
		sendLog := new(MarshalJson)
		sendLog.TemplateManageId = id
		for _, m := range subList {
			if strings.Compare(month, m) == 0 || strings.Compare(month, "0"+m) == 0 {
				sendLog.MemberId = v.MemberId
				sendLog.CorpId = v.CorpId
				sendLog.Mobile = v.Mobile
				sendLog.MallId = new(module.Member).GetMallId(v.MemberId)
				if sendLog.MallId ==0{
					continue
				}
				result,err:=sendLog.marshalJson(msg)
				if err !=nil{
					common.Log.Errorln(err)
					continue
				}
				Push("myPusher","rmq_test",result)
				continue
			}
		}
	}
	return nil
}
//根据会员等级
func (l *LogicService) SendMessageForGrade(id int, typeDate, msg string) error {
	list, err := new(module.MemberCard).GetMessageOfGrade(typeDate)
	if err != nil {
		return err
	}
	sendLog := new(MarshalJson)
	sendLog.TemplateManageId = id
	for _, v := range list {
		sendLog.MemberId = v.MemberId
		sendLog.CorpId = v.CorpId
		sendLog.Mobile = v.Mobile
		sendLog.MallId = new(module.Member).GetMallId(v.MemberId)
		if sendLog.MallId == 0 {
			continue
		}

		result, err := sendLog.marshalJson(msg)
		if err != nil {
			common.Log.Errorln(err)
			continue
		}
		Push("myPusher", "rmq_test", result)
	}
	return nil
}

//根据性别
func (l *LogicService) SendMessageForSex(id int, typeDae, msg string) error {
	list, err := new(module.MemberInfo).GetMessageOfSex(typeDae)
	if err != nil {
		return err
	}
	//发送短信
	sendLog := new(MarshalJson)
	sendLog.TemplateManageId = id
	for _, v := range list {
		sendLog.MemberId = v.MemberId
		sendLog.CorpId = v.CorpId
		sendLog.Mobile = v.Mobile
		sendLog.MallId = new(module.Member).GetMallId(v.MemberId)
		if sendLog.MallId == 0 {
			continue
		}

		result, err := sendLog.marshalJson(msg)
		if err != nil {
			common.Log.Errorln(err)
			continue
		}
		Push("myPusher", "rmq_test", result)
	}
	return nil
}


//发送所有人
func(l* LogicService)SendMessageEveryOne(id int,msg string)error{
	list,err:=new(module.MemberInfo).GetAllMember()
	if err!=nil{
		return err
	}
	sendLog := new(MarshalJson)
	sendLog.TemplateManageId = id
	for _, v := range list {
		sendLog.MemberId = v.MemberId
		sendLog.CorpId = v.CorpId
		sendLog.Mobile = v.Mobile
		sendLog.MallId = new(module.Member).GetMallId(v.MemberId)
		if sendLog.MallId == 0 {
			continue
		}

		result, err := sendLog.marshalJson(msg)
		if err != nil {
			common.Log.Errorln(err)
			continue
		}
		Push("myPusher", "rmq_test", result)
	}

	return nil
}


//指定手机号码发送
func(l*LogicService)SendMessageOfMobile(id int, typeDate,msg string)error{
	ob:=new(module.MemberInfo)
	err:=ob.GetMessageOfPhone(typeDate)
	if err!=nil{
		return err
	}
	sendLog:=new(MarshalJson)
	sendLog.TemplateManageId =id
	sendLog.MemberId = ob.MemberId
	sendLog.CorpId = ob.CorpId
	sendLog.Mobile = ob.Mobile
	sendLog.MallId = new(module.Member).GetMallId(ob.MemberId)
	if sendLog.MallId ==0{
		return fmt.Errorf("会员%d不存在！！",ob.MemberId)
	}
	result,err:=sendLog.marshalJson(msg)
	if err!=nil{
		common.Log.Errorln(err)
		return err
	}
	Push("myPusher","rmq_test",result)
	return nil
}


func (l *LogicService) AboutIdInfo(id int) {
	//再此处理逻辑业务
	//判断怎么发送 发送那些人
	tsm:=new(module.TemplateSmsManage)
	err := tsm.GetManageOfId(id)
	if err != nil {
		common.Log.Errorln(err)
		return
	}
	if tsm.AcceptUserType == 1 { //指定会员
		if tsm.SendType == 2 { //即时发
			fmt.Println("指定会员——即时发送")
			//new(TemplateSmsType).SearchOfManageId(t.Id, t.TemplateId)
			l.SendMessage(tsm.Id, tsm.TemplateId)

		} else if tsm.SendType == 1 { //定时发
			//启动定时器
			str := "2006-01-02 15:04:05"
			str = tsm.SendTime.Format(str)
			fmt.Println("定时发送短息时间是" + str)
		ConditionUser:
			for {
				time.Sleep(5 * time.Second)
				tim := time.Now()
				err:=tsm.GetManageOfId(id)
				if err!=nil{
					common.Log.Errorln(err)
					break ConditionUser
				}
				//判断是否到发送时间
				if tsm.SendTime.Unix() <= tim.Unix() {

					fmt.Printf("开始发送时间为%s\r\n",time.Now().Format(str))
					l.SendMessage(tsm.Id, tsm.TemplateId)
					break ConditionUser
				}
			}

		}

	} else if tsm.AcceptUserType == 2 { //全部会员
		 fmt.Println("全员发送")
		//获取短息模板
		message, err := new(module.TemplateSms).GetText(tsm.TemplateId)
		if err != nil {
			common.Log.Infoln(err)
			return
		}
		if tsm.SendType == 2 { //即时发
			fmt.Println("全员即时发送")
			l.SendMessageEveryOne(tsm.Id, message)
		} else if tsm.SendType == 1 { //定时发
			//启动定时器
			fmt.Println("全会员——定时发送")
			str := "2006-01-02 15:04:05"
			common.Log.Infoln("全会员——定时发送")
		EveryMark:
			for {
				time.Sleep(5 * time.Second)
				tim := time.Now()
				err:=tsm.GetManageOfId(id)
				if err!=nil{
					common.Log.Errorln(err)
					break EveryMark
				}
				//判断是否到发送时间
				if tsm.SendTime.Unix() <= tim.Unix() {
					fmt.Printf("全员定时发送时间为%q\r\n",time.Now().Format(str))
					l.SendMessageEveryOne(tsm.Id, message)
					break EveryMark
				}
			}
		}

	} else if tsm.AcceptUserType == 3 { //指定手机号发送
		 fmt.Println("指定电话号码发送")
		message, err := new(module.TemplateSms).GetText(tsm.TemplateId)
		if err != nil {
			common.Log.Infoln(err)
			return
		}
		if tsm.SendType == 1 { //即时发
			fmt.Println("指定电话号码即时发送")
			err:=l.SendMessageOfMobile(tsm.Id, tsm.Mobile, message)
			if err!=nil{
				common.Log.Errorln(err)
				return
			}
		} else if tsm.SendType == 2 { //定时发
		PhoneMar:
			for {
				time.Sleep(5 * time.Second)
				tim := time.Now()
				//再次查询数据库是否已经被取消掉了
				err:=tsm.GetManageOfId(id)
				if err!=nil{
					common.Log.Errorln(err)
					break PhoneMar
				}
				//判断是否到发送时间
				if tsm.SendTime.Unix() <= tim.Unix() {
					fmt.Printf("指定电话号码定时发送 发送时间%q\r\n",time.Now().Format("2006-01-02 15:04:05"))
					l.SendMessageOfMobile(tsm.Id, tsm.Mobile, message)
					break PhoneMar
				}
			}

		}
	}

}

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

func (l *LogicService) SendMessage(id, templateId int)int {
	count:=0
	list, err := new(module.TemplateSmsType).GetSmsTypeOfId(id)
	if err != nil {
		common.Log.Errorln(err)
		return count
	}
	message, err := new(module.TemplateSms).GetText(templateId)
	if err != nil {
		return count
	}
	for _, v := range list {
		if v.Type == 1 { //姓别
			fmt.Println("根据会员性别发送")
			count,err = l.SendMessageForSex(id, v.TypeData, message)
			if err != nil {
				common.Log.Errorln(err)
				continue
			}
		} else if v.Type == 2 { //会员等级
			fmt.Println("根据会员等级发送")
			count=l.SendMessageForGrade(id, v.TypeData, message)
		} else if v.Type == 3 { //会员生日
			fmt.Println("根据会员生日发送")
			count,err=l.SendMessageOfBirthDay(id, v.TypeData, message)
			if err!=nil{
				common.Log.Errorln(err)
				continue
			}
		}
	}
	return count
}
//根据会员生日
func (l*LogicService)SendMessageOfBirthDay(id int,typeDate,msg string) (int,error)  {
	count:=0
	list,err:=new(module.MemberInfo).GetAllMember()
	if err!=nil{
		return count, err
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

				err=Push("Pusher","",result)
				if err!=nil{
					common.Log.Errorln(err)
					continue
				}
				count++
				continue
			}
		}
	}
	return count,nil
}
//根据会员等级
func (l *LogicService) SendMessageForGrade(id int, typeDate, msg string) (int) {
	count:=0
	listMemberId,err:= new(module.Member).GetMemberId(typeDate)
	if err!=nil{
		common.Log.Errorln(err)
		return count
	}
	if len(listMemberId)==0{
		err= fmt.Errorf("member 表中没有 MemberCarId=%v的记录",typeDate)
		common.Log.Errorln(err)
		return count
	}
	list, err := new(module.MemberInfo).GetMemberIdList(listMemberId)
	if err != nil {
		common.Log.Errorln(err)
		return count
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

		Push("Pusher", "", result)
		if err!=nil{
			common.Log.Errorln(err)
			continue
		}
		count++
	}
	return count
}

//根据性别
func (l *LogicService) SendMessageForSex(id int, typeDae, msg string) (int,error) {
	count:=0
	list, err := new(module.MemberInfo).GetMessageOfSex(typeDae)
	if err != nil {
		return count,err
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

		err=Push("Pusher", "", result)
		if err!=nil{
			common.Log.Errorln(err)
			continue
		}
		count++
	}
	return count,nil
}


//发送所有人
func(l* LogicService)SendMessageEveryOne(id int,msg string) (int,error){
	count :=0
	list,err:=new(module.MemberInfo).GetAllMember()
	if err!=nil{
		return count,err
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

		err=Push("Pusher", "", result)
		if err!=nil{
			common.Log.Errorln(err)
			continue
		}
		count++
	}

	return count,nil
}


//指定手机号码发送
//返回值 1 为发送条数
func(l*LogicService)SendMessageOfMobile(id int, typeDate,msg string)(int,error){
	ob:=new(module.MemberInfo)
	err:=ob.GetMessageOfPhone(typeDate)
	if err!=nil{
		return 0,err
	}
	sendLog:=new(MarshalJson)
	sendLog.TemplateManageId =id
	sendLog.MemberId = ob.MemberId
	sendLog.CorpId = ob.CorpId
	sendLog.Mobile = ob.Mobile
	sendLog.MallId = new(module.Member).GetMallId(ob.MemberId)
	if sendLog.MallId ==0{
		return 0,fmt.Errorf("会员%d不存在！！",ob.MemberId)
	}
	result,err:=sendLog.marshalJson(msg)
	if err!=nil{
		common.Log.Errorln(err)
		return 0,err
	}
	err=Push("Pusher","",result)
	if err!=nil{
		common.Log.Errorln(err)
		return 0,err
	}
	return 1,nil
}


//入口函数
func (l *LogicService) AboutIdInfo(id int) {
	//再此处理逻辑业务
	//判断怎么发送 发送那些人
	count:=0
	tsm:=new(module.TemplateSmsManage)
	err := tsm.GetManageOfId(id)
	if err != nil {
		common.Log.Errorln(err)
		if err:=tsm.UpdateSendStatus(tsm.Id);err!=nil{
			common.Log.Errorln(err)
		}
		return
	}
	if tsm.AcceptUserType == 1 { //指定会员
		if tsm.SendType == 2 { //即时发
			fmt.Println("指定会员——即时发送")
			count=l.SendMessage(tsm.Id, tsm.TemplateId)
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

					fmt.Printf("开始发送时间为%v\r\n",time.Now().Format(str))
					count=l.SendMessage(tsm.Id, tsm.TemplateId)
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
			if err:=tsm.UpdateSendStatus(tsm.Id);err!=nil{
				common.Log.Errorln(err)
			}
			return
		}
		if tsm.SendType == 2 { //即时发
			fmt.Println("全员即时发送")
			count,err=l.SendMessageEveryOne(tsm.Id, message)
			if err!=nil{
				common.Log.Errorln(err)
				if err:=tsm.UpdateSendStatus(tsm.Id);err!=nil{
					common.Log.Errorln(err)
				}
				return
			}
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
					count,err=l.SendMessageEveryOne(tsm.Id, message)
					if err!=nil{
						common.Log.Errorln(err)
						break EveryMark
					}
					break EveryMark
				}
			}
		}

	} else if tsm.AcceptUserType == 3 { //指定手机号发送
		 fmt.Println("指定电话号码发送")
		message, err := new(module.TemplateSms).GetText(tsm.TemplateId)
		if err != nil {
			common.Log.Infoln(err)
			if err:=tsm.UpdateSendStatus(tsm.Id);err!=nil{
				common.Log.Errorln(err)
			}
			return
		}
		if tsm.SendType == 2 { //即时发
			fmt.Println("指定电话号码即时发送")
			count,err=l.SendMessageOfMobile(tsm.Id, tsm.Mobile, message)
			if err!=nil{
				common.Log.Errorln(err)
				if err:=tsm.UpdateSendStatus(tsm.Id);err!=nil{
					common.Log.Errorln(err)
				}
				return
			}
		} else if tsm.SendType == 1 { //定时发
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
					count,err=l.SendMessageOfMobile(tsm.Id, tsm.Mobile, message)
					if err!=nil{
						common.Log.Errorln(err)
						break PhoneMar
					}
					break PhoneMar
				}
			}

		}
	}
	//更新状态
	err=tsm.UpdateCount(tsm.Id,count)
	if err!=nil{
		common.Log.Errorln(err)
	}
}

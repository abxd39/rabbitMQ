package module

import "time"

type TemplateSmsLog struct {
	Id               int       `xorm:"not null pk autoincr INT(10)"`
	MemberId         int       `xorm:"default 0 INT(10)"`
	TemplateManageId int       `xorm:"default 0 comment('template_sms_manage表的主键id') INT(10)"`
	Mobile           string    `xorm:"not null default '' comment(' 手机号') VARCHAR(20)"`
	Status           int       `xorm:"default 0 comment('发送状态：1-成功；2-失败') TINYINT(1)"`
	Created          time.Time `xorm:"comment('备注时间') DATETIME"`
}


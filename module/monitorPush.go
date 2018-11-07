package module

type MonitorPush struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	CorpId    int    `xorm:"not null comment('机构id') INT(11)"`
	MallId    int    `xorm:"not null comment('分店id') INT(11)"`
	Emails    string `xorm:"comment('邮箱,多个邮件用分号隔开') VARCHAR(200)"`
	Wechatids string `xorm:"comment(''微信用户id，多个用分号隔开'') VARCHAR(200)"`
}

func (a MonitorPush) TableName() string {
	return "monitor_push"
}

package module

type MemberCardActivateSetting struct {
	Id            int    `xorm:"not null pk autoincr INT(11)"`
	CardSettingId int    `xorm:"not null default 0 INT(11)"`
	Type          int    `xorm:"not null default 0 comment('选项类型（1：文本输入，2：日期选择，3：下拉选择：4：地区选择）') TINYINT(4)"`
	Name          string `xorm:"not null default '' comment('字段名') VARCHAR(255)"`
	Title         string `xorm:"not null default '' comment('标题') VARCHAR(255)"`
	Desc          string `xorm:"not null default '' comment('提示信息') VARCHAR(255)"`
	IsRequired    int    `xorm:"not null default 0 comment('是否必填（0：否，1：是）') TINYINT(1)"`
	IsDisabled    int    `xorm:"not null default 0 comment('是否允许修改（0：允许，1：不允许）') TINYINT(1)"`
}

func (a MemberCardActivateSetting) TableName() string {
	return "member_card_activate_setting"
}

package module

import (
	Log "github.com/sirupsen/logrus"
	"sctek.com/typhoon/th-platform-gateway/common"
)

type MemberInfo struct {
	MemberId int    `xorm:"not null pk comment('会员ID') INT(11)"`
	CorpId   int    `xorm:"not null default 0 INT(11)"`
	CardNo   string `xorm:"not null default '' comment('会员卡号') index VARCHAR(50)"`
	Name     string `xorm:"not null default '' comment('会员姓名') VARCHAR(255)"`
	Mobile   string `xorm:"not null default '' comment(' 手机号') index VARCHAR(20)"`
	Sex      int    `xorm:"not null default 0 comment('性别（0：未知，1：男，2：女）') TINYINT(1)"`
	Birthday string `xorm:"not null comment(' 生日') DATE"`
	ProvId   int    `xorm:"not null default 0 INT(11)"`
	CityId   int    `xorm:"not null default 0 INT(11)"`
	AreaId   int    `xorm:"not null default 0 INT(11)"`
	Address  string `xorm:"not null comment('详细地址') TEXT"`
	Email    string `xorm:"not null default '' comment('邮箱') VARCHAR(50)"`
}

func (m *MemberInfo) TableName() string {
	return "member_info"
}

func (m *MemberInfo) GetMemberIdList(idList []int) ([]MemberInfo, error) {
	Log.Infof("会员id%v\r\n", idList)
	list := make([]MemberInfo, 0)
	err := common.DB.Cols("member_id", "corp_id", "mobile").In("member_id", idList).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *MemberInfo) GetMessageOfSex(sex string) ([]MemberInfo, error) {
	Log.Infof("性别为%q的会员", sex)
	engine := common.DB
	list := make([]MemberInfo, 0)
	err := engine.Where("sex=?", sex).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *MemberInfo) GetAllMember() ([]MemberInfo, error) {
	Log.Infoln("即时全员发送")
	engine := common.DB
	list := make([]MemberInfo, 0)
	err := engine.Select(" * ").Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *MemberInfo) GetMessageOfPhone(phone string) (bool, error) {
	Log.Infoln("指定电话号码发送短息")
	engine := common.DB
	has, err := engine.Where("mobile=?", phone).Get(m)
	if err != nil {
		return true, err
	}
	return has, nil
}

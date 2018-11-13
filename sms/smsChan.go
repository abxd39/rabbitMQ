package sms



var ChanPushMsgs chan *SmsMsgStruct // = make(chan *SmsMsgStruct, 200)

type SmsWorkStruct struct {
}

type SmsMsgStruct struct {
	QueueName string
	Msg       string
}

func InitWorker(maxWork int, maxQueueSize int) {
	ChanPushMsgs = make(chan *SmsMsgStruct, maxQueueSize)
	//初始工作
	for i := 0; i < maxWork; i++ {
		smsWork := &SmsWorkStruct{}
		go smsWork.DoWork()
	}
}

func (p *SmsWorkStruct) DoWork() {
	for {
		select {
		case pMsg := <-ChanPushMsgs:
			if pMsg != nil {
				Send(pMsg.QueueName, pMsg.Msg)
			} else {
				return
			}
		}
	}
}

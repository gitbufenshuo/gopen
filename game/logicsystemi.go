package game

import "fmt"

type LogicSystemI interface {
	GetLogicByName(gi *GlobalInfo, name string) LogicSupportI
	BindLogicByName(gi *GlobalInfo, logicName string, f func(gi *GlobalInfo) LogicSupportI)
}
type LogicBind struct {
	logicNewFunc map[string]func(gi *GlobalInfo) LogicSupportI
}

func NewLogicBind() LogicSystemI {
	logicBindIns := new(LogicBind)
	logicBindIns.logicNewFunc = make(map[string]func(gi *GlobalInfo) LogicSupportI)
	//
	return logicBindIns
}

func (lb *LogicBind) BindLogicByName(gi *GlobalInfo, logicName string, f func(gi *GlobalInfo) LogicSupportI) {
	lb.logicNewFunc[fmt.Sprintf("logic_%s", logicName)] = f
}

func (lb *LogicBind) GetLogicByName(gi *GlobalInfo, logicName string) LogicSupportI {
	return lb.logicNewFunc[logicName](gi)
}

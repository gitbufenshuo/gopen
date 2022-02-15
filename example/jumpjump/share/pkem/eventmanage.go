package pkem

/*
	本地逻辑事件管理器
	不用提交到远程的
*/

type EventManager struct {
	evlist []*Event
}

func NewEventManager() *EventManager {
	res := new(EventManager)
	return res
}

// 生成一个事件
func (em *EventManager) FireEvent(ev *Event) {
	em.evlist = append(em.evlist, ev)
}

func (em *EventManager) GetEvList() []*Event {
	return em.evlist
}
func (em *EventManager) Clear() {
	em.evlist = nil
}

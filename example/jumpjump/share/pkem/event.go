package pkem

type EventKind int64

const (
	EK_DoAtt = iota + 1 // 攻击事件
)

type Event struct {
	PID              int64     // playerid
	TargetPID        int64     // 目标
	EK               EventKind // 事件类型
	SubKind          int64     // 子类型
	PosX, PosY, PosZ int64     // 事件发生的坐标
	DirX, DirY, DirZ int64     // 事件发生的朝向
}

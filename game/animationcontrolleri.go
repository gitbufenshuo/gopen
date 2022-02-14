package game

type AnimationControllerI interface {
	BindBoneNode(name string, transform *Transform)
	ChangeMode(mode string)
	GetModeList() []string
	NowMode() string
	RecordInitFrame()
	Update()
}

type AniMoving struct {
	GBID     int    // 每一个gameobject 只能绑定一个 AC
	BoneName string // 节点名称
}

type AnimationSystemI interface {
	CreateAC(amname string) AnimationControllerI
}

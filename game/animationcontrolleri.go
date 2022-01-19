package game

type AnimationControllerI interface {
	BindBoneNode(name string, transform *Transform)
	ChangeMode(mode string)
	GetModeList() []string
	NowMode() string
	RecordInitFrame()
}

type AniMoving struct {
	GBID     int    // 每一个gameobject 只能绑定一个 AC
	BoneName string // 节点名称
}

type AnimationSystemI interface {
	CreateAC(amname string, gbid int) AnimationControllerI
	GetAC(gbid int) AnimationControllerI
	GetMoving(gbid int) []*AniMoving
	GameobjectDel(gbid int)
	GameobjectDetach(gbid int)
	// gbid 是主节点gameobject id
	BindBoneNode(gbid int, bonename string, transform *Transform)
	CloneAC(oldgbid, newgbid int) AnimationControllerI
	Update()
}

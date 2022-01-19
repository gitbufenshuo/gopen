package game

type AnimationControllerI interface {
	BindBoneNode(name string, transform *Transform)
	ChangeMode(mode string)
	GetModeList() []string
	NowMode() string
	RecordInitFrame()
}
type AnimationSystemI interface {
	CreateAC(amname string, gbid int) AnimationControllerI
	GetAC(gbid int) AnimationControllerI
	GameobjectDel(gbid int)
	GameobjectDetach(gbid int)
	Update()
}

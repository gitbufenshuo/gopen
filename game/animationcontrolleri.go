package game

type AnimationControllerI interface {
	BindBoneNode(name string, transform *Transform)
	ChangeMode(mode string)
	GetModeList() []string
	NowMode() string
	RecordInitFrame()
}
type AnimationSystemI interface {
	CreateAnimationController(amname string, gbid int) AnimationControllerI
	GetAnimationController(gbid int) AnimationControllerI
	GameobjectDel(gbid int)
	Update()
}

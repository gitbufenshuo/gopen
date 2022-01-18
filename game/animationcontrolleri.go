package game

type AnimationControllerI interface {
	BindBoneNode(name string, transform *Transform)
	RecordInitFrame()
}
type AnimationSystemI interface {
	CreateAnimationController(amname string, gbid int) AnimationControllerI
	GameobjectDel(gbid int)
	Update()
}

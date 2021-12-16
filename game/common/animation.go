package common

import "github.com/gitbufenshuo/gopen/matmath"

type BoneSatus struct {
	Position *matmath.VECX
	Rotation *matmath.VECX
}

func NewBoneSatus(px, py, pz, rx, ry, rz float32) *BoneSatus {
	var pos matmath.VECX
	pos.Init3()
	pos.SetValue3(px, py, pz)
	var rot matmath.VECX
	rot.Init3()
	rot.SetValue3(rx, ry, rz)
	return &BoneSatus{
		Position: &pos,
		Rotation: &rot,
	}
}

type AnimationFrame struct {
	HeadStatus      *BoneSatus
	BodyStatus      *BoneSatus
	HandLeftStatus  *BoneSatus
	HandRightStatus *BoneSatus
	WheelStatus     *BoneSatus
}

type AnimationController struct {
	InitFrame *AnimationFrame
	AniMode   map[string][]*AnimationFrame
	ModeList  []string
	CurMode   string
	CurIndex  int
	CurDir    int
	////////////////////////////////////////////////
	headNode      *Transform
	bodyNode      *Transform
	handLeftNode  *Transform
	handRightNode *Transform
	wheelNode     *Transform
}

func NewAnimationController() *AnimationController {
	res := new(AnimationController)
	res.AniMode = make(map[string][]*AnimationFrame)
	res.CurDir = 1
	return res
}

func (ac *AnimationController) AddMode(mode string, frameList []*AnimationFrame) {
	ac.AniMode[mode] = frameList
}

func (ac *AnimationController) ChangeMode(mode string) {
	ac.CurMode = mode
	ac.CurIndex = 0
}

// generate the init frame
func (ac *AnimationController) BindBoneNode(head, body, handLeft, handRight, wheel *Transform) {
	ac.headNode = head
	ac.bodyNode = body
	ac.handLeftNode = handLeft
	ac.handRightNode = handRight
	ac.wheelNode = wheel
	ac.RecordInitFrame()
}
func (ac *AnimationController) RecordInitFrame() {
	initFrame := new(AnimationFrame)
	{
		position := ac.headNode.Postion.Clone()
		rotation := ac.headNode.Rotation.Clone()
		initFrame.HeadStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.bodyNode.Postion.Clone()
		rotation := ac.bodyNode.Rotation.Clone()
		initFrame.BodyStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.handLeftNode.Postion.Clone()
		rotation := ac.handLeftNode.Rotation.Clone()
		initFrame.HandLeftStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.handRightNode.Postion.Clone()
		rotation := ac.handRightNode.Rotation.Clone()
		initFrame.HandRightStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.wheelNode.Postion.Clone()
		rotation := ac.wheelNode.Rotation.Clone()
		initFrame.WheelStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	ac.InitFrame = initFrame
}

func (ac *AnimationController) Update() {
	initFrame := ac.InitFrame
	list := ac.AniMode[ac.CurMode]
	if len(list) == 0 {
		return
	}
	curFrame := list[ac.CurIndex]
	if ac.headNode != nil {
		ac.headNode.Postion.Add2_InPlace(initFrame.HeadStatus.Position, curFrame.HeadStatus.Position)
		ac.headNode.Rotation.Add2_InPlace(initFrame.HeadStatus.Rotation, curFrame.HeadStatus.Rotation)
	}
	if ac.bodyNode != nil {
		ac.bodyNode.Postion.Add2_InPlace(initFrame.BodyStatus.Position, curFrame.BodyStatus.Position)
		ac.bodyNode.Rotation.Add2_InPlace(initFrame.BodyStatus.Rotation, curFrame.BodyStatus.Rotation)
	}
	if ac.handLeftNode != nil {
		ac.handLeftNode.Postion.Add2_InPlace(initFrame.HandLeftStatus.Position, curFrame.HandLeftStatus.Position)
		ac.handLeftNode.Rotation.Add2_InPlace(initFrame.HandLeftStatus.Rotation, curFrame.HandLeftStatus.Rotation)
	}
	//
	if ac.handRightNode != nil {
		ac.handRightNode.Postion.Add2_InPlace(initFrame.HandRightStatus.Position, curFrame.HandRightStatus.Position)
		ac.handRightNode.Rotation.Add2_InPlace(initFrame.HandRightStatus.Rotation, curFrame.HandRightStatus.Rotation)
	}
	if ac.wheelNode != nil {
		ac.wheelNode.Postion.Add2_InPlace(initFrame.WheelStatus.Position, curFrame.WheelStatus.Position)
		ac.wheelNode.Rotation.Add2_InPlace(initFrame.WheelStatus.Rotation, curFrame.WheelStatus.Rotation)
	}
	//
	ac.CurIndex++
	ac.CurIndex %= len(list)
}

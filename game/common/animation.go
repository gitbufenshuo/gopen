package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
)

type BoneSatus struct {
	Position *matmath.Vec4
	Rotation *matmath.Vec4
}

func (bs *BoneSatus) ToByte() []byte {
	return []byte(fmt.Sprintf("%f|%f|%f|%f|%f|%f", bs.Position.GetIndexValue(0),
		bs.Position.GetIndexValue(1),
		bs.Position.GetIndexValue(2),
		bs.Rotation.GetIndexValue(0),
		bs.Rotation.GetIndexValue(1),
		bs.Rotation.GetIndexValue(2),
	))
}

// 0|1|2|3|4|5
func NewBoneStatusFromData(data string) *BoneSatus {
	if data == "NONE" {
		data = "0.000000|0.000000|0.000000|0.000000|0.000000|0.000000"
	}
	segs := strings.Split(data, "|")
	px, py, pz, rx, ry, rz := help.Str2Float32(segs[0]), help.Str2Float32(segs[1]), help.Str2Float32(segs[2]), help.Str2Float32(segs[3]), help.Str2Float32(segs[4]), help.Str2Float32(segs[5])
	return NewBoneSatus(px, py, pz, rx, ry, rz)
}

func NewBoneSatus(px, py, pz, rx, ry, rz float32) *BoneSatus {
	var pos matmath.Vec4
	pos.SetValue3(px, py, pz)
	var rot matmath.Vec4
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
	LegLeftStatus   *BoneSatus
	LegRightStatus  *BoneSatus
}

func (af *AnimationFrame) ToByte() []byte {
	var list = [][]byte{
		af.HeadStatus.ToByte(),
		af.BodyStatus.ToByte(),
		af.HandLeftStatus.ToByte(),
		af.HandRightStatus.ToByte(),
		af.LegLeftStatus.ToByte(),
		af.LegRightStatus.ToByte(),
	}
	return bytes.Join(list, []byte(","))
}

// 0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5
func NewAnimationFrameFromData(data string) *AnimationFrame {
	af := new(AnimationFrame)
	//
	segs := strings.Split(data, ",")
	if len(segs) != 6 {
		fmt.Println(segs)
		os.Exit(1)
	}
	af.HeadStatus = NewBoneStatusFromData(segs[0])
	af.BodyStatus = NewBoneStatusFromData(segs[1])
	af.HandLeftStatus = NewBoneStatusFromData(segs[2])
	af.HandRightStatus = NewBoneStatusFromData(segs[3])
	af.LegLeftStatus = NewBoneStatusFromData(segs[4])
	af.LegRightStatus = NewBoneStatusFromData(segs[5])
	return af
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
	legLeftNode   *Transform
	legRightNode  *Transform
}

func NewAnimationController() *AnimationController {
	res := new(AnimationController)
	res.AniMode = make(map[string][]*AnimationFrame)
	res.CurDir = 1
	return res
}

/*
STOP
head,body,left,right,wheel
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,
--------------------------
JUMP
head,body,left,right,wheel
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,
--------------------------
FIRE
head,body,left,right,wheel
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,
,       ,    ,    ,     ,


*/
func (ac *AnimationController) LoadFromFile(filename string) {
	data, _ := ioutil.ReadFile(filename)
	ac.LoadFromData(data)
}

func (ac *AnimationController) LoadFromData(data []byte) {
	split := []byte("\n--------------------------\n")
	modeList := bytes.Split(data, split)
	////////////////////////////////////
	for _, onemodeData := range modeList {
		buffer := bytes.NewBuffer(onemodeData)
		scanner := bufio.NewScanner(buffer)
		mode := ""
		linenum := 0
		aflist := []*AnimationFrame{}
		for scanner.Scan() { // 扫描一行
			content := scanner.Text()
			////////////////////////
			if linenum == 0 {
				mode = content
				linenum++
				continue
			}
			////////////////////////
			if linenum == 1 {
				linenum++
				continue
			}
			if len(content) < 10 {
				linenum++
				continue
			}
			////////////////////////
			linenum++
			af := NewAnimationFrameFromData(content)
			aflist = append(aflist, af)
		}
		if !strings.HasPrefix(mode, "//") {
			ac.ModeList = append(ac.ModeList, mode)
			ac.AniMode[mode] = aflist
		}
	}
	ac.CurMode = ac.ModeList[0]
}

func (ac *AnimationController) Write2Data() []byte {
	split := []byte("\n--------------------------\n")
	////////
	modeDataList := [][]byte{}
	for _, modename := range ac.ModeList {
		modebuffer := bytes.NewBuffer(nil)
		modebuffer.Reset()
		modebuffer.WriteString(modename)
		modebuffer.WriteString("\n")
		modebuffer.WriteString("head,body,left,right,wheel\n")
		afList := ac.AniMode[modename]
		dataList := [][]byte{}
		for _, oneaf := range afList {
			dataList = append(dataList, oneaf.ToByte())
		}
		modebuffer.Write(bytes.Join(dataList, []byte("\n")))
		modeDataList = append(modeDataList, modebuffer.Bytes())
	}
	return bytes.Join(modeDataList, split)
}

func (ac *AnimationController) AddMode(mode string, frameList []*AnimationFrame) {
	ac.AniMode[mode] = frameList
}

func (ac *AnimationController) ChangeMode(mode string) {
	ac.CurMode = mode
	ac.CurIndex = 0
}

// generate the init frame
func (ac *AnimationController) BindBoneNode(head, body, handLeft, handRight, legLeft, legRight *Transform) {
	ac.headNode = head
	ac.bodyNode = body
	ac.handLeftNode = handLeft
	ac.handRightNode = handRight
	ac.legLeftNode = legLeft
	ac.legRightNode = legRight
	ac.RecordInitFrame()
}
func (ac *AnimationController) RecordInitFrame() {
	initFrame := new(AnimationFrame)
	{
		position := ac.headNode.Postion
		rotation := ac.headNode.Rotation
		initFrame.HeadStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.bodyNode.Postion
		rotation := ac.bodyNode.Rotation
		initFrame.BodyStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.handLeftNode.Postion
		rotation := ac.handLeftNode.Rotation
		initFrame.HandLeftStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.handRightNode.Postion
		rotation := ac.handRightNode.Rotation
		initFrame.HandRightStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.legLeftNode.Postion
		rotation := ac.legLeftNode.Rotation
		initFrame.LegLeftStatus = &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		}
	}
	{
		position := ac.legRightNode.Postion
		rotation := ac.legRightNode.Rotation
		initFrame.LegRightStatus = &BoneSatus{
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
	if ac.legLeftNode != nil {
		ac.legLeftNode.Postion.Add2_InPlace(initFrame.LegLeftStatus.Position, curFrame.LegLeftStatus.Position)
		ac.legLeftNode.Rotation.Add2_InPlace(initFrame.LegLeftStatus.Rotation, curFrame.LegLeftStatus.Rotation)
	}
	if ac.legRightNode != nil {
		ac.legRightNode.Postion.Add2_InPlace(initFrame.LegRightStatus.Position, curFrame.LegRightStatus.Position)
		ac.legRightNode.Rotation.Add2_InPlace(initFrame.LegRightStatus.Rotation, curFrame.LegRightStatus.Rotation)
	}
	//
	ac.CurIndex++
	ac.CurIndex %= len(list)
}

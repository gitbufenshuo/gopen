package animationsystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gitbufenshuo/gopen/game"
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
	StatusList []*BoneSatus
}

// 0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5,0|1|2|3|4|5
func NewAnimationFrameFromData(data string, boneNum int) *AnimationFrame {
	af := new(AnimationFrame)
	//
	segs := strings.Split(data, ",")
	if len(segs) != boneNum {
		fmt.Println(segs)
		os.Exit(1)
	}
	for idx := range segs {
		af.StatusList = append(af.StatusList, NewBoneStatusFromData(segs[idx]))
	}
	return af
}

type AnimationMeta struct {
	AniMode  map[string][]*AnimationFrame
	ModeList []string
	IndexMap map[string]int
}

func LoadAnimationMetaFromData(data []byte) *AnimationMeta {
	am := new(AnimationMeta)
	am.AniMode = make(map[string][]*AnimationFrame)
	am.IndexMap = make(map[string]int)
	split := []byte("\n--------------------------\n")
	modeList := bytes.Split(data, split)
	////////////////////////////////////
	for _, onemodeData := range modeList {
		buffer := bytes.NewBuffer(onemodeData)
		scanner := bufio.NewScanner(buffer)
		descline := "" // __init 6
		mode := ""
		boneNum := 6
		linenum := 0
		aflist := []*AnimationFrame{}
		for scanner.Scan() { // 扫描一行
			content := scanner.Text()
			////////////////////////
			if linenum == 0 {
				descline = content
				if strings.Contains(descline, " ") {
					descsegs := strings.Split(descline, " ")
					mode = descsegs[0]
					boneNum = help.Str2Int(descsegs[1])
				} else {
					mode = descline
					boneNum = 6
				}
				linenum++
				continue
			}
			////////////////////////
			if linenum == 1 { // head,body,handleft,handright,legleft,legright
				if len(am.IndexMap) == 0 {
					namesegs := strings.Split(content, ",")
					for idx := range namesegs {
						am.IndexMap[namesegs[idx]] = idx
					}
				}
				linenum++
				continue
			}
			if len(content) < 10 {
				linenum++
				continue
			}
			////////////////////////
			linenum++
			af := NewAnimationFrameFromData(content, boneNum)
			aflist = append(aflist, af)
		}
		if !strings.HasPrefix(descline, "//") {
			am.ModeList = append(am.ModeList, mode)
			am.AniMode[mode] = aflist
		}
	}
	return am
}

func LoadAnimationMetaFromFile(filename string) *AnimationMeta {
	data, _ := ioutil.ReadFile(filename)
	return LoadAnimationMetaFromData(data)
}

type AnimationControlSpec struct {
	Name      string
	Index     int
	transform *game.Transform
}

type AnimationController struct {
	InitFrame *AnimationFrame
	AM        *AnimationMeta
	CurMode   string
	CurIndex  int
	////////////////////////////////////////////////
	NodeList []*AnimationControlSpec
}

func NewAnimationController() *AnimationController {
	res := new(AnimationController)
	return res
}
func (ac *AnimationController) UseAimationMeta(am *AnimationMeta) {
	ac.AM = am
	ac.CurMode = am.ModeList[0]
	ac.CurIndex = 0
}
func (ac *AnimationController) ChangeMode(mode string) {
	ac.CurMode = mode
	ac.CurIndex = 0
}

func (ac *AnimationController) BindBoneNode(name string, transform *game.Transform) {
	newAnimationControlSpec := new(AnimationControlSpec)
	newAnimationControlSpec.Name = name
	newAnimationControlSpec.Index = ac.AM.IndexMap[name]
	newAnimationControlSpec.transform = transform

	ac.NodeList = append(ac.NodeList, newAnimationControlSpec)
	// ac.RecordInitFrame()
}

func (ac *AnimationController) RecordInitFrame() {
	initFrame := new(AnimationFrame)
	for idx := range ac.NodeList {
		position := ac.NodeList[idx].transform.Postion
		rotation := ac.NodeList[idx].transform.Rotation
		initFrame.StatusList = append(initFrame.StatusList, &BoneSatus{
			Position: &position,
			Rotation: &rotation,
		})
	}
	ac.InitFrame = initFrame
}

func (ac *AnimationController) Update() {
	initFrame := ac.InitFrame
	list := ac.AM.AniMode[ac.CurMode]
	if len(list) == 0 {
		return
	}
	curFrame := list[ac.CurIndex]
	for idx := range ac.NodeList {
		if ac.NodeList[idx] != nil {
			ac.NodeList[idx].transform.Postion.Add2_InPlace(
				initFrame.StatusList[idx].Position,
				curFrame.StatusList[ac.NodeList[idx].Index].Position,
			)
			ac.NodeList[idx].transform.Rotation.Add2_InPlace(
				initFrame.StatusList[idx].Rotation,
				curFrame.StatusList[ac.NodeList[idx].Index].Rotation,
			)
		}
	}
	//
	ac.CurIndex++
	ac.CurIndex %= len(list)
}

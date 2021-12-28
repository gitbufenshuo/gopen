package inputsystem

import (
	"fmt"
	"strconv"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const doubleClickInterval float64 = 500

type InputSystemManager struct {
	ID      int
	gi      *game.GlobalInfo
	keyList map[int]*InputListenerQueue
}

type InputListenerQueue struct {
	key        int
	stateNow   glfw.Action
	stateTime  float64
	frameState []*InputAction // 记录十帧状态，最多支持 五击-检测
	nowFrame   int            // 现在是哪一帧
	//
	justPress    bool    // 是否刚刚按下 // 只会持续一帧
	justRelease  bool    // 是否放开 // 只会持续一帧
	justDBClick  bool    // 是否双击 // 只会持续一帧
	justTriClick bool    // 是否三击 // 只会持续一帧 // todolist
	holdValue    float64 // 持续按下时，积攒蓄力

	action   *InputAction
	eeAction *InputAction

	list   []*InputListener
	eeList []*InputListener
}

func NewInputListenerQueue(key int) *InputListenerQueue {
	newilq := new(InputListenerQueue)
	//
	newilq.key = key
	newilq.stateNow = 0
	newilq.action = new(InputAction)
	newilq.eeAction = new(InputAction)
	newilq.list = make([]*InputListener, 0, 5)
	newilq.eeList = make([]*InputListener, 0, 5)
	//
	newilq.frameState = make([]*InputAction, 10)
	for idx := range newilq.frameState {
		newilq.frameState[idx] = new(InputAction)
	}
	newilq.nowFrame = 100
	return newilq
}

func (ilq *InputListenerQueue) CurFrame() int {
	return ilq.nowFrame % 10
}

func (ilq *InputListenerQueue) CalcjustPress(ism *InputSystemManager) {
	ilq.justPress = false // 先清空
	//
	if ilq.stateNow == glfw.Release {
		return
	}
	prevFrameState := ilq.frameState[(ilq.nowFrame-1)%10]
	if prevFrameState.state != glfw.Release {
		return
	}
	ilq.justPress = true
}

func (ilq *InputListenerQueue) CalcjustRelease(ism *InputSystemManager) {
	ilq.justRelease = false // 先清空
	//
	if ilq.stateNow == glfw.Press {
		return
	}
	prevFrameState := ilq.frameState[(ilq.nowFrame-1)%10]
	if prevFrameState.state != glfw.Press {
		return
	}
	ilq.justRelease = true
}

func (ilq *InputListenerQueue) CalcjustDBClick(ism *InputSystemManager) {
	ilq.justDBClick = false // 先清空
	if ilq.stateNow != glfw.Press {
		return
	}
	// 前一帧必须是 release
	prevFrameState := ilq.frameState[(ilq.nowFrame-1)%10]
	if prevFrameState.state != glfw.Release {
		return
	}
	// 往前搜，第一个 press 状态, 跟现在的时间差不超过 doubleClickInterval 毫秒
	var targetIdx int = -1
	for idx := 2; idx <= 9; idx++ {
		frameState := ilq.frameState[(ilq.nowFrame-idx)%10]
		if frameState.state == glfw.Release {
			continue
		}
		//
		targetIdx = idx
		break
	}
	if targetIdx == -1 {
		return
	}
	//
	frameState := ilq.frameState[(ilq.nowFrame-targetIdx)%10]
	if ilq.stateTime-frameState.stateTime < doubleClickInterval {
		ilq.justDBClick = true
	}
}
func (ilq *InputListenerQueue) CalcholdValue(ism *InputSystemManager) {
	ilq.holdValue = 0
	if ilq.stateNow == glfw.Press {
		return
	}
	// 前面连续9帧都是Press
	for idx := 1; idx <= 9; idx++ {
		frameState := ilq.frameState[(ilq.nowFrame-idx)%10]
		if frameState.state == glfw.Release {
			ilq.holdValue = 0
			return
		}
		ilq.holdValue = ilq.stateTime - frameState.stateTime
	}
}

func (ilq *InputListenerQueue) CheckListener(ism *InputSystemManager) {
	keyNow := ism.gi.Window().GetKey(glfw.Key(ilq.key))
	var stateSwitched bool = keyNow != ilq.stateNow
	{

		ilq.stateNow = keyNow
		ilq.stateTime = ism.gi.NowMS
		// base info
		frameState := ilq.frameState[ilq.CurFrame()]
		frameState.stateTime = ism.gi.NowMS
		frameState.state = keyNow
		//
		prevFrameState := ilq.frameState[(ilq.nowFrame-1)%10]
		frameState.value = frameState.stateTime - prevFrameState.stateTime // delta value
		if frameState.value > 10000 {
			frameState.value = 1
		}
	}
	{
		// justPress check
		ilq.CalcjustPress(ism)
		// justRelease check
		ilq.CalcjustRelease(ism)
		// double click check
		ilq.CalcjustDBClick(ism)
		// hold release check
		ilq.CalcholdValue(ism)
	}
	ilq.nowFrame++
	if stateSwitched {
		// 按键状态切换
		if ilq.stateNow == glfw.Press {
			ilq.action.phase = Pressing
			ilq.stateTime = ism.gi.NowMS
		} else {
			ilq.action.phase = Releasing
			ilq.stateTime = ism.gi.NowMS
		}
		for i := 0; i < len(ilq.list); i++ {

			ilq.action.keyType = ilq.list[i].keyType
			ilq.action.value = 1
			ilq.action.stateTime = ilq.stateTime
			ilq.list[i].callback(ilq.action)
		}
	} else {
		// 持续中
		if len(ilq.eeList) > 0 && keyNow == glfw.Press {
			for i := 0; i < len(ilq.eeList); i++ {
				ilq.eeAction.phase = Holding
				ilq.eeAction.keyType = ilq.eeList[i].keyType
				ilq.eeAction.value = ism.gi.NowMS - ilq.stateTime
				ilq.eeList[i].callback(ilq.eeAction)
			}
		}
	}
}

type InputListener struct {
	keyType  KeyType
	callback func(action *InputAction)
}

type KeyType int

const (
	KeyStatus KeyType = 1
	KeyValue  KeyType = 2
)

type phaseType int

const (
	Disabled  phaseType = 0
	Pressing  phaseType = 1
	Releasing phaseType = 2
	Holding   phaseType = 3
)

type InputAction struct {
	keyType   KeyType
	phase     phaseType
	value     float64     // delta time
	stateTime float64     // timestamp
	state     glfw.Action // press or release
}

func (action *InputAction) GetKeyType() KeyType {
	return action.keyType
}

func (action *InputAction) GetPhase() phaseType {
	return action.phase
}

func (action *InputAction) GetValue() float64 {
	return action.value
}

var ins *InputSystemManager

func (ism *InputSystemManager) Start() {
	fmt.Println("InputSystem Start")
}

func (ism *InputSystemManager) Update() {
	for _, ilq := range ism.keyList {
		ilq.CheckListener(ism)
	}
}

func (ism *InputSystemManager) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return ism.ID
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	ism.ID = _id[0]
	return ism.ID
}

func (ism *InputSystemManager) KeyDown(key int) bool {
	return ism.keyList[key].justPress
}
func (ism *InputSystemManager) KeyUp(key int) bool {
	return ism.keyList[key].justRelease
}
func (ism *InputSystemManager) KeyDoubleClick(key int) bool {
	return ism.keyList[key].justDBClick
}
func (ism *InputSystemManager) KeyHoldRelease(key int) float64 {
	return ism.keyList[key].holdValue
}

func GetInputSystem() *InputSystemManager {
	if ins == nil {
		fmt.Println("InputSystem New ")
		ins = new(InputSystemManager)
		ins.keyList = make(map[int]*InputListenerQueue)

	}
	return ins
}

func InitInputSystem(gi *game.GlobalInfo) {
	GetInputSystem().gi = gi
	gi.AddManageObject(ins)
	ins.TestId()
}

func (ism *InputSystemManager) TestId() {
	fmt.Println("InputSystem Test " + strconv.Itoa(ism.ID_sg()))
}

func (ism *InputSystemManager) AddKeyListener(keyType KeyType, key int, callback func(action *InputAction)) {
	var ilq *InputListenerQueue
	if _ilq, found := ism.keyList[key]; !found {
		_ilq := NewInputListenerQueue(key)
		ism.keyList[key] = _ilq
		ilq = _ilq
	} else {
		ilq = _ilq
	}
	inputListener := new(InputListener)
	inputListener.keyType = keyType
	inputListener.callback = callback
	if keyType == KeyStatus {
		ilq.list = append(ilq.list, inputListener)
	} else {
		ilq.eeList = append(ilq.eeList, inputListener)
	}
}

package inputsystem

import (
	"fmt"
	"strconv"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type InputSystemManager struct {
	ID      int
	gi      *game.GlobalInfo
	keyList map[int]*InputListenerQueue
}

type InputListenerQueue struct {
	uniqueIndex int
	stateNow    glfw.Action
	stateTime   float64
	action      *InputAction
	eeAction    *InputAction

	list   []*InputListener
	eeList []*InputListener
}

func (ilq *InputListenerQueue) CheckListener(key int, ism *InputSystemManager) {
	keyNow := ism.gi.Window().GetKey(glfw.Key(key))

	if keyNow != ilq.stateNow {
		// 按键状态切换
		ilq.stateNow = keyNow

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
	value     float64
	stateTime float64
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

	for k, ilq := range ism.keyList {
		ilq.CheckListener(k, ism)
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
	if ist == nil {
		fmt.Println("InputSystemTools New ")
		ist = new(InputSystemTools)
		ist.DoubleClickList = make(map[string]*DoubleClickInfo)
	}
}

func (ism *InputSystemManager) TestId() {
	fmt.Println("InputSystem Test " + strconv.Itoa(ism.ID_sg()))
}

func (ism *InputSystemManager) AddKeyListener(keyType KeyType, key int, callback func(action *InputAction)) {
	_, ok := ism.keyList[key]
	if !ok {
		ism.keyList[key] = new(InputListenerQueue)
		ism.keyList[key].stateNow = 0
		ism.keyList[key].action = new(InputAction)
		ism.keyList[key].eeAction = new(InputAction)
		ism.keyList[key].list = make([]*InputListener, 0, 5)
		ism.keyList[key].eeList = make([]*InputListener, 0, 5)
	}
	inputListener := new(InputListener)
	inputListener.keyType = keyType
	inputListener.callback = callback
	if keyType == KeyStatus {
		ism.keyList[key].list = append(ism.keyList[key].list, inputListener)
	} else {
		ism.keyList[key].eeList = append(ism.keyList[key].eeList, inputListener)
	}

}

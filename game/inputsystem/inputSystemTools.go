package inputsystem

import (
	"fmt"

	"github.com/google/uuid"
)

type InputSystemTools struct {
	DoubleClickList map[string]*DoubleClickInfo
}

type DoubleClickInfo struct {
	Id            string
	LastClickTime float64
	Callback      func(action *InputAction)
}

const (
	doubleClickInterval float64 = 500
)

var ist *InputSystemTools

func (ism *InputSystemManager) KeyDown(key int, callback func(action *InputAction)) {
	fmt.Println("InputSystemTools KeyDown")
	onKeyCallback := func(action *InputAction) {
		if action.GetPhase() == Pressing {
			callback(action)
		}
	}
	GetInputSystem().AddKeyListener(KeyStatus, key, onKeyCallback)
}

func (ism *InputSystemManager) KeyUp(key int, callback func(action *InputAction)) {
	fmt.Println("InputSystemTools KeyUp")
	onKeyCallback := func(action *InputAction) {
		if action.GetPhase() == Releasing {
			callback(action)
		}
	}
	GetInputSystem().AddKeyListener(KeyStatus, key, onKeyCallback)
}

func (ism *InputSystemManager) KeyHold(key int, callback func(action *InputAction)) {
	fmt.Println("InputSystemTools KeyHold")
	onKeyCallback := func(action *InputAction) {
		if action.GetPhase() == Releasing {
			callback(action)
		}
	}
	GetInputSystem().AddKeyListener(KeyValue, key, onKeyCallback)
}

func (ism *InputSystemManager) DoubleClick(key int, callback func(action *InputAction)) {
	fmt.Println("InputSystemTools DoubleClick")
	info := new(DoubleClickInfo)
	info.Id = uuid.NewString()
	info.LastClickTime = 0
	info.Callback = func(action *InputAction) {
		if action.GetPhase() == Releasing {
			if ist.DoubleClickList[info.Id].LastClickTime == 0 {
				ist.DoubleClickList[info.Id].LastClickTime = action.stateTime
			} else if action.stateTime-ist.DoubleClickList[info.Id].LastClickTime < doubleClickInterval {
				callback(action)
				ist.DoubleClickList[info.Id].LastClickTime = 0
			} else {
				ist.DoubleClickList[info.Id].LastClickTime = action.stateTime
			}
		}
	}
	ist.DoubleClickList[info.Id] = info

	GetInputSystem().AddKeyListener(KeyStatus, key, info.Callback)
}

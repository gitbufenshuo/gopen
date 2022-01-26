package manage_main

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/sceneloader"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ManageMain struct {
	gi *game.GlobalInfo
	*gameobjects.NilManageObject
	UserMap map[string]int
	UID     string // 自己的uid
	//
	which          int64 // which player is this client
	MainPlayer     game.GameObjectI
	MainPlayerJump *logic_jump.LogicJump
	SubPlayer      game.GameObjectI
	SubPlayerJump  *logic_jump.LogicJump
	//
	InMsgChan    chan commmsg.JumpMSGTurn
	OutMsgChan   chan commmsg.JumpMSGTurn
	wsadDone     bool
	turnMsgLocal commmsg.JumpMSGTurn
	frame        int
	Turn         int64 // 回合
}

func NewManageMain(gi *game.GlobalInfo) *ManageMain {
	res := new(ManageMain)
	//
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	res.UserMap = make(map[string]int)
	res.which = -1
	res.InMsgChan = make(chan commmsg.JumpMSGTurn, 100)
	res.OutMsgChan = make(chan commmsg.JumpMSGTurn, 100)
	go res.fakeServer()
	return res
}

func (lm *ManageMain) fakeServer() {
	for msglist := range lm.OutMsgChan {
		lm.InMsgChan <- msglist
	}
}

func (lm *ManageMain) Start() {
	inputsystem.InitInputSystem(lm.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeySpace))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyA))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyD))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyW))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyP))
	lm.gi.SetInputSystem(inputsystem.GetInputSystem())
	//
	lm.MainPlayer = sceneloader.FindGameobjectByName("scenespec", "MainPlayer")
	logiclist := lm.MainPlayer.GetLogicSupport()
	for idx := range logiclist {
		if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
			lm.MainPlayerJump = v
			lm.MainPlayerJump.Chosen = true
		}
	}
}

func (lm *ManageMain) clonePlayer() {
	if lm.SubPlayerJump != nil {
		return
	}
	lm.SubPlayer = lm.gi.InstantiateGameObject(lm.MainPlayer)
	logiclist := lm.SubPlayer.GetLogicSupport()
	for idx := range logiclist {
		if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
			lm.SubPlayerJump = v
		}
	}
}

func (lm *ManageMain) Update() {
	lm.frame++
	lm.clonePlayer()
	////////////////
	// 收集本地指令
	lm.Local_Login_Collect()
	lm.Local_WSAD_Collect()
	////////////////
	if lm.frame%3 == 0 {
		//
		lm.turnMsgLocal.Turn = lm.Turn
		if len(lm.turnMsgLocal.List) == 0 {
			// idle
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, commmsg.JumpMSGOne{
				Kind: "idle",
				UID:  lm.UID,
			})
		}
		// 将本地指令先发出去
		lm.OutMsgChan <- lm.turnMsgLocal
		lm.turnMsgLocal.List = nil
		// 接收服务器指令
		inmsglist := <-lm.InMsgChan
		fmt.Println("inmsglist", inmsglist.Turn, lm.Turn)
		for _, onemsg := range inmsglist.List {
			lm.MSG_Update(onemsg)
		}
		// 对本地程序步进
		lm.MainPlayerJump.OnForce()
		lm.SubPlayerJump.OnForce()
		// 回合收尾
		lm.wsadDone = false
		lm.Turn++ // 回合加1
	}
}

func (lm *ManageMain) Local_WSAD_Collect() {
	if lm.wsadDone {
		return
	}
	apressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyA))
	dpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyD))
	wpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyW))
	spressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyS))
	var mx, mz int64
	if apressed {
		mx = -50

	} else if dpressed {
		mx = 50

	}
	if wpressed {
		mz = -50

	} else if spressed {
		mz = 50

	}
	lm.wsadDone = true
	if lm.UID != "" {
		//fmt.Printf("{Collect}, (%s)(%d %d)\n", lm.UID, mx, mz)
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, commmsg.JumpMSGOne{
			Kind:     "move",
			UID:      lm.UID,
			MoveValX: mx,
			MoveValZ: mz,
		})
	}
	return
}

func (lm *ManageMain) Local_Login_Collect() {
	if lm.which != -1 {
		return
	}
	if inputsystem.GetInputSystem().KeyDown(int(glfw.KeyP)) {
		// 按下P键就是登陆选角色
		if lm.UID == "" {
			lm.UID = fmt.Sprintf("%f", lm.gi.NowMS)
		}
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, commmsg.JumpMSGOne{
			Kind: "login",
			UID:  lm.UID,
		})
	}
}

func (lm *ManageMain) MSG_Update(msg commmsg.JumpMSGOne) {
	if msg.Kind == "login" {
		if _, found := lm.UserMap[msg.UID]; found {
			return
		} else {
			lm.UserMap[msg.UID] = int(lm.which) + 1
			lm.which++
		}
		which := lm.UserMap[msg.UID]
		fmt.Printf("{login}, (%s:%d)\n", msg.UID, which)
	} else if msg.Kind == "move" {
		// 通过 uid 找到 which
		if len(lm.UserMap) != 2 {
			return // two player login then begin the game
		}
		if which, found := lm.UserMap[msg.UID]; found {
			if which == 0 {
				fmt.Printf("{Collect}, (%s)(%d %d)\n", msg.UID, msg.MoveValX, msg.MoveValZ)
				lm.MainPlayerJump.Velx = msg.MoveValX
				lm.MainPlayerJump.Velz = msg.MoveValZ
			}
		}
	}
}

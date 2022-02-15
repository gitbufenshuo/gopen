package manage_main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/server/imple"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ManageMain struct {
	gi *game.GlobalInfo
	*gameobjects.NilManageObject
	MaxCount int
	UserMap  map[string]int64
	UID      string // 自己的uid
	Login    bool
	//
	which          int64 // which player is this client
	MainPlayer     game.GameObjectI
	MainPlayerJump *logic_jump.LogicJump
	SubPlayer      game.GameObjectI
	SubPlayerJump  *logic_jump.LogicJump
	//
	localPlayerJump *logic_jump.LogicJump
	//
	InMsgChan    chan *jump.JumpMSGTurn
	OutMsgChan   chan *jump.JumpMSGTurn
	turnMsgLocal *jump.JumpMSGTurn
	frame        int
	Turn         int64 // 回合
	//
	apressed bool
	dpressed bool
	wpressed bool
	spressed bool
	mpressed bool
	jpressed bool
	kpressed bool
	ppressed bool
	//
	serverConn net.Conn
	//
	cameraX, cameraY, cameraZ float32
	cameraFollow              *game.Transform
}

func NewManageMain(gi *game.GlobalInfo) *ManageMain {
	res := new(ManageMain)
	//
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	res.UserMap = make(map[string]int64)
	res.which = -1
	res.InMsgChan = make(chan *jump.JumpMSGTurn, 100)
	res.OutMsgChan = make(chan *jump.JumpMSGTurn, 100)
	// camera
	res.cameraX, res.cameraY, res.cameraZ = gi.MainCamera.Transform.Postion.GetValue3()
	return res
}

func (lm *ManageMain) sendToServer() {
	for msglist := range lm.OutMsgChan {
		commmsg.WriteJumpMSGTurn([]net.Conn{lm.serverConn}, msglist)
	}
}

func (lm *ManageMain) readFromServer() {
	for {
		var msg jump.JumpMSGTurn
		commmsg.ReadOnePack(lm.serverConn, &msg)
		//fmt.Println("read from server turn:", msg.Turn, msg.List, "lenuser", len(lm.UserMap))
		lm.InMsgChan <- &msg
	}
}

func (lm *ManageMain) Start() {
	inputsystem.InitInputSystem(lm.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeySpace))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyA))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyD))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyW))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyM))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyP))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyJ))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyK))
	lm.gi.SetInputSystem(inputsystem.GetInputSystem())
	//
	{
		lm.MainPlayer = modelcustom.SceneSystemIns.GetSceneOb("main", "mainplayer")
		logiclist := lm.MainPlayer.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
				lm.MainPlayerJump = v
				lm.MainPlayerJump.Chosen = true
			}
		}
	}
	{
		lm.SubPlayer = modelcustom.SceneSystemIns.GetSceneOb("main", "subplayer")
		logiclist := lm.SubPlayer.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
				lm.SubPlayerJump = v
				lm.SubPlayerJump.Chosen = true
			}
		}
	}
	lm.connect()
}

func (lm *ManageMain) connect() {
	var remoteadrr string
	if os.Args[1] == "local" {
		remoteadrr = "127.0.0.1:9090"
		lm.MaxCount = 1
		go imple.Main(1, remoteadrr)
	} else {
		lm.MaxCount = 2
		remoteadrr = os.Args[1]
	}
	time.Sleep(time.Millisecond * 10)
	conn, err := net.Dial("tcp", remoteadrr)
	if err != nil {
		panic(err)
	}
	lm.serverConn = conn
	go lm.sendToServer()
	go lm.readFromServer()
}

func (lm *ManageMain) cameraControl() {
	return
	if lm.cameraFollow == nil {
		return
	}
	mainCamera := lm.gi.MainCamera
	fox, foy, foz := lm.cameraFollow.Postion.GetValue3()
	mainCamera.Transform.Postion.SetValue3(
		lm.cameraX+fox,
		lm.cameraY+foy,
		lm.cameraZ+foz,
	)
}

func (lm *ManageMain) fromWhichGetLogic(which int64) *logic_jump.LogicJump {
	if which == 0 {
		return lm.MainPlayerJump
	}
	return lm.SubPlayerJump
}
func (lm *ManageMain) fromWhichOtherLogic(which int64) *logic_jump.LogicJump {
	if which == 0 {
		return lm.SubPlayerJump
	}
	return lm.MainPlayerJump
}

func (lm *ManageMain) Update() {
	defer func() {
		lm.frame++
	}()
	lm.cameraControl()
	////////////////
	// 收集本地指令
	if lm.frame%3 == 0 {
		lm.turnMsgLocal = new(jump.JumpMSGTurn)
	}
	lm.Local_Total_Collect()
	////////////////
	if lm.frame%3 == 2 {
		//
		lm.Local_Total_Merge()
		lm.turnMsgLocal.Turn = lm.Turn
		if len(lm.turnMsgLocal.List) == 0 {
			// idle
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
				Kind: "idle",
				Uid:  lm.UID,
			})
		}
		// 将本地指令先发出去
		lm.OutMsgChan <- lm.turnMsgLocal
		// 接收服务器指令
		//fmt.Println("准备接受服务器指令")
		inmsglist := <-lm.InMsgChan
		//fmt.Println("inmsglist", inmsglist.Turn, lm.Turn, len(inmsglist.List))
		for _, onemsg := range inmsglist.List {
			lm.MSG_Update(onemsg)
		}
		// 对本地程序步进
		lm.MainPlayerJump.OutterUpdate()
		lm.SubPlayerJump.OutterUpdate()
		// 回合收尾
		lm.Local_Collect_End()
		lm.Turn++ // 回合加1
	}
}

func (lm *ManageMain) Local_Total_Collect() {
	apressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyA))
	dpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyD))
	wpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyW))
	spressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyS))
	mpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyM))
	jpressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyJ))
	kpressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyK))
	ppressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyP))
	if !lm.apressed {
		lm.apressed = apressed
	}
	if !lm.dpressed {
		lm.dpressed = dpressed
	}
	if !lm.wpressed {
		lm.wpressed = wpressed
	}
	if !lm.spressed {
		lm.spressed = spressed
	}
	if !lm.mpressed {
		lm.mpressed = mpressed
	}
	if !lm.jpressed {
		lm.jpressed = jpressed
	}
	if !lm.ppressed {
		lm.ppressed = ppressed
	}
	if !lm.kpressed {
		lm.kpressed = kpressed
	}
	return
}

func (lm *ManageMain) Local_Collect_End() {
	lm.apressed = false
	lm.dpressed = false
	lm.wpressed = false
	lm.spressed = false
	lm.mpressed = false
	lm.jpressed = false
	lm.kpressed = false
	lm.ppressed = false
}
func (lm *ManageMain) Local_Total_Merge() {
	lm.Login_Merge()
	lm.Action_Merge()
	return
}

func (lm *ManageMain) Login_Merge() {
	if lm.Login {
		return
	}
	if lm.ppressed {
		// 按下P键就是登陆选角色
		if lm.UID == "" {
			lm.UID = fmt.Sprintf("%f", lm.gi.NowMS)
		}
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
			Kind: "login",
			Uid:  lm.UID,
		})
	}
}

func (lm *ManageMain) Action_Merge() {
	if !lm.Login {
		return
	}
	if len(lm.UserMap) != lm.MaxCount {
		return
	}
	var mx, mz int64
	if lm.apressed {
		mx = -500
	} else if lm.dpressed {
		mx = 500
	}
	if lm.wpressed {
		mz = -500
	} else if lm.spressed {
		mz = 500
	}
	if lm.jpressed { // 发起攻击
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
			Kind:     "doatt",
			Uid:      lm.UID,
			MoveValX: mx,
			MoveValZ: mz,
			M:        lm.mpressed,
		})
	} else { // 普通移动
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
			Kind:     "move",
			Uid:      lm.UID,
			MoveValX: mx,
			MoveValZ: mz,
			M:        lm.mpressed,
		})
		if lm.kpressed { // 技能 影分身
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
				Kind: "skill-yingfenshen",
				Uid:  lm.UID,
			})
		}
	}
}

func (lm *ManageMain) MSG_Update(msg *jump.JumpMSGOne) {
	if msg.Kind == "login" {
		if _, found := lm.UserMap[msg.Uid]; found {
			fmt.Println("msg.UID found ", msg.Uid)
			return
		} else {
			lm.UserMap[msg.Uid] = lm.which + 1
			lm.which++
		}
		which := lm.UserMap[msg.Uid]
		if msg.Uid == lm.UID {
			lm.Login = true
			if which == 0 {
				lm.localPlayerJump = lm.MainPlayerJump
			} else {
				lm.localPlayerJump = lm.SubPlayerJump
			}
		}
		fmt.Printf("{login}, (%s:%d) 设置相机绑定\n", msg.Uid, which)
		return
		// lm.cameraFollow = lm.localPlayerJump.Transform
	}
	switch msg.Kind {
	case "move", "doatt", "skill-yingfenshen":
		if which, found := lm.UserMap[msg.Uid]; found {
			//fmt.Printf("{Collect}, (%s)(%d %d)\n", msg.UID, msg.MoveValX, msg.MoveValZ)
			logijump := lm.fromWhichGetLogic(which)
			logijump.ProcessMSG(msg)
			break
		}
	}
}

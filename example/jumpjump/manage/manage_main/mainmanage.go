package manage_main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/server/imple"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_bullet"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
	"github.com/gitbufenshuo/gopen/help"
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
	PlayerLogicList []*logic_jump.LogicJump
	BulletLogicList []*logic_bullet.LogicBullet
	//
	auto            bool
	localPlayerJump *logic_jump.LogicJump
	//
	InMsgChan    chan *jump.JumpMSGTurn
	OutMsgChan   chan *jump.JumpMSGTurn
	turnMsgLocal *jump.JumpMSGTurn
	frame        int
	Turn         int64 // 回合
	evmanager    *pkem.EventManager
	//
	apressed  bool
	dpressed  bool
	wpressed  bool
	spressed  bool
	mpressed  bool
	jpressed  bool
	kpressed  bool
	lpressed  bool
	ppressed  bool
	opressed  bool
	f1pressed bool
	//
	serverConn net.Conn
	//
	cameraX, cameraY, cameraZ float32
	cameraRotateCounter       int64
}

func NewManageMain(gi *game.GlobalInfo) *ManageMain {
	res := new(ManageMain)
	//
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	res.UserMap = make(map[string]int64)
	res.which = -1
	res.evmanager = pkem.NewEventManager()
	res.InMsgChan = make(chan *jump.JumpMSGTurn, 100)
	res.OutMsgChan = make(chan *jump.JumpMSGTurn, 100)
	// camera
	res.cameraX, res.cameraY, res.cameraZ = gi.MainCamera.Transform.Postion.GetValue3()
	res.LandInit()
	return res
}

// 地形初始化
func (lm *ManageMain) LandInit() {
	// return
	prefab := modelcustom.PrefabSystemIns.GetPrefab("bullet")

	for x := int64(-3); x <= 3; x++ {
		for y := int64(-3); y <= 3; y++ {
			newgb := prefab.Instantiate(lm.gi)
			newgb.GetTransform().Scale.SetValue4(0.5, 0.5, 0.5, 1)
			logiclist := newgb.GetLogicSupport()
			for idx := range logiclist {
				if v, ok := logiclist[idx].(*logic_bullet.LogicBullet); ok {
					v.LogicPosX = x * 10 * 1000
					v.LogicPosY = y * 10 * 1000
					v.LogicPosZ = -5 * 10 * 1000
				}
			}

		}
	}
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
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyL))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyO))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyF1))
	lm.gi.SetInputSystem(inputsystem.GetInputSystem())
	//
	{
		lm.MainPlayer = modelcustom.SceneSystemIns.GetSceneOb("main", "mainplayer")
		logiclist := lm.MainPlayer.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
				lm.MainPlayerJump = v
			}
		}
		lm.PlayerLogicList = append(lm.PlayerLogicList, lm.MainPlayerJump)
	}
	{
		lm.SubPlayer = modelcustom.SceneSystemIns.GetSceneOb("main", "subplayer")
		logiclist := lm.SubPlayer.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_jump.LogicJump); ok {
				lm.SubPlayerJump = v
			}
		}
		lm.PlayerLogicList = append(lm.PlayerLogicList, lm.SubPlayerJump)
	}
	lm.connect()
	lm.cameraInit()
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

func (lm *ManageMain) cameraInit() {

	mainCamera := lm.gi.MainCamera
	x, y, z := mainCamera.Transform.Postion.GetValue3()
	fo := mainCamera.GetForward()
	fx, fy, fz := fo.GetValue3()
	fmt.Println("MainCameraPos:", x, y, z)
	fmt.Println("MainCameraForward:", fx, fy, fz)

	// mainCamera.SetForward(1, help.Sin(frame), help.Cos(frame))
}

func (lm *ManageMain) cameraControl() {
	if lm.f1pressed {
		lm.cameraRotateCounter++
		mainCamera := lm.gi.MainCamera
		frame := float32(lm.cameraRotateCounter) / 100
		sinvalue := help.Sin(frame)
		// -1 1 --> -PI/4 PI/4
		sinvalue *= (3.141592653 / 4)
		/////////////////////////////////////////
		R := help.Sqrt(lm.cameraY*lm.cameraY + lm.cameraZ*lm.cameraZ)
		lastz := R * help.Cos(sinvalue)
		lasty := R * help.Sin(sinvalue)
		mainCamera.Transform.Postion.SetValue3(0, lasty, lastz)
		mainCamera.SetTarget(0, 0, 0)
	}
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
		lm.Outter_Update()
		// 事件
		lm.Event_Update()
		// 回合收尾
		lm.Local_Collect_End()
		lm.Turn++ // 回合加1
	}
}

func (lm *ManageMain) Local_Total_Collect() {
	opressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyO))
	if lm.auto {
		if lm.localPlayerJump.GetLogicPosX() < 0 {
			lm.dpressed = true
		} else {
			lm.dpressed = false
		}
		if lm.localPlayerJump.GetLogicPosX() > 0 {
			lm.apressed = true
		} else {
			lm.apressed = false
		}
		//
		if lm.localPlayerJump.GetLogicPosZ() < 0 {
			lm.spressed = true
		} else {
			lm.spressed = false
		}
		if lm.localPlayerJump.GetLogicPosZ() > 0 {
			lm.wpressed = true
		} else {
			lm.wpressed = false
		}
		if !lm.opressed {
			lm.opressed = opressed
		}
		return
	}
	apressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyA))
	dpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyD))
	wpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyW))
	spressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyS))
	mpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyM))
	jpressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyJ))
	kpressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyK))
	lpressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyL))
	ppressed := inputsystem.GetInputSystem().KeyDown(int(glfw.KeyP))
	f1pressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyF1))
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
	if !lm.lpressed {
		lm.lpressed = lpressed
	}
	if !lm.opressed {
		lm.opressed = opressed
	}
	if !lm.f1pressed {
		lm.f1pressed = f1pressed
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
	lm.lpressed = false
	lm.ppressed = false
	lm.opressed = false
	lm.f1pressed = false
}
func (lm *ManageMain) Local_Total_Merge() {
	lm.Login_Merge()
	lm.Action_Merge()
	lm.Auto_Merge()
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

func (lm *ManageMain) Auto_Merge() {
	if lm.opressed {
		// 按下O键就是切换自动或者不自动
		lm.auto = !lm.auto
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
		mx = -100
	} else if lm.dpressed {
		mx = 100
	}
	if lm.wpressed {
		mz = -100
	} else if lm.spressed {
		mz = 100
	}
	if lm.jpressed { // 发起攻击
		lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
			Kind:     "doatt",
			Uid:      lm.UID,
			MoveValX: mx,
			MoveValZ: mz,
			M:        lm.mpressed,
			WhichAtt: int64(logic_jump.Att_Pugong),
		})
	} else { // 普通移动
		if mx*mx+mz*mz > 0 {
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
				Kind:     "move",
				Uid:      lm.UID,
				MoveValX: mx,
				MoveValZ: mz,
				M:        lm.mpressed,
			})
		}
		if lm.kpressed { // 技能 影分身
			fmt.Println("lm.kpressed")
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
				Kind:     "doatt",
				Uid:      lm.UID,
				WhichAtt: int64(logic_jump.Att_Skill1),
			})
		} else if lm.lpressed {
			lm.turnMsgLocal.List = append(lm.turnMsgLocal.List, &jump.JumpMSGOne{
				Kind: "jiasu", // 隐身
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
		thelogic := lm.fromWhichGetLogic(which)
		thelogic.SetPID(which)
		thelogic.SetEVM(lm.evmanager)
		if msg.Uid == lm.UID {
			lm.Login = true
			lm.localPlayerJump = thelogic
		}
		fmt.Printf("{用户登录}, (uid:%s which:%d)\n", msg.Uid, which)
		return
		// lm.cameraFollow = lm.localPlayerJump.Transform
	}
	switch msg.Kind {
	case "move", "doatt", "jiasu":
		if which, found := lm.UserMap[msg.Uid]; found {
			//fmt.Printf("{Collect}, (%s)(%d %d)\n", msg.UID, msg.MoveValX, msg.MoveValZ)
			logijump := lm.fromWhichGetLogic(which)
			logijump.ProcessMSG(msg)
			break
		}
	}
}

// 两个 player 步进
func (lm *ManageMain) Outter_Update() {
	for _, oneplayer := range lm.PlayerLogicList {
		oneplayer.OutterUpdate()
	}
	for _, onebullet := range lm.BulletLogicList {
		if onebullet.ShouldDel {

		}
		targetlogic := lm.fromWhichGetLogic(onebullet.TargetPID)
		if targetlogic == nil {
			continue
		}
		onebullet.TargetPosX = targetlogic.GetLogicPosX()
		onebullet.TargetPosY = targetlogic.GetLogicPosY()
		onebullet.TargetPosZ = targetlogic.GetLogicPosZ()
		onebullet.OutterUpdate()
	}
}

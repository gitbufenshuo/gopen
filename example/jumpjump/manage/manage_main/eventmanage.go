package manage_main

import (
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_bullet"
	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
)

func (lm *ManageMain) Event_Update() {
	evlist := lm.evmanager.GetEvList()
	for _, oneev := range evlist {
		if oneev.EK == pkem.EK_DoAtt { // 如果是攻击事件
			if oneev.SubKind == int64(logic_jump.Att_Pugong) {
				for _, oneplayer := range lm.PlayerLogicList {
					if oneplayer.GetPID() == oneev.TargetPID {
						newmsg := new(jump.JumpMSGOne)
						newmsg.Kind = "underatt"
						newmsg.PosX, newmsg.PosY, newmsg.PosZ = oneev.PosX, oneev.PosY, oneev.PosZ
						oneplayer.ProcessMSG(newmsg)
					}
				}
			} else if oneev.SubKind == int64(logic_jump.Att_Skill1) {
				prefab := modelcustom.PrefabSystemIns.GetPrefab("bullet")
				newgb := prefab.Instantiate(lm.gi)
				logiclist := newgb.GetLogicSupport()
				for idx := range logiclist {
					if v, ok := logiclist[idx].(*logic_bullet.LogicBullet); ok {
						lm.BulletLogicList = append(lm.BulletLogicList, v)
						v.LogicPosX = oneev.PosX
						v.LogicPosZ = oneev.PosZ
						v.TargetPID = oneev.TargetPID
						v.SetEVM(lm.evmanager)
						lm.BulletLogicList = append(lm.BulletLogicList, v)
					}
				}
			}
		}
	}
	// 清空事件
	lm.evmanager.Clear()
}

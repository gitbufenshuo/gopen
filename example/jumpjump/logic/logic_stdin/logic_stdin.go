package logic_stdin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/game/supports/logicinner"
	"github.com/gitbufenshuo/gopen/help"
)

type LogicStdin struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	//
	bootok  bool
	msgchan chan string
	nowgbid int
}

func newLogicStdin(gi *game.GlobalInfo) *LogicStdin {
	res := new(LogicStdin)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func NewLogicStdin(gi *game.GlobalInfo) game.LogicSupportI {
	return newLogicStdin(gi)
}

func (ls *LogicStdin) boot(gb game.GameObjectI) {
	if ls.bootok {
		return
	}
	fmt.Println("[LogicStdin boot]")
	ls.bootok = true
	ls.msgchan = make(chan string, 100)
	//
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ls.msgchan <- scanner.Text()
		}
	}()
}
func (ls *LogicStdin) jianbian(gb game.GameObjectI) {

	var thisgbid int
	select {
	case newmsg := <-ls.msgchan:
		fmt.Println("[STDIN]", newmsg)
		gbid := help.Str2Int(newmsg)
		if gbid > 0 {
			thisgbid = gbid
			ls.nowgbid = thisgbid
		}
		break
	default:
		ls.nowgbid++
		ls.nowgbid %= 100
		thisgbid = ls.nowgbid
	}
	if gb := ls.gi.GetGameObject(thisgbid); gb != nil {
		fmt.Println("[STDIN GBID]", thisgbid)
		logiclist := gb.GetLogicSupport()
		var colorcontrol bool
		var rawcc *logicinner.LogicColorControl
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logicinner.LogicColorControl); ok {
				rawcc = v
				colorcontrol = true
				break
			}
		}
		if !colorcontrol {
			newcc := logicinner.NewLogicColorControl()
			gb.AddLogicSupport(newcc)
			rawcc = newcc
		} else {
			if rawcc.Color[0] > 0.5 {
				rawcc.Color = [3]float32{0.1, 0.1, 0.1}
			} else {
				rawcc.Color = [3]float32{1, 1, 1}
			}
		}
	}

}
func (ls *LogicStdin) Update(gb game.GameObjectI) {
	ls.boot(gb)
	//
	ls.jianbian(gb)
}

func (ls *LogicStdin) Clone() game.LogicSupportI {
	fmt.Println("<><><><><><><><><")
	return nil
}

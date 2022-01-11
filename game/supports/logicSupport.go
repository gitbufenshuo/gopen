package supports

import "github.com/gitbufenshuo/gopen/game"

type NilLogic struct {
}

func NewNilLogic() *NilLogic {
	return new(NilLogic)
}

func (nl *NilLogic) Start() {

}
func (nl *NilLogic) Update() {

}
func (nl *NilLogic) OnDraw(game.GameObjectI) {

}
func (nl *NilLogic) OnDrawFinish(game.GameObjectI) {

}

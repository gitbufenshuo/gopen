package main

import (
	"fmt"
	"reflect"

	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_bullet"
	"github.com/gitbufenshuo/gopen/game"
)

func main() {
	shiti := (*logic_bullet.LogicBullet)(nil)
	shitiif := interface{}(shiti)
	shititype := reflect.TypeOf(shitiif)
	{
		var logic game.LogicSupportI
		logic = shiti
		logictpe := reflect.TypeOf(logic)
		fmt.Println("refname:", logictpe.ConvertibleTo(shititype))
	}
}

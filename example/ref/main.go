package main

import (
	"fmt"
	"reflect"

	"github.com/gitbufenshuo/gopen/example/jumpjump/logic/logic_bullet"
	"github.com/gitbufenshuo/gopen/game"
)

func main() {
	var ifa game.LogicSupportI
	var shiti logic_bullet.LogicBullet
	ifa = &shiti
	{
		a := reflect.TypeOf(ifa)
		fmt.Println("refname:", a)
		fmt.Println("  refname kind:", a.Kind())
		fmt.Println("  refname elem kind:", a.Elem().Kind())
	}
	{
		shitiaddr := &shiti
		a := reflect.TypeOf(shitiaddr)
		fmt.Println("refname:", a)
		fmt.Println("  refname kind:", a.Kind())
		fmt.Println("  refname elem kind:", a.Elem().Kind())
	}
	{
		var zhi = func() {}
		a := reflect.TypeOf(zhi)
		fmt.Println("refname:", a)
		fmt.Println("  refname kind:", a.Kind())
		// fmt.Println("  refname elem kind:", a.Elem().Kind())
	}
	{
		var ifa game.LogicSupportI
		ifa = &shiti
		var ifb *game.LogicSupportI
		ifb = &ifa
		a := reflect.TypeOf(ifb)
		fmt.Println("refname:", a)
		fmt.Println("  refname kind:", a.Kind())
		fmt.Println("  refname elem kind:", a.Elem().Kind())
		fmt.Println("  implements:", a.Elem().Implements(a.Elem()))
	}
}

package help

import "reflect"

// 前面是if
// 后面是具体
func TypeTheSame(onelogic interface{}, spLogic interface{}) bool {
	onelogicTpe := reflect.TypeOf(onelogic)
	spLogicTpe := reflect.TypeOf(spLogic)
	return onelogicTpe.ConvertibleTo(spLogicTpe)
}

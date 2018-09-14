package game

import (
	"github.com/gitbufenshuo/gopen/matmath"
)

type GameObject struct {
	set          bool
	drawEnable   bool
	drawPrepared bool
	Position     *matmath.VECX
	Rotation     *matmath.VECX
}

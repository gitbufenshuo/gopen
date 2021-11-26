package test

import (
	"testing"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TestMATX_ToIdentity(t *testing.T) {
	mat4 := matmath.GetMATX(4)
	mat4.ToIdentity()
	mat4.PrettyShow()
}

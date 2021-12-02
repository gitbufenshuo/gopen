package test

import (
	"testing"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TestMATX_ToIdentity(t *testing.T) {
	var mat4 matmath.MATX
	mat4.Init4()
	mat4.ToIdentity()
	mat4.PrettyShow()
}

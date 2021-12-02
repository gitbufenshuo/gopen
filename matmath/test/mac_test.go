package test

import (
	"fmt"
	"testing"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TestMATPrettyShow(t *testing.T) {
	var mat3 matmath.MATX
	mat3.Init3()
	mat3.PrettyShow()
	var mat4 matmath.MATX
	mat4.Init4()
	mat4.PrettyShow()
	mat4.SetEleByRowAndCol(3, 2, 1)
	mat4.SetEleByRowAndCol(2, 3, 1.5)
	mat4.PrettyShow()
}
func TestMATSetCol(t *testing.T) {
	vec4_1 := matmath.GetVECX(4)
	vec4_1.SetIndexValue(0, 1)
	vec4_1.SetIndexValue(1, 2)
	vec4_1.SetIndexValue(2, 3)

	var mat4_1 matmath.MATX
	mat4_1.Init4()
	mat4_1.SetByCol_InPlace(vec4_1, 3)

	mat4_1.PrettyShow()
}
func TestMatMatEqual(t *testing.T) {
	vec4_1 := matmath.GetVECX(4)
	vec4_1.SetIndexValue(0, 1)
	vec4_1.SetIndexValue(1, 2)
	vec4_1.SetIndexValue(2, 3)

	var mat4_1 matmath.MATX
	mat4_1.Init4()
	mat4_1.SetByCol_InPlace(vec4_1, 3)

	var mat4_2 matmath.MATX
	mat4_2.Init4()
	mat4_2.SetByCol_InPlace(vec4_1, 3)

	if !mat4_1.EqualsTo(&mat4_2) {
		t.Error("wrong")
	}
}
func TestMatMat(t *testing.T) {
	vec4_1 := matmath.GetVECX(4)
	vec4_1.SetIndexValue(0, 1)
	vec4_1.SetIndexValue(1, 2)
	vec4_1.SetIndexValue(2, 3)
	vec4_1.SetIndexValue(3, 3)

	var mat4_1 matmath.MATX
	mat4_1.Init4()
	mat4_1.SetByCol_InPlace(vec4_1, 3)
	mat4_1.SetByCol_InPlace(vec4_1, 1)

	var mat4_2 matmath.MATX
	mat4_2.Init4()
	mat4_2.SetByCol_InPlace(vec4_1, 3)
	mat4_2.SetByCol_InPlace(vec4_1, 2)

	fmt.Println("--")
	mat4_1.PrettyShow()
	mat4_1.RightMul_InPlace(&mat4_2)
	mat4_1.PrettyShow()
	mat4_2.PrettyShow()
}
func TestMatVec(t *testing.T) {
	vec4_1 := matmath.GetVECX(4)
	vec4_1.SetIndexValue(0, 1)
	vec4_1.SetIndexValue(1, 2)
	vec4_1.SetIndexValue(2, 3)
	vec4_1.SetIndexValue(3, 4)
	fmt.Println("--")
	vec4_1.PrettyShow()
	var mat4_1 matmath.MATX
	mat4_1.Init4()
	mat4_1.SetByCol_InPlace(vec4_1, 1)
	mat4_1.SetByCol_InPlace(vec4_1, 2)
	mat4_1.SetByCol_InPlace(vec4_1, 3)
	mat4_1.SetByCol_InPlace(vec4_1, 4)
	mat4_1.PrettyShow()
	vec4_1.RightMul_InPlace(&mat4_1)
	vec4_1.PrettyShow()
}

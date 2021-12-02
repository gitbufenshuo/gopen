package test

import (
	"math"
	"testing"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TestOne(t *testing.T) {
	var vec3 matmath.VECX
	vec3.Init3()
	vec3.SetIndexValue(0, 1)
	vec3.SetIndexValue(1, 2)
	vec3.SetIndexValue(2, 3)
	vec3.PrettyShow()
}
func TestVecEqual(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 2)
	vec3_1.SetIndexValue(2, 3)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 1)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 3)

	if !vec3_1.EqualsTo(&vec3_2) {
		t.Error("wrong")
	}
	vec3_2.SetIndexValue(2, 4)
	if vec3_1.EqualsTo(&vec3_2) {
		t.Error("wrong")
	}
}
func TestVecScale(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 1)
	vec3_1.SetIndexValue(2, 1)
	vec3_1.Scale_InPlace(2)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 2)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 2)

	if !vec3_1.EqualsTo(&vec3_2) {
		t.Error("wrong")
	}
}

func TestVecDiff(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 1)
	vec3_1.SetIndexValue(2, 1)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 2)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 2)

	if vec3_1.EqualsTo(&vec3_2) {
		t.Error("wrong")
	}
}
func TestVecAdd(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 1)
	vec3_1.SetIndexValue(2, 1)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 2)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 2)

	var shouldRes matmath.VECX
	shouldRes.Init3()
	shouldRes.SetIndexValue(0, 3)
	shouldRes.SetIndexValue(1, 3)
	shouldRes.SetIndexValue(2, 3)

	res := vec3_1.Add(&vec3_2)

	if !shouldRes.EqualsTo(&res) {
		t.Error("wrong")
	}
	vec3_1.Add_InPlace(&vec3_2)

	if !shouldRes.EqualsTo(&vec3_1) {
		t.Error("wrong")
	}
}
func TestVecInterpolation(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 1)
	vec3_1.SetIndexValue(2, 1)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 2)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 2)

	var vec3_3 matmath.VECX
	vec3_3.Init3()
	vec3_3.SetIndexValue(0, 1.5)
	vec3_3.SetIndexValue(1, 1.5)
	vec3_3.SetIndexValue(2, 1.5)

	begin := vec3_1.Interpolation(&vec3_2, 0)
	if !begin.EqualsTo(&vec3_1) {
		t.Error("wrong")
	}
	end := vec3_1.Interpolation(&vec3_2, 1)
	if !end.EqualsTo(&vec3_2) {
		t.Error("wrong")
	}
	mid := vec3_1.Interpolation(&vec3_2, 0.5)
	if !mid.EqualsTo(&vec3_3) {
		t.Error("wrong")
	}
}

func TestVecDot(t *testing.T) {
	var vec3_1 matmath.VECX
	vec3_1.Init3()
	vec3_1.SetIndexValue(0, 1)
	vec3_1.SetIndexValue(1, 1)
	vec3_1.SetIndexValue(2, 1)

	var vec3_2 matmath.VECX
	vec3_2.Init3()
	vec3_2.SetIndexValue(0, 2)
	vec3_2.SetIndexValue(1, 2)
	vec3_2.SetIndexValue(2, 2)

	cases := []struct {
		v1   *matmath.VECX
		v2   *matmath.VECX
		want float32
	}{
		{&vec3_1, &vec3_2, 6},
	}

	for _, c := range cases {
		got := c.v1.Dot(c.v2)
		if math.Abs(float64(got-c.want)) > 0.00001 {
			t.Error("wrong")
		}
	}
}

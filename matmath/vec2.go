package matmath

import (
	"fmt"
	"math"

	"github.com/gitbufenshuo/gopen/help"
)

type Vec2 struct {
	data [2]float32
}

func (vec2 *Vec2) PrettyShow(reason string) {
	fmt.Printf("%s:(", reason)
	for i := 0; i != 2; i++ {
		fmt.Printf(" %v ", vec2.data[i])
	}
	fmt.Printf(")\n")
}

func CreateVec2(x, y float32) Vec2 {
	return Vec2{
		data: [2]float32{x, y},
	}
}
func CreateVec2FromVec4(input Vec4) Vec2 {
	return Vec2{
		data: [2]float32{
			input.data[0],
			input.data[1],
		},
	}
}

// 计算两个vec2的夹角
func Vec2Angle(one, two *Vec2) float32 {
	return help.ArcCos(
		Vec2InnerProduct(one, two) / (Vec2Length(one) * Vec2Length(two)),
	)
}

// 计算两个vec2的内积
func Vec2InnerProduct(one, two *Vec2) float32 {
	return one.data[0]*two.data[0] + one.data[1]*two.data[1]
}

// 计算vec2的模
func Vec2Length(one *Vec2) float32 {
	return help.Sqrt(one.data[0]*one.data[0] + one.data[1]*one.data[1])
}

// 任意四个点组成的四边形 注意一定是 凸四边形
// 检测另一个点是否在这个四边形内
// 不一定是矩形，也不一定是和xy轴对齐的
func Vec2BoundCheck(bounds []*Vec2, target *Vec2) bool {
	tempvec2_1 := CreateVec2(0, 0)
	tempvec2_2 := CreateVec2(0, 0)
	//
	var totalangle float32
	for idx := 0; idx != len(bounds); idx++ {
		tempvec2_1.data[0], tempvec2_1.data[1] = bounds[idx].data[0]-target.data[0], bounds[idx].data[1]-target.data[1]
		tempvec2_2.data[0], tempvec2_2.data[1] = bounds[(idx+1)%4].data[0]-target.data[0], bounds[(idx+1)%4].data[1]-target.data[1]
		totalangle += Vec2Angle(&tempvec2_1, &tempvec2_2)
	}
	return help.TheSame(totalangle, math.Pi*2)
}

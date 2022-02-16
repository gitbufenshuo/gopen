package resource

import (
	"github.com/gitbufenshuo/gopen/matmath"
)

var modelSphereJSON = `
{
	"Vertices":[
		-0.33333,-1,1, 0.33333,0, 0,0,1,
		0.33333,-1,1, 0.66666,0, 0,0,1,
		1,-0.33333,1, 1,0.33333, 0,0,1,
		1,0.33333,1, 1,0.6666, 0,0,1,
		0.33333,1,1, 0.66666,1, 0,0,1,
		-0.33333,1,1, 0.33333,1, 0,0,1,
		-1,0.33333,1, 0,0.6666, 0,0,1,
		-1,-0.33333,1, 0,0.33333, 0,0,1,
		0,0,1, 0.5,0.5, 0,0,1
		],
		"Indices":[8,0,1, 8,1,2, 8,2,3, 8,3,4, 8,4,5, 8,5,6, 8,6,7, 8,7,0],
		"Stripes":[3,2,3]
}`

func NewSphereModel_BySpec(pivot, size matmath.Vec4) *Model {
	res := NewModel()
	//
	newvers := []float32{}
	newins := []uint32{}
	inAdd := uint32(9)
	res.ReadFromContent(modelSphereJSON)
	for ry := float32(90); ry < 360; ry += 90 {
		for vidx := 0; vidx != 9; vidx++ {
			x, y, z := res.Vertices[vidx*8+0], res.Vertices[vidx*8+1], res.Vertices[vidx*8+2]
			u, v := res.Vertices[vidx*8+3], res.Vertices[vidx*8+4]
			nx, ny, nz := res.Vertices[vidx*8+5], res.Vertices[vidx*8+6], res.Vertices[vidx*8+7]
			xyz := matmath.CreateVec4(x, y, z, 1)
			nxyz := matmath.CreateVec4(nx, ny, nz, 1)
			xyzbian, nxyzbian := matmath.RotateY(xyz, ry), matmath.RotateY(nxyz, ry)
			bianx, biany, bianz := xyzbian.GetValue3()
			biannx, bianny, biannz := nxyzbian.GetValue3()
			//
			newvers = append(newvers, bianx, biany, bianz, u, v, biannx, bianny, biannz)
		}
		for idx := 0; idx != 24; idx++ {
			newins = append(newins, inAdd+res.Indices[idx])
		}
		inAdd += 9
	}
	{
		var xrlist = []float32{90, -90}
		for _, onerx := range xrlist {
			for vidx := 0; vidx != 9; vidx++ {
				x, y, z := res.Vertices[vidx*8+0], res.Vertices[vidx*8+1], res.Vertices[vidx*8+2]
				u, v := res.Vertices[vidx*8+3], res.Vertices[vidx*8+4]
				nx, ny, nz := res.Vertices[vidx*8+5], res.Vertices[vidx*8+6], res.Vertices[vidx*8+7]
				xyz := matmath.CreateVec4(x, y, z, 1)
				nxyz := matmath.CreateVec4(nx, ny, nz, 1)
				xyzbian, nxyzbian := matmath.RotateX(xyz, onerx), matmath.RotateX(nxyz, onerx)
				bianx, biany, bianz := xyzbian.GetValue3()
				biannx, bianny, biannz := nxyzbian.GetValue3()
				//
				newvers = append(newvers, bianx, biany, bianz, u, v, biannx, bianny, biannz)
				// fmt.Println(onerx, "qian", x, y, z, "Zuihou", bianx, biany, bianz)
			}
			for idx := 0; idx != 24; idx++ {
				newins = append(newins, inAdd+res.Indices[idx])
			}
			inAdd += 9
		}
	}
	res.Vertices = append(res.Vertices, newvers...)
	res.Indices = append(res.Indices, newins...)
	// fmt.Println("SSSSSSSSSS", len(res.Vertices)/8)
	// for zu := 0; zu != 6; zu++ {
	// 	fmt.Println("zuIndex", zu, res.Indices[zu*24:(zu+1)*24])
	// }
	//
	basever := []float32{
		1, 0.33333, 1, 0.5, 0.5, 0.5773, 0.5773, 0.5773,
		1, 1, 0.33333, 0.5, 0.5, 0.5773, 0.5773, 0.5773,
		0.33333, 1, 1, 0.5, 0.5, 0.5773, 0.5773, 0.5773,
	}
	baseindex := []uint32{54, 55, 56}
	for ry := float32(90); ry < 360; ry += 90 {
		for vidx := 0; vidx != 3; vidx++ { // 三个点
			x, y, z := basever[vidx*8+0], basever[vidx*8+1], basever[vidx*8+2]
			u, v := basever[vidx*8+3], basever[vidx*8+4]
			nx, ny, nz := basever[vidx*8+5], basever[vidx*8+6], basever[vidx*8+7]
			xyz := matmath.CreateVec4(x, y, z, 1)
			nxyz := matmath.CreateVec4(nx, ny, nz, 1)
			xyzbian, nxyzbian := matmath.RotateY(xyz, ry), matmath.RotateY(nxyz, ry)
			bianx, biany, bianz := xyzbian.GetValue3()
			biannx, bianny, biannz := nxyzbian.GetValue3()
			basever = append(basever, bianx, biany, bianz, u, v, biannx, bianny, biannz)
		}
		// fmt.Println(baseindex[len(baseindex)-1]+1, baseindex[len(baseindex)-1]+2, baseindex[len(baseindex)-1]+3)
		baseindex = append(baseindex, baseindex[len(baseindex)-1]+1, baseindex[len(baseindex)-1]+2, baseindex[len(baseindex)-1]+3)
	}
	res.Vertices = append(res.Vertices, basever...)
	basever = nil
	basever = []float32{
		1, -0.33333, 1, 0.5, 0.5, 0.5773, -0.5773, 0.5773,
		0.33333, -1, 1, 0.5, 0.5, 0.5773, -0.5773, 0.5773,
		1, -1, 0.33333, 0.5, 0.5, 0.5773, -0.5773, 0.5773,
	}
	for ry := float32(90); ry < 360; ry += 90 {
		for vidx := 0; vidx != 3; vidx++ { // 三个点
			x, y, z := basever[vidx*8+0], basever[vidx*8+1], basever[vidx*8+2]
			u, v := basever[vidx*8+3], basever[vidx*8+4]
			nx, ny, nz := basever[vidx*8+5], basever[vidx*8+6], basever[vidx*8+7]
			xyz := matmath.CreateVec4(x, y, z, 1)
			nxyz := matmath.CreateVec4(nx, ny, nz, 1)
			xyzbian, nxyzbian := matmath.RotateY(xyz, ry), matmath.RotateY(nxyz, ry)
			bianx, biany, bianz := xyzbian.GetValue3()
			biannx, bianny, biannz := nxyzbian.GetValue3()
			basever = append(basever, bianx, biany, bianz, u, v, biannx, bianny, biannz)
		}
		baseindex = append(baseindex, baseindex[len(baseindex)-1]+1, baseindex[len(baseindex)-1]+2, baseindex[len(baseindex)-1]+3)
	}
	baseindex = append(baseindex, baseindex[len(baseindex)-1]+1, baseindex[len(baseindex)-1]+2, baseindex[len(baseindex)-1]+3)
	res.Vertices = append(res.Vertices, basever...)
	res.Indices = append(res.Indices, baseindex...)
	///////////////////////////////////////////////
	// 1. 考虑 size
	var maxx, maxy, maxz float32
	var minx, miny, minz float32
	sizex, sizey, sizez := size.GetValue3()
	for idx := 0; idx != 78; idx++ {
		res.Vertices[idx*8+0] *= sizex
		res.Vertices[idx*8+1] *= sizey
		res.Vertices[idx*8+2] *= sizez
		if maxx < res.Vertices[idx*8+0] {
			maxx = res.Vertices[idx*8+0]
		}
		if maxy < res.Vertices[idx*8+1] {
			maxy = res.Vertices[idx*8+1]
		}
		if maxz < res.Vertices[idx*8+2] {
			maxz = res.Vertices[idx*8+2]
		}
	}
	minx, miny, minz = -maxx, -maxy, -maxz
	///////////////////////////////////////////////
	// 2. 考虑 pivot
	pivotx, pivoty, pivotz := pivot.GetValue3()
	xoff, yoff, zoff := (minx-maxx)*pivotx/2, (miny-maxy)*pivoty/2, (minz-maxz)*pivotz/2
	//
	for idx := 0; idx != 78; idx++ {
		res.Vertices[idx*8+0] += xoff
		res.Vertices[idx*8+1] += yoff
		res.Vertices[idx*8+2] += zoff
	}
	// fmt.Println("sphere model", "vercount:", len(res.Vertices)/8)
	return res
}

package logic_jump

func (lj *LogicJump) GetVelX() int64 {
	return (lj.movex + lj.underattx) * lj.moveSpeed / 100
}

func (lj *LogicJump) GetVelZ() int64 {
	return (lj.movez + lj.underattz) * lj.moveSpeed / 100
}

func (lj *LogicJump) OnForce() {
	var upForce int64
	if lj.logicposy <= 0 {
		lj.logicposy = 0
		upForce = -lj.gravity // 如果在地面，向上的弹力应该正好与重力相反
		lj.Vely = 0
	}
	//
	//deltams := float32(lj.gi.FrameElapsedMS / 1000) // 单位变成秒
	mergeforce := lj.gravity + upForce // 合力
	lj.Vely += (mergeforce) * 10
	lj.logicposy += lj.Vely
	lj.logicposx += lj.GetVelX()
	lj.logicposz += lj.GetVelZ()
	lj.factor = 5
	{
		// clamp x and z
		if lj.logicposx < -16*1000 {
			lj.logicposx = -16 * 1000
		}
		if lj.logicposx > 16*1000 {
			lj.logicposx = 16 * 1000
		}

		if lj.logicposz < -10*1000 {
			lj.logicposz = -10 * 1000
		}
		if lj.logicposz > 10*1000 {
			lj.logicposz = 10 * 1000
		}
	}

	// fmt.Printf("lj.logicposy:%f lj.vel:%f imp:%f mode:%v\n", lj.logicposy, lj.vel, mergeforce*deltams, lj.PlayerMode)
}

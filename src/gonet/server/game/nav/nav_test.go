package detour_test

import (
	"gonet/base"
	"gonet/server/zone/game/lmath"
	detour "gonet/server/zone/game/nav"
	"testing"
)

func Test66(t *testing.T)  {
	t.Log("测试巡径")
	dt := detour.NewDetour(1000)
	dt.Load("../../../../bin/nav/scene1.obj.tile.bin")
	for j := 0; j < 10000; j++{
		dt.FindPath(lmath.Point3F{-500, 0, 0}, lmath.Point3F{0,0,0}, base.NewVector())
	}
}



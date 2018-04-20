package game

const(
	TM_NORMAL = 0
	TM_BLOCK = 1
	TM_SHADOW = 2
	TM_MINE = 3
	TM_PLAYER_BLOCK = 6
	TM_PLAYER_BLOCK_SHADOW = 7
)//TileMask

type(
	Tile struct {
		
	}
)

func getEntitiesInTile(simsMap map [int] *SimObject, x int, y int, resultsims []*SimObject)  {
	for _,v := range simsMap{
		v.Type = 1
	}
}
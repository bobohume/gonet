package social

import (
	"server/world"
	"actor"
	"database/sql"
	"base"
	"server/world/player"
	"db"
)

const(
	Temp = iota			// 临时
	Friend	= iota			// 好友
	Consort = iota			// 配偶
	Master	= iota			// 师傅
	Prentice = iota		// 徒弟
	Enemy = iota			// 仇人
	Mute = iota			// 屏蔽(黑名单)
	Count = iota

	sqlTable = "tbl_social"
)

//enSocialError
const (
	SocialError_None 		= iota
	SocialError_Unknown	= iota			// 未知
	SocialError_Self		= iota				// 不能和自身成为社会关系
	SocialError_MaxCount	= iota			// 此社会关系人数达到最大上限
	SocialError_NotFound	= iota			// 目标玩家不存在
	SocialError_Existed	= iota			// 目标玩家已是此在社会关系列表中
	SocialError_Unallowed	= iota	    // 该操作不允许
	SocialError_DbError	= iota	       	// 数据库操作错误
)

type (
	SocialItem struct {
		PlayerId int	`primary`
		TargetId int	`primary`
		Type	int8
		FriendValue	int
	}

	SOCIALITEMMAP map[int] *SocialItem

	SocialMgr struct{
		actor.Actor
		m_db *sql.DB
		m_Log *base.CLog
		m_SocialMap map[int]	SOCIALITEMMAP
	}

	ISocialMgr interface {
		actor.IActor
		makeLink(int, int, int8) int//加好友
		destoryLink(int, int) int//删除好友
		addFreindValue(int, int , int) int//增加好友度

		isFriendType(int8) bool
		isBestFriendType(int8) bool
		hasMakeLink(int8, int8)	bool
	}
)

var(
	SOCIALMGR SocialMgr
)

func getSocialTypeMaxCount(Type int8) int{
	if Type >= Count ||  Type < 0{
		return 0
	}

	SocialTypeMaxCount	:= [Count]int{
		50, 50, 1, 1, 5, 20, 100,
	}

	return  SocialTypeMaxCount[Type]
}

func (this *SocialMgr) Init(num int) {
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("social", this)
	this.m_SocialMap = make(map[int] SOCIALITEMMAP)

	this.RegisterCall("C_W_MakeLinkRequest", func(caller *actor.Caller, PlayerId, TargetId int, Type int8) {
		pPlayer := player.PLAYERSIMPLEMGR.GetPlayerDataById(PlayerId)
		pTarget	:= player.PLAYERSIMPLEMGR.GetPlayerDataById(TargetId)
		if pPlayer == nil || pTarget == nil{
			this.m_Log.Printf("查询玩家id[%d][%d]数据为空", PlayerId, TargetId)
			return
		}
	})

	this.Actor.Start()
}

func (this *SocialMgr) isFreindType(Type int8) bool{
	if Type != Temp && Type != Mute && Type != Enemy{
		return	true
	}
	return  false
}

func (this *SocialMgr) isBestFriendType(Type int8) bool{
	if Type != Friend && Type != Temp && Type != Mute && Type != Enemy{
		return true
	}
	return	false
}

func (this *SocialMgr) hasMakeLink(oldType, newType int8) bool{
	if oldType == newType{
		return false
	}

	if oldType != Temp && newType == Temp{
		return false
	}

	if (this.isBestFriendType(oldType) || this.isBestFriendType(newType)){
		return false
	}

	if (oldType == Friend || oldType == Mute) && newType == Enemy{
		return false
	}

	return true
}

func (this *SocialMgr) 	makeLink(PlayerId, TargetId int, Type int8) int{
	if Type >= Count{
		return SocialError_Unknown
	}

	if this.isBestFriendType(Type){
		return SocialError_Unallowed
	}

	if PlayerId == TargetId{
		return SocialError_Self
	}

	SocialMap, exist := this.m_SocialMap[PlayerId]
	if exist{
		nCount := 0
		firstPlayerId := 0
		var Item *SocialItem

		for _, v := range SocialMap{
			if Item == nil && v.TargetId == TargetId{
				Item = v
			}

			if v.Type == Type{
				nCount++
				if firstPlayerId == 0{
					firstPlayerId = v.TargetId
				}
			}
		}

		// type is full, send error
		if nCount > getSocialTypeMaxCount(Type){
			if Type != Enemy && Type != Temp{
				return SocialError_MaxCount
			}else{
				// 当仇人或临时好友添加满时，删除一个仇人或临时好友
				if firstPlayerId != 0 && SocialError_None != this.destoryLink( PlayerId, firstPlayerId ){
					return SocialError_DbError
				}
			}
		}

		if Item != nil{
			if Item.Type == Type{
				return SocialError_Existed
			}else{
				if !this.hasMakeLink(Item.Type, Type){
					return SocialError_Unallowed
				}
			}

			Item.PlayerId = PlayerId
			Item.TargetId = TargetId
			Item.Type = Type
			Item.FriendValue = 0
			this.m_db.Exec(db.UpdateSql(Item, sqlTable))
			this.m_Log.Printf("更新社会关系playerId=%d,destPlayerId=%d,newType=%d", PlayerId, TargetId, Item.Type)
			return SocialError_None
		}
	}else{
		Item := &SocialItem{}
		Item.PlayerId = PlayerId
		Item.TargetId = TargetId
		Item.Type = Type
		Item.FriendValue = 0

		SocialMap = make(SOCIALITEMMAP)
		SocialMap[TargetId] = Item
		this.m_SocialMap[PlayerId] = SocialMap
		this.m_db.Exec(db.InsertSql(Item, sqlTable))
		this.m_Log.Printf("新增社会关系playerId=%d,destPlayerId=%d,type=%d", PlayerId, TargetId, Item.Type)
		return SocialError_None
	}

	return SocialError_DbError
}

func (this *SocialMgr) destoryLink(PlayerId, TargetId int) int{
	if PlayerId == TargetId{
		return  SocialError_Self
	}

	SocialMap, exist := this.m_SocialMap[PlayerId]
	if !exist{
		return SocialError_NotFound
	}

	Item, exist := SocialMap[TargetId]
	if(!exist){
		return SocialError_NotFound
	}

	if Item.Type == Friend || Item.Type == Mute || Item.Type == Temp || Item.Type == Enemy{
		this.m_db.Exec(db.DeleteSql(Item, sqlTable))
		delete(SocialMap, TargetId)
		return SocialError_None
	}

	return SocialError_Unallowed
}

func (this *SocialMgr) addFriendValue(PlayerId, TargetId, Value int) int{
	SocialMap1, exist1 := this.m_SocialMap[PlayerId]
	SocialMap2, exist2 := this.m_SocialMap[TargetId]
	if !exist1 || !exist2{
		return 0
	}

	Item1, exist1 := SocialMap1[TargetId]
	Item2, exist2 := SocialMap2[PlayerId]
	if !exist1 || !exist2{
		return 0
	}

	if !this.isFreindType(Item1.Type) || !this.isFreindType(Item2.Type){
		return 0
	}

	Item1.FriendValue +=Value
	Item2.FriendValue = Item1.FriendValue
	this.m_db.Exec(db.UpdateSqlEx(Item1, sqlTable, "FriendValue"))
	this.m_db.Exec(db.UpdateSqlEx(Item2, sqlTable, "FriendValue"))
	return Value
}
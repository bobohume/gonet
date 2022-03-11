package gate

import "gonet/base"

//消息防火墙
var(
	s_clientCheckFilters map[string] bool//use for no check playerid
	s_clientCheckFilterInit bool
)

func IsCheckClient(msg string) bool {
	if !s_clientCheckFilterInit{
		s_clientCheckFilters = make(map[string] bool)
		s_clientCheckFilters[base.ToLower("LoginAccountRequest")] = true
		s_clientCheckFilters[base.ToLower("CreatePlayerRequest")] = true
		s_clientCheckFilters[base.ToLower("LoginPlayerRequset")] = true
		s_clientCheckFilterInit = true
	}

	_,exist := s_clientCheckFilters[msg]
	return exist
}

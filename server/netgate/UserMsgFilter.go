package netgate

import "gonet/base"

//消息防火墙
var(
	s_clientCheckFilters map[string] bool//use for no check accountid
	s_clientCheckFilterInit bool
)

func IsCheckClient(msg string) bool {
	if !s_clientCheckFilterInit{
		s_clientCheckFilters = make(map[string] bool)
		s_clientCheckFilters[base.ToLower("C_A_LoginRequest")] = true
		s_clientCheckFilters[base.ToLower("C_A_RegisterRequest")] = true
		s_clientCheckFilterInit = true
	}

	_,exist := s_clientCheckFilters[msg]
	return exist
}

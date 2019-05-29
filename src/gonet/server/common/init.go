package common

import "gonet/message"

//调用在 全局变量之后，为了防止有些全局变量依赖，比如a依赖b，但是a先创建了，导致出问题

func Init(){
	message.Init()
}

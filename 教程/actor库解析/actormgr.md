actor我怎么知道我要往那个actor发送呢

    对于actor绑定network,在network里面分析。
    
    那actor对actor呢 ActorMgr.go里面有,通过AddActor类型名字注册actor,这个适用于全局的actor,考虑性能
    
    没用锁模式,直接map多线程访问,主要还是map,作为全局定义,在多线程不增加删除的情况下是安全的
    
    那么不是全局的怎么办,这种需求少,目前只有玩家才有。比如玩家管理统筹所有玩家actor,登录和离开过playermgr,
    
    其他消息通过玩家id分发到玩家actor

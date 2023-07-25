package cache

//CacheKey 缓存标识
type CacheKey string

//CacheDataModelInterface 缓存数据模板接口
type CacheDataModelInterface interface {
	GetCacheKey(interface{}, interface{}) (CacheKey, error) //获取缓存标识 参数1:缓存key 参数2:缓存数据 响应1:缓存标识 响应2:错误信息
	EliminatePolicyFunc(CacheKey, interface{}) bool         //淘汰策略 参数1:缓存key 参数2:缓存数据 响应1:是否淘汰
	PollUpdatePolicyFunc(interface{}) interface{}           //巡检更新策略  参数1： 之前的缓存数据列表  响应1：需要缓存的新增数据列表（差异数据）
}

//CacheInterface 缓存接口
type CacheInterface interface {
	AddCacheData(key CacheKey, data interface{}) bool                         //AddCacheData 手动强制添加缓存数据，不会触发巡检更新，仅触发有效性检测逻辑 参数1:缓存key 参数2:缓存数据 响应1:是否添加成功
	BatchAddCacheData(list interface{}) error                                 //BatchAddCacheData 手动强制批量添加缓存数据，不会触发巡检更新，仅触发有效性检测逻辑  参数1:缓存数据列表 响应1:错误信息
	GetCacheInfo(key CacheKey, hasIgnorePoll bool) (interface{}, bool, error) //GetCacheInfo 获取缓存数据 参数1:缓存key 参数2:是否忽略巡检维护新数据 响应1:缓存数据 响应2:是否存在该缓存 响应3:错误信息
	GetCacheData(hasIgnorePoll bool) (map[CacheKey]interface{}, error)        //GetCacheData 获取缓存数据列表 参数1:是否忽略巡检维护新数据 响应1:缓存数据列表 响应2:错误信息
	MustPollUpdatePolicyFunc() error                                          //MustPollUpdatePolicyFunc 强制巡检维护新数据 响应1:错误信息
	DelCache(key CacheKey)                                                    //DelCache 删除缓存数据 参数1:缓存key
}

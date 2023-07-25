package cache

import (
	"errors"
	"reflect"
	"sync"
)

//MemCache 内存缓存
type MemCache struct {
	sync.RWMutex
	cacheDataList map[CacheKey]interface{} //缓存数据列表
	dataInterface CacheDataModelInterface  //数据接口
}

//NewMemCache 创建缓存
func NewMemCache(dataInterface CacheDataModelInterface) *MemCache {
	return &MemCache{
		cacheDataList: make(map[CacheKey]interface{}),
		dataInterface: dataInterface,
	}
}

//BatchAddCacheData 批量添加缓存数据
func (_memCache *MemCache) BatchAddCacheData(list interface{}) error {
	typeObj := reflect.ValueOf(list)
	if !typeObj.IsValid() {
		return nil
	}
	switch typeObj.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < typeObj.Len(); i++ {
			info := typeObj.Index(i).Interface()
			//fmt.Println()
			//获取缓存key
			cacheKey, err := _memCache.dataInterface.GetCacheKey(i, info)
			if err != nil {
				return err
			}
			//添加缓存数据
			_memCache.AddCacheData(cacheKey, info)
		}
	case reflect.Map:
		for _, cacheKey := range typeObj.MapKeys() {
			info := typeObj.MapIndex(cacheKey).Interface()
			//获取缓存key
			cacheKey, err := _memCache.dataInterface.GetCacheKey(cacheKey.Interface(), info)
			if err != nil {
				return err
			}
			//添加缓存数据
			_memCache.AddCacheData(cacheKey, info)
		}
	//case reflect.:

	default:
		return errors.New("只支持数组、切片、map类型")
	}
	return nil
}

//AddCacheData 添加缓存数据
func (_memCache *MemCache) AddCacheData(key CacheKey, data interface{}) bool {
	//判断是否需要淘汰
	if _memCache.dataInterface.EliminatePolicyFunc(key, data) {
		return false
	}
	_memCache.cacheDataList[key] = data
	return true
}

//pollUpdatePolicyFunc 巡检维护新数据
func (_memCache *MemCache) pollUpdatePolicyFunc() error {
	waitAddData := _memCache.dataInterface.PollUpdatePolicyFunc(_memCache.cacheDataList)
	err := _memCache.BatchAddCacheData(waitAddData)
	return err
}

//GetCacheInfo 获取缓存数据
func (_memCache *MemCache) GetCacheInfo(key CacheKey, hasIgnorePoll bool) (interface{}, bool, error) {
	_memCache.Lock()
	defer _memCache.Unlock()
	//判断是否需要巡检维护新数据
	if !hasIgnorePoll {
		//巡检维护新数据
		err := _memCache.pollUpdatePolicyFunc()
		if err != nil {
			return nil, false, err
		}
	}

	//获取缓存数据
	data, ok := _memCache.cacheDataList[key]
	if ok {
		if _memCache.expiredCheck(key, data) {
			return nil, false, nil
		}
		return data, true, nil
	}
	return nil, false, nil
}

//GetCacheData 获取缓存数据列表
func (_memCache *MemCache) GetCacheData(hasIgnorePoll bool) (map[CacheKey]interface{}, error) {
	_memCache.Lock()
	defer _memCache.Unlock()
	//判断是否需要巡检维护新数据
	if !hasIgnorePoll {
		//巡检维护新数据
		err := _memCache.pollUpdatePolicyFunc()
		if err != nil {
			return nil, err
		}
	}
	//淘汰检测
	for cacheKey, cacheData := range _memCache.cacheDataList {
		//淘汰检测 符合淘汰规则会删除当前遍历的数据
		_memCache.expiredCheck(cacheKey, cacheData)
	}
	return _memCache.cacheDataList, nil
}

//MustPollUpdatePolicyFunc 强制巡检维护新数据
func (_memCache *MemCache) MustPollUpdatePolicyFunc() error {
	//巡检维护新数据
	return _memCache.pollUpdatePolicyFunc()
}

//DelCache 删除缓存数据
func (_memCache *MemCache) DelCache(key CacheKey) {
	_memCache.Lock()
	defer _memCache.Unlock()
	delete(_memCache.cacheDataList, key)
}

//expiredCheck 过期检查
func (_memCache *MemCache) expiredCheck(key CacheKey, data interface{}) bool {
	if _memCache.dataInterface.EliminatePolicyFunc(key, data) {
		delete(_memCache.cacheDataList, key)
		return true
	}
	return false
}

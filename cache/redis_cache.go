package cache

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"sync"
)

//RedisCache redis缓存
type RedisCache struct {
	sync.RWMutex
	dataInterface CacheDataModelInterface //数据接口
	redisConn     redis.Conn              //redis连接
	hashMapKey    string                  //redis hashmap Key
}

//NewRedisCache 创建redis缓存实例
func NewRedisCache(dataInterface CacheDataModelInterface, redisConn redis.Conn, hashMapKey string) *RedisCache {
	return &RedisCache{
		redisConn:     redisConn,
		hashMapKey:    hashMapKey,
		dataInterface: dataInterface,
	}
}

//BatchAddCacheData 批量添加缓存数据
func (_redisCache *RedisCache) BatchAddCacheData(list interface{}) error {
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
			cacheKey, err := _redisCache.dataInterface.GetCacheKey(i, info)
			if err != nil {
				return err
			}
			//添加缓存数据
			_redisCache.AddCacheData(cacheKey, info)
		}
	case reflect.Map:
		for _, cacheKey := range typeObj.MapKeys() {
			info := typeObj.MapIndex(cacheKey).Interface()
			//获取缓存key
			cacheKey, err := _redisCache.dataInterface.GetCacheKey(cacheKey.Interface(), info)
			if err != nil {
				return err
			}
			//添加缓存数据
			_redisCache.AddCacheData(cacheKey, info)
		}
	default:
		return errors.New("只支持数组、切片、map类型")
	}
	return nil
}

//AddCacheData 添加缓存数据
func (_redisCache *RedisCache) AddCacheData(key CacheKey, data interface{}) bool {
	//判断是否需要淘汰
	if _redisCache.dataInterface.EliminatePolicyFunc(key, data) {
		return false
	}
	reflectDataObj := reflect.ValueOf(data)
	if reflectDataObj.Kind() == reflect.Ptr {
		reflectDataObj = reflectDataObj.Elem()
	}
	var addData interface{}
	var err error
	switch reflectDataObj.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
		addData, err = json.Marshal(reflectDataObj.Interface())
		if err != nil {
			return false
		}

	default:
		addData = reflectDataObj.Interface()
	}
	rel, err := redis.Int(_redisCache.redisConn.Do("HSET", _redisCache.hashMapKey, key, addData))
	if err != nil {
		return false
	}
	if rel == 0 {
		return false
	}
	return true
}

//pollUpdatePolicyFunc 巡检维护新数据
func (_redisCache *RedisCache) pollUpdatePolicyFunc() error {
	//获取巡检维护新数据
	dataList, err := redis.StringMap(_redisCache.redisConn.Do("HGETALL", _redisCache.hashMapKey))
	if err != nil {
		return err
	}
	waitAddData := _redisCache.dataInterface.PollUpdatePolicyFunc(dataList)
	err = _redisCache.BatchAddCacheData(waitAddData)
	return err
}

//GetCacheInfo 获取缓存数据
func (_redisCache *RedisCache) GetCacheInfo(key CacheKey, hasIgnorePoll bool) (interface{}, bool, error) {
	_redisCache.Lock()
	defer _redisCache.Unlock()
	//判断是否需要巡检维护新数据
	if !hasIgnorePoll {
		//巡检维护新数据
		err := _redisCache.pollUpdatePolicyFunc()
		if err != nil {
			return nil, false, err
		}
	}

	//获取缓存数据
	data, err := redis.String(_redisCache.redisConn.Do("HGET", _redisCache.hashMapKey, key))
	if err != nil {
		return nil, false, err
	}
	if _redisCache.expiredCheck(key, data) {
		return nil, false, nil
	}
	return data, true, nil
}

//GetCacheData 获取缓存数据列表
func (_redisCache *RedisCache) GetCacheData(hasIgnorePoll bool) (map[CacheKey]interface{}, error) {
	_redisCache.Lock()
	defer _redisCache.Unlock()
	//判断是否需要巡检维护新数据
	if !hasIgnorePoll {
		//巡检维护新数据
		err := _redisCache.pollUpdatePolicyFunc()
		if err != nil {
			return nil, err
		}
	}
	//获取缓存数据
	dataList, err := redis.StringMap(_redisCache.redisConn.Do("HGETALL", _redisCache.hashMapKey))
	if err != nil {
		return nil, err
	}
	newDataList := make(map[CacheKey]interface{})
	//淘汰检测
	for cacheKey, cacheData := range dataList {
		//淘汰检测 符合淘汰规则会删除当前遍历的数据
		if _redisCache.expiredCheck(CacheKey(cacheKey), cacheData) == true {
			delete(dataList, cacheKey)
		} else {
			newDataList[CacheKey(cacheKey)] = cacheData
		}
	}
	return newDataList, nil
}

//MustPollUpdatePolicyFunc 强制巡检维护新数据
func (_redisCache *RedisCache) MustPollUpdatePolicyFunc() error {
	//巡检维护新数据
	return _redisCache.pollUpdatePolicyFunc()
}

//DelCache 删除缓存数据
func (_redisCache *RedisCache) DelCache(key CacheKey) {
	_redisCache.Lock()
	defer _redisCache.Unlock()
	_redisCache.redisConn.Do("HDEL", _redisCache.hashMapKey, key)
}

//expiredCheck 过期检查
func (_redisCache *RedisCache) expiredCheck(key CacheKey, data interface{}) bool {
	if _redisCache.dataInterface.EliminatePolicyFunc(key, data) {
		_redisCache.redisConn.Do("HDEL", _redisCache.hashMapKey, key)
		return true
	}
	return false
}

package data

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/rstat1/frost/common"
	"github.com/mediocregopher/radix.v2/pool"
)

type CacheService struct {
	redisClient *pool.Pool
}

func NewCacheService() *CacheService {
	client, err := pool.New("tcp", common.CurrentConfig.RedisServerAddr, 3)
	if err != nil {
		common.CreateFailureResponse(err, "NewCacheService", 500)
	}

	return &CacheService{
		redisClient: client,
	}
}

//PutString Caches a string with the given key+field
func (c *CacheService) PutString(key, field, value string) {
	if resp := c.redisClient.Cmd("SET", key+":"+field, value); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "PutString", 500)
	}
}

//PutObject Caches an object (as a string) with the given key+field
func (c *CacheService) PutObject(key, field string, value interface{}) {
	newValue, _ := json.Marshal(value)
	if resp := c.redisClient.Cmd("SET", key+":"+field, newValue); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "PutObject", 500)
	}
}

//PutStringWithExpiration Caches a string with the given key+field that expires in expiresIn seconds.
func (c *CacheService) PutStringWithExpiration(key, field, value string, expiresIn int) {
	if resp := c.redisClient.Cmd("SET", key+":"+field, value, "EX", expiresIn); resp.Err != nil {
		common.CreateFailureResponseWithFields(resp.Err, 500, logrus.Fields{
			"func":  "PutStringWithExpiration",
			"key":   key,
			"field": field,
		})
	}
}

//AddStringToSet Adds a string to a set with the given key+field
func (c *CacheService) AddStringToSet(setname, key, value string) {
	if resp := c.redisClient.Cmd("SADD", setname+":"+key, value); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "AddStringToSet", 500)
	}
}

//AddStringToGlobalSet Add a string to a "global" set (ie one that's not specific to a user) with the given key
func (c *CacheService) AddStringToGlobalSet(key, value string) {
	if resp := c.redisClient.Cmd("SADD", key, value); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "AddStringToGlobalSet", 500)
	}
}

//SetTTLOnKey Sets a TTL (time-to-live) on the given key+fieldS
func (c *CacheService) SetTTLOnKey(key, field string, ttl int) {
	if resp := c.redisClient.Cmd("EXPIRE", key+":"+field, ttl); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "SetTTLOnKey", 500)
	}
}

//GetString Returns a cached string with the given key+field
func (c *CacheService) GetString(key, field string) string {
	if c.DoesKeyExist(key, field) == false {
		return ""
	}
	if resp := c.redisClient.Cmd("GET", key+":"+field); resp.Err != nil {
		common.CreateFailureResponseWithFields(resp.Err, 500, logrus.Fields{
			"func":  "GetString",
			"key":   key,
			"field": field,
		})
	} else {
		if value, err := resp.Str(); err == nil {
			return value
		} else {
			return ""
		}
	}
	return ""
}

//GetSet ...
func (c *CacheService) GetSet(key, field string) []string {
	if c.DoesKeyExist(key, field) == false {
		return nil
	}
	if resp := c.redisClient.Cmd("SMEMBERS", key+":"+field); resp.Err != nil {
		common.CreateFailureResponseWithFields(resp.Err, 500, logrus.Fields{
			"func":  "GetSet",
			"key":   key,
			"field": field,
		})
	} else {
		if set, err := resp.List(); err == nil {
			return set
		} else {
			common.CreateFailureResponse(err, "GetSet", 500)
			return nil
		}
	}
	return nil
}

//IsInSet ...
func (c *CacheService) IsInSet(key, field string) bool {
	if resp := c.redisClient.Cmd("SISMEMBER", key+":"+field); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "IsInSet", 500)
	} else {
		if ret, err := resp.Int(); err == nil {
			if ret == 1 {
				return true
			} else {
				return false
			}
		} else {
			common.CreateFailureResponse(resp.Err, "IsInSet", 500)
			return false
		}
	}
	return false
}

//DeleteString ...
func (c *CacheService) DeleteString(key, field string) {
	if resp := c.redisClient.Cmd("DEL", key+":"+field); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "DeleteString", 500)
	}
}

//DeleteSetItem ...
func (c *CacheService) DeleteSetItem(setname, key, member string) {
	if resp := c.redisClient.Cmd("SREM", key+":"+setname, member); resp.Err != nil {
		common.CreateFailureResponse(resp.Err, "DeleteString", 500)
	}
}

//DoesKeyExist ...
func (c *CacheService) DoesKeyExist(key, field string) bool {
	if resp := c.redisClient.Cmd("EXISTS", key+":"+field); resp.Err != nil {
		common.CreateFailureResponseWithFields(resp.Err, 500, logrus.Fields{
			"func":  "DoesKeyExist",
			"key":   key,
			"field": field,
		})
	} else {
		if ret, err := resp.Int(); err == nil {
			if ret == 1 {
				return true
			} else {
				return false
			}
		} else {
			common.CreateFailureResponse(resp.Err, "DoesKeyExist", 500)
			return false
		}
	}
	return false
}

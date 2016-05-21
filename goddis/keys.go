package goddis

import (
	"strconv"
	"time"
)

func (g *Goddis) Exists(key string) bool {
	_, existsInStrings := g.State.StringMap[key]
	_, existsInHash := g.State.HashMap[key]
	if existsInStrings || existsInHash {
		return true
	}
	return false
}

func (g *Goddis) Keys(pattern string) []string {
	var keys []string
	return keys
}

func (g *Goddis) Expire(key, seconds string) bool {
	sec, err := strconv.Atoi(seconds)
	if err == nil {
		if sec <= 0 {
			g.Del(key)
			return true
		} else {
			expireAt := time.Now().Add(time.Duration(sec) * time.Second)
			g.State.Expiry[key] = expireAt
			return true
		}
	}
	return false
}

func (g *Goddis) Del(keys ...string) int {
	var removed int
	for _, key := range keys {
		_, existsInStrings := g.State.StringMap[key]
		_, existsInHash := g.State.HashMap[key]
		if existsInStrings || existsInHash {
			delete(g.State.StringMap, key)
			delete(g.State.HashMap, key)
			delete(g.State.Expiry, key)
			removed += 1
		}
	}
	return removed
}

func (g *Goddis) Rename(key, newkey string) {

}

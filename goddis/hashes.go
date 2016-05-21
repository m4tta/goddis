package goddis

import (
	"errors"
	"strconv"
	"time"
)

// var HashMap = make(map[string]map[string]string)

func (g *Goddis) HSet(key, field, value string) (wasNew bool) {
	// ok := g.Get(key)
	// if ok {
	//   return false, errors.New("Operation against a key holding the wrong kind of value")
	// }
	_, o := g.State.HashMap[key]
	if o {
		if _, exists := g.State.HashMap[key][field]; exists {
			g.State.HashMap[key][field] = value
		} else {
			g.State.HashMap[key][field] = value
			wasNew = true
		}
	} else {
		g.State.HashMap[key] = make(map[string]string)
		g.State.HashMap[key][field] = value
		wasNew = true
	}
	delete(g.State.Expiry, key)
	return
}

func (g *Goddis) HGet(key, field string) (s string, ok bool) {
	if expireTime, exists := g.State.Expiry[key]; exists {
		if expireTime.After(time.Now()) {
			g.Del(key)
			delete(g.State.Expiry, key)
			ok = false
			return
		}
	}
	keyVal, exists := g.State.HashMap[key]
	if exists {
		val, exists := keyVal[field]
		if exists {
			s = val
			ok = true
		}
	}
	return
}

func (g *Goddis) HMGet(key string, fields ...string) []string {
	values := make([]string, 0)
	for _, f := range fields {
		val, _ := g.HGet(key, f)
		values = append(values, val)
	}
	return values
}

func (g *Goddis) HIncrBy(key, field, value string) (int, error) {
	var storedInt int
	newValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("value is not an integer or out of range")
	}
	val, exists := g.HGet(key, field)
	if exists {
		oldInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.New("hash value is not an integer")
		}
		storedInt = oldInt + newValue
		g.HSet(key, field, strconv.Itoa(storedInt))
	} else {
		storedInt = newValue
		g.HSet(key, field, strconv.Itoa(storedInt))
	}
	return storedInt, nil
}

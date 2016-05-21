package goddis

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func (g *Goddis) Set(key, value string, modifiers ...string) (ok bool) {
	var hasEx bool
	var EX int
	var hasPx bool
	var PX int
	var NX bool
	var XX bool
	for i, mod := range modifiers {
		mod = strings.ToUpper(mod)
		if mod == "NX" {
			NX = true
		} else if mod == "XX" {
			XX = true
		} else if mod == "EX" {
			hasEx = true
			if i+1 < len(modifiers) {
				EX, _ = strconv.Atoi(modifiers[i+1])
			} else {
				// missing seconds arg
				ok = false
				return
			}
		} else if mod == "PX" {
			hasPx = true
			if i+1 < len(modifiers) {
				PX, _ = strconv.Atoi(modifiers[i+1])
			} else {
				// missing milliseconds arg
				ok = false
				return
			}
		}
	}
	if (XX && NX) || (hasEx && EX == 0) || (hasPx && PX == 0) {
		return
	}
	_, exists := g.State.StringMap[key]
	if (!XX && !NX) || (NX && !exists) || (XX && exists) {
		g.State.StringMap[key] = value
		delete(g.State.Expiry, key)
		ok = true
	}
	return
}

func (g *Goddis) Get(key string) (value string, ok bool) {
	if expireTime, exists := g.State.Expiry[key]; exists {
		if expireTime.After(time.Now()) {
			g.Del(key)
			delete(g.State.Expiry, key)
		}
	} else if val, exists := g.State.StringMap[key]; exists {
		value = val
		ok = exists
	}
	return
}

func (g *Goddis) IncrBy(key, value string) (int, error) {
	var storedInt int
	newValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("value is not an integer or out of range")
	}
	val, exists := g.Get(key)
	if exists {
		oldInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.New("key value is not an integer")
		}
		storedInt = oldInt + newValue
		g.Set(key, strconv.Itoa(storedInt))
	} else {
		storedInt = newValue
		g.Set(key, strconv.Itoa(storedInt))
	}
	return storedInt, nil
}

func (g *Goddis) Incr(key string) (int, error) {
	var storedInt int
	val, exists := g.Get(key)
	if exists {
		oldInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.New("key value is not an integer")
		}
		storedInt = oldInt + 1
		g.Set(key, strconv.Itoa(storedInt))
	} else {
		storedInt = 1
		g.Set(key, strconv.Itoa(storedInt))
	}
	return storedInt, nil
}

func (g *Goddis) MGet(keys ...string) []string {
	values := make([]string, 0)
	for _, key := range keys {
		val, _ := g.Get(key)
		values = append(values, val)
	}
	return values
}

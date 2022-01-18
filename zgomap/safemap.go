package zgomap

import (
  "fmt"
  "sort"
  "sync"
)

/*
@Time : 2019-03-15 11:29
@Author : rubinus.chu
@File : safeMap
@project: zgo
*/

var Map Maper

type Maper interface {
  New() *safeMap
  Get(k interface{}) interface{}
  Set(k interface{}, v interface{}) bool
  IsExists(k interface{}) bool
  IsEmpty() bool
  Delete(k interface{})
  Size() int
  Range() chan *Sma
  Keys() []string
  Values() []string
  Join(s string) string
}

// safeMap is concurrent security map
type safeMap struct {
  lock *sync.RWMutex
  sm   map[interface{}]interface{}
}

// NewsafeMap get a new concurrent security map
func GetMap() *safeMap {
  return &safeMap{
    lock: new(sync.RWMutex),
    sm:   make(map[interface{}]interface{}),
  }
}

func (m *safeMap) New() *safeMap {
  return GetMap()
}

// Get used to get a value by key
func (m *safeMap) Get(k interface{}) interface{} {
  m.lock.RLock()
  defer m.lock.RUnlock()
  if val, ok := m.sm[k]; ok {
    return val
  }
  return nil
}

// Set used to set value with key
func (m *safeMap) Set(k interface{}, v interface{}) bool {
  m.lock.Lock()
  defer m.lock.Unlock()
  //if val, ok := m.sm[k]; !ok {
  //	m.sm[k] = v
  //} else if val != v {
  //	m.sm[k] = v
  //} else {
  //	return false
  //}
  m.sm[k] = v
  return true
}

// IsExists determine whether k exists
func (m *safeMap) IsExists(k interface{}) bool {
  m.lock.RLock()
  defer m.lock.RUnlock()
  if _, ok := m.sm[k]; !ok {
    return false
  }
  return true
}

// Delete used to delete a key
func (m *safeMap) Delete(k interface{}) {
  m.lock.Lock()
  defer m.lock.Unlock()
  delete(m.sm, k)
}

// Len长度
func (m *safeMap) Size() int {
  m.lock.RLock()
  defer m.lock.RUnlock()
  return len(m.sm)
}

// IsEmpty
func (m *safeMap) IsEmpty() bool {
  m.lock.RLock()
  defer m.lock.RUnlock()
  return len(m.sm) == 0
}

type Sma struct {
  Key interface{}
  Val interface{}
}

func (m *safeMap) Range() chan *Sma {
  //m.lock.RLock()
  //defer m.lock.RUnlock()
  out := make(chan *Sma)
  go func() {
    m.lock.RLock()
    defer m.lock.RUnlock()
    for k, v := range m.sm {
      c := &Sma{
        Key: k,
        Val: v,
      }
      out <- c
    }
    close(out)
  }()
  return out
}
func (m *safeMap) Keys() []string {
  m.lock.RLock()
  defer m.lock.RUnlock()
  keys := make([]string, 0, len(m.sm))
  for k := range m.sm {
    keys = append(keys, k.(string))
  }
  sort.Strings(keys)
  return keys
}

func (m *safeMap) Values() []string {
  m.lock.RLock()
  defer m.lock.RUnlock()
  values := make([]string, 0, len(m.sm))
  for _, v := range m.sm {
    values = append(values, v.(string))
  }
  sort.Strings(values)
  return values
}
func (m *safeMap) Join(s string) string {
  m.lock.RLock()
  defer m.lock.RUnlock()
  str := ""
  for k, v := range m.sm {
    str += fmt.Sprintf("%v=%v%v", k, v, s)
  }
  return str[:len(str)-len(s)]
}

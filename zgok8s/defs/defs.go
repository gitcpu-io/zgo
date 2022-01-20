package defs

import (
  "k8s.io/apimachinery/pkg/version"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "sync"
)

var (
  K8sContextMap  = make(map[string]*K8sInfo)
  CurrentContext string
  mutex          sync.RWMutex
)

type K8sInfo struct {
  Config    *rest.Config
  ClientSet *kubernetes.Clientset
}

type K8s struct {
  host string
}

func New() *K8s {
  return &K8s{}
}

func oneOfMap(m map[string]*K8sInfo, keys []string, curhost string) string {
  var one string
  if len(keys) > 0 {
    one = keys[0]
  } else if curhost != "" {
    one = curhost
  } else {
    for key, _ := range m {
      if key != "" {
        one = key
        break
      }
    }
  }
  return one
}

func (k *K8s) UseContext(host string) *rest.Config {
  mutex.RLock()
  defer mutex.RUnlock()
  CurrentContext = host
  return K8sContextMap[CurrentContext].Config
}

func (k *K8s) GetContext(host ...string) *rest.Config {
  mutex.RLock()
  defer mutex.RUnlock()
  one := oneOfMap(K8sContextMap, host, CurrentContext)
  return K8sContextMap[one].Config
}

func (k *K8s) SetContext(host string, config *rest.Config) {
  mutex.Lock()
  defer mutex.Unlock()
  if _, ok := K8sContextMap[host]; !ok {
    K8sContextMap[host] = &K8sInfo{}
  }
  K8sContextMap[host].Config = config
  CurrentContext = host //设置当前的host
  return
}

func (k *K8s) GetClientSet(host ...string) *kubernetes.Clientset {
  mutex.RLock()
  defer mutex.RUnlock()
  one := oneOfMap(K8sContextMap, host, CurrentContext)
  return K8sContextMap[one].ClientSet
}

func (k *K8s) SetClientSet(host string, cs *kubernetes.Clientset) {
  mutex.Lock()
  defer mutex.Unlock()
  if _, ok := K8sContextMap[host]; !ok {
    K8sContextMap[host] = &K8sInfo{}
  }
  K8sContextMap[host].ClientSet = cs
  return
}

func (k *K8s) ServerVersion(host ...string) (*version.Info, error) {
  mutex.RLock()
  defer mutex.RUnlock()
  one := oneOfMap(K8sContextMap, host, CurrentContext)
  return K8sContextMap[one].ClientSet.ServerVersion()
}

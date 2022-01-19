package zgok8s

import (
  "github.com/gitcpu-io/zgo/zgok8s/buildclient"
  "github.com/gitcpu-io/zgo/zgok8s/configoption"
)

type K8sContexter interface {
  ConfigOption() configoption.ConfigOptioner
  Client() buildclient.BuildClienter
}

type K8s struct {}

func New() K8sContexter  {
  return &K8s{}
}

func (k *K8s) ConfigOption() configoption.ConfigOptioner  {
  return configoption.New()
}

func (k *K8s) Client() buildclient.BuildClienter  {
 return buildclient.New()
}

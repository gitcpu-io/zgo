package zgok8s

import (
  "github.com/gitcpu-io/zgo/zgok8s/builder"
  "github.com/gitcpu-io/zgo/zgok8s/configoption"
  "github.com/gitcpu-io/zgo/zgok8s/defs"
  "github.com/gitcpu-io/zgo/zgok8s/deployment"
  "k8s.io/apimachinery/pkg/version"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)

type K8sContexter interface {
  UseContext(host string) *rest.Config
  GetContext(host ...string) *rest.Config
  GetClientSet(host ...string) *kubernetes.Clientset
  ServerVersion(host string) (*version.Info, error)

  // ConfigOption create config and option
  ConfigOption() configoption.ConfigOptioner
  // Builder create client and config to zgo.K8s
  Builder() builder.Builder

  // Deployment define
  Deployment() deployment.Deploymenter
}

func New() K8sContexter {
  return &K8sContext{res: defs.New()}
}

type K8sContext struct {
  res *defs.K8s
}

func (k *K8sContext) UseContext(host string) *rest.Config {
  return k.res.UseContext(host)
}

func (k *K8sContext) GetContext(host ...string) *rest.Config {
  return k.res.GetContext(host...)
}

func (k *K8sContext) GetClientSet(host ...string) *kubernetes.Clientset {
  return k.res.GetClientSet(host...)
}

func (k *K8sContext) ServerVersion(host string) (*version.Info, error) {
  return k.res.ServerVersion(host)
}

func (k *K8sContext) ConfigOption() configoption.ConfigOptioner {
  return configoption.New()
}

// Builder 创建config 和 clientset
func (k *K8sContext) Builder() builder.Builder {
  return builder.New()
}

// Deployment 开始build其它接口
func (k *K8sContext) Deployment() deployment.Deploymenter {
  return deployment.New()
}

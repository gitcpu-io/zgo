package builder

import (
  "fmt"
  "github.com/gitcpu-io/zgo/zgok8s/configoption"
  "github.com/gitcpu-io/zgo/zgok8s/defs"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/klog"
)

// Builder is interface
type Builder interface {
  BuildConfig(configOption *configoption.ConfigOption) (config *rest.Config, err error)
  BuildClientSet(host string, config *rest.Config) (clientSet *kubernetes.Clientset, err error)
}

type BuildClient struct {
  K8s *defs.K8s
}

func New() Builder {
  return &BuildClient{
    K8s: defs.New(),
  }
}

func (bc *BuildClient) BuildConfig(co *configoption.ConfigOption) (config *rest.Config, err error) {
  if co == nil || (co.MasterUrl == "" && co.KubeConfig == "") {
    //use cluster
    config, err = rest.InClusterConfig()
  } else {
    //use builder master url and kube config
    config, err = clientcmd.BuildConfigFromFlags(co.MasterUrl, co.KubeConfig)
  }
  if err != nil {
    errStr := fmt.Sprintf("Fail to load k8s config option: %v,\nerr: %v", co, err)
    klog.Error(errStr)
    return
  }
  //config.NegotiatedSerializer = scheme.Codecs
  bc.K8s.SetContext(config.Host, config)
  return
}

func (bc *BuildClient) BuildClientSet(host string, config *rest.Config, ) (clientSet *kubernetes.Clientset, err error) {
  if clientSet, err = kubernetes.NewForConfig(config); err != nil {
    klog.Errorf("Fail to create k8s client config: %v,\nerr: %v", config, err)
    return
  }
  bc.K8s.SetClientSet(host, clientSet)
  return
}

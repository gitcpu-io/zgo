package buildclient

import (
  "fmt"
  "github.com/gitcpu-io/zgo/zgok8s/configoption"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/kubernetes/scheme"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/klog"
)

type BuildClienter interface {
  BuildConfig(configOption *configoption.ConfigOption) (config *rest.Config,err error)
  ClientSet(config *rest.Config) (clientSet *kubernetes.Clientset, err error)
}

type buildClient struct {
}

func New() BuildClienter {
  return &buildClient{}
}

func (bc *buildClient) BuildConfig(co *configoption.ConfigOption) (config *rest.Config,err error) {
  if co.MasterUrl == "" && co.KubeConfig == "" {
    //use cluster
    config, err = rest.InClusterConfig()
  } else {
    //use build master url and kube config
    config, err = clientcmd.BuildConfigFromFlags(co.MasterUrl, co.KubeConfig)
  }
  if err != nil {
    errStr := fmt.Sprintf("Fail to load k8s config option: %v,\nerr: %v", co, err)
    klog.Error(errStr)
    return
  }
  config.NegotiatedSerializer = scheme.Codecs
  return
}

func (bc *buildClient) ClientSet(config *rest.Config,) (clientSet *kubernetes.Clientset, err error) {
  if clientSet, err = kubernetes.NewForConfig(config); err != nil {
    klog.Errorf("Fail to create k8s client config: %v,\nerr: %v", config, err)
    return
  }
  return
}

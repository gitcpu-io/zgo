package configoption

import (
  "errors"
)

type configOptionFunc func(kco *ConfigOption) error

func (kcf configOptionFunc) WithMasterUrl(masterUrl string) configOptionFunc {
  return func(kco *ConfigOption) error {
    if masterUrl == "" {
      return errors.New("master url must be have")
    }
    kco.MasterUrl = masterUrl
    return nil
  }
}

func (kcf configOptionFunc) WithKubeConfig(kubeConfigFile string) configOptionFunc {
  return func(kco *ConfigOption) error {
    if kubeConfigFile == "" {
      return errors.New("kube config file must be have")
    }
    kco.KubeConfig = kubeConfigFile
    return nil
  }
}

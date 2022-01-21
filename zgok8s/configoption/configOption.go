package configoption


type ConfigOptioner interface {
  // GetFunc 取configOption 的func函数
  GetFunc() configOptionFunc
  // Build 通过func函数为 config option赋值
  Build(opts ...configOptionFunc) (k8sco *ConfigOption,err error)
}

type ConfigOption struct {
  MasterUrl string
  KubeConfig string
}

func New() ConfigOptioner {
  return &ConfigOption{}
}


func (kc *ConfigOption) GetFunc() configOptionFunc {
  return func(kco *ConfigOption) error {
    return nil
  }
}

func (kc *ConfigOption) Build(opts ...configOptionFunc) (*ConfigOption,error) {
  k8sco := &ConfigOption{}
  for _, opt := range opts {
    err := opt(k8sco)
    if err != nil {
      return nil,err
    }
  }
  return k8sco,nil
}

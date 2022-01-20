package deployment

import (
  "context"
  "github.com/gitcpu-io/zgo/zgok8s/defs"
  appv1 "k8s.io/api/apps/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deploymenter interface {
  List() (*appv1.DeploymentList, error)
}

type deployment struct {
  K8s *defs.K8s
}

func New() Deploymenter {
  return &deployment{
    K8s: defs.New(),
  }
}

func (d *deployment) List() (*appv1.DeploymentList, error) {
  return d.K8s.GetClientSet().AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
}

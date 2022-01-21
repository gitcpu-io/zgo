package deployment

import (
  "context"
  "github.com/gitcpu-io/zgo/zgok8s/defs"
  appv1 "k8s.io/api/apps/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deploymenter interface {
  List(ctx context.Context, ns, ls, fs string, limit int64, watch bool) (*appv1.DeploymentList, error)
}

type deployment struct {
  K8s *defs.K8s
}

func New() Deploymenter {
  return &deployment{
    K8s: defs.New(),
  }
}

func (d *deployment) List(ctx context.Context, ns, ls, fs string, limit int64, watch bool) (*appv1.DeploymentList, error) {
  opts := metav1.ListOptions{
    LabelSelector: ls,
    FieldSelector: fs,
    Watch:         watch,
    Limit:         limit,
  }
  return d.K8s.GetClientSet().AppsV1().Deployments(ns).List(ctx, opts)
}

package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// List installed Helm charts
func (h *Helm) List() ([]*release.Release, error) {
	client := action.NewList(h.ActionConfig)
	return client.Run()
}

package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

// Install a specific Helm chart
func (h *Helm) Install(chartPath string) (*release.Release, error) {
	client := action.NewInstall(h.ActionConfig)
	client.Namespace = h.Namespace
	client.CreateNamespace = false

	validatedChartPath, err := client.ChartPathOptions.LocateChart(chartPath, CliSettings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(validatedChartPath)
	if err != nil {
		return nil, err
	}

	client.ReleaseName = chart.Name()

	rls, err := client.Run(chart, nil)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

// InstallNamed a named helm chart
func (h *Helm) InstallNamed(releaseName, chartPath string) (*release.Release, error) {
	client := action.NewInstall(h.ActionConfig)
	client.ReleaseName = releaseName
	client.Namespace = h.Namespace
	client.CreateNamespace = false

	validatedChartPath, err := client.ChartPathOptions.LocateChart(chartPath, CliSettings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(validatedChartPath)
	if err != nil {
		return nil, err
	}

	rls, err := client.Run(chart, nil)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

// InstallWithValues a specific Helm chart with value overwrite
func (h *Helm) InstallWithValues(chartPath string, values map[string]interface{}) (*release.Release, error) {
	client := action.NewInstall(h.ActionConfig)
	client.Namespace = h.Namespace
	client.CreateNamespace = false

	validatedChartPath, err := client.ChartPathOptions.LocateChart(chartPath, CliSettings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(validatedChartPath)
	if err != nil {
		return nil, err
	}

	client.ReleaseName = chart.Name()

	rls, err := client.Run(chart, values)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

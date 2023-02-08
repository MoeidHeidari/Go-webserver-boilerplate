package kubes

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

// Gets release in Helm
func (u KubeRequest) HGetRelease() ([]*release.Release, error) {
	u.ActionConfiguration.Init(u.Settings.RESTClientGetter(), u.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewList(u.ActionConfiguration)
	results, err := client.Run()
	if err != nil || results == nil {
		return nil, err
	}
	return results, nil
}

func (u KubeRequest) HCreateRelease(chartBody ChartBody) (*release.Release, error) {

	u.ActionConfiguration.Init(u.Settings.RESTClientGetter(), u.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewInstall(u.ActionConfiguration)
	client.Namespace = chartBody.Namespace
	client.ReleaseName = chartBody.ReleaseName

	locatedChart, err := client.LocateChart(chartBody.ChartPath, u.Settings)
	if err != nil {
		return nil, err
	}

	newChart, err := loader.Load(locatedChart)
	if err != nil {
		return nil, err
	}

	release, err := client.Run(newChart, nil)
	if err != nil {
		return nil, err
	}

	return release, nil
}

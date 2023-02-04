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
	if err := u.actionConfiguration.Init(u.settings.RESTClientGetter(), u.settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	client := action.NewList(u.actionConfiguration)
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		return nil, err
	}
	return results, nil

}

func (u KubeRequest) HCreateRelease(chartBody ChartBody) (*release.Release, error) {

	if err := u.actionConfiguration.Init(u.settings.RESTClientGetter(), u.settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	client := action.NewInstall(u.actionConfiguration)
	client.Namespace = chartBody.Namespace
	client.ReleaseName = chartBody.ReleaseName

	locatedChart, err := client.LocateChart(chartBody.ChartPath, u.settings)
	if err != nil {
		return nil, err
	}

	newChart, err := loader.Load(locatedChart)
	if err != nil {
		return nil, err
	}

	release, err := client.Run(newChart, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return release, nil
}

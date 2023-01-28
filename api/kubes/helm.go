package kubes

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func helmInit() error {

	return nil
}

// Gets release in Helm
func (u KubeRequest) HGetRelease(release_name string) (release.Release, error) {
	panic("Not implemented exception")
}

func (u KubeRequest) HCreateRelease(chartBody ChartBody) (*release.Release, error) {

	settings := cli.New()
	actionConfiguration := new(action.Configuration)

	if err := actionConfiguration.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	client := action.NewInstall(actionConfiguration)
	client.Namespace = chartBody.Namespace
	client.ReleaseName = chartBody.ReleaseName

	locatedChart, err := client.LocateChart(chartBody.ChartPath, settings)
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

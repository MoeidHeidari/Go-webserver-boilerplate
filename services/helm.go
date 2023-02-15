package services

import (
	"context"
	"io/ioutil"
	"log"
	"main/models"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

// Gets release in Helm
func (u KubernetesService) HGetRelease() ([]*release.Release, error) {
	u.Repository.ActionConfiguration.Init(u.Repository.Settings.RESTClientGetter(), u.Repository.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewList(u.Repository.ActionConfiguration)
	results, err := client.Run()
	if err != nil || results == nil {
		return nil, err
	}
	return results, nil
}

func (u KubernetesService) HCreateRelease(chartBody models.ChartBody) (*release.Release, error) {

	u.Repository.ActionConfiguration.Init(u.Repository.Settings.RESTClientGetter(), u.Repository.Settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	client := action.NewInstall(u.Repository.ActionConfiguration)
	client.Namespace = chartBody.Namespace
	client.ReleaseName = chartBody.ReleaseName

	locatedChart, err := client.LocateChart(chartBody.ChartPath, u.Repository.Settings)
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

func (u KubernetesService) HelmRepoAdd(body models.RepositoryBody) error {
	name := body.Name
	url := body.Url
	repoFile := u.Repository.Settings.RepositoryConfig
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {
		return errors.New("Repo already exists")
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(u.Repository.Settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart Repository or cannot be reached", url)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

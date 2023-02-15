package lib

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	Clientset           *kubernetes.Clientset
	Settings            *cli.EnvSettings
	ActionConfiguration *action.Configuration
	logger              Logger
}

func NewKubernetesClient(logger Logger) KubernetesClient {
	Settings := cli.New()
	ActionConfiguration := new(action.Configuration)
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()

	if err != nil {
		logger.Error(err)
	}

	Clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err)
	}

	return KubernetesClient{
		logger:              logger,
		Clientset:           Clientset,
		Settings:            Settings,
		ActionConfiguration: ActionConfiguration,
	}
}

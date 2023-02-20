package repository

import "main/lib"

type KubernetesRepository struct {
	lib.KubernetesClient
	logger lib.Logger
}

func NewKubernetesRepository(kubeclient lib.KubernetesClient, logger lib.Logger) KubernetesRepository {
	return KubernetesRepository{
		KubernetesClient: kubeclient,
		logger:           logger,
	}
}

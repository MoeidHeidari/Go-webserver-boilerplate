package kubes

type ChartBody struct {
	ChartPath   string `json:"chart_path"`
	Namespace   string `json:"namespace"`
	ReleaseName string `json:"release_name"`
}

type RepositoryBody struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type PodBody struct {
	Name          string   `json:"name"`
	Namespace     string   `json:"namespace"`
	ContainerName string   `json:"container_name"`
	Image         string   `json:"image"`
	Command       []string `json:"command"`
	ConfigmapName string   `json:"configmap_name"`
	SecretName    string   `json:"secret_name"`
	Port          int32    `json:"port"`
}

type ConfigMapBody struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"env"`
}

type SecretBody struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"env"`
}

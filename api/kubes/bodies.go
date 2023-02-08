package kubes

type ChartBody struct {
	ChartPath   string `json:"chart_path"`
	Namespace   string `json:"namespace"`
	ReleaseName string `json:"release_name"`
	Reponame    string `json:"repo_name"`
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
	ClaimName     string   `json:"claim_name"`
	VolumeName    string   `json:"volume_name"`
	MountPath     string   `json:"mountpath"`
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

type PV struct {
	Name    string `json:"name"`
	Storage string `json:"storage"`
	Path    string `json:"path"`
}

type PVC struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Storage   string `json:"storage"`
}

type Nodeport struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Port         int32  `json:"port"`
	RedirectPort int32  `json:"redirect_port"`
}

package config

type Config struct {
	Provider      string
	Kubeconfig    string
	Namespace     string
	LabelSelector string
	FieldSelector string
	HostNetwork   bool
}

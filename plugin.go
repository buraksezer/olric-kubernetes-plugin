package main

import "github.com/buraksezer/olric-kubernetes-plugin/lib"

// ServiceDiscovery defines a service discovery plugin
// for Olric, backed by Kubernetes.
var ServiceDiscovery lib.KubernetesDiscovery

func main() {
	_ = ServiceDiscovery
}

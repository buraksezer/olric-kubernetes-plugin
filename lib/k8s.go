package lib

import (
	"fmt"
	"log"

	"github.com/buraksezer/olric-kubernetes-plugin/config"
	"github.com/buraksezer/olric-kubernetes-plugin/internal/k8s"
	"github.com/mitchellh/mapstructure"
)

type KubernetesDiscovery struct {
	config *config.Config
	log    *log.Logger
	k8s    *k8s.K8S
}

func (k *KubernetesDiscovery) checkErrors() error {
	if k.config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if k.log == nil {
		return fmt.Errorf("logger cannot be nil")
	}
	return nil
}

func (k *KubernetesDiscovery) Initialize() error {
	if err := k.checkErrors(); err != nil {
		return err
	}

	k.k8s = &k8s.K8S{}
	k.log.Printf("[INFO] Service discovery plugin is enabled, provider: %s", k.config.Provider)
	return nil
}

func (k *KubernetesDiscovery) SetLogger(l *log.Logger) {
	k.log = l
}

func (k *KubernetesDiscovery) SetConfig(cfg map[string]interface{}) error {
	var cg config.Config
	err := mapstructure.Decode(cfg, &cg)
	if err != nil {
		return err
	}
	k.config = &cg
	return nil
}

func (k *KubernetesDiscovery) DiscoverPeers() ([]string, error) {
	peers, err := k.k8s.Addresses(k.config, k.log)
	if err != nil {
		return nil, err
	}
	if len(peers) == 0 {
		return nil, fmt.Errorf("no peer found")
	}
	return peers, nil
}

func (k *KubernetesDiscovery) Register() error { return nil }

func (k *KubernetesDiscovery) Deregister() error { return nil }

func (k *KubernetesDiscovery) Close() error { return nil }

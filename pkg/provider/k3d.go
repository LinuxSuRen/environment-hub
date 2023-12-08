/*
Copyright 2023 Rick.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"os"
	"strconv"

	"fmt"

	"github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/config"
	k3dtypes "github.com/k3d-io/k3d/v5/pkg/config/types"
	k3dconf "github.com/k3d-io/k3d/v5/pkg/config/v1alpha4"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	dockerutils "github.com/k3d-io/k3d/v5/pkg/runtimes/docker"
	"github.com/k3d-io/k3d/v5/pkg/types"
	k3d "github.com/k3d-io/k3d/v5/pkg/types"
	k3dversion "github.com/k3d-io/k3d/v5/version"
	"github.com/linuxsuren/environment-hub/pkg/model"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type k3dCluster struct {
}

func (c *k3dCluster) Create(cluster model.Cluster) (err error) {
	cfg := k3dconf.SimpleConfig{
		TypeMeta: k3dtypes.TypeMeta{
			APIVersion: config.DefaultConfigApiVersion,
			Kind:       "Simple",
		},
		ObjectMeta: k3dtypes.ObjectMeta{
			Name: cluster.Name,
		},
		Servers: cluster.Servers,
		Agents:  cluster.Agents,
		// ClusterToken: p.Token,
		Image: fmt.Sprintf("%s:%s", k3d.DefaultK3sImageRepo, k3dversion.K3sVersion),
		// ExposeAPI: k3dconf.SimpleExposureOpts{
		// 	HostIP:   ipPorts[0],
		// 	HostPort: ipPorts[1],
		// },
		Options: k3dconf.SimpleConfigOptions{
			Runtime:    k3dconf.SimpleConfigOptionsRuntime{},
			K3dOptions: k3dconf.SimpleConfigOptionsK3d{
				// DisableImageVolume:  p.NoImageVolume,
				// DisableLoadbalancer: p.NoLB,
			},
			K3sOptions: k3dconf.SimpleConfigOptionsK3s{},
		},
	}
	if cluster.Port > 0 {
		cfg.Ports = []k3dconf.PortWithNodeFilters{{
			Port: strconv.Itoa(cluster.Port),
			NodeFilters: []string{
				"server:*",
			},
		}}
	}

	cc, err := config.TransformSimpleToClusterConfig(context.Background(), runtimes.SelectedRuntime, cfg)
	if err == nil {
		err = client.ClusterRun(context.Background(), runtimes.SelectedRuntime, cc)
	}
	return
}

func (c *k3dCluster) List() (clusters []model.Cluster, err error) {
	var rawClusters []*k3d.Cluster
	if rawClusters, err = client.ClusterList(context.Background(), runtimes.SelectedRuntime); err == nil {
		for _, item := range rawClusters {

			cfg := &k3d.Cluster{
				Name: item.Name,
			}
			var cluster *k3d.Cluster
			if cluster, err = client.ClusterGet(context.Background(), runtimes.SelectedRuntime, cfg); err == nil {
				clusters = append(clusters, model.Cluster{
					Name:        cluster.Name,
					PortBinding: getMappingPort(cluster),
					Nodes:       convertToNode(cluster.Nodes),
				})
			}
		}
	}
	return
}

func convertToNode(items []*types.Node) (nodes []model.Node) {
	for _, item := range items {
		nodes = append(nodes, model.Node{
			Name:   item.Name,
			Role:   string(item.Role),
			Status: item.State.Status,
		})
	}
	return
}

func (c *k3dCluster) Delete(name string) error {
	cfg := &k3d.Cluster{
		Name: name,
	}

	return client.ClusterDelete(context.Background(), runtimes.SelectedRuntime, cfg, k3d.ClusterDeleteOpts{})
}

func getMappingPort(cluster *k3d.Cluster) (portMap map[string]string) {
	if cluster.ServerLoadBalancer == nil || cluster.ServerLoadBalancer.Node == nil {
		return
	}
	node := cluster.ServerLoadBalancer.Node
	portMap = make(map[string]string)

	docker, err := dockerutils.GetDockerClient()
	if err == nil {
		if container, err := docker.ContainerInspect(context.Background(), node.Name); err == nil {
			for port, binds := range container.NetworkSettings.Ports {
				for _, bind := range binds {
					portMap[port.Port()] = bind.HostPort
					break
				}
			}
		}
	}
	return
}

func (c *k3dCluster) Get(name string) (result model.Cluster, err error) {
	cfg := &k3d.Cluster{
		Name: name,
	}

	ctx := context.Background()
	var cluster *k3d.Cluster
	if cluster, err = client.ClusterGet(ctx, runtimes.SelectedRuntime, cfg); err == nil {
		result = model.Cluster{
			Name:        cluster.Name,
			PortBinding: getMappingPort(cluster),
			Nodes:       convertToNode(cluster.Nodes),
		}

		var configFile *os.File
		if configFile, err = os.CreateTemp(os.TempDir(), "kubeconfig"); err == nil {
			defer os.Remove(configFile.Name())

			var config *clientcmdapi.Config
			if config, err = client.KubeconfigGet(ctx, runtimes.SelectedRuntime, cluster); err == nil {
				if err = client.KubeconfigWrite(ctx, config, configFile.Name()); err == nil {
					var data []byte
					if data, err = os.ReadFile(configFile.Name()); err == nil {
						result.KubeConfig = string(data)
					}
				}
			}
		}

	}
	return
}

func (c *k3dCluster) Update(cluster model.Cluster) error {
	ctx := context.Background()
	cfg := &k3d.Cluster{
		Name: cluster.Name,
	}

	simpleConfig := &k3dconf.SimpleConfig{
		Servers: cluster.Servers,
		Agents:  cluster.Agents,
	}
	return client.ClusterEditChangesetSimple(ctx, runtimes.SelectedRuntime, cfg, simpleConfig)
}

func (c *k3dCluster) Start(name string) error {
	cfg := &k3d.Cluster{
		Name: name,
	}

	return client.ClusterStart(context.Background(), runtimes.SelectedRuntime, cfg, k3d.ClusterStartOpts{})
}

func (c *k3dCluster) Stop(name string) error {
	cfg := &k3d.Cluster{
		Name: name,
	}

	return client.ClusterStop(context.Background(), runtimes.SelectedRuntime, cfg)
}

func init() {
	Register("k3d", &k3dCluster{})
}

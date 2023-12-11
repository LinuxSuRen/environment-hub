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
	"errors"

	"github.com/linuxsuren/environment-hub/pkg/model"
	kCluster "sigs.k8s.io/kind/pkg/cluster"
)

type kindCluster struct {
	serverAddress string
}

func (c *kindCluster) Create(cluster model.Cluster) (err error) {
	opt := kCluster.ProviderWithDocker()
	p := kCluster.NewProvider(opt)
	err = p.Create(cluster.Name)
	return
}

func (c *kindCluster) List() (clusters []model.Cluster, err error) {
	opt := kCluster.ProviderWithDocker()
	p := kCluster.NewProvider(opt)

	var items []string
	if items, err = p.List(); err == nil {
		for _, item := range items {
			clusters = append(clusters, model.Cluster{
				Name: item,
			})
		}
	}
	return
}

func (c *kindCluster) Delete(name string) error {
	opt := kCluster.ProviderWithDocker()
	p := kCluster.NewProvider(opt)

	return p.Delete(name, "")
}

func (c *kindCluster) Get(name string) (result model.Cluster, err error) {
	opt := kCluster.ProviderWithDocker()
	p := kCluster.NewProvider(opt)

	result.Name = name

	kubeconfig, _ := p.KubeConfig(name, false)
	result.KubeConfig = kubeconfig

	nodes, _ := p.ListNodes(name)
	for _, node := range nodes {
		role, _ := node.Role()

		result.Nodes = append(result.Nodes, model.Node{
			Name: node.String(),
			Role: role,
		})
	}
	return
}

func (c *kindCluster) Update(model.Cluster) error {
	return errors.New("not support yet")
}

func (c *kindCluster) Start(name string) error {
	return errors.New("not support yet")
}

func (c *kindCluster) Stop(name string) error {
	return errors.New("not support yet")
}

func (c *kindCluster) WithServerAddress(serverAddress string) {
	c.serverAddress = serverAddress
}

func init() {
	Register("kind", &kindCluster{})
}

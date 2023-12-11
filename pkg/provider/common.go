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

import "github.com/linuxsuren/environment-hub/pkg/model"

type Provider interface {
	Create(model.Cluster) error
	List() ([]model.Cluster, error)
	Delete(string) error
	Get(string) (model.Cluster, error)
	Update(model.Cluster) error
	Start(string) error
	Stop(string) error

	// TODO having a better way to pass the server address
	WithServerAddress(string)
}

var providers = make(map[string]Provider, 0)

func Register(kind string, provider Provider) {
	providers[kind] = provider
}

func GetProvider(kind string) Provider {
	return providers[kind]
}

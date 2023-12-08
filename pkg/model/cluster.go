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

package model

import "time"

type Cluster struct {
	ID             string
	Name           string `json:"name"`
	Description    string
	Servers        int    `json:"servers"`
	Agents         int    `json:"agents"`
	Nodes          []Node `json:"nodes"`
	Port           int
	PortBinding    map[string]string `json:"portBinding"`
	KubeConfig     string            `json:"kubeconfig"`
	Kind           string
	K8sVersion     string
	CreatedTime    time.Time
	LastUpdateTime time.Time
}

type Node struct {
	Name   string `json:"name"`
	Role string `json:"role"`
	Status string `json:"status"`
}

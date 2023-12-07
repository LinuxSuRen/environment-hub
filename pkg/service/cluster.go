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

package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/environment-hub/pkg/model"
	"github.com/linuxsuren/environment-hub/pkg/provider"
)

func listCluster(c *gin.Context) {
	kind := c.Param("kind")
	pro := provider.GetProvider(kind)
	if clusters, err := pro.List(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, clusters)
	}
}

func createCluster(c *gin.Context) {
	var cluster model.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kind := c.Param("kind")
	pro := provider.GetProvider(kind)
	if err := pro.Create(cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func deleteCluster(c *gin.Context) {
	kind := c.Param("kind")
	name := c.Param("name")
	pro := provider.GetProvider(kind)
	if err := pro.Delete(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func getCluster(c *gin.Context) {
	kind := c.Param("kind")
	name := c.Param("name")
	pro := provider.GetProvider(kind)
	if cluster, err := pro.Get(name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}

func init() {
	RegisterHandler(listCluster, "/v1/:kind/clusters", http.MethodGet)
	RegisterHandler(createCluster, "/v1/:kind/clusters", http.MethodPost)
	RegisterHandler(deleteCluster, "/v1/:kind/clusters/:name", http.MethodDelete)
	RegisterHandler(getCluster, "/v1/:kind/clusters/:name", http.MethodGet)
}

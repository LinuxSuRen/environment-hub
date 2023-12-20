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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	helmLoader "helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

type helmChart struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	RepoURL   string   `json:"repoURL"`
	Version   string   `json:"version"`
	Values    []string `json:"values"`
}

func helmInstall(c *gin.Context) {
	var chart helmChart
	if err := c.ShouldBindJSON(&chart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if chart.Namespace == "" {
		chart.Namespace = "default"
	}

	cluster, err := getClusterByName(c, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write kubeconfig to a temporary file
	kubeconfigFile, err := os.CreateTemp(os.TempDir(), "kubeconfig")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to create temporary file": err.Error()})
		return
	} else {
		if err = os.WriteFile(kubeconfigFile.Name(), []byte(cluster.KubeConfig), 0600); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"failed to create kubeconfig file": err.Error()})
			return
		}
		defer os.RemoveAll(kubeconfigFile.Name())
	}

	actionConfig := new(action.Configuration)
	var settings = cli.New()
	helmDriver := os.Getenv("HELM_DRIVER")
	settings.KubeConfig = kubeconfigFile.Name()

	if err = actionConfig.Init(settings.RESTClientGetter(), chart.Namespace, helmDriver, func(format string, v ...interface{}) {
		log.Printf(format, v)
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valueOpts := &values.Options{}
	valueOpts.Values = chart.Values

	client := action.NewInstall(actionConfig)
	client.CreateNamespace = true
	client.Namespace = chart.Namespace
	client.RepoURL = chart.RepoURL
	client.Version = chart.Version
	client.ReleaseName = chart.Name
	client.Timeout = time.Minute * 2
	//client.InsecureSkipTLSverify = true

	chartPath, err := client.ChartPathOptions.LocateChart(chart.Name, settings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to load chart": err.Error()})
		return
	}

	v, err := valueOpts.MergeValues(getter.All(settings))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ch, err := helmLoader.Load(chartPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to load chart": err.Error()})
		return
	}

	_, err = client.RunWithContext(c.Request.Context(), ch, v)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to install chart": err.Error()})
	} else {
		c.String(http.StatusOK, "ok")
	}
}

func init() {
	RegisterHandler(helmInstall, "/v1/:kind/clusters/:name/install", http.MethodPost)
}

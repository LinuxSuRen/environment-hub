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

package cmd

import (
	_ "embed"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/environment-hub/pkg/service"
	"github.com/spf13/cobra"
)

type serverOption struct {
	address string
}

func newServerCmd() (c *cobra.Command) {
	opt := &serverOption{}
	c = &cobra.Command{
		Use:  "server",
		RunE: opt.runE,
	}
	flags := c.Flags()
	flags.StringVarP(&opt.address, "address", "", "0.0.0.0:8080", "The listen address")
	return
}

func (o *serverOption) runE(c *cobra.Command, args []string) (err error) {
	r := gin.Default()

	handlers := service.GetAllHandlers()
	for _, handler := range handlers {
		switch method := handler.Method; method {
		case http.MethodGet:
			r.GET(handler.Path, handler.Process)
		case http.MethodPost:
			r.POST(handler.Path, handler.Process)
		case http.MethodPut:
			r.PUT(handler.Path, handler.Process)
		case http.MethodDelete:
			r.DELETE(handler.Path, handler.Process)
		default:
			log.Println("unknown handler method:", method, "path:", handler.Path)
		}
	}

	r.GET("", func(c *gin.Context) {
		c.Header("content-type", "text/html;charset=utf-8")
		c.String(http.StatusOK, staticIndex)
	})
	r.GET("/assets/:file", func(c *gin.Context) {
		file := c.Param("file")

		if strings.HasSuffix(file, ".js") {
			c.Header("content-type", "application/javascript")
			c.String(http.StatusOK, staticJs)
		} else if strings.HasSuffix(file, ".css") {
			c.Header("content-type", "text/css; charset=utf-8")
			c.String(http.StatusOK, staticCss)
		}
	})
	r.Run(o.address)
	return
}

//go:embed data/index.html
var staticIndex string

//go:embed data/assets/*.js
var staticJs string

//go:embed data/assets/*.css
var staticCss string

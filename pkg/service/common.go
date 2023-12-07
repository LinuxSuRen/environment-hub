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

import "github.com/gin-gonic/gin"

type Handler struct {
	Process func(c *gin.Context)
	Path    string
	Method  string
}

var handlers []Handler

func GetAllHandlers() []Handler {
	return handlers
}

func RegisterHandler(process func(c *gin.Context), path, method string) {
	handlers = append(handlers, Handler{
		Process: process,
		Path:    path,
		Method:  method,
	})
}

/*
 * Copyright 2019 VMware, Inc..
 * SPDX-License-Identifier: Apache-2.0
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	//"github.com/labstack/gommon/log"
	"log"
	//"net/http"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/server"
)

func main() {
	confDirStr := flag.String("confDir", "", "Configuration directory")
	apiPortStr := flag.String("apiPort", "1323", "REST API port")
	flag.Parse()
	if *confDirStr == "" {
		log.Println("confDir is not defined")
		flag.Usage()
		return
	}
	apiPort, err := strconv.Atoi(*apiPortStr)
	if err != nil {
		fmt.Errorf("apiPort %s is not an integer", *apiPortStr)
		os.Exit(1)
	}
	arachneRestAPI := server.NewArachne(*confDirStr, apiPort)
	e := echo.New()
	err = arachneRestAPI.ConnectArachneAPIToEcho(e)
	if err != nil {
		e.Logger.Fatal(err)
	}
	err = arachneRestAPI.ConnectMiniS3ToEcho(e)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":" + *apiPortStr))
}


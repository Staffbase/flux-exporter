/*
Copyright 2020, Staffbase GmbH and contributors.

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

package exporter

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Staffbase/flux-exporter/pkg/api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

const namespace = "flux_exp"

var (
	imageMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "image",
	}, []string{"name", "namespace", "current", "new"})
)

type Exporter struct {
	Endpoint string
}

func New() (*Exporter, error) {
	if os.Getenv("ENDPOINT") == "" {
		return nil, errors.New("ENDPOINT is missing")
	}

	return &Exporter{
		Endpoint: os.Getenv("ENDPOINT"),
	}, nil
}

func Run(interval int64, exporter *Exporter) {
	for {
		images, err := api.GetImages(exporter.Endpoint)
		if err != nil {
			log.WithError(err).Error("Could not get images")
		}

		imageMetric.Reset()

		for _, image := range images {
			log.Debugf("Process image: %#v", image)
			for _, container := range image.Containers {
				log.Debugf("Process container: %#v", container)
				var newImageAvailable float64

				namespace := "-"
				currentVersion := "-"
				newVersion := "-"

				if split := strings.Split(image.ID, ":"); len(split) == 2 {
					namespace = split[0]
				}

				if split := strings.Split(container.Current.ID, ":"); len(split) == 2 {
					currentVersion = split[1]
				}

				if split := strings.Split(container.LatestFiltered.ID, ":"); len(split) == 2 {
					newVersion = split[1]
				}

				if currentVersion != newVersion {
					newImageAvailable = 1
				}

				imageMetric.With(prometheus.Labels{
					"name":      container.Name,
					"namespace": namespace,
					"current":   currentVersion,
					"new":       newVersion,
				}).Set(newImageAvailable)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

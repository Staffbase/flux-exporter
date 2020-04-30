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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Staffbase/flux-exporter/pkg/exporter"
	"github.com/Staffbase/flux-exporter/pkg/version"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	interval      int64
	listenAddress string
	logLevel      string
	logOutput     string
	metricsPath   string
)

var rootCmd = &cobra.Command{
	Use:   "Flux Exporter",
	Short: "Flux Exporter - export Prometheus metrics for available updates via Flux.",
	Long:  "Flux Exporter - export Prometheus metrics for available updates via Flux.",
	Run: func(cmd *cobra.Command, args []string) {
		if logOutput == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}

		log.SetReportCaller(true)
		lvl, err := log.ParseLevel(logLevel)
		if err != nil {
			log.WithError(err).Fatal("Could not set log level")
		}
		log.SetLevel(lvl)

		log.Infof(version.Info())
		log.Infof(version.BuildContext())

		exp, err := exporter.New()
		if err != nil {
			log.WithError(err).Fatal("Could not create exporter")
		}

		go exporter.Run(interval, exp)

		router := chi.NewRouter()
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Flux Exporter</title></head>
			<body>
			<h1>Flux Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			<p>
			<ul>
			<li>version: ` + version.Version + `</li>
			<li>branch: ` + version.Branch + `</li>
			<li>revision: ` + version.Revision + `</li>
			<li>go version: ` + version.GoVersion + `</li>
			<li>build user: ` + version.BuildUser + `</li>
			<li>build date: ` + version.BuildDate + `</li>
			</ul>
			</p>
			</body>
			</html>`))
		})
		router.Mount(metricsPath, promhttp.Handler())

		server := &http.Server{
			Addr:    listenAddress,
			Handler: router,
		}

		// Listen for SIGINT and SIGTERM signals and try to gracefully shutdown
		// the HTTP server. This ensures that enabled connections are not
		// interrupted.
		go func() {
			term := make(chan os.Signal, 1)
			signal.Notify(term, os.Interrupt, syscall.SIGTERM)

			<-term
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := server.Shutdown(ctx)
			if err != nil {
				log.WithError(err).Fatalf("Failed to shutdown Flux Exporter gracefully")
			}

			log.Infof("Shutdown Flux Exporter...")
			os.Exit(0)
		}()

		log.Infof("Server listen on: %s", listenAddress)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Fatal("HTTP server died unexpected")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for Flux Exporter.",
	Long:  "Print version information for Flux Exporter.",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := version.Print("Flux Exporter")
		if err != nil {
			log.WithError(err).Fatal("Failed to print version information")
		}

		fmt.Fprintln(os.Stdout, v)
		return
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().Int64Var(&interval, "interval", 3600, "Set interval for fetching the resource quota and usage.")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log.level", "info", "Set the log level. Must be one of the follwing values: trace, debug, info, warn, error, fatal or panic.")
	rootCmd.PersistentFlags().StringVar(&logOutput, "log.output", "plain", "Set the output format of the log line. Must be plain or json.")
	rootCmd.PersistentFlags().StringVar(&listenAddress, "web.listen-address", ":8080", "Address to listen on for web interface and telemetry.")
	rootCmd.PersistentFlags().StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Failed to initialize Flux Exporter")
	}
}

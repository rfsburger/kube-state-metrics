/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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
	"os"
	"strings"

	"github.com/prometheus/common/version"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"

	"k8s.io/kube-state-metrics/v2/pkg/customresource"
	"k8s.io/kube-state-metrics/v2/pkg/customresourcestate"

	"k8s.io/kube-state-metrics/v2/pkg/app"
	"k8s.io/kube-state-metrics/v2/pkg/options"
)

func main() {
	opts := options.NewOptions()
	opts.AddFlags()

	if err := opts.Parse(); err != nil {
		klog.Fatalf("Parsing flag definitions error: %v", err)
	}

	if opts.Version {
		fmt.Printf("%s\n", version.Print("kube-state-metrics"))
		os.Exit(0)
	}

	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}

	var factories []customresource.RegistryFactory
	if config, set := resolveCustomResourceConfig(); set {
		crf, err := customresourcestate.FromConfig(config)
		if err != nil {
			klog.Fatal(err)
		}
		factories = append(factories, crf...)
	}

	ctx := context.Background()
	if err := app.RunKubeStateMetrics(ctx, opts, factories...); err != nil {
		klog.Fatalf("Failed to run kube-state-metrics: %v", err)
	}
}

func resolveCustomResourceConfig() (customresourcestate.ConfigDecoder, bool) {
	if s, set := os.LookupEnv("KSM_CUSTOM_RESOURCE_METRICS"); set {
		return yaml.NewDecoder(strings.NewReader(s)), true
	}
	if file, set := os.LookupEnv("KSM_CUSTOM_RESOURCE_METRICS_FILE"); set {
		f, err := os.Open(file)
		if err != nil {
			klog.Fatal(err)
		}
		return yaml.NewDecoder(f), true
	}
	return nil, false
}

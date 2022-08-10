package main

import (
	"os"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"

	// Ensure gvk gets loaded in wrangler/pkg/gvk cache
	_ "github.com/rancher/wrangler/pkg/generated/controllers/apiextensions.k8s.io/v1"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/harvester/vm-import-controller/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"source.harvesterhci.io": {
				Types: []interface{}{
					"./pkg/apis/source.harvesterhci.io/v1beta1",
				},
				GenerateTypes: true,
			},
			"importjob.harvesterhci.io": {
				Types: []interface{}{
					"./pkg/apis/importjob.harvesterhci.io/v1beta1",
				},
				GenerateTypes: true,
			},
		},
	})
}

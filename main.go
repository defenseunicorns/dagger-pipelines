package main

import (
	"context"
	"fmt"
	"os"

	"main/pkg/container"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// Create Dagger client
	client, err := dagger.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	// Create new Zarf container
	ctr := container.NewZarfContainer(client)

	// Login to OCI registry
	ctr, err = container.ZarfRegistryLogin(ctx, ctr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create Zarf package
	pkgCreateArgs := []string{".", "-a=arm64"}
	ctr, err = container.CreateZarfPackage(ctx, ctr, pkgCreateArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Deploy Zarf package
	pkgDeployArgs := []string{"zarf-package-podinfo-arm64-1.0.0.tar.zst"}
	_, err = container.DeployZarfPackage(ctx, ctr, pkgDeployArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Publish Zarf package
	pkgPublishArgs := []string{"zarf-package-podinfo-arm64-1.0.0.tar.zst", "oci://host.docker.internal:5000"}
	_, err = container.PublishZarfPackage(ctx, ctr, pkgPublishArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

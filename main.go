package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"main/pkg/container"

	"dagger.io/dagger"
)

func main() {
	auth := &container.RegistryAuth{}

	auth.URL = os.Getenv("REGISTRY_URL")
	auth.Username = os.Getenv("REGISTRY_USERNAME")
	auth.Password = os.Getenv("REGISTRY_PASSWORD")

	arch := runtime.GOARCH
	pkgPath := fmt.Sprintf("zarf-package-podinfo-%s-1.0.0.tar.zst", arch)

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
	ctr, err = container.ZarfRegistryLogin(ctx, ctr, auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create Zarf package
	pkgCreateArgs := []string{".", "-a", arch}
	ctr, err = container.CreateZarfPackage(ctx, ctr, pkgCreateArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Deploy Zarf package
	pkgDeployArgs := []string{pkgPath}
	_, err = container.DeployZarfPackage(ctx, ctr, pkgDeployArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Publish Zarf package
	pkgPublishArgs := []string{pkgPath, "oci://" + auth.URL}
	_, err = container.PublishZarfPackage(ctx, ctr, pkgPublishArgs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

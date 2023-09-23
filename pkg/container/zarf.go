package container

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"dagger.io/dagger"
)

type RegistryAuth struct {
	URL      string
	Username string
	Password string
}

func NewRegistryAuth() *RegistryAuth {
	auth := &RegistryAuth{
		URL:      os.Getenv("REGISTRY_URL"),
		Username: os.Getenv("REGISTRY_USERNAME"),
		Password: os.Getenv("REGISTRY_PASSWORD"),
	}
	return auth
}

// NewZarfContainer creates a new container to run Zarf commands from.
func NewZarfContainer(client *dagger.Client) *dagger.Container {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	homeDir := currentUser.HomeDir
	kubePath := filepath.Join(homeDir, ".kube", "config")

	return client.Container().
		From("cgr.dev/chainguard/wolfi-base:latest").
		WithFile("/tmp/kubeconfig.yaml", client.Host().File(kubePath)).
		WithEnvVariable("KUBECONFIG", "/tmp/kubeconfig.yaml").
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithDirectory("/ci", client.Host().Directory(".")).
		WithWorkdir("/ci").
		WithExec([]string{"apk", "add", "zarf"}).
		WithExec([]string{"./update-kubeconfig.sh"}).
		WithEntrypoint([]string{"/usr/bin/zarf"})
}

func ZarfRegistryLogin(ctx context.Context, ctr *dagger.Container, auth *RegistryAuth, args ...string) (*dagger.Container, error) {
	base := []string{"tools", "registry", "login", "--username", auth.Username, "--password", auth.Password, auth.URL}
	cmd := append(base, args...)

	ctr = ctr.WithExec(cmd).Pipeline("Create Zarf Package")

	if _, err := ctr.Stderr(ctx); err != nil {
		return nil, err
	}

	return ctr, nil
}

func CreateZarfPackage(ctx context.Context, ctr *dagger.Container, args ...string) (*dagger.Container, error) {
	base := []string{"package", "create", "--confirm"}
	cmd := append(base, args...)

	ctr = ctr.WithExec(cmd).Pipeline("Create Zarf Package")

	if _, err := ctr.Stderr(ctx); err != nil {
		return nil, err
	}

	return ctr, nil
}

func DeployZarfPackage(ctx context.Context, ctr *dagger.Container, args ...string) (*dagger.Container, error) {
	base := []string{"package", "deploy", "--confirm"}
	cmd := append(base, args...)

	ctr = ctr.WithExec(cmd).Pipeline("Deploy Zarf Package")

	if _, err := ctr.Stderr(ctx); err != nil {
		return nil, err
	}

	return ctr, nil
}

func PublishZarfPackage(ctx context.Context, ctr *dagger.Container, args ...string) (*dagger.Container, error) {
	base := []string{"package", "publish", "--insecure"}
	cmd := append(base, args...)

	ctr = ctr.WithExec(cmd).Pipeline("Deploy Zarf Package")

	if _, err := ctr.Stderr(ctx); err != nil {
		return nil, err
	}

	return ctr, nil
}

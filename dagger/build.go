package main

import (
	"runtime"
)

const (
	// cgr.dev/chainguard/wolfi-base:latest 6/26/2024
	wolfiBase = "cgr.dev/chainguard/wolfi-base@sha256:7a5b796ae54f72b78b7fc33c8fffee9a363af2c6796dac7c4ef65de8d67d348d"
)

// Build Goreleaser
func (g *Goreleaser) Build(
	// Target OS to build
	// +default="linux"
	os string,
	// Target architecture to build
	// +optional
	arch string,
) *File {
	if arch == "" {
		arch = runtime.GOARCH
	}
	return g.BuildEnv().
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{"go", "build", "-o", "/src/dist/goreleaser"}).
		File("/src/dist/goreleaser")
}

// Container to build Goreleaser
func (g *Goreleaser) BuildEnv() *Container {
	// Base image with Go
	env := dag.Container().
		From(wolfiBase).
		WithExec([]string{"apk", "add", "go"})

	// Mount the Go cache
	env = env.
		WithMountedCache(
			"/go",
			dag.CacheVolume("goreleaser-goroot"),
			ContainerWithMountedCacheOpts{
				Owner: "nonroot",
			}).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod")

	// Mount the Go build cache
	env = env.
		WithMountedCache(
			"/gocache",
			dag.CacheVolume("goreleaser-gobuild"),
			ContainerWithMountedCacheOpts{
				Owner: "nonroot",
			}).
		WithEnvVariable("GOCACHE", "/gocache")

	// Mount the source code
	env = env.
		WithMountedDirectory("/src", g.Source, ContainerWithMountedDirectoryOpts{
			Owner: "nonroot",
		}).
		WithWorkdir("/src")

	return env
}

//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aserto-dev/mage-loot/buf"
	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/aserto-dev/sver/pkg/sver"
	"github.com/magefile/mage/mg"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.19")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

var (
	oras       = deps.BinDep("oras")
	mediaType  = "application/vnd.unknown.layer.v1+txt"
	pluginName = "aserto-idp-plugin-json"
	ghName     = "ghcr.io/aserto-dev/aserto-idp-plugins_"
	osMap      = map[string][]string{
		"linux":   {"arm64", "amd64_v1"},
		"darwin":  {"arm64", "amd64_v1"},
		"windows": {"amd64_v1"},
	}

	extensions = map[string]string{
		"linux":   "",
		"darwin":  "",
		"windows": ".exe",
	}
)

// Generate generates all code.
func Generate() error {
	return common.Generate()
}

// TODO: Will be moved to Proto repo
func GenerateProto() error {
	// Build generate
	return buf.Run(
		buf.AddArg("generate"),
	)
}

// Build builds all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Cleans the bin director
func Clean() error {
	return os.RemoveAll("dist")
}

// Release releases the project.
func Release() error {
	return common.Release()
}

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAllReleaser("--snapshot")
}

func Deps() {
	deps.GetAllDeps()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

// All runs all targets in the appropriate order.
// The targets are run in the following order:
// deps, generate, lint, test, build, dockerImage
func All() error {
	mg.SerialDeps(Deps, Generate, Lint, Test, Build)
	return nil
}

func Publish() error {

	username := os.Getenv("DOCKER_USERNAME")
	if username == "" {
		return errors.New("env var DOCKER_USERNAME is not set")
	}
	password := os.Getenv("DOCKER_PASSWORD")
	if password == "" {
		return errors.New("env var DOCKER_PASSWORD is not set")
	}

	version, err := sver.CurrentVersion(true, true)
	if err != nil {
		return fmt.Errorf("couldn't calculate current version: %w", err)
	}

	pwd := os.Getenv("PWD")
	defer os.Chdir(pwd)

	for operatingSystem, archs := range osMap {
		for _, arch := range archs {
			buildPath := filepath.Join(pwd, "dist", pluginName+"_"+operatingSystem+"_"+arch)
			os.Chdir(buildPath)
			grName := fmt.Sprintf("%s%s_%s:%s-%s", ghName, operatingSystem, iff(arch == "amd64_v1", "amd64", arch), "json", version)
			location := fmt.Sprintf("%s%s:%s", pluginName, extensions[operatingSystem], mediaType)

			err = oras("push", "-u", username, "-p", password, grName, location)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func iff[T any](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	}
	return valFalse
}

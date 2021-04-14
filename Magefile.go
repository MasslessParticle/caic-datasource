//+build mage

package main

import (
	"os"
	"path/filepath"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Coverage runs backend tests and makes a coverage report.
func CoverageRace() error {
	// Create a coverage file if it does not already exist
	if err := os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm); err != nil {
		return err
	}

	if err := sh.RunV("go", "test", "-race", "./pkg/...", "-v", "-cover", "-coverprofile=coverage/backend.out"); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.out", "-o", "coverage/backend.html"); err != nil {
		return err
	}

	return nil
}

// Default configures the default target.
var Default = build.BuildAll

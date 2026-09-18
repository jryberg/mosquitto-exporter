package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version = "dev"
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()

func versionString() string {
	return fmt.Sprintf("%s (%s), %s", Version, Commit, runtime.Version())
}

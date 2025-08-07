package main

import (
	"github.com/reflet-devops/go-media-resizer/cli"
	"github.com/reflet-devops/go-media-resizer/context"
)

func main() {
	ctx := context.DefaultContext()
	rootCmd := cli.GetRootCmd(ctx)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

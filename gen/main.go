package main

import (
	"github.com/adlerhurst/protoc-gen-go-cli/gen/proto/example"
	"github.com/spf13/cobra"
)

func main() {
	err := example.ExampleCmd.Execute()
	cobra.CheckErr(err)
}

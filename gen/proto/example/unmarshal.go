package example

import (
	os "os"

	cobra "github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
)

// func UnmarshalExampleMyCallRequest(cmd *cobra.Command, args []string) {
func Unmarshal(cmd *cobra.Command, args []string) {
	cmd.DisableFlagParsing = false
	indexes := messageIndexes(args, "nested")

	if len(indexes) == 0 {
		if err := cmd.ParseFlags(args); err != nil {
			DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
			os.Exit(1)
		}
		return
	}

	// cmd.ParseFlags(args[:remaining]) remaining args
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "nested" {
			DefaultConfig.Logger.Info("visit", "name", f.Name, "anno", f.Annotations)
		}
	})

	if err := cmd.ParseFlags(args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}
}

func unmarshalFlag(cmd *cobra.Command, args []string)

type index struct {
	index int
	flag  string
}

func messageIndexes(args []string, flags ...string) map[string][]int {
	indexes := make(map[string][]int, len(flags))

	for i, arg := range args {
		for _, flag := range flags {
			if arg != flag {
				continue
			}
			if _, ok := indexes[flag]; !ok {
				indexes[flag] = []int{}
			}
			indexes[flag] = append(indexes[flag], i)
			break
		}
	}

	return indexes
}

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var out string

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Output shell completion code",
	Long: `Output shell completion code.
To configure your shell to load completions for each session

# bash
echo '. <(frgm completion bash)' > ~/.bashrc

# zsh
frgm completion zsh > $fpath[1]/_frgm
`,
	ValidArgs: []string{"bash", "zsh"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			o   *os.File
			err error
		)
		sh := args[0]
		if out == "" {
			o = os.Stdout
		} else {
			o, err = os.Create(out)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}

		switch sh {
		case "bash":
			if err := rootCmd.GenBashCompletion(o); err != nil {
				_ = o.Close()
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		case "zsh":
			if err := rootCmd.GenZshCompletion(o); err != nil {
				_ = o.Close()
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
		if err := o.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringVarP(&out, "out", "o", "", "output file path")
}

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
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/k1LoW/ghput/gh"
	"github.com/spf13/cobra"
)

// prCommentCmd represents the prComment command
var prCommentCmd = &cobra.Command{
	Use:   "pr-comment",
	Short: "Put comment to pull request",
	Long:  `Put comment to pull request.`,
	Args: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stdin.Stat()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			return errors.New("ghput need STDIN. Please use pipe")
		}
		if owner == "" || repo == "" {
			return errors.New("ghput need `--owner` AND `--repo` flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runPrComment(os.Stdin, os.Stdout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runPrComment(stdin io.Reader, stdout io.Writer) (int, error) {
	g, err := gh.New(owner, repo)
	if err != nil {
		return 1, err
	}
	c, err := ioutil.ReadAll(stdin)
	if err != nil {
		return 1, err
	}
	ctx := context.Background()
	if err := g.PutPrComment(ctx, number, string(c)); err != nil {
		return 1, err
	}
	return 0, nil
}

func init() {
	rootCmd.AddCommand(prCommentCmd)
	prCommentCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	prCommentCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	prCommentCmd.Flags().IntVarP(&number, "number", "", 0, "pull request number")
}

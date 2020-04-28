/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/k1LoW/ghput/gh"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Put commit to branch",
	Long:  `Put commit to branch.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if owner == "" || repo == "" || branch == "" {
			return errors.New("`ghput commit` need `--owner` AND `--repo` AND `--branch` flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runCommit(os.Stdin, os.Stdout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func runCommit(stdin io.Reader, stdout io.Writer) error {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	return g.CommitAndPush(ctx, branch, file, path, message)
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	commitCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	commitCmd.Flags().StringVarP(&branch, "branch", "", "master", "branch")
	commitCmd.Flags().StringVarP(&file, "file", "", "", "commit target file")
	commitCmd.Flags().StringVarP(&path, "path", "", "", "commit path")
	commitCmd.Flags().StringVarP(&message, "message", "", "commit by ghput", "commit message")
}

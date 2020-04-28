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
		if owner == "" || repo == "" || number == 0 {
			return errors.New("`ghput pr-comment` need `--owner` AND `--repo` AND `--number` flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runPrComment(os.Stdin, os.Stdout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func runPrComment(stdin io.Reader, stdout io.Writer) error {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	b, err := g.IsPullRequest(ctx, number)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("#%d is not pull request", number)
	}
	comment, err := g.MakeComment(ctx, stdin, header, footer)
	if err != nil {
		return err
	}
	if err := g.DeleteCurrentIssueComment(ctx, number); err != nil {
		return err
	}
	if err := g.PutIssueComment(ctx, number, comment); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(prCommentCmd)
	prCommentCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	prCommentCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	prCommentCmd.Flags().IntVarP(&number, "number", "", 0, "pull request number")
	prCommentCmd.Flags().StringVarP(&header, "header", "", "", "comment header")
	prCommentCmd.Flags().StringVarP(&footer, "footer", "", "", "comment footer")
	prCommentCmd.Flags().StringVarP(&key, "key", "", "", "key for uniquely identifying the comment")
}

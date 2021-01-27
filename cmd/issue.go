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

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Put new issue to repo",
	Long:  `Put new issue to repo.`,
	Args: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			return errors.New("ghput need STDIN. Please use pipe")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIssue(os.Stdin, os.Stdout)
	},
}

func runIssue(stdin io.Reader, stdout io.Writer) error {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	comment, err := g.MakeComment(ctx, stdin, header, footer)
	if err != nil {
		return err
	}
	n, err := g.CreateIssue(ctx, title, comment, assignees)
	if err != nil {
		return err
	}
	if err := g.CloseIssuesUsingTitleMatch(ctx, closeTitle, n); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(stdout, "%d\n", n)
	return nil
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	if err := issueCmd.MarkFlagRequired("owner"); err != nil {
		issueCmd.PrintErrln(err)
		os.Exit(1)
	}
	issueCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	if err := issueCmd.MarkFlagRequired("repo"); err != nil {
		issueCmd.PrintErrln(err)
		os.Exit(1)
	}
	issueCmd.Flags().StringVarP(&title, "title", "", "", "issue title")
	if err := issueCmd.MarkFlagRequired("title"); err != nil {
		issueCmd.PrintErrln(err)
		os.Exit(1)
	}
	issueCmd.Flags().StringVarP(&header, "header", "", "", "comment header")
	issueCmd.Flags().StringVarP(&footer, "footer", "", "", "comment footer")
	issueCmd.Flags().StringSliceVarP(&assignees, "assignee", "a", []string{}, "issue assignee")
	issueCmd.Flags().StringVarP(&closeTitle, "close-issues-using-title-match", "", "", "close current open issues using title match")
}

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

// issueCommentCmd represents the issueComment command
var issueCommentCmd = &cobra.Command{
	Use:   "issue-comment",
	Short: "Put comment to issue",
	Long:  `Put comment to issue.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := setOwnerRepo(); err != nil {
			return err
		}
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
		return runIssueComment(os.Stdin, os.Stdout)
	},
}

func runIssueComment(stdin io.Reader, stdout io.Writer) error {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	b, err := g.IsIssue(ctx, number)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("#%d is not issue", number)
	}
	c, err := getStdin(ctx, stdin)
	if err != nil {
		return err
	}
	body := string(c)
	comment, err := g.MakeComment(ctx, body, header, footer)
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
	rootCmd.AddCommand(issueCommentCmd)
	issueCommentCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	issueCommentCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	issueCommentCmd.Flags().IntVarP(&number, "number", "", 0, "issue number")
	if err := issueCommentCmd.MarkFlagRequired("number"); err != nil {
		issueCommentCmd.PrintErrln(err)
		os.Exit(1)
	}
	issueCommentCmd.Flags().StringVarP(&header, "header", "", "", "comment header")
	issueCommentCmd.Flags().StringVarP(&footer, "footer", "", "", "comment footer")
	issueCommentCmd.Flags().StringVarP(&key, "key", "", "", "key for uniquely identifying the comment")
}

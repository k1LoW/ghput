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
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/k1LoW/ghput/gh"
	"github.com/mattn/go-colorable"
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
	ctx := context.Background()
	g, err := gh.New(owner, repo)
	if err != nil {
		return 1, err
	}
	c, err := getStdin(ctx, stdin)
	if err != nil {
		return 1, err
	}
	if !allowMulti {
		if err := g.DeleteCurrentPrComment(ctx, number); err != nil {
			return 1, err
		}
	}
	if err := g.PutPrComment(ctx, number, string(c)+gh.Footer); err != nil {
		return 1, err
	}
	return 0, nil
}

func getStdin(ctx context.Context, stdin io.Reader) (string, error) {
	in := bufio.NewReader(stdin)
	out := new(bytes.Buffer)
	nc := colorable.NewNonColorable(out)
	for {
		s, err := in.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		select {
		case <-ctx.Done():
			break
		default:
			_, err = nc.Write(s)
			if err != nil {
				return "", err
			}
		}
	}
	return out.String(), nil
}

func init() {
	rootCmd.AddCommand(prCommentCmd)
	prCommentCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	prCommentCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	prCommentCmd.Flags().IntVarP(&number, "number", "", 0, "pull request number")
	prCommentCmd.Flags().BoolVarP(&allowMulti, "allow-multiple-comments", "", false, "allow multiple comments")
}

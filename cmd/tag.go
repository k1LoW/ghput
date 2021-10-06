/*
Copyright © 2021 Ken'ichiro Oyama <k1lowxb@gmail.com>

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
	"io"
	"os"
	"time"

	"github.com/itchyny/timefmt-go"
	"github.com/k1LoW/ghput/gh"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Put tag to branch",
	Long:  `Put tag to branch.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := setOwnerRepo(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTag(os.Stdin, os.Stdout)
	},
}

func runTag(stdin io.Reader, stdout io.Writer) error {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	if branch == "" {
		branch, err = g.GetDefaultBranch(ctx)
		if err != nil {
			return err
		}
	}
	if tag == "" {
		tag = timefmt.Format(time.Now(), tagTimeFormat)
	}
	if err := g.CreateTag(ctx, branch, tag); err != nil {
		return err
	}
	if release {
		if releaseBody == "" {
			fi, err := os.Stdin.Stat()
			if err != nil {
				return err
			}
			if (fi.Mode() & os.ModeCharDevice) == 0 {
				c, err := getStdin(ctx, stdin)
				if err != nil {
					return err
				}
				releaseBody = string(c)
			}
		}
		if err := g.CreateRelease(ctx, tag, releaseTitle, releaseBody); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().StringVarP(&owner, "owner", "", "", "owner")
	tagCmd.Flags().StringVarP(&repo, "repo", "", "", "repo")
	tagCmd.Flags().StringVarP(&branch, "branch", "", "", "branch (default: default branch of repository)")
	tagCmd.Flags().StringVarP(&tag, "tag", "", "", "tag")
	tagCmd.Flags().StringVarP(&tagTimeFormat, "tag-time-format", "", "%Y%m%d-%H%M%S%z", "time format of tag")

	tagCmd.Flags().BoolVarP(&release, "release", "", false, "create a tag as a release.")
	tagCmd.Flags().StringVarP(&releaseTitle, "release-title", "", "", "release title")
	tagCmd.Flags().StringVarP(&releaseBody, "release-body", "", "", "release body")
}

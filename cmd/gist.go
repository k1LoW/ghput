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
	"io"
	"os"
	"path/filepath"

	"github.com/k1LoW/ghput/gh"
	"github.com/spf13/cobra"
)

// gistCmd represents the gist command
var gistCmd = &cobra.Command{
	Use:   "gist",
	Short: "Put gist",
	Long:  `Put gist.`,
	Args: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if (fi.Mode()&os.ModeCharDevice) != 0 && file == "" {
			return errors.New("`ghput gist` need `--file` OR STDIN")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGist(os.Stdin, os.Stdout)
	},
}

func runGist(stdin io.Reader, stdout io.Writer) (err error) {
	ctx := context.Background()
	g, err := gh.New(owner, repo, key)
	if err != nil {
		return err
	}
	var (
		r io.Reader
	)
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer func() {
			err = f.Close()
		}()
		r = f
		if filename == "" {
			filename = filepath.Base(file)
		}
	} else {
		r = stdin
		if filename == "" {
			filename = "stdin"
		}
	}
	return g.CreateGist(ctx, filename, public, r, stdout)
}

func init() {
	rootCmd.AddCommand(gistCmd)
	gistCmd.Flags().StringVarP(&file, "file", "", "", "target file")
	gistCmd.Flags().StringVarP(&filename, "filename", "", "", "filename")
	gistCmd.Flags().BoolVarP(&public, "public", "", false, "public gist")
}

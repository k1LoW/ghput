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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/k1LoW/ghput/gh"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

var (
	owner      string
	repo       string
	number     int
	header     string
	footer     string
	allowMulti bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghput",
	Short: "ghput",
	Long:  `ghput.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {}

func makeComment(ctx context.Context, stdin io.Reader, header, footer string) (string, error) {
	c, err := getStdin(ctx, stdin)
	if err != nil {
		return "", err
	}
	body := string(c)
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	if header != "" && !strings.HasSuffix(header, "\n") {
		header += "\n"
	}
	if footer != "" && !strings.HasSuffix(footer, "\n") {
		footer += "\n"
	}
	return fmt.Sprintf("%s%s%s%s\n", header, body, footer, gh.Footer), nil
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

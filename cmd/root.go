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
	"errors"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

var (
	owner         string
	repo          string
	number        int
	header        string
	footer        string
	key           string
	branch        string
	file          string
	path          string
	message       string
	public        bool
	filename      string
	title         string
	assignees     []string
	closeTitle    string
	latestMerged  bool
	tag           string
	tagTimeFormat string
	release       bool
	releaseTitle  string
	releaseBody   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "ghput",
	Short:        "ghput is a CI-friendly tool that puts * on GitHub.",
	Long:         `ghput is a CI-friendly tool that puts * on GitHub.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setOwnerRepo() error {
	if owner == "" && repo == "" && os.Getenv("GITHUB_REPOSITORY") != "" {
		splitted := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
		if len(splitted) == 2 {
			owner = splitted[0]
			repo = splitted[1]
		}
	}
	if owner == "" || repo == "" {
		return errors.New("--owner and --repo are not set")
	}
	return nil
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

# ghput [![Build Status](https://github.com/k1LoW/ghput/workflows/build/badge.svg)](https://github.com/k1LoW/ghput/actions) [![GitHub release](https://img.shields.io/github/release/k1LoW/ghput.svg)](https://github.com/k1LoW/ghput/releases)

:octocat: ghput is a CI-friendly tool that puts `*` on GitHub.

## Usage

**Put comment to issue:**

``` console
$ echo 'This is comment !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput issue-comment --owner k1LoW --repo myrepo --number 1
```

**Put comment to issue on GitHub Actions:**

ghput get `owner` and `repo` from [`GITHUB_REPOSITORY`](https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables).

``` yaml
- name: Put comment
  run: echo 'This is comment !!' | ghput issue-comment --number 1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Put comment to pull request:**

``` console
$ echo 'This is comment !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput pr-comment --owner k1LoW --repo myrepo --number 2
```

**Put comment to latest merged pull request:**

``` console
$ echo 'Hello merged pull request !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput pr-comment --owner k1LoW --repo myrepo --latest-merged
```

**Put issue to repo:**

``` console
$ echo 'This is new isssue !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput issue --owner k1LoW --repo myrepo --title 'New Issue !!!!!'
```

**Put commit to branch:**

``` console
$ GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput commit --owner k1LoW --repo myrepo --branch main --file file.txt --path path/to/file.txt --message 'Commit file !!'
```

**Put tag to branch:**


``` console
$ GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput tag --owner k1LoW --repo myrepo --branch main --tag v0.0.1
```

Put a time based tag when the tag option is omitted ( default time format: `%Y%m%d-%H%M%S%z` )

``` console
$ GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput tag --owner k1LoW --repo myrepo --tag-time-format '%Y%m%d%H%M'
```

**Put tag as release:**


``` console
$ GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput tag --owner k1LoW --repo myrepo --tag v0.0.1 --release --release-title 'Code name: Elvis Juice' --release-body 'Release !!'
```

or

``` console
$ echo 'Release !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput tag --owner k1LoW --repo myrepo --tag v0.0.1 --release --release-title 'Code name: Elvis Juice'
```

**Put file to Gist:**

``` console
$ GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput gist --file file.txt
```

``` console
$ cat file.txt | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput gist
```

**Use on GitHub Enterprise:**

``` console
$ export GITHUB_API_URL=https://git.my-company.com/api/v3/
```

## Install

**deb:**

Use [dpkg-i-from-url](https://github.com/k1LoW/dpkg-i-from-url)

``` console
$ export GHPUT_VERSION=X.X.X
$ curl -L https://git.io/dpkg-i-from-url | bash -s -- https://github.com/k1LoW/ghput/releases/download/v$GHPUT_VERSION/ghput_$GHPUT_VERSION-1_amd64.deb
```

**RPM:**

``` console
$ export GHPUT_VERSION=X.X.X
$ yum install https://github.com/k1LoW/ghput/releases/download/v$GHPUT_VERSION/ghput_$GHPUT_VERSION-1_amd64.rpm
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/ghput
```

**manually:**

Download binany from [releases page](https://github.com/k1LoW/ghput/releases)

**go get:**

```console
$ go get github.com/k1LoW/ghput
```

## Alternatives

- [github-commenter](https://github.com/cloudposse/github-commenter): Command line utility for creating GitHub comments on Commits, Pull Request Reviews or Issues

# ghput

ghput is a CI-Friendly tool for put comment on GitHub.

## Usage

``` console
$ echo 'This is comment message !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput pr-comment --owner k1LoW --repo ghput --number 2
```

**GitHub Enterprise:**

``` console
$ export GITHUB_BASE_URL=https://git.my-company.com/api/v3/
```

## Install

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

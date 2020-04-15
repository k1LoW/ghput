# ghput

ghput is a CI-Friendly tool for put comment on GitHub.

## Usage

``` console
$ echo 'This is comment !!' | GITHUB_TOKEN=XXXXXxxxxxXXxxxx ghput pr-comment --owner k1LoW --repo ghput --number 2
```

**GitHub Enterprise:**

``` console
$ export GITHUB_BASE_URL=https://git.my-company.com/api/v3/
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

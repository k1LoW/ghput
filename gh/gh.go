package gh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/k1LoW/go-github-client/v39/factory"
)

const (
	footerFormat = "<!-- Put by ghput %s-->"
)

var rep = strings.NewReplacer("\\n", "\n", "\\t", "\t")

type Gh struct {
	client *github.Client
	owner  string
	repo   string
	key    string
}

// New return Gh
func New(owner, repo, key string) (*Gh, error) {
	client, err := factory.NewGithubClient()
	if err != nil {
		return nil, err
	}
	return &Gh{
		client: client,
		owner:  owner,
		repo:   repo,
		key:    key,
	}, nil
}

func (g *Gh) MakeComment(ctx context.Context, body, header, footer string) (string, error) {
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	header = rep.Replace(header)
	if header != "" && !strings.HasSuffix(header, "\n") {
		header += "\n"
	}
	footer = rep.Replace(footer)
	if footer != "" && !strings.HasSuffix(footer, "\n") {
		footer += "\n"
	}
	return fmt.Sprintf("%s%s%s%s\n", header, body, footer, g.CommentFooter()), nil
}

func (g *Gh) CommentFooter() string {
	if g.key == "" {
		return fmt.Sprintf(footerFormat, g.key)
	}
	key := fmt.Sprintf("[key:%s] ", g.key)
	return fmt.Sprintf(footerFormat, key)
}

func (g *Gh) FetchLatestMergedPullRequest(ctx context.Context) (int, error) {
	commits, _, err := g.client.Repositories.ListCommits(ctx, g.owner, g.repo, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err != nil {
		return 0, err
	}
	for _, c := range commits {
		m := c.GetCommit().GetMessage()
		if strings.HasPrefix(m, "Merge pull request #") {
			splitted := strings.Split(strings.TrimPrefix(m, "Merge pull request #"), " ")
			if len(splitted) < 1 {
				break
			}
			n, err := strconv.Atoi(splitted[0])
			if err != nil {
				break
			}
			return n, nil
		}
	}
	// fallback
	q := fmt.Sprintf("type:pr is:merged sort:updated-desc repo:%s/%s", g.owner, g.repo)
	prs, _, err := g.client.Search.Issues(ctx, q, &github.SearchOptions{
		Sort:  "updated",
		Order: "desc",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	})
	if err != nil {
		return 0, err
	}
	if len(prs.Issues) == 0 {
		return 0, err
	}
	return prs.Issues[0].GetNumber(), nil
}

func (g *Gh) GetDefaultBranch(ctx context.Context) (string, error) {
	r, _, err := g.client.Repositories.Get(ctx, g.owner, g.repo)
	if err != nil {
		return "", err
	}
	return r.GetDefaultBranch(), nil
}

func (g Gh) IsPullRequest(ctx context.Context, n int) (bool, error) {
	i, _, err := g.client.Issues.Get(ctx, g.owner, g.repo, n)
	if err != nil {
		return false, err
	}
	return i.IsPullRequest(), nil
}

func (g Gh) IsIssue(ctx context.Context, n int) (bool, error) {
	b, err := g.IsPullRequest(ctx, n)
	if err != nil {
		return false, err
	}
	return !b, nil
}

func (g *Gh) CreateIssue(ctx context.Context, title string, comment string, assignees []string) (int, error) {
	// trim assignees
	as := []string{}
	for _, a := range assignees {
		splitted := strings.Split(a, " ")
		for _, s := range splitted {
			if s == "" {
				continue
			}
			trimed := strings.Trim(s, "@")
			if !strings.Contains(trimed, "/") {
				as = append(as, trimed)
				continue
			}
			splitted := strings.Split(trimed, "/")
			org := splitted[0]
			slug := splitted[1]
			opts := &github.TeamListTeamMembersOptions{}
			users, _, err := g.client.Teams.ListTeamMembersBySlug(ctx, org, slug, opts)
			if err != nil {
				return 0, err
			}
			for _, u := range users {
				as = append(as, *u.Login)
			}
		}
	}
	as = unique(as)

	r := &github.IssueRequest{Title: &title, Body: &comment, Assignees: &as}
	i, _, err := g.client.Issues.Create(ctx, g.owner, g.repo, r)
	if err != nil {
		return 0, err
	}
	return *i.Number, nil
}

func (g *Gh) PutIssueComment(ctx context.Context, n int, comment string) error {
	c := &github.IssueComment{Body: &comment}
	if _, _, err := g.client.Issues.CreateComment(ctx, g.owner, g.repo, n, c); err != nil {
		return err
	}
	return nil
}

func (g *Gh) DeleteCurrentIssueComment(ctx context.Context, n int) error {
	listOptions := &github.IssueListCommentsOptions{}
	comments, _, err := g.client.Issues.ListComments(ctx, g.owner, g.repo, n, listOptions)
	if err != nil {
		return err
	}
	for _, c := range comments {
		if strings.Contains(*c.Body, g.CommentFooter()) {
			_, err = g.client.Issues.DeleteComment(ctx, g.owner, g.repo, *c.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Gh) CommitAndPush(ctx context.Context, branch, content, rPath, message string) error {
	srv := g.client.Git

	dRef, _, err := srv.GetRef(ctx, g.owner, g.repo, path.Join("heads", branch))
	if err != nil {
		return err
	}

	parent, _, err := srv.GetCommit(ctx, g.owner, g.repo, *dRef.Object.SHA)
	if err != nil {
		return err
	}

	var tree *github.Tree

	if rPath != "" {
		blob := &github.Blob{
			Content:  github.String(content),
			Encoding: github.String("utf-8"),
			Size:     github.Int(len(content)),
		}

		resB, _, err := srv.CreateBlob(ctx, g.owner, g.repo, blob)
		if err != nil {
			return err
		}

		entry := &github.TreeEntry{
			Path: github.String(rPath),
			Mode: github.String("100644"),
			Type: github.String("blob"),
			SHA:  resB.SHA,
		}

		entries := []*github.TreeEntry{entry}

		tree, _, err = srv.CreateTree(ctx, g.owner, g.repo, *dRef.Object.SHA, entries)
		if err != nil {
			return err
		}
	} else {
		tree, _, err = srv.GetTree(ctx, g.owner, g.repo, *parent.Tree.SHA, false)
	}

	commit := &github.Commit{
		Message: github.String(message),
		Tree:    tree,
		Parents: []*github.Commit{parent},
	}
	resC, _, err := srv.CreateCommit(ctx, g.owner, g.repo, commit)
	if err != nil {
		return err
	}

	nref := &github.Reference{
		Ref: github.String(path.Join("refs", "heads", branch)),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  resC.SHA,
		},
	}
	if _, _, err := srv.UpdateRef(ctx, g.owner, g.repo, nref, false); err != nil {
		return err
	}

	return nil
}

func (g *Gh) CommitAndPushFile(ctx context.Context, branch, file, rPath, message string) error {
	content := ""
	if file != "" {
		f, err := os.Stat(file)
		if err != nil {
			return err
		}
		if f.IsDir() {
			return errors.New("'ghput commit' does not yet support directory commit.")
		}
		b, err := ioutil.ReadFile(filepath.Clean(file))
		if err != nil {
			return err
		}
		content = string(b)
		if rPath == "" {
			rPath = filepath.Base(file)
		}
	}
	return g.CommitAndPush(ctx, branch, content, rPath, message)
}

func (g *Gh) CreateGist(ctx context.Context, fname string, public bool, in io.Reader, out io.Writer) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	content := string(b)
	files := make(map[github.GistFilename]github.GistFile, 1)
	files[github.GistFilename(fname)] = github.GistFile{
		Size:     github.Int(len(content)),
		Filename: github.String(fname),
		Content:  github.String(content),
	}

	input := &github.Gist{
		Description: github.String("Put by ghput"),
		Public:      github.Bool(public),
		Files:       files,
	}
	gist, _, err := g.client.Gists.Create(ctx, input)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "%s\n", *gist.HTMLURL)
	return nil
}

func (g *Gh) CloseIssuesUsingTitle(ctx context.Context, closeTitle string, ignoreNumber int) error {
	if closeTitle == "" {
		return nil
	}
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	r, _, err := g.client.Search.Issues(ctx, fmt.Sprintf("%s state:open type:issue in:title repo:%s/%s", closeTitle, g.owner, g.repo), opts)
	if err != nil {
		return err
	}
	for _, i := range r.Issues {
		if *i.Number == ignoreNumber {
			continue
		}
		if err := g.PutIssueComment(ctx, *i.Number, fmt.Sprintf("Closed when ghput created #%d.", ignoreNumber)); err != nil {
			return err
		}
		closed := "closed"
		if _, _, err := g.client.Issues.Edit(ctx, g.owner, g.repo, *i.Number, &github.IssueRequest{
			State: &closed,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gh) CreateTag(ctx context.Context, branch, tag string) error {
	ref, _, err := g.client.Git.GetRef(ctx, g.owner, g.repo, path.Join("heads", branch))
	if err != nil {
		return err
	}
	tref := fmt.Sprintf("refs/tags/%s", tag)
	if _, _, err := g.client.Git.CreateRef(ctx, g.owner, g.repo, &github.Reference{
		Ref:    &tref,
		Object: ref.GetObject(),
	}); err != nil {
		return err
	}
	return nil
}

func (g *Gh) CreateRelease(ctx context.Context, tag, title, body string) error {
	if _, _, err := g.client.Repositories.CreateRelease(ctx, g.owner, g.repo, &github.RepositoryRelease{
		TagName: &tag,
		Name:    &title,
		Body:    &body,
	}); err != nil {
		return err
	}
	return nil
}

func unique(in []string) []string {
	m := map[string]struct{}{}
	u := []string{}
	for _, s := range in {
		if _, ok := m[s]; !ok {
			u = append(u, s)
		}
		m[s] = struct{}{}
	}
	return u
}

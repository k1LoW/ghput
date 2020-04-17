package gh

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/mattn/go-colorable"
)

const (
	defaultBaseURL = "https://api.github.com/"
	uploadBaseURL  = "https://uploads.github.com/"
	footerFormat   = "<!-- Generated by ghput [key:%s]-->"
)

type Gh struct {
	client *github.Client
	owner  string
	repo   string
	key    string
}

// New return Gh
func New(owner, repo, key string) (*Gh, error) {
	c := github.NewClient(httpClient())
	if baseURL := os.Getenv("GITHUB_BASE_URL"); baseURL != "" {
		baseEndpoint, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(baseEndpoint.Path, "/") {
			baseEndpoint.Path += "/"
		}
		c.BaseURL = baseEndpoint
	}
	if uploadURL := os.Getenv("GITHUB_UPLOAD_URL"); uploadURL != "" {
		uploadEndpoint, err := url.Parse(uploadURL)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(uploadEndpoint.Path, "/") {
			uploadEndpoint.Path += "/"
		}
		c.UploadURL = uploadEndpoint
	}
	return &Gh{
		client: c,
		owner:  owner,
		repo:   repo,
		key:    key,
	}, nil
}

func (g *Gh) MakeComment(ctx context.Context, stdin io.Reader, header, footer string) (string, error) {
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
	return fmt.Sprintf("%s%s%s%s\n", header, body, footer, g.CommentFooter()), nil
}

func (g *Gh) CommentFooter() string {
	return fmt.Sprintf(footerFormat, g.key)
}

func (g *Gh) PutPrComment(ctx context.Context, n int, comment string) error {
	c := &github.IssueComment{Body: &comment}
	if _, _, err := g.client.Issues.CreateComment(ctx, g.owner, g.repo, n, c); err != nil {
		return err
	}
	return nil
}

func (g *Gh) DeleteCurrentPrComment(ctx context.Context, n int) error {
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

type roundTripper struct {
	transport   *http.Transport
	accessToken string
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", rt.accessToken))
	return rt.transport.RoundTrip(r)
}

func httpClient() *http.Client {
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	rt := roundTripper{
		transport:   t,
		accessToken: os.Getenv("GITHUB_TOKEN"),
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: rt,
	}
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

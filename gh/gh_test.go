package gh

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeComment(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		body   string
		header string
		footer string
		want   string
	}{
		{"", "", "", "<!-- Put by ghput -->\n"},
		{"body", "header", "footer", "header\nbody\nfooter\n<!-- Put by ghput -->\n"},
		{"body\n", "header\n", "footer\n", "header\nbody\nfooter\n<!-- Put by ghput -->\n"},
	}
	for _, tt := range tests {
		gh, err := New("o", "r", "")
		if err != nil {
			t.Fatal(err)
		}
		got, err := gh.MakeComment(ctx, tt.body, tt.header, tt.footer)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("got\n%v\nwant\n%v", got, tt.want)
		}
	}
}

func TestCommentFooter(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{
			key:  "",
			want: "<!-- Put by ghput -->",
		},
		{
			key:  "value",
			want: "<!-- Put by ghput [key:value] -->",
		},
	}

	for _, tt := range tests {
		gh, err := New("o", "r", tt.key)
		if err != nil {
			t.Fatal(err)
		}
		got := gh.CommentFooter()
		if got != tt.want {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		{[]string{}, []string{}},
		{[]string{"a", "c", "b"}, []string{"a", "c", "b"}},
		{[]string{"b", "c", "b"}, []string{"b", "c"}},
		{[]string{"a", "a", "b"}, []string{"a", "b"}},
	}
	for _, tt := range tests {
		got := unique(tt.in)
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

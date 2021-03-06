package gh

import (
	"testing"
)

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

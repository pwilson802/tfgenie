package accessanalyzer_test

import (
	"testing"

	"github.com/pwilson802/tfgenie/pkg/accessanalyzer"
)

func TestClientCreation(t *testing.T) {
	t.Parallel()
	testAnalyzer := accessanalyzer.AnalyzerClient{
		UserName: "testuser-1234b",
	}
	want := "testuser-1234b"
	got := testAnalyzer.UserName
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

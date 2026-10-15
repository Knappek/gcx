package output_test

import (
	"runtime"
	"testing"

	"github.com/grafana/gcx/internal/output"
)

func TestBuildPaginationCommandPreservesArgvAndReplacesPaginationFlags(t *testing.T) {
	got := output.BuildPaginationCommand(
		[]string{"gcx", "--context", "dev", "dashboards", "list", "--api-version=dashboard.grafana.app/v2", "--limit", "50", "--continue=old"},
		100,
		"next-token",
	)
	want := "gcx --context dev dashboards list --api-version=dashboard.grafana.app/v2 --limit 100 --continue next-token"
	if runtime.GOOS == "windows" {
		want = "& 'gcx' '--context' 'dev' 'dashboards' 'list' '--api-version=dashboard.grafana.app/v2' '--limit' '100' '--continue' 'next-token'"
	}
	if got != want {
		t.Fatalf("BuildPaginationCommand() = %q, want %q", got, want)
	}
}

func TestBuildPaginationCommandQuotesUnsafeArgs(t *testing.T) {
	got := output.BuildPaginationCommand(
		[]string{"/tmp/my gcx", "dashboards", "list", "--limit=50"},
		50,
		"token with spaces",
	)
	want := "'/tmp/my gcx' dashboards list --limit 50 --continue 'token with spaces'"
	if runtime.GOOS == "windows" {
		want = "& '/tmp/my gcx' 'dashboards' 'list' '--limit' '50' '--continue' 'token with spaces'"
	}
	if got != want {
		t.Fatalf("BuildPaginationCommand() = %q, want %q", got, want)
	}
}

func TestBuildPaginationCommandQuotesApostrophe(t *testing.T) {
	got := output.BuildPaginationCommand([]string{"gcx", "dashboards", "list"}, 50, "next'page")
	want := "gcx dashboards list --limit 50 --continue 'next'\\''page'"
	if runtime.GOOS == "windows" {
		want = "& 'gcx' 'dashboards' 'list' '--limit' '50' '--continue' 'next''page'"
	}
	if got != want {
		t.Fatalf("BuildPaginationCommand() = %q, want %q", got, want)
	}
}

func TestBuildListLimitCommandQuotesPowerShellSyntax(t *testing.T) {
	got := output.BuildListLimitCommand([]string{"C:\\Program Files\\gcx.exe", "dashboards", "list", "--selector=a,b", "--context", "team's prod"}, 25)
	if runtime.GOOS != "windows" {
		return
	}
	want := "& 'C:\\Program Files\\gcx.exe' 'dashboards' 'list' '--selector=a,b' '--context' 'team''s prod' '--limit' '25'"
	if got != want {
		t.Fatalf("BuildListLimitCommand() = %q, want %q", got, want)
	}
}

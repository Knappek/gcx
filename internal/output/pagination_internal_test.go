package output

import "testing"

func TestShellJoinForOS(t *testing.T) {
	args := []string{`C:\Program Files\gcx.exe`, "--selector=a,b", "--context", "team's prod"}
	tests := []struct {
		name string
		goos string
		want string
	}{
		{
			name: "PowerShell",
			goos: "windows",
			want: `& 'C:\Program Files\gcx.exe' '--selector=a,b' '--context' 'team''s prod'`,
		},
		{
			name: "POSIX shell",
			goos: "linux",
			want: `'C:\Program Files\gcx.exe' --selector=a,b --context 'team'\''s prod'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellJoinForOS(args, tt.goos); got != tt.want {
				t.Fatalf("shellJoinForOS() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShellJoinForOSPreservesLegacyNativeArguments(t *testing.T) {
	args := []string{`C:\Program Files\gcx.exe`, `a"b`, "", `tail\`, "team's prod"}
	want := `& { $PSNativeCommandArgumentPassing = 'Legacy'; & 'C:\Program Files\gcx.exe' '"a\"b"' '""' '"tail\\"' '"team''s prod"' }`
	if got := shellJoinForOS(args, "windows"); got != want {
		t.Fatalf("shellJoinForOS() = %q, want %q", got, want)
	}
}

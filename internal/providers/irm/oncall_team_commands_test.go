//nolint:testpackage // white-box tests require access to unexported IRM command builders
package irm

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// fakeTeamAPI stubs the team working-hours surface. Unimplemented OnCallAPI
// methods panic via the embedded nil interface.
type fakeTeamAPI struct {
	OnCallAPI

	hours    map[string]any
	gotID    string
	gotHours map[string]any
}

func (f *fakeTeamAPI) GetTeamWorkingHours(_ context.Context, id string) (map[string]any, error) {
	f.gotID = id
	return f.hours, nil
}

func (f *fakeTeamAPI) UpdateTeamWorkingHours(_ context.Context, id string, h map[string]any) (map[string]any, error) {
	f.gotID = id
	f.gotHours = h
	return h, nil
}

func runTeamsCmd(t *testing.T, fake *fakeTeamAPI, args ...string) (string, error) {
	t.Helper()
	cmd := newTeamsCmd(&fakeLoader{client: fake})
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs(args)
	err := cmd.ExecuteContext(context.Background())
	return out.String(), err
}

func TestTeamGetWorkingHoursCommand(t *testing.T) {
	resetAgentMode(t)

	fake := &fakeTeamAPI{hours: map[string]any{
		"irm_working_hours": []any{map[string]any{"name": "Office hours"}},
	}}

	out, err := runTeamsCmd(t, fake, "get-working-hours", "T1", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	if fake.gotID != "T1" {
		t.Errorf("expected hours of T1, got %q", fake.gotID)
	}

	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if intervals, _ := got["irm_working_hours"].([]any); len(intervals) != 1 {
		t.Errorf("unexpected intervals: %v", got["irm_working_hours"])
	}
}

func TestTeamUpdateWorkingHoursCommand(t *testing.T) {
	resetAgentMode(t)

	file := filepath.Join(t.TempDir(), "hours.yaml")
	doc := `irm_working_hours:
  - name: Office hours
    irm_working_hours:
      monday: [{start: "09:00:00", end: "17:00:00"}]
      sunday: []
`
	if err := os.WriteFile(file, []byte(doc), 0o600); err != nil {
		t.Fatal(err)
	}

	fake := &fakeTeamAPI{}
	if _, err := runTeamsCmd(t, fake, "update-working-hours", "T1", "-f", file, "-o", "json"); err != nil {
		t.Fatal(err)
	}
	if fake.gotID != "T1" {
		t.Errorf("expected update of T1, got %q", fake.gotID)
	}

	intervals, _ := fake.gotHours["irm_working_hours"].([]any)
	if len(intervals) != 1 {
		t.Fatalf("expected one interval, got %v", fake.gotHours)
	}
	interval, _ := intervals[0].(map[string]any)
	if interval["name"] != "Office hours" {
		t.Errorf("unexpected interval: %v", interval)
	}
}

func TestTeamUpdateWorkingHoursClear(t *testing.T) {
	resetAgentMode(t)

	file := filepath.Join(t.TempDir(), "clear.yaml")
	if err := os.WriteFile(file, []byte("irm_working_hours: null\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fake := &fakeTeamAPI{}
	if _, err := runTeamsCmd(t, fake, "update-working-hours", "T1", "-f", file, "-o", "json"); err != nil {
		t.Fatal(err)
	}
	if v, ok := fake.gotHours["irm_working_hours"]; !ok || v != nil {
		t.Errorf("expected an explicit null to reach the backend, got %v", fake.gotHours)
	}
}

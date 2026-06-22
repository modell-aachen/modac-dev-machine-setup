package orbstack

import "testing"

func TestWaitForRunning(t *testing.T) {
	tests := []struct {
		name      string
		statuses  []string
		attempts  int
		wantErr   bool
		wantWaits int
	}{
		{"already running", []string{statusRunning}, 3, false, 0},
		{"stopped then running", []string{"Stopped", statusRunning}, 3, false, 1},
		{"empty status then running", []string{"", statusRunning}, 3, false, 1},
		{"running on last attempt", []string{"Stopped", "Stopped", "Stopped", statusRunning}, 3, false, 3},
		{"never running", []string{"Stopped", "Stopped", "Stopped", "Stopped"}, 3, true, 3},
		{"never running with empty status", []string{"", "", "", ""}, 3, true, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			status := func() string {
				// Repeat the final scripted value if polled beyond the script,
				// mirroring a status that keeps reporting the same state.
				idx := calls
				if idx >= len(tt.statuses) {
					idx = len(tt.statuses) - 1
				}
				calls++
				return tt.statuses[idx]
			}

			waits := 0
			wait := func() { waits++ }

			err := waitForRunning(status, wait, tt.attempts)

			if (err != nil) != tt.wantErr {
				t.Errorf("waitForRunning() error = %v, wantErr %v", err, tt.wantErr)
			}
			if waits != tt.wantWaits {
				t.Errorf("waitForRunning() waited %d times, want %d", waits, tt.wantWaits)
			}
		})
	}
}

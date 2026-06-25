package devboxupdate

import "testing"

func TestAlreadyUpdated(t *testing.T) {
	t.Run("unset means not yet updated", func(t *testing.T) {
		t.Setenv(RestartGuardEnv, "")
		if AlreadyUpdated() {
			t.Error("AlreadyUpdated() = true with the guard unset, want false")
		}
	})

	t.Run("set to 1 means already updated", func(t *testing.T) {
		t.Setenv(RestartGuardEnv, "1")
		if !AlreadyUpdated() {
			t.Error("AlreadyUpdated() = false with the guard set, want true")
		}
	})
}

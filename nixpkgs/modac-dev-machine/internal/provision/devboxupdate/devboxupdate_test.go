package devboxupdate

import "testing"

func TestAlreadyUpdated(t *testing.T) {
	t.Run("unset means not yet updated", func(t *testing.T) {
		t.Setenv(selfUpdatedEnv, "")
		if alreadyUpdated() {
			t.Error("alreadyUpdated() = true with the guard unset, want false")
		}
	})

	t.Run("set to 1 means already updated", func(t *testing.T) {
		t.Setenv(selfUpdatedEnv, "1")
		if !alreadyUpdated() {
			t.Error("alreadyUpdated() = false with the guard set, want true")
		}
	})
}

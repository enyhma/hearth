package hearth

import "testing"

func TestShouldColor_Precedence(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("FORCE_COLOR", "")

	if ShouldColor(ColorConfig{NoColor: true, ForceColor: true, IsTTY: true}) {
		t.Error("--no-color must win over everything else")
	}
	if !ShouldColor(ColorConfig{ForceColor: true, IsTTY: false}) {
		t.Error("--color should force on even without a TTY")
	}
	if ShouldColor(ColorConfig{IsTTY: false}) {
		t.Error("no flags, no TTY: should stay off")
	}
	if !ShouldColor(ColorConfig{IsTTY: true}) {
		t.Error("no flags, a TTY: should turn on")
	}
}

func TestShouldColor_EnvVars(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	if !ShouldColor(ColorConfig{IsTTY: false}) {
		t.Error("FORCE_COLOR env should force on")
	}
	t.Setenv("FORCE_COLOR", "")

	t.Setenv("NO_COLOR", "1")
	if ShouldColor(ColorConfig{IsTTY: true}) {
		t.Error("NO_COLOR env should force off even on a TTY")
	}
}

func TestColorize(t *testing.T) {
	if got := Colorize(ColorGreen, "ok", false); got != "ok" {
		t.Errorf("disabled: got %q, want unchanged", got)
	}
	got := Colorize(ColorGreen, "ok", true)
	if got == "ok" {
		t.Error("enabled: expected ANSI wrapping, got plain text")
	}
	if visLen(got) != 2 {
		t.Errorf("visLen(got) = %d, want 2", visLen(got))
	}
}

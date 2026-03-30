package validate

import "testing"

func TestPassword_valid(t *testing.T) {
	valid := []string{
		"password1",
		"P4ssword",
		"12345abc",
		"Ab1!@#$%",
		"longpassword99",
	}
	for _, p := range valid {
		if msg := Password(p); msg != "" {
			t.Errorf("Password(%q) = %q, want empty (valid)", p, msg)
		}
	}
}

func TestPassword_tooShort(t *testing.T) {
	if msg := Password("abc1"); msg == "" {
		t.Error("Password(\"abc1\") should be invalid (too short)")
	}
}

func TestPassword_noLetter(t *testing.T) {
	if msg := Password("12345678"); msg == "" {
		t.Error("Password(\"12345678\") should be invalid (no letter)")
	}
}

func TestPassword_noNumber(t *testing.T) {
	if msg := Password("abcdefgh"); msg == "" {
		t.Error("Password(\"abcdefgh\") should be invalid (no number)")
	}
}

func TestPassword_empty(t *testing.T) {
	if msg := Password(""); msg == "" {
		t.Error("Password(\"\") should be invalid")
	}
}

// --- ProjectTag ---

func TestProjectTag_valid(t *testing.T) {
	valid := []string{"KB", "MKB", "PROJ"}
	for _, tag := range valid {
		if msg := ProjectTag(tag); msg != "" {
			t.Errorf("ProjectTag(%q) = %q, want empty (valid)", tag, msg)
		}
	}
}

func TestProjectTag_tooShort(t *testing.T) {
	if msg := ProjectTag("K"); msg == "" {
		t.Error("single character tag should be invalid")
	}
}

func TestProjectTag_tooLong(t *testing.T) {
	if msg := ProjectTag("ABCDE"); msg == "" {
		t.Error("5-character tag should be invalid")
	}
}

func TestProjectTag_lowercase(t *testing.T) {
	if msg := ProjectTag("kb"); msg == "" {
		t.Error("lowercase tag should be invalid")
	}
}

func TestProjectTag_numbers(t *testing.T) {
	if msg := ProjectTag("K1"); msg == "" {
		t.Error("tag with numbers should be invalid")
	}
}

func TestProjectTag_empty(t *testing.T) {
	if msg := ProjectTag(""); msg == "" {
		t.Error("empty tag should be invalid")
	}
}

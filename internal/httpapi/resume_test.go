package httpapi

import "testing"

func TestSlugValidation(t *testing.T) {
	good := []string{"elahe-roshandel", "abc", "a1-b2"}
	bad := []string{"ab", "Admin", "admin", "has space", "-lead", "trail-", "a--b", "api"}
	for _, s := range good {
		if !validSlug(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
	for _, s := range bad {
		if validSlug(s) {
			t.Errorf("expected %q to be invalid", s)
		}
	}
}

func TestSlugify(t *testing.T) {
	if got := slugify("Elahe Roshandel"); got != "elahe-roshandel" {
		t.Errorf("got %q", got)
	}
	if got := slugify("Jo"); got != "user-jo" {
		t.Errorf("got %q", got)
	}
}

func TestLocalizeFallsBack(t *testing.T) {
	if got := pick("en", "Architektin", ""); got != "Architektin" {
		t.Errorf("fallback failed: %q", got)
	}
	if got := pick("en", "Architektin", "Architect"); got != "Architect" {
		t.Errorf("en not picked: %q", got)
	}
	if got := pick("de", "Architektin", "Architect"); got != "Architektin" {
		t.Errorf("de not picked: %q", got)
	}
}

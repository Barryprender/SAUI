package saas

import "strings"

type PreseededMember struct {
	UserID string
	Name   string
	Email  string
	Role   string
}

var preseededMembers = []PreseededMember{
	{"alice-morgan", "Alice Morgan", "alice@demo.io", "admin"},
	{"bob-keller", "Bob Keller", "bob@demo.io", "admin"},
	{"carol-lim", "Carol Lim", "carol@demo.io", "member"},
	{"david-park", "David Park", "david@demo.io", "member"},
	{"elena-ruiz", "Elena Ruiz", "elena@demo.io", "member"},
	{"frank-wong", "Frank Wong", "frank@demo.io", "viewer"},
	{"grace-hall", "Grace Hall", "grace@demo.io", "viewer"},
	{"henry-mills", "Henry Mills", "henry@demo.io", "viewer"},
}

func PreseededMembers() []PreseededMember { return preseededMembers }

// Slugify converts a name like "John Smith" to "john-smith".
func Slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	result := strings.TrimRight(b.String(), "-")
	if result == "" {
		return "user"
	}
	return result
}

// Initials returns up to two initials from a full name.
func Initials(name string) string {
	words := strings.Fields(name)
	if len(words) == 0 {
		return "?"
	}
	if len(words) == 1 {
		return strings.ToUpper(string([]rune(words[0])[:1]))
	}
	return strings.ToUpper(string([]rune(words[0])[:1]) + string([]rune(words[len(words)-1])[:1]))
}

// RoleAvatarClass returns a CSS class for the role's avatar colour.
func RoleAvatarClass(role string) string {
	switch role {
	case "admin":
		return "avatar--amber"
	case "member":
		return "avatar--blue"
	default:
		return "avatar--teal"
	}
}

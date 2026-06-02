package saas

import "strconv"

func activeClass(current, target string) string {
	if current == target {
		return " filter-tab--active"
	}
	return ""
}

func ariaCurrent(current, target string) string {
	if current == target {
		return "page"
	}
	return "false"
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func countLabel(n int) string {
	if n == 1 {
		return "1 active member"
	}
	return strconv.Itoa(n) + " active members"
}

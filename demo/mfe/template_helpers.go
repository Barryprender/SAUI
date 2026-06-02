package mfe

import "fmt"

func taskRowClass(t Task, justToggled string) string {
	cls := "task-row"
	if t.Completed {
		cls += " task-row--done"
	}
	if t.ID == justToggled {
		cls += " task-row--flash"
	}
	return cls
}

func feedRowClass(e FeedEntry) string {
	return "feed-row " + e.ColorClass
}

func progressWidth(pct int) string {
	return fmt.Sprintf("width:%d%%", pct)
}

func statLabel(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

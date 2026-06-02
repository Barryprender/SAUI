package mfe

import "fmt"

var preseedTitles = []string{
	"Set up CI pipeline",
	"Review API documentation",
	"Implement rate limiting",
	"Write integration tests",
	"Deploy to staging",
}

func taskID(seq int) string {
	return fmt.Sprintf("task-%d", seq)
}

package distributed

import "fmt"

func (n PipelineNode) StatusClass() string {
	switch n.Status {
	case StageDone:
		return "pipeline-node--done"
	case StageActive:
		return "pipeline-node--active"
	case StageFailed:
		return "pipeline-node--failed"
	default:
		return "pipeline-node--pending"
	}
}

func (n PipelineNode) CSSClasses() string {
	classes := "pipeline-node " + n.StatusClass() + " " + n.ColorClass
	if n.JustCompleted {
		classes += " pipeline-node--pulse"
	}
	return classes
}

func (n PipelineNode) IconText() string {
	switch n.Status {
	case StageDone:
		return "✓"
	case StageFailed:
		return "✕"
	default:
		return ""
	}
}

func FormatEventID(id int64) string {
	return fmt.Sprintf("EVT-%06d", id)
}

func ConnectorClass(left, right PipelineNode) string {
	if left.Status == StageDone {
		return "pipeline-connector--done"
	}
	return "pipeline-connector--pending"
}

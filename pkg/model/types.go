package model

type HunkOperation string

const (
	ADD    = "add"
	DELETE = "delete"
	MODIFY = "modify"
)

type GitDiff struct {
	FileDiffs []FileDiff
}

type FileDiff struct {
	OldFilename string
	NewFilename string
	Hunks       []Hunk
}

type Hunk struct {
	HunkOperation    HunkOperation
	OldFileLineStart int
	OldFileLineCount int
	NewFileLineStart int
	NewFileLineCount int
	ChangedLines     []ChangedLine
}

type ChangedLine struct {
	Content    string
	IsDeletion bool
}

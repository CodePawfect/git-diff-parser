package diffparser

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

func Parse(gitDiff string) (GitDiff, error) {
	fileDiffsRaw := strings.Split(gitDiff, "diff --git")
	fileDiffsRaw = fileDiffsRaw[1:]

	var fd []FileDiff
	for _, fileDiffRaw := range fileDiffsRaw {
		h, err := extractHunks(fileDiffRaw)
		if err != nil {
			return GitDiff{}, fmt.Errorf("failed to extract h: %w", err)
		}

		fileDiff := FileDiff{
			OldFilename: extractOldFilename(fileDiffRaw),
			NewFilename: extractNewFilename(fileDiffRaw),
			Hunks:       h,
		}

		fd = append(fd, fileDiff)
	}

	return GitDiff{
		FileDiffs: fd,
	}, nil
}

func extractOldFilename(str string) string {
	i := strings.Index(str, "--- a/")
	i += len("--- a/")

	if i == -1 {
		return ""
	}

	j := strings.Index(str[i:], "\n")

	if j == -1 {
		return ""
	}

	return str[i : j+i]
}

func extractNewFilename(str string) string {
	i := strings.Index(str, "+++ b/")
	i += len("+++ b/")

	if i == -1 {
		return ""
	}

	j := strings.Index(str[i:], "\n")

	if j == -1 {
		return ""
	}

	return str[i : j+i]
}

var hunkHeaderRegex = regexp.MustCompile(`(?m)^\s*@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)

func extractHunks(str string) ([]Hunk, error) {
	var hunks []Hunk
	matches := hunkHeaderRegex.FindAllStringSubmatchIndex(str, -1)

	for i := 0; i < len(matches); i++ {
		oldLineStart, err := strconv.Atoi(str[matches[i][2]:matches[i][3]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse old line start: %w", err)
		}
		oldLineCount, err := strconv.Atoi(str[matches[i][4]:matches[i][5]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse old line count: %w", err)
		}
		newLineStart, err := strconv.Atoi(str[matches[i][6]:matches[i][7]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse new line start: %w", err)
		}
		newLineCount, err := strconv.Atoi(str[matches[i][8]:matches[i][9]])
		if err != nil {
			return nil, fmt.Errorf("failed to parse new line count: %w", err)
		}

		var hunkContent string
		if i+1 < len(matches) {
			hunkContent = str[matches[i][9]:matches[i+1][9]]
		} else {
			hunkContent = str[matches[i][0]:len(str)]
		}

		cl := extractChangedLines(hunkContent)
		ho := determineHunkOperation(cl)

		h := Hunk{
			HunkOperation:    ho,
			OldFileLineStart: oldLineStart,
			OldFileLineCount: oldLineCount,
			NewFileLineStart: newLineStart,
			NewFileLineCount: newLineCount,
			ChangedLines:     cl,
		}
		hunks = append(hunks, h)
	}

	return hunks, nil
}

func determineHunkOperation(cl []ChangedLine) HunkOperation {
	hasAdditions := false
	hasDeletions := false

	for _, l := range cl {
		if l.IsDeletion {
			hasDeletions = true
		} else {
			hasAdditions = true
		}

		if hasAdditions && hasDeletions {
			return MODIFY
		}
	}

	if hasAdditions {
		return ADD
	}
	if hasDeletions {
		return DELETE
	}

	return MODIFY
}

func extractChangedLines(str string) []ChangedLine {
	var cls []ChangedLine
	s := bufio.NewScanner(strings.NewReader(str))

	for s.Scan() {
		l := s.Text()

		if !strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "-") {
			continue
		}
		if strings.HasPrefix(l, "+") {
			cl := ChangedLine{
				IsDeletion: false,
				Content:    strings.TrimSpace(l[1:]),
			}
			cls = append(cls, cl)
		} else {
			cl := ChangedLine{
				IsDeletion: true,
				Content:    strings.TrimSpace(l[1:]),
			}
			cls = append(cls, cl)
		}
	}

	return cls
}

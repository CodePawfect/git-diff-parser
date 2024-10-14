package parse

import (
	"bufio"
	"fmt"
	"github.com/codepawfect/git-diff-parser/pkg/model"
	"regexp"
	"strconv"
	"strings"
)

func Parse(gitDiff string) (model.GitDiff, error) {
	fileDiffsRaw := strings.Split(gitDiff, "diff --git")
	fileDiffsRaw = fileDiffsRaw[1:]

	var fileDiffs []model.FileDiff
	for _, fileDiffRaw := range fileDiffsRaw {
		hunks, err := extractHunks(fileDiffRaw)
		if err != nil {
			return model.GitDiff{}, fmt.Errorf("failed to extract hunks: %w", err)
		}

		fileDiff := model.FileDiff{
			OldFilename: extractOldFilename(fileDiffRaw),
			NewFilename: extractNewFilename(fileDiffRaw),
			Hunks:       hunks,
		}

		fileDiffs = append(fileDiffs, fileDiff)
	}

	return model.GitDiff{
		FileDiffs: fileDiffs,
	}, nil
}

func extractOldFilename(str string) string {
	startIndex := strings.Index(str, "--- a/")
	startIndex += len("--- a/")

	if startIndex == -1 {
		return ""
	}

	endIndex := strings.Index(str[startIndex:], "\n")

	if endIndex == -1 {
		return ""
	}

	return str[startIndex : endIndex+startIndex]
}

func extractNewFilename(str string) string {
	startIndex := strings.Index(str, "+++ b/")
	startIndex += len("+++ b/")

	if startIndex == -1 {
		return ""
	}

	endIndex := strings.Index(str[startIndex:], "\n")

	if endIndex == -1 {
		return ""
	}

	return str[startIndex : endIndex+startIndex]
}

func extractHunks(str string) ([]model.Hunk, error) {
	var hunks []model.Hunk
	hunkHeaderRegex := regexp.MustCompile(`(?m)^\s*@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	//hunkHeaderRegex := regexp.MustCompile(`(?m)^@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
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

		changedLines := extractChangedLines(hunkContent)
		hunkOperation := determineHunkOperation(changedLines)

		hunk := model.Hunk{
			HunkOperation:    hunkOperation,
			OldFileLineStart: oldLineStart,
			OldFileLineCount: oldLineCount,
			NewFileLineStart: newLineStart,
			NewFileLineCount: newLineCount,
			ChangedLines:     changedLines,
		}
		hunks = append(hunks, hunk)
	}

	return hunks, nil
}

func determineHunkOperation(changedLines []model.ChangedLine) model.HunkOperation {
	hasAdditions := false
	hasDeletions := false

	for _, line := range changedLines {
		if line.IsDeletion {
			hasDeletions = true
		} else {
			hasAdditions = true
		}

		if hasAdditions && hasDeletions {
			return model.MODIFY
		}
	}

	if hasAdditions {
		return model.ADD
	}
	if hasDeletions {
		return model.DELETE
	}

	return model.MODIFY
}

func extractChangedLines(str string) []model.ChangedLine {
	var changedLines []model.ChangedLine
	scanner := bufio.NewScanner(strings.NewReader(str))

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "-") {
			continue
		}
		if strings.HasPrefix(line, "+") {
			changedLine := model.ChangedLine{
				IsDeletion: false,
				Content:    strings.TrimSpace(line[1:]),
			}
			changedLines = append(changedLines, changedLine)
		} else {
			changedLine := model.ChangedLine{
				IsDeletion: true,
				Content:    strings.TrimSpace(line[1:]),
			}
			changedLines = append(changedLines, changedLine)
		}
	}

	return changedLines
}

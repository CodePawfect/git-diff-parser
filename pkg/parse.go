package parse

import (
	"bufio"
	"git-diff-parser/pkg/model"
	"regexp"
	"strconv"
	"strings"
)

func parse(gitDiff string) model.GitDiff {
	fileDiffsRaw := strings.Split(gitDiff, "diff --git")
	fileDiffsRaw = fileDiffsRaw[1:]

	var fileDiffs []model.FileDiff
	for _, fileDiffRaw := range fileDiffsRaw {
		fileDiff := model.FileDiff{
			OldFilename: extractOldFilename(fileDiffRaw),
			NewFilename: extractNewFilename(fileDiffRaw),
			Hunks:       extractHunks(fileDiffRaw),
		}
		fileDiffs = append(fileDiffs, fileDiff)
	}

	return model.GitDiff{
		FileDiffs: fileDiffs,
	}
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

func extractHunks(str string) []model.Hunk {
	var hunks []model.Hunk
	hunkHeaderRegex := regexp.MustCompile(`(?m)^\s*@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	//hunkHeaderRegex := regexp.MustCompile(`(?m)^@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	matches := hunkHeaderRegex.FindAllStringSubmatchIndex(str, -1)

	for i := 0; i < len(matches); i++ {
		oldLineStart, _ := strconv.Atoi(str[matches[i][2]:matches[i][3]])
		oldLineCount, _ := strconv.Atoi(str[matches[i][4]:matches[i][5]])
		newLineStart, _ := strconv.Atoi(str[matches[i][6]:matches[i][7]])
		newLineCount, _ := strconv.Atoi(str[matches[i][8]:matches[i][9]])

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

	return hunks
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
	}
	if hasAdditions && hasDeletions {
		return model.MODIFY
	} else if hasAdditions {
		return model.ADD
	} else if hasDeletions {
		return model.DELETE
	} else {
		return model.MODIFY
	}
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

package shell

import (
	"strconv"
	"strings"
)

type Feature byte

const (
	Jobs Feature = iota
	Azure
	PoshGit
	LineError
	Tooltips
	Transient
	FTCSMarks
	Upgrade
	Notice
	PromptMark
	RPrompt
	CursorPositioning
	Async
)

type Features []Feature

func (f Features) Lines(shell string) Lines {
	var lines Lines

	for _, feature := range f {
		var code Code

		switch shell {
		case PWSH, PWSH5:
			code = feature.Pwsh()
		case ZSH:
			code = feature.Zsh()
		case BASH:
			code = feature.Bash()
		case ELVISH:
			code = feature.Elvish()
		case FISH:
			code = feature.Fish()
		case CMD:
			code = feature.Cmd()
		case NU:
			code = feature.Nu()
		case XONSH:
			code = feature.Xonsh()
		}

		if len(code) > 0 {
			lines = append(lines, code)
		}
	}

	return lines
}

func (f Features) Hash() string {
	var sb strings.Builder

	for _, feature := range f {
		sb.WriteString(strconv.Itoa(int(feature)))
	}

	return sb.String()
}

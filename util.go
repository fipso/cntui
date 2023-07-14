package main

import "strings"

func escapeShell(s string) string {
	return strings.Replace(s, "'", "'\\''", -1)
}

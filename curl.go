package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func exportAsCurl(req *request) {
	headers := ""

	headersMap, err := req.req.Headers.Map()
	if err != nil {
		panic(err)
	}
	for key, value := range headersMap {
		if strings.HasPrefix(key, "sec-") || strings.HasPrefix(key, ":") {
			continue
		}
		headers += fmt.Sprintf(" -H '%s: %s'", escapeShell(key), escapeShell(value))
	}

	body := ""
	if req.req.PostData != nil {
		body = fmt.Sprintf(" -d '%s'", escapeShell(*req.req.PostData))
	}

	method := ""
	if req.req.Method != "GET" && req.req.Method != "POST" {
		method = fmt.Sprintf(" -X %s", req.req.Method)
	}

	cmd := fmt.Sprintf("curl%s%s%s '%s'", headers, body, method, req.req.URL)
	clipboard.WriteAll(cmd)
}

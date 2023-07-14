package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func exportAsCurl(req *request) {
        // add -H
	headers := ""
	headersMap, err := req.req.Headers.Map()
	if err != nil {
		panic(err)
	}
	for key, value := range headersMap {
		// Skip accept-encoding header to prefer plaintext response
		if strings.ToLower(key) == "accept-encoding" {
			continue
		}
                // Drop these
		if strings.HasPrefix(key, "sec-") || strings.HasPrefix(key, ":") {
			continue
		}
		headers += fmt.Sprintf(" -H '%s: %s'", escapeShell(key), escapeShell(value))
	}

        // add -d
	body := ""
	if req.req.PostData != nil {
		body = fmt.Sprintf(" -d '%s'", escapeShell(*req.req.PostData))
	}

        // add -X
	method := ""
	if req.req.Method != "GET" && req.req.Method != "POST" {
		method = fmt.Sprintf(" -X %s", req.req.Method)
	}

        // assemble command
	cmd := fmt.Sprintf("curl%s%s%s '%s'", headers, body, method, req.req.URL)
	clipboard.WriteAll(cmd)
}

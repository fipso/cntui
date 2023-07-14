package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/rpcc"
)

var requests []request
var requestsUpdated chan bool

type request struct {
	time      time.Time
	id        network.RequestID
	initiator network.Initiator
	req       network.Request
	res       *network.Response
}

func main() {
	requestsUpdated = make(chan bool)

	connectChrome()
	setupTUI()
}

func connectChrome() {
	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(context.Background(), devtool.Page)
	if err != nil {
		pt, err = devt.Create(context.Background())
		if err != nil {
			panic(err)
		}
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(context.Background(), pt.WebSocketDebuggerURL)
	if err != nil {
		panic(err)
	}
	//defer conn.Close() // Leaving connections open will leak memory.

	c := cdp.NewClient(conn)
	c.Network.Enable(context.Background(), &network.EnableArgs{})

	reqWillBeSent, err := c.Network.RequestWillBeSent(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			reply, err := reqWillBeSent.Recv()
			if err != nil {
				panic(err)
			}

			if !strings.HasPrefix(reply.Request.URL, "http") {
				continue
			}

			requests = append(
				requests,
				request{
					id:        reply.RequestID,
					time:      reply.WallTime.Time(),
					initiator: reply.Initiator,
					req:       reply.Request,
				},
			)

			requestsUpdated <- true
		}
	}()

	resReceived, err := c.Network.ResponseReceived(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			reply, err := resReceived.Recv()
			if err != nil {
				panic(err)
			}

			for i, req := range requests {
				if req.id == reply.RequestID {
					requests[i].res = &reply.Response
					requestsUpdated <- true
					break
				}

			}

		}
	}()

	reqWillBeSentExtraInfo, err := c.Network.RequestWillBeSentExtraInfo(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			reply, err := reqWillBeSentExtraInfo.Recv()
			if err != nil {
				panic(err)
			}

			newHeaders, err := reply.Headers.Map()
			if err != nil {
				panic(err)
			}

			for i, req := range requests {
				if req.id == reply.RequestID {
					headers, err := req.req.Headers.Map()
					if err != nil {
						panic(err)
					}

				outer:
					for key, value := range newHeaders {
						for existingKey := range headers {
							existingKeyLower := strings.ToLower(existingKey)
							if existingKeyLower == key {
								headers[existingKey] = value
								continue outer
							}
						}
						headers[key] = value
					}

					newJson, err := json.Marshal(headers)
					if err != nil {
						panic(err)
					}

					requests[i].req.Headers = newJson

					requestsUpdated <- true
					break
				}

			}
		}
	}()
}

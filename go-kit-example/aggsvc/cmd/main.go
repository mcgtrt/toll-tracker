package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/mcgtrt/toll-tracker/go-kit-example/aggsvc/aggendpoint"
	"github.com/mcgtrt/toll-tracker/go-kit-example/aggsvc/aggservice"
	"github.com/mcgtrt/toll-tracker/go-kit-example/aggsvc/aggtransport"
)

func main() {
	httpAddr := flag.String("httpAddr", ":3031", "http serv er listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	var (
		svc         = aggservice.New()
		endpoints   = aggendpoint.New(svc, logger)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)

	logger.Log("transport", "HTTP", "addr", *httpAddr)
	http.Serve(httpListener, httpHandler)
}

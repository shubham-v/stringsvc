package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"stringsvc/middlewares/instruments"
	logging "stringsvc/middlewares/loggers"
	"stringsvc/middlewares/proxys"
	"stringsvc/services"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"

	"stringsvc/transports"
	"stringsvc/transports/requests"
	transportUtils "stringsvc/transports/utils"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	var Logger log.Logger
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "listen", *listen, "caller", log.DefaultCaller)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var svc services.StringService
	svc = services.StringServiceType{}
	svc = proxys.ProxyingMiddleware(context.Background(), *proxy, Logger)(svc)
	svc = logging.LoggingMiddleware{Logger, svc}
	svc = instruments.InstrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	uppercaseHandler := httptransport.NewServer(
		transports.MakeUppercaseEndpoint(svc),
		requests.DecodeUppercaseRequest,
		transportUtils.EncodeResponse,
	)
	countHandler := httptransport.NewServer(
		transports.MakeCountEndpoint(svc),
		requests.DecodeCountRequest,
		transportUtils.EncodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	Logger.Log("msg", "HTTP", "addr", *listen)
	Logger.Log("err", http.ListenAndServe(*listen, nil))
}
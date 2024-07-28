package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/time/rate"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()
	// set default logger to json
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	limitDuration := time.Minute * 1
	if debug {
		limitDuration = time.Second * 1
		slogger.Debug("Debug mode enabled.")
	}
	limiter := rate.NewLimiter(rate.Every(limitDuration), 1)
	for {
		if err := limiter.Wait(context.Background()); err != nil {
			slog.Error("Failed to wait.", "error", err)
		}
		pinger, err := probing.NewPinger("8.8.4.4")
		if err != nil {
			slogger.Error("Failed while constructing new pinger.", "error", err)
		}
		pinger.Count = 1
		if err := pinger.Run(); err != nil {
			slog.Error("Failed to ping.", "error", err)
		}
		stats := pinger.Statistics()

		slogger.Info("ping successful",
			"address", stats.Addr,
			"packets_transmitted", stats.PacketsSent,
			"received", stats.PacketsRecv,
			"packet_loss", stats.PacketLoss,
			"time_ns", stats.AvgRtt)
	}
}

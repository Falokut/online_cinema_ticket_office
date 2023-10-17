package shutdown

import (
	"io"
	"os"
	"os/signal"

	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/logging"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	logger.Infof("Caught signal %s. Shutting down...", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
	logger.Writer().Close()
}

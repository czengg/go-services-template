package jobs

import "template/internal/logger"

func FakeJob(logger logger.Logger) {
	logger.Info("fake job running")
}

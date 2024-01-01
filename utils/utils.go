package utils

import (
	"os"

	"github.com/benharmonics/personal-site-backend/logging"
)

func Must(err error) {
	if err != nil {
		logging.Error("Fatal error:", err)
		os.Exit(1)
	}
}

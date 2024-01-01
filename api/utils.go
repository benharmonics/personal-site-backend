package api

import (
	"net/http"

	"github.com/benharmonics/backend/logging"
	"github.com/benharmonics/backend/utils/web"
)

func logAndEmitHTTPError(w http.ResponseWriter, r *http.Request, statusCode int, messages ...string) {
	logging.HTTPError(r, statusCode)
	web.HTTPError(w, statusCode, messages...)
}

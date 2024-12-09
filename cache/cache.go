package cache

import (
	"net/http"
	"time"
)

type CacheDB struct {
	Response     *http.Response
	ResponseBody []byte
	Created      time.Time
}

package server

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/airmap/tegola/cache"
)

//	TileCacheHandler implements a request cache for tiles on requests when the URLs
//	have a /:z/:x/:y scheme suffix (i.e. /osm/1/3/4.pbf)
func TileCacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		//	check if a cache backend exists
		cacher := Atlas.GetCache()
		if cacher == nil {
			//	nope. move on
			next.ServeHTTP(w, r)
			return
		}

		//	parse our URI into a cache key structure (pop off the "maps/" prefix)
		//	5 is the value of len("maps/")
		key, err := cache.ParseKey(r.URL.Path[5:])
		if err != nil {
			log.Println("cache middleware: ParseKey err: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		//	use the URL path as the key
		cachedTile, hit, err := cacher.Get(key)
		if err != nil {
			log.Printf("cache middleware: error reading from cache: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		//	cache miss
		if !hit {
			//	buffer which will hold a copy of the response for writing to the cache
			var buff bytes.Buffer

			//	ovewrite our current responseWriter with a tileCacheResponseWriter
			w = newTileCacheResponseWriter(w, &buff)

			next.ServeHTTP(w, r)

			//	check if our request context has been canceled
			if r.Context().Err() != nil {
				return
			}

			if err := cacher.Set(key, buff.Bytes()); err != nil {
				log.Println("cache response writer err: %v", err)
			}
			return
		}

		//	TODO: how configurable do we want the CORS policy to be?
		//	set CORS header
		w.Header().Add("Access-Control-Allow-Origin", "*")

		//	mimetype for protocol buffers
		w.Header().Add("Content-Type", "application/x-protobuf")

		//	communicate the cache is being used
		w.Header().Add("Tegola-Cache", "HIT")

		w.Write(cachedTile)
		return
	})
}

func newTileCacheResponseWriter(resp http.ResponseWriter, w io.Writer) http.ResponseWriter {
	return &tileCacheResponseWriter{
		resp:  resp,
		multi: io.MultiWriter(w, resp),
	}
}

//	tileCacheResponsWriter wraps http.ResponseWriter (https://golang.org/pkg/net/http/#ResponseWriter)
//	to additionally write the response to a cache when there is a cache MISS
type tileCacheResponseWriter struct {
	resp  http.ResponseWriter
	multi io.Writer
}

func (w *tileCacheResponseWriter) Header() http.Header {
	//	communicate the tegola cache is being used
	w.resp.Header().Set("Tegola-Cache", "MISS")

	return w.resp.Header()
}

func (w *tileCacheResponseWriter) Write(b []byte) (int, error) {
	//	write to our multi writer
	return w.multi.Write(b)
}

func (w *tileCacheResponseWriter) WriteHeader(i int) {
	w.resp.WriteHeader(i)
}

package server

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func createProxy(target string, s *Server, r *http.Request) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	reqUrl := fmt.Sprintf("%s:%s", url.Host, r.URL.RequestURI())
	// slog.Info("redis key", "reqUrl", reqUrl)
	if err != nil {
		slog.Error("failed parsing url", "error", err.Error())
	}

	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()
			pr.SetURL(url)
		},
		ModifyResponse: HandleResponse(s, reqUrl),
	}
}

// DecompressBody handles gzip and zstd decompression
func DecompressBody(r *http.Response) (io.ReadCloser, error) {
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		return gzipReader, nil
	case "zstd":
		zstdReader, err := zstd.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(zstdReader)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewBuffer(data)), nil
		// return zstdReader.IOReadCloser(), nil
	default:
		return r.Body, nil
	}
}

// HandleResponse handles the modification of the response
func HandleResponse(s *Server, reqUrl string) func(*http.Response) error {
	return func(r *http.Response) error {

		if r.StatusCode != http.StatusOK {
			log.Printf("Error: Received status code %d", r.StatusCode)
			return errors.New("failed call")
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			log.Printf("Error: Expected JSON response but got Content-Type: %s", r.Header.Get("Content-Type"))
			return errors.New("invalid content type")
		}
		reader, err := DecompressBody(r)
		if err != nil {
			log.Printf("Error decompressing response body: %v", err)
			return err
		}
		defer reader.Close()
		data, err := io.ReadAll(reader)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return err
		}
		// log.Printf("Raw response body: %s", string(data))

		// Remove Content-Encoding header to indicate decompression
		r.Header.Del("Content-Encoding")
		r.Header.Set("Content-Length", fmt.Sprint(len(data)))

		// Set the body to the decompressed data
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		s.db.SetString(reqUrl, string(data))
		return nil
	}

}

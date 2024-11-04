package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/link_generator"
	"main/internal/redis_storage"
	"main/internal/storage"
	"main/internal/utils"
	"main/lib/api"
	"main/lib/api/response"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	store         storage.Storage               // DI
	linkGenerator *link_generator.LinkGenerator //DI
	linkLength    = 10                          // default TODO: move to configuration default value
	publicURL     = "http://localhost:3000"     // default
)

// TODO: move info cmd to separate file (?)
func init() {
	if lenEnv, err := strconv.Atoi(os.Getenv("LINK_LENGTH")); err == nil {
		linkLength = lenEnv
	}

	if publicUrlEnv := os.Getenv("PUBLIC_URL"); publicUrlEnv != "" {
		publicURL = publicUrlEnv
	}

	store = redis_storage.NewRedisStorage(os.Getenv("REDIS_URL"))
	linkGenerator = link_generator.NewLinkGenerator(store)
}

// Check that alias for url is not used already.
// POST https://clck.knsrg.com/check with body {"alias": "alias"}
func fetchHandler(w http.ResponseWriter, r *http.Request) {
	requestID := utils.RandomRequestID()
	clientIP := utils.GetClientIP(r)

	log.Printf("[INFO] [%s] Received request from IP: %s for fetching", requestID, clientIP)

	ctx := context.Background()

	var req struct {
		Alias string `json:"alias"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Alias == "" {
		errMsg := fmt.Sprintf("Invalid request payload: %v", err)
		resp := response.Error(errMsg)

		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	code := req.Alias

	aliasExists, err := store.Exists(ctx, code)

	if err != nil {
		respErrorMsg := fmt.Sprintf("Failed to check alias: %v", err)
		resp := response.Error(respErrorMsg)

		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	// if exists return 200 else 404 without body
	if aliasExists {
		resp := response.OK()

		api.RespondWithOK(w, requestID, resp, http.StatusOK)
	} else {
		resp := response.Error("Alias not found")

		api.RespondWithError(w, requestID, resp, http.StatusNotFound)
	}
}

// Create short link handler.
// POST https://clck.knsrg.com
// TODO: add basic auth (middleware?)
func encodeHandler(w http.ResponseWriter, r *http.Request) {
	requestID := utils.RandomRequestID()
	clientIP := utils.GetClientIP(r)

	log.Printf("[INFO] [%s] Received request from IP: %s for encoding", requestID, clientIP)

	var req struct {
		Link         string `json:"link"`
		LifetimeDays *int   `json:"lifetime_days,omitempty"`
		Alias        string `json:"alias"`
	}

	// TODO: add custom validator for validate request schema (???)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Link == "" {
		errMsg := fmt.Sprintf("Invalid request payload: %v", err)
		resp := response.Error(errMsg)
		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	var lifetimeSeconds int
	if req.LifetimeDays != nil {
		seconds := *req.LifetimeDays * 24 * 60 * 60
		lifetimeSeconds = seconds
	} else {
		seconds := 0
		lifetimeSeconds = seconds
	}

	ctx := context.Background()

	code, err := linkGenerator.GenerateLinkAlias(ctx, req.Link, lifetimeSeconds, req.Alias, linkLength)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate link code: %v", err)
		resp := response.Error(errMsg)

		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	encodedLink, err := url.JoinPath(publicURL, code)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate link: %v", err)
		resp := response.Error(errMsg)

		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	resp := response.OKWithEncodedLink(encodedLink)
	api.RespondWithOK(w, requestID, resp, http.StatusOK)

	log.Printf("[INFO] [%s] Link encoded successfully. Encoded link: %s", requestID, encodedLink)
}

// Redirect to original link handler.
// GET https://clck.knsrg.com/{code}
func decodeHandler(w http.ResponseWriter, r *http.Request) {
	requestID := utils.RandomRequestID()
	clientIP := utils.GetClientIP(r)

	log.Printf("[INFO] [%s] Received request from IP: %s for decoding", requestID, clientIP)

	code := r.URL.Path[len("/"):]
	if code == "" {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()
	link, err := store.Get(ctx, code)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to get link: %v", err)
		resp := response.Error(errMsg)
		api.RespondWithError(w, requestID, resp, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}

// TODO: add external http router for better routing (???)
func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler

	p := r.URL.Path

	switch {
	case r.Method == `POST`:
		if match(p, `/`) {
			h = http.HandlerFunc(encodeHandler)
		}

		if match(p, `/check`) {
			h = http.HandlerFunc(fetchHandler)
		}
	case r.Method == `GET`:
		h = http.HandlerFunc(decodeHandler)
	default:
		http.NotFound(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", Serve)

	// TODO: add configuration pattern ?
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

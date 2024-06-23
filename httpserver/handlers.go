package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/chackett/dating-service/datingservice"
	"github.com/chackett/dating-service/rankingservice"
	"github.com/chackett/dating-service/repository"
	"log/slog"
	"net/http"
	"os"
)

const (
	maxRequestBodySizeBytes = 1048576
	ctxKeySessionUserID     = "session_user_id"
)

// handler defines functionality for exposing routes via HTTP and also parsing the messages before passing onto the relevant
// service.
type handler struct {
	dateService *datingservice.DateService
	logger      *slog.Logger
	mux         http.Handler
	routes      map[string]routeConfig
}

// routeConfig stores an HTTP route and any config related to it. i.e. Authenticate it or not.
// For a basic ACL you could add access level here.
type routeConfig struct {
	authUser bool
	handler  func(http.ResponseWriter, *http.Request)
}

// newHandler creates and initialises the handler/routes.
func newHandler(ds *datingservice.DateService) (*handler, error) {
	if ds == nil {
		return nil, errors.New("datingservice is nil")
	}

	result := &handler{
		dateService: ds,
		logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	result.routes = map[string]routeConfig{
		"POST /user/create": {
			authUser: false,
			handler:  result.handlePOSTCreateUser,
		},
		"POST /user/preferences": {
			authUser: true,
			handler:  result.handlePOSTUserPreferences,
		},
		"POST /login": {
			authUser: false,
			handler:  result.handlePOSTLogin,
		},
		"GET /discover": {
			authUser: true,
			handler:  result.handleGETDiscover,
		},
		"POST /swipe": {
			authUser: true,
			handler:  result.handlePOSTSwipe,
		},
	}
	return result, nil
}

// setupRoutes applies the route configs to a HTTP mux/handler for serving.
func (h *handler) setupRoutes(middlewares []func(h http.Handler) http.Handler) {
	mux := http.NewServeMux()

	for rp, rc := range h.routes {
		h.logger.Debug("set up route: %s", rp)
		mux.HandleFunc(rp, rc.handler)
	}

	h.mux = http.NewServeMux()

	for _, mw := range middlewares {
		h.mux = mw(mux)
	}
}

// handlePOSTCreateUser handles requests to create new user
func (h *handler) handlePOSTCreateUser(w http.ResponseWriter, r *http.Request) {
	u := repository.User{}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySizeBytes)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.logger.Error("decode create user message: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	createdUser, err := h.dateService.CreateUser(r.Context(), u)
	if err != nil {
		h.logger.Error("create user: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	btsUser, err := json.Marshal(createdUser)
	if err != nil {
		h.logger.Error("marshal created user: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writePlainResponse(w, http.StatusCreated, string(btsUser))
}

// handlePOSTUserPreferences handle request to set user preferences
func (h *handler) handlePOSTUserPreferences(w http.ResponseWriter, r *http.Request) {
	input := repository.UserPreferences{}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySizeBytes)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decode create user preferences message: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionUserID, ok := r.Context().Value(ctxKeySessionUserID).(int)
	if !ok {
		h.writePlainResponse(w, http.StatusBadRequest, "invalid user")
		return
	}

	if sessionUserID != input.UserID {
		h.logger.Info("unauthorized attempt by user `%d` to update preferences for other user `%d`", sessionUserID, input.UserID)
		h.writePlainResponse(w, http.StatusBadRequest, "logged in user mismatch with preference request")
		return
	}

	err = h.dateService.SetUserPreferences(r.Context(), input)
	if err != nil {
		h.logger.Error("create user preferences: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writePlainResponse(w, http.StatusCreated, "")
}

// handlePOSTLogin handle requests to create authenticated session (i.e. Login)
// Once this is successfully called with a valid username/password combination, then a token is returned which can be used
// against subsequent authenticated HTTP calls.
func (h *handler) handlePOSTLogin(w http.ResponseWriter, r *http.Request) {
	// Use anonymous struct as login messages are pretty much isolated to this function
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decode login message", err)
		h.writePlainResponse(w, http.StatusUnauthorized, "incorrect email / password combination")
		return
	}

	token, err := h.dateService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		h.logger.Error("date service login attempt", err)
		h.writePlainResponse(w, http.StatusUnauthorized, "incorrect email / password combination")
		return
	}

	// Use anonymous struct type as this message type is not used elsewhere.
	tokenResponse := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	btsResp, err := json.Marshal(tokenResponse)
	if err != nil {
		h.logger.Error("marshal login token message to JSON", err)
		h.writePlainResponse(w, http.StatusUnauthorized, "incorrect email / password combination")
		return
	}

	h.writeJSONResponse(w, http.StatusAccepted, string(btsResp))
}

// handleGETDiscover a handler for requests to discover matched candidates
func (h *handler) handleGETDiscover(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value(ctxKeySessionUserID).(int)
	if !ok {
		h.writePlainResponse(w, http.StatusBadRequest, "invalid user")
		return
	}

	matches, err := h.dateService.Discover(r.Context(), sessionUserID)
	if err != nil {
		h.writePlainResponse(w, http.StatusInternalServerError, "an error has occurred")
		return
	}

	resp := struct {
		Results []rankingservice.RankedMatch `json:"results"`
	}{
		Results: matches.Matches,
	}

	btsResp, err := json.Marshal(resp)
	if err != nil {
		h.writePlainResponse(w, http.StatusInternalServerError, "an error has occurred")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, string(btsResp))
}

// handlePOSTSwipe handle requests from users where they are voting on a candidate.
func (h *handler) handlePOSTSwipe(w http.ResponseWriter, r *http.Request) {
	input := repository.Swipe{}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decode swipe message", err)
		h.writePlainResponse(w, http.StatusBadRequest, "unable to parse swipe message")
		return
	}

	sessionUserID, ok := r.Context().Value(ctxKeySessionUserID).(int)
	if !ok {
		h.writePlainResponse(w, http.StatusBadRequest, "invalid user")
		return
	}

	if sessionUserID != input.UserID {
		h.logger.Info("unauthorized swipe attempt by user `%d` for other user `%d`", sessionUserID, input.CandidateID)
		h.writePlainResponse(w, http.StatusBadRequest, "logged in user mismatch with swipe message")
		return
	}

	match, err := h.dateService.Swipe(r.Context(), input)
	if err != nil {
		if errors.As(datingservice.ErrDuplicateSwipe, &err) {
			h.writePlainResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writePlainResponse(w, http.StatusInternalServerError, "unable to submit swipe message")
		return
	}

	// TODO: This is all a bit janky, admittedly. Tidy up.
	type subResults struct {
		Matched bool `json:"matched"`
		MatchID int  `json:"matchID,omitempty"`
	}

	var matchedCandidateID int
	if match {
		matchedCandidateID = input.CandidateID
	}

	result := struct {
		Results subResults `json:"results"`
	}{
		Results: subResults{
			Matched: match,
			MatchID: matchedCandidateID,
		},
	}

	btsResult, err := json.Marshal(result)
	if err != nil {
		h.writePlainResponse(w, http.StatusInternalServerError, "an error has occurred")
		return
	}
	h.writeJSONResponse(w, http.StatusCreated, string(btsResult))
}

// writeJSONResponse a helper function to reduce duplicated code to return a JSON message.
func (h *handler) writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		h.logger.Error("unable to write http response", err)
	}
}

// writePlainResponse a helper function to reduce duplicated code to return a plain text message.
func (h *handler) writePlainResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		h.logger.Error("unable to write http response", err)
	}
}

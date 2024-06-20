package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/chackett/dating-service/datingservice"
	"github.com/chackett/dating-service/repository"
	"log/slog"
	"net/http"
	"os"
)

const (
	maxRequestBodySize  = 1048576
	ctxKeySessionUserID = "session_user_id"
)

type handler struct {
	dateService *datingservice.DateService
	logger      *slog.Logger
	mux         http.Handler
	routes      map[string]routeConfig
}

type routeConfig struct {
	authUser bool
	handler  func(http.ResponseWriter, *http.Request)
}

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

func (h *handler) handlePOSTCreateUser(w http.ResponseWriter, r *http.Request) {
	u := repository.User{}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.logger.Error("decode create user message: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	createdUser, err := h.dateService.CreateUser(r.Context(), u)
	if err != nil {
		h.logger.Error("create user: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, "unable to create user")
		return
	}

	btsUser, err := json.Marshal(createdUser)
	if err != nil {
		h.logger.Error("marshal created user: %w", err)
		h.writePlainResponse(w, http.StatusInternalServerError, "user created, but error returning to caller")
		return
	}

	h.writePlainResponse(w, http.StatusCreated, string(btsUser))
}
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
		Results []repository.User `json:"results"`
	}{
		Results: matches,
	}

	btsResp, err := json.Marshal(resp)
	if err != nil {
		h.writePlainResponse(w, http.StatusInternalServerError, "an error has occurred")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, string(btsResp))
}
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
		h.logger.Error("date service swipe: %w", err)
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
		Results subResults
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

func (h *handler) writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		h.logger.Error("unable to write http response", err)
	}
}

func (h *handler) writePlainResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		h.logger.Error("unable to write http response", err)
	}
}

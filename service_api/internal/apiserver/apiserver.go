package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SedovSG/zaplog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"service_api/internal/config"
	"service_api/internal/metrics"
	"service_api/internal/model"
	"service_api/internal/storage"
	"service_api/internal/storage/cache"
	"service_api/internal/storage/pgsql"
	"strconv"
)

// APIServer структура компонентов приложения
type APIServer struct {
	config    *config.Config
	router    *mux.Router
	routerApi *mux.Router
	storage   storage.IStorageSQL
	cache     storage.ICache
}

func New(ctx context.Context, cfg *config.Config) (*APIServer, error) {
	postgres, err := pgsql.NewPostgresStorage(ctx, cfg.DBSQLConnection)
	if err != nil {
		return nil, errors.Wrap(err, "fail to init postgresql database")
	}
	c, err := cache.NewCacheInMemory(cfg.Cache.Size, cfg.Cache.TtlSecond)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize cache")
	}
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("").Subrouter()
	server := APIServer{
		config:    cfg,
		router:    router,
		routerApi: apiRouter,
		storage:   postgres,
		cache:     c,
	}
	return &server, nil
}

func (s *APIServer) Start() error {
	s.configureRouterApi()
	zaplog.Throw().Info("starting api server")
	return http.ListenAndServe(":"+s.config.Server.Port, s.router)
}

func (s *APIServer) MiddlewareHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		xApiKey := r.Header.Get("X-API-KEY")
		w.Header().Set("Content-Type", "application/json")
		if xApiKey == s.config.Server.ApiKey {
			next.ServeHTTP(w, r)
			return
		}
		errorResponse(w, http.StatusUnauthorized, "Invalid api key")
		return
	})
}

func (s *APIServer) MiddlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health/readiness" || r.URL.Path == "/api/health/live" {
			next.ServeHTTP(w, r)
			return
		}
		defer metrics.ObserveRequestDurationSeconds(r.URL.Path)()
		metrics.IncRequestsTotal(metrics.OkStatus, r.URL.Path)
		next.ServeHTTP(w, r)

	})
}

func (s *APIServer) configureRouterApi() {
	s.router.Use(s.MiddlewareMetrics)
	s.routerApi.Use(s.MiddlewareHeaders)

	s.routerApi.HandleFunc("", s.getProducts).Methods("GET")
	s.router.HandleFunc("/api/health/live", s.Health).Methods("GET")
	s.router.HandleFunc("/api/health/readiness", s.Health).Methods("GET")
	s.router.Handle("/metrics", promhttp.Handler()).Methods("GET")

}

// getProducts метод обработки GET запроса на получение списка продуктов
func (s *APIServer) getProducts(w http.ResponseWriter, r *http.Request) {
	var queryData model.Params
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		lim, err := strconv.Atoi(limit)
		queryData.Limit = lim
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		errorResponse(w, http.StatusBadRequest, "missing parameter limit")
		return
	}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		off, err := strconv.Atoi(offset)
		queryData.Offset = off
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		errorResponse(w, http.StatusBadRequest, "missing parameter offset")
		return
	}
	result, err := s.cache.Get(r.Context(), queryData, s.storage)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	var res []model.Product
	for _, val := range result {
		res = append(res, model.Product{Id: val.Id, Price: val.Price})
	}

	successResponse(w, res)
	return

}

// Health метод обработки GET запроса на health приложения
func (s *APIServer) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Healthy"))
	if err != nil {
		zaplog.Throw().Error(err.Error())
	}
}

// errorResponse функция формирования ответа об ошибке
func errorResponse(w http.ResponseWriter, code int, errorText string) {
	w.WriteHeader(code)
	zaplog.Throw().Error(fmt.Sprintf("error response code %d, %s", code, errorText))
	jsonResponse, jsonError := json.Marshal(model.Error{Error: errorText})
	if jsonError != nil {
		zaplog.Throw().Error(jsonError.Error())
		return
	}
	_, err := w.Write(jsonResponse)
	if err != nil {
		zaplog.Throw().Error(err.Error())
	}
	return
}

// successResponse функция формирования успешного ответа
func successResponse(w http.ResponseWriter, body any) {
	jsonResponse, jsonError := json.Marshal(body)
	if jsonError != nil {
		zaplog.Throw().Error(jsonError.Error())
		return
	}
	_, err := w.Write(jsonResponse)
	if err != nil {
		zaplog.Throw().Error(err.Error())
	}
}

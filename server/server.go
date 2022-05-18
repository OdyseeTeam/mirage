package http

import (
	"context"
	"net/http"
	"time"

	"github.com/OdyseeTeam/mirage/internal/metrics"
	"github.com/OdyseeTeam/mirage/metadata"
	"github.com/OdyseeTeam/mirage/optimizer"
	"github.com/bluele/gcache"

	"github.com/OdyseeTeam/gody-cdn/store"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/lbryio/lbry.go/v2/extras/stop"
	log "github.com/sirupsen/logrus"
)

// Server is an instance of a peer server that houses the listener and store.
type Server struct {
	grp             *stop.Group
	optimizer       *optimizer.Optimizer
	cache           store.ObjectStore
	metadataManager *metadata.Manager
	errorCache      gcache.Cache
}

// NewServer returns an initialized Server pointer.
func NewServer(optimizer *optimizer.Optimizer, cache store.ObjectStore, metadataManager *metadata.Manager) *Server {
	return &Server{
		grp:             stop.New(),
		optimizer:       optimizer,
		cache:           cache,
		metadataManager: metadataManager,
		errorCache:      gcache.New(10000).Expiration(2 * time.Minute).Build(),
	}
}

// Shutdown gracefully shuts down the peer server.
func (s *Server) Shutdown() {
	log.Debug("shutting down HTTP server")
	s.grp.StopAndWait()
	log.Debug("HTTP server stopped")
}

// Start starts the server listener to handle connections.
func (s *Server) Start(address string) error {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(s.ErrorHandle())
	router.Use(nice.Recovery(s.recoveryHandler))
	router.Use(s.addCSPHeaders)
	metrics.InstallRoute(router)
	//https://thumbnails.odycdn.com/optimize/s:100:0/quality:85/plain/https://thumbnails.lbry.com/UCX_t3BvnQtS5IHzto_y7tbw
	router.GET("/optimize/:dimensions/quality:quality/plain/*url", s.optimizeHandler)
	router.GET("/optimize/plain/*url", s.simpleRedirect)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}
	go s.listenForShutdown(srv)
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	s.grp.Add(1)
	go func() {
		defer s.grp.Done()
		log.Println("HTTP server listening on " + address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return nil
}

func (s *Server) listenForShutdown(listener *http.Server) {
	<-s.grp.Ch()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := listener.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

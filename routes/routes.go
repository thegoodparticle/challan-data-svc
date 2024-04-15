package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	grpcclient "github.com/thegoodparticle/challan-data-svc/grpc-client"
	handler "github.com/thegoodparticle/challan-data-svc/rest-handler"
	"github.com/thegoodparticle/challan-data-svc/store"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter(timeout int) *Router {
	return &Router{
		config: NewConfig().SetTimeout(timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository store.Interface, grpcClient *grpcclient.GRPCClient) *chi.Mux {
	r.setConfigsRouters()

	r.RouterVehicle(repository, grpcClient)

	return r.router
}

func (r *Router) setConfigsRouters() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()
}

func (r *Router) RouterVehicle(repository store.Interface, grpcClient *grpcclient.GRPCClient) {
	handler := handler.NewHandler(repository, grpcClient)

	r.router.Route("/", func(route chi.Router) {
		route.Get("/health-check", handler.HealthCheck)
		route.Get("/", handler.HealthCheck)
	})

	r.router.Route("/challan-info", func(route chi.Router) {
		// route.Post("/", handler.Post)
		// route.Get("/", handler.Get)
		route.Get("/{RegID}", handler.Get)
		// route.Put("/{RegID}", handler.Put)
		// route.Delete("/{RegID}", handler.Delete)
		// route.Options("/", handler.Options)
	})
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}

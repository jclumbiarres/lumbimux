package lumbimux

import (
	"log"
	"net/http"
)

// Middleware es una función que toma un http.HandlerFunc y devuelve otro http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// RouteKey es un tipo personalizado para la clave del mapa de rutas.
type RouteKey struct {
	Method string
	Path   string
}

// Router es una interfaz que define los métodos para registrar rutas y controladores HTTP.
type Router interface {
	GET(path string, handler http.HandlerFunc, middlewares ...Middleware)
	POST(path string, handler http.HandlerFunc, middlewares ...Middleware)
	PUT(path string, handler http.HandlerFunc, middlewares ...Middleware)
	DELETE(path string, handler http.HandlerFunc, middlewares ...Middleware)
	PATCH(path string, handler http.HandlerFunc, middlewares ...Middleware)
	OPTIONS(path string, handler http.HandlerFunc, middlewares ...Middleware)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// LumbiMuxRouter es una estructura que implementa la interfaz Router.
type LumbiMuxRouter struct {
	routes      map[RouteKey]http.HandlerFunc
	middlewares map[RouteKey][]Middleware
}

// NewLumbiMux crea una nueva instancia de LumbiMuxRouter.
func NewLumbiMux() Router {
	return &LumbiMuxRouter{
		routes:      make(map[RouteKey]http.HandlerFunc),
		middlewares: make(map[RouteKey][]Middleware),
	}
}

// GET registra una ruta y controlador HTTP GET.
func (r *LumbiMuxRouter) GET(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("GET", path, handler, middlewares...)
}

// POST registra una ruta y controlador HTTP POST.
func (r *LumbiMuxRouter) POST(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("POST", path, handler, middlewares...)
}

// PUT registra una ruta y controlador HTTP PUT.
func (r *LumbiMuxRouter) PUT(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("PUT", path, handler, middlewares...)
}

// DELETE registra una ruta y controlador HTTP DELETE.
func (r *LumbiMuxRouter) DELETE(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("DELETE", path, handler, middlewares...)
}

// PATCH registra una ruta y controlador HTTP PATCH.
func (r *LumbiMuxRouter) PATCH(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("PATCH", path, handler, middlewares...)
}

// OPTIONS registra una ruta y controlador HTTP OPTIONS.
func (r *LumbiMuxRouter) OPTIONS(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.registerRoute("OPTIONS", path, handler, middlewares...)
}

// ServeHTTP maneja las solicitudes HTTP entrantes mediante la coincidencia del método de solicitud y la ruta de URL con un controlador registrado.
func (r *LumbiMuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := RouteKey{Method: req.Method, Path: req.URL.Path}
	if handler, ok := r.routes[key]; ok {
		for _, middleware := range r.middlewares[key] {
			handler = middleware(handler)
		}
		handler(w, req)
		return
	}
	http.NotFound(w, req)
}

// registerRoute registra una nueva ruta y controlador con el método HTTP especificado.
func (r *LumbiMuxRouter) registerRoute(method string, path string, handler http.HandlerFunc, middlewares ...Middleware) {
	key := RouteKey{Method: method, Path: path}
	r.routes[key] = anadeMiddleware(handler, middlewares...)
	r.middlewares[key] = middlewares
}

// anadeMiddleware toma un http.HandlerFunc y una lista de Middleware y devuelve un nuevo http.HandlerFunc que aplica los middlewares a la función original.
func anadeMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// helloHandler es un controlador HTTP simple que escribe "¡Hola, mundo!" en la respuesta.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("¡Hola, mundo!"))
}

// loggingMiddleware es un middleware que registra información sobre la solicitud HTTP.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	}
}

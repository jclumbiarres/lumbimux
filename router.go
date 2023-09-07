package lumbimux

import (
	"log"
	"net/http"
)

// Middleware es una función que toma un http.HandlerFunc y devuelve otro http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Router es una interfaz que define los métodos para registrar rutas y handlers HTTP.
type Router interface {
	GET(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	POST(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	PUT(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	DELETE(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	PATCH(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	OPTIONS(ruta string, handler http.HandlerFunc, middlewares ...Middleware)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// LumbiMuxRouter es una estructura que implementa la interfaz Router.
type LumbiMuxRouter struct {
	reglas      map[string]map[string]http.HandlerFunc
	middlewares map[string][]Middleware
}

// NewLumbiMux crea una nueva instancia de LumbiMuxRouter.
func NewLumbiMux() Router {
	return &LumbiMuxRouter{
		reglas:      make(map[string]map[string]http.HandlerFunc),
		middlewares: make(map[string][]Middleware),
	}
}

// GET registra una ruta y handler HTTP GET.
func (r *LumbiMuxRouter) GET(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["GET"] == nil {
		r.reglas["GET"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["GET"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["GET:"+ruta] = middlewares
}

// POST registra una ruta y handler HTTP POST.
func (r *LumbiMuxRouter) POST(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["POST"] == nil {
		r.reglas["POST"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["POST"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["POST:"+ruta] = middlewares
}

// PUT registra una ruta y handler HTTP PUT.
func (r *LumbiMuxRouter) PUT(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["PUT"] == nil {
		r.reglas["PUT"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["PUT"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["PUT:"+ruta] = middlewares
}

// DELETE registra una ruta y handler HTTP DELETE.
func (r *LumbiMuxRouter) DELETE(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["DELETE"] == nil {
		r.reglas["DELETE"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["DELETE"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["DELETE:"+ruta] = middlewares
}

// PATCH registra una ruta y handler HTTP PATCH.
func (r *LumbiMuxRouter) PATCH(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["PATCH"] == nil {
		r.reglas["PATCH"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["PATCH"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["PATCH:"+ruta] = middlewares
}

// OPTIONS registra una ruta y handler HTTP OPTIONS.
func (r *LumbiMuxRouter) OPTIONS(ruta string, handler http.HandlerFunc, middlewares ...Middleware) {
	if r.reglas["OPTIONS"] == nil {
		r.reglas["OPTIONS"] = make(map[string]http.HandlerFunc)
	}
	r.reglas["OPTIONS"][ruta] = anadeMiddleware(handler, middlewares...)
	r.middlewares["OPTIONS:"+ruta] = middlewares
}

// ServeHTTP maneja las solicitudes HTTP entrantes mediante la coincidencia del método de solicitud y la ruta de URL con un handler registrado.
func (r *LumbiMuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.reglas[req.Method][req.URL.Path]; ok {
		for _, middleware := range r.middlewares[req.Method+":"+req.URL.Path] {
			handler = middleware(handler)
		}
		handler(w, req)
		return
	}
	http.NotFound(w, req)
}

// anadeMiddleware toma un http.HandlerFunc y una lista de Middleware y devuelve un nuevo http.HandlerFunc que aplica los middleware a la función original.
func anadeMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// helloHandler es un handler HTTP simple que escribe "¡Hola, mundo!" en la respuesta.
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

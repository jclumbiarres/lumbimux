# LumbiMux
---
Sistema simple de routing para Go usando solo net/http para proyectos locales.

Ejemplo de uso:
```go
package main

import (
	"net/http"
	"github.com/jclumbiarres/lumbimux"
)

func main() {
	r := lumbimux.NewLumbiMux()
	r.GET("/", lumbimux.HelloHandler, lumbimux.LoggingMiddleware)
	http.ListenAndServe(":8000", r)

}
```
* HelloHandler es un handler para pruebas simple que responde con un "¡Hola Mundo!"
* LoggingMiddleware es opcional.
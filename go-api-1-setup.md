# API con Go Gin y SQLite

Notas sobre el tutorial [Building a web app with Go and SQLite, Jeremy Morgan (PluralSight Tutorial)](https://www.allhandsontech.com/programming/golang/web-app-sqlite-go/)

Gin es un *framework* para construir aplicaciones web basadas en API de forma sencilla.

Como *backend* usamos SQLite, un *motor* de base de datos ligero, rápido y autocontenido, que no requiere un servidor separado, sino que forma parte de la aplicación desarrollada.

## Inicializar el módulo

Creo una carpeta `go-api/`, donde ubicaremos todos los ficheros de la aplicación.

En `go-api/`, inicializamos el módulo:

```shell
go mod init github.com/onthedock/go-api
```

Creamos el fichero `go-api/main.go` y declaramos la dependencia de `gin`:

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {

}
```

> VSCode puede configurarse para gestionar los *imports* automáticamente. En este caso, es posible que al guardar se elimine el *import* de Gin, ya que no hay código que lo requiera. Para mantener el *import*, puedes añadir el *default router* en Gin `r := gin.Default()` dentro de la función `main`.

Como hemos añadido una dependencia, ejecutamos `go mod tidy` para que se descargue Gin (y sus dependencias).

## Crear un *router*

Un *router* establece la relación entre las peticiones que recibe el servidor web y las funciones que las gestionan.

Creamos una instance del *router* `Default()` proporcionado por Gin:

```go
// main.go -> func main()
r := gin.Default()
```

A continuación se define un [grupo](https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.Group), de manera que todas las rutas definidas en el grupo tienen *middlewares* comunes o el mismo prefijo. En este caso, lo usamos para versionar la API. Tras crear las rutas, usamos [r.Run()](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.Run) que asocia el *router* al servidor HTTP y empiza a escuchar y a responder peticiones (*listen & serve*):

```go
// main.go -> func main()
r := gin.Default()

// API v1
v1 := r.Group("/api/v1")
{
  v1.GET("person", getPersons)
  v1.GET("person/:id", getPersonById)
  v1.POST("person", addPerson)
  v1.PUT("person/:id", updatePerson)
  v1.DELETE("person/:id", deletePerson)
  v1.OPTIONS("person", options)
}

// By default it serves on :8080 unless a
// PORT environment variable was defined.
r.Run()
```

Toddas las peticiones que llegan al servidor `www.example.org/api/v1` se gestionan por este grupo de *routes*.

Gestionamos cada petición en función del **verbo** HTTP, la ruta a la que va dirigida y en función de ello, identificamos la función a llamar. Por ejemplo:

```go
v1.GET("person", getPersons)
```

- Verbo: `GET`
- Ruta: `person` (en realidad, `/api/v1/person`)
- Función a llamar: `getPersons`

Gin también gestiona rutas *variables*, como en:

```go
v1.GET("person/:id", getPersonById)
```

- Verbo: `GET`
- Ruta: `person/:id` (en realidad, `/api/v1/person`). En este caso, `:id` permite capturar esta parte de la ruta y almacenarla en la variable `id`.
- Función a llamar: `getPersons`

## *Handlers*

El siguiente paso consiste en definir las funciones que van a ejecutarse cuando se reciba una petición en una determinada ruta.

Para devolver el código HTTP correspondiente en cada caso, importamos el paquete `net/http` (en vez de devolver [*números mágicos*](https://en.wikipedia.org/wiki/Magic_number_(programming))) y definimos funciones sencillas que nos permiten validar que todo funciona.

```go
// main.go

import (
    "net/http"

    "github.com/gin-gonic/gin"
)
```

> El paquete `net/http` forma parte de la biblioteca standard de Go, por lo que no es necesario ejecutar `go mod tidy`.

Las funciones *dummy* usadas como *handlers* temporales:

```go
// main.go

func getPersons(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "getPersons Called"})
}

func getPersonById(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"message": "getPersonById " + id +" Called"})
}

func addPerson(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "addPerson Called"})
}

func updatePerson(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"message": "deletePerson " + id + " Called"})
}

func options(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "options Called"})
}
```

Cada una de estas funciones acepta como parámetro un puntero a la estructura `Context` definida en Gin: [`c *gin.Context`](https://pkg.go.dev/github.com/gin-gonic/gin#Context). La estructura `Context` permite pasar variables entre el *middleware*, gestionar el flujo y validar los JSON de las peticiones y las respuestas.

Usamos la función [func (c *Context) JSON(code int, obj any)](https://pkg.go.dev/github.com/gin-gonic/gin#Context.JSON) para *serializar* el objeto JSON en el *body* de la respuesta, indicando a Gin que vamos a devolver un objeto JSON.

El primer parámetro es un *int* que corresponde con el *status code* (para el que usamos los códigos definidos en `net/http`). El segundo, es el objeto JSON que devolvemos como respuesta (usando `gin.H()`: *H is a shortcut for map[string]interface{}*):

```go
c.JSON(http.StatusOK, gin.H{"message": "getPersons Called"})
```

En aquellas funciones en la que hemos definido una parte de la ruta de la petición como variable, la asignamos a una variable obteniéndola del contexto ([Param](https://pkg.go.dev/github.com/gin-gonic/gin#Param)) y la incluimos concatenando su valor en la respuesta:

```go
id := c.Param("id")
c.JSON(http.StatusOK, gin.H{"message": "getPersonById " + id +" Called"})
```

## Primera prueba

Construimos y ejecutamos la aplicación:

```shell
go build
./go-api
```

Desde otro terminal, usamos `curl` para lanzar peticiones a la aplicación. Por defecto, el servidor escucha en el puerto `8080`:

> En mi caso, la aplicación está "publicada" en `http://go.dev.vm`

```shell
$ curl http://go.dev.vm:8080/api/v1/
404 page not found
```

No hemos definido ninguna ruta para `/`, por lo que obtenemos un mensaje de error (no se encuentra la página).

Por defecto, `curl` realiza peticiones `GET`, pero podemos indicar el *verbo* a usar mediante `-X`:

```shell
$ curl -X GET http://go.dev.vm:8080/api/v1/
404 page not found
```

En la salida del servidor, vemos el registro que genera Gin por defecto:

```shell
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2022/06/12 - 10:48:04 | 404 |         387ns |   192.168.1.139 | GET      "/api/v1/"
[GIN] 2022/06/12 - 10:50:29 | 404 |       2.479µs |   192.168.1.139 | GET      "/api/v1/"
```

Si lanzamos la petición a una de las rutas para las que sí tenemos definido un *handler*:

```shell
$ curl -X GET http://go.dev.vm:8080/api/v1/person
{"message":"getPersons Called"}
```

En este caso, en el registro de Gin:

```shell
[GIN] 2022/06/12 - 10:52:51 | 200 |      94.815µs |   192.168.1.139 | GET      "/api/v1/person"
```

Del mismo modo, si pasamos un `id`:

```shell
$ curl -X GET http://go.dev.vm:8080/api/v1/person/151 
{"message":"getPersonById 151 Called"}
```

Y en el *log* de Gin:

```shell
[GIN] 2022/06/12 - 10:54:47 | 200 |     267.503µs |   192.168.1.139 | GET      "/api/v1/person/151"
```

Cambiando el *verbo* en `curl`, podemos probar el resto de rutas; por ejemplo, `DELETE`:

```shell
$ url -X DELETE http://go.dev.vm:8080/api/v1/person/13 
{"message":"deletePerson 13 Called"}
```

Y en el registro:

```shell
[GIN] 2022/06/12 - 10:56:17 | 200 |     242.367µs |   192.168.1.139 | DELETE   "/api/v1/person/13"
```

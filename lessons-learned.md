# Lecciones aprendidas

## Importar el paquete `gin-gonic/gin`

```go
package main

import (
    "net/http" // Opcional, pero interesante para las constantes http.StatusOK
    "github.com/gin-gonic/gin"
)
```

## *Router* por defecto

```go
func main() {
    r := gin.Default()
}
```

El *router* por defecto en Gin incluye dos *middlewares*; uno permite recuperarse de un *panic* (llamado *Recovery*) y el *Logger*, que genera logs.

## *Router* con *middleware* personalizado

Usando `gin.New()` creamos un *router* sin ningún *middleware*.

Para añadir un *middleware custom* llamado `Logger()`, usamos:

```go
r := gin.New()
r.Use(Logger())
```

El *middleware* tiene la *signature*

```go
func Logger() gin.HandlerFunc {
    // code
}
```

El *router* creado mediante `gin.Default()` es equivalente a:

```go
r := gin.New()
r.Use(gin.Revocery())
r.Use(gin.Logger())
```


## Rutas

Para añadir una ruta, usamos el método asociado al *verbo* HTTP; por ejemplo:

```go
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
```

### Grupos de rutas

Puede resultar interesante usar [*grupos de rutas*](https://gin-gonic.com/docs/examples/grouping-routes/), por ejemplo mediante `r.Group("/api/v1"){}`:

```go
v1 := r.Group("/api/v1")
{
    v1.GET("personas", getPersonas) // example.org/api/v1/personas
    // ...
}
```

Todas las rutas en un grupo comparten el *path* base (`/api/v1` en mi caso) y los *middlewares*.

### Parámetros en el *path*

Para definir una parte variable del *path*  que va a ser capturada como una variable, usamos `:id` (para guardar el contenido del *path* en la variable `id`):

```go
v1.GET("personas/:id", getPersonasById) // example.org/api/v1/personas/123
```

Podemos acceder al valor en la variable `id` mediante `c.Param("id")`.

Ver [Parameters in path](https://gin-gonic.com/docs/examples/param-in-path/).

### Parámetros como parte de la *query string*

En este caso, la URL de la petición es de la forma `/welcome?firstname=John&lastname=Doe`.

Definimos la ruta como (obtenido del ejempo en [Query string parameters](https://gin-gonic.com/docs/examples/querystring-param/)):

```go
router.GET("/welcome", func (c *gin.Context") {
    firstname: = c.DefaultQuery("firstname", "Guest") // Provides a default value
    lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname") 
})
```

## Ejecutar el servidor

```go
r.Run() // equivalente a ListenAndServe 0.0.0.0:8080
```

Se puede especificar un puerto diferente especificándolo `r.Run(":8910")`.

## Ejemplo completo

De la documentación oficial en el [Quickstart](https://gin-gonic.com/docs/quickstart/):

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

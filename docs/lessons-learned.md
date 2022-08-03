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
    firstname: = c.DefaultQuery("firstname", "Guest") // Provides a default value 'Guest'
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

## Bases de datos

Para una introducción de Go y bases de datos SQL, revisa el [Go database/sql tutorial](http://go-database-sql.org/).

### Importar paquete y *driver*

El paquete `database/sql` proporciona un interfaz genérico para todas las bases de datos SQL (aunque el *driver* específico para cada una debe instalarse por separado).

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)
```

> Tras añadir paquetes no incluidos en la biblioteca estándar, ejecutar `go mod tidy`.

### Conectar con la base de datos

[Open](https://pkg.go.dev/database/sql#Open) devuelve un puntero a [`sql.DB`](https://pkg.go.dev/database/sql#DB), que se puede usar tanto en conexiones concurrentes como reutilizarse de forma segura (mantiene un *pool* de conexiones).

Open sólo debe usarse una vez.

Se puede definir `DB` como una variable global:

```go
var DB *sql.DB
```

Y usar `sql.Open()` en una función `OpenDB()`, por ejemplo:

```go
func OpenDB() error {
    db, err := sql.Open("sqlite3", "./names.db")
    if err != nil {
        log.Printf("Error opening database ./names.db: %s\n", err.Error())
        return err
    }
    DB = db
    return nil
}
```

> Hay otras opciones, ya que en general no es una buena idea usar variables globales. Puedes ver diferentes alternativas en el artículo [Organising Database Access in Go](https://www.alexedwards.net/blog/organising-database-access) de Alex Edwards.

### Crear tablas si no existen (desde el código)

Usar la función `CREATE TABLE IF NOT EXISTS <nombre_tabla> (col1, col2, ...)` y crear la(s) tabla(s) necesarias desde el código.

### Usar consultas *preparadas*

Ver [Using prepared statements](https://go.dev/doc/database/prepared-statements) y la diferencia entre una consulta y una transacción en [Executing transactions](https://go.dev/doc/database/execute-transactions).

> If a function name includes `Query`, it is designed to ask a question of the database, and will return a set of rows, even if it’s empty. Statements that don’t return rows should not use `Query` functions; they should use `Exec()`.

Por lo que he entendido, usar `db.Query()`, obtener `rows`, recorrer las filas obtenidas en un bucle (con `rows.Next()`, `rows.Scan()`) y finalmente cerrar la conexión (con un `defer rows.Clode()`) es el patrón recomendado cuando se realizan consultas; en este contexto, *consultas* significa obtener datos de la base de datos, sin modificar su contenido.

La modificación de datos es preferible realizarla usando *prepared statements* (ver [Statements that Modify Data](http://go-database-sql.org/modifying.html)).

En el tutorial, se usan *transacciones* a la hora de insertar, actualizar y eliminar datos. No veo la diferencia entre usar *prepared statements* o *transacciones*.

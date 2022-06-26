# Conexión con la base de datos

En la primera parte del tutorial hemos configurado Gin y hemos creado las rutas que vamos a *publicar* en nuestra aplicación.

En esta segunda parte nos centramos en conectar con una base de datos, en SQLite.

Usamos la base de datos `names.db` proporcionada por Jeremy Morgan en la segunda parte de su tutorial, [Building a Web App with Go and SQLite](https://www.allhandsontech.com/programming/golang/web-app-sqlite-go-2/).

Esta base de datos SQLite 3 contiene 1000 entradas con el siguiente esquema:

```sql
CREATE TABLE "people" (
    "id" INTEGER,
    "first_name" TEXT,
    "last_name" TEXT,
    "email" TEXT,
    "ip_address" TEXT,
    PRIMARY KEY("id" AUTOINCREMENT)
)
```

## *Package* `models`

Creamos la carpeta `models/` y en ella, el fichero `models/person.go`.

Este fichero contendrá el paquete `models`:

```go
// models/person.go
package models

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
```

Antes de continuar, ejecutamos `go mod tidy` para importar los paquetes necesarios.

```shell
$ go mod tidy
go: finding module for package github.com/mattn/go-sqlite3
go: found github.com/mattn/go-sqlite3 in github.com/mattn/go-sqlite3 v1.14.13
```

El paquete [`database/sql`](https://pkg.go.dev/database/sql) forma parte de la biblioteca standard. Proporciona una intefaz genérica para interaccionar con bases de datos SQL **pero no incluye ningún driver**.

Por ello, es necesario importar el segundo paquete; en realidad, lo que hacemos es aprovechar que al importar un paquete se ejecutan unas funciones de inicialización (mediante la función `init()`) y éstas realizan configuraciones que permiten usar el paquete, aunque no usemos el código importado *directamente*; por ello necesitamos prefijar el `import` del *blank identifier* (`_`) y así evitar los avisos del compilador indicando que se ha importado un paquete pero que no se utiliza el código.

Finalmente, declaramos una variable global `DB` que contiene un puntero a la conexión con la base de datos.

## Conectar a la base de datos

Creamos una función para realizar la conexión a la base de datos:

> En el tutorial original el autor no muestra ningún mensaje cuando se produce un error. IMHO, esto dificulta saber qué es lo que está fallando y dónde, por lo que prefiero añadir algo de contexto en los logs.
>
> Requiere importar el paquete `log` (de la biblioteca standard).

```go
// models/person.go
func ConnectDatabase() error {
    db, err := sql.Open("sqlite3", "./names.db")
    if err != nil {
        log.Printf("Error opening database ./names.db: %s\n", err.Error())
        return err
    }
    DB = db
    return nil
}
```

## Crear el modelo

Creamos un modelo que refleja la estructura de la base de datos. Decoramos los nombres de la `struct` con etiquetas que permiten especificar cómo referirse a los campos en el objeto JSON equivalente:

```go
type Person struct {
    Id          int    `json:"id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    Email       string `json:"email"`
    IpAddress   string `json:"ip_address"`
}
```

Esta es la `struct` que usaremos para pasar los registros de la base de datos a la web y viceversa.

## Obtener una lista de personas

```go
func GetPersons(count int) ([]Person, error) {

    rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address from people LIMIT " + strconv.Itoa(count))
    if err != nil {
        log.Printf("Error querying the database: %s\n", err.Error())
        return nil, err
    }
    defer rows.Close()

    people := make([]Person, 0)

    for rows.Next() {
        singlePerson := Person{}
        err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IpAddress)

        if err != nil {
            log.Printf("Error getting the Next record: %s", err.Error())
            return nil, err
        }

        people = append(people, singlePerson)
    }

    err = rows.Err()
    if err != nil {
        log.Printf("Error getting the Next: %s", err.Error())
        return nil, err
    }

    return people, err
}
```

Declaramos la función `GetPersons`, a la que pasamos un número (el de registros que queremos obtener de la base de datos) y que devuelve un *slice* de `Person` y un error. La idea es que, si cualquier cosa en la función falla, devolvemos `nil` y el error.

Consultamos la base de datos con [`DB.Query()`](https://pkg.go.dev/database/sql#DB.Query) para obtener el número de registros indicado.

Si se produce un error en la *query*, devolvemos el error; si no, postponemos el cierre de la conexión con la base de datos al finalizar la función con `defer rows.Close()`.

La consulta devuelve varias filas; las recorremos mediante [`rows.Next()`](https://pkg.go.dev/database/sql#Rows.Next). `rows.Next()` contiene `true` si tenemos una fila con datos que *escanear* con `rows.Scan()`; esto nos permite recorrer los resultados usando un bucle.

El resultado de cada fila obtenida de la base de datos la insertamos en la variable de tipo `Person` (`singlePerson`) y añadimos la *singlePerson* al *slice* de `Person` `people`, donde vamos acumulando los resultados.

Si no hemos tenido ningún error, devolvemos el *slice* `people` (también devolvemos `err`, aunque si hemos llegado hasta aquí significa que `err=nil`)...

## Importar el paquete `models` en `main.go`

Volvemos al fichero `main.go` para importar el paquete `models`. En mi caso, el módulo creado para el proyecto se llama `go-api` de manera que lo añado a los `import`:

```go
// main.go
import (
    "net/http"
    "github.com/onthedock/go-api/models"
    "github.com/gin-gonic/gin"
)
```

Modificamos la función `getPersons(c *gin.Context)` para devolver vía web los valores obtenidos de la consulta a la base de datos:

```go
func getPersons(c *gin.Context) {

    persons, err := models.GetPersons(10)
    checkErr(err)

    if persons == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
        return
    } else {
        c.JSON(http.StatusOK, gin.H{"data": persons})
    }
}
```

> El autor del tutorial original define una función `checkError(err error)` que registra el error mediante `log.Fatal()` y que a continuación finaliza la ejecución del programa:
>
> ```go
> func checkErr(err error) {
>   if err != nil {
>       log.Fatal(err)
>   }
> }
> ```
>
> En mi opinión, esto diluye el contexto en el que se produce el error, lo que dificulta identificar dónde se ha originado el error. Por ello, he introducido mensajes *contextuales* antes de devolver el error, cuando se produce.

La función `getPersons` ejecuta la función `models.GetPersons(10)`; si se produce cualquier error con la base de datos, `checkErr` finaliza el programa.

Si no se ha producido ningún error, `persons` no puede ser `nil`; por tanto, no me parece correcto devolver `http.StatusBadRequest` (400), sino que lo más adecuado creo que sería `http.StatusNotFound` (404), que además es coherente con el mensaje devuelto de `No records found`.

Si `persons` no es `nil`, devolvemos el objeto JSON con los resultados obtenidos en `persons`.

### Realizar la conexión con la base de datos

Hasta ahora hemos estado enfocados en actualizar la llamada `getPersons(..)` para llamar a `models.GetPersons(10)` y procesar el resultado obtenido desde la base de datos y devolverlo en la web.

Antes de poder probar la aplicación, tenemos que añadir la conexión con la base de datos; al principio de la función `main.go`, añadimos:

```go
// main.go --> main()
func main() {

    err := models.ConnectDatabase()
    checkErr(err)
    ...
```

La primera línea realiza la conexión con SQLite e inicializa la variable `DB` con la que el *driver* gestionará toda la interacción con la base de datos.

## Probar la consulta

Compilamos la aplicación y la ejecutamos.

La función `getPerson(10)` tiene el valor de registros **fijo** en el código de la aplicación, por lo que lanzamos una petición a `/api/v1/person`:

> Uso **jq** para mostrar de forma más inteligible el JSON devuelto por la aplicación.

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person | jq
{
  "data": [
    {
      "id": 1,
      "first_name": "Glyn",
      "last_name": "Quaife",
      "email": "gquaife0@edublogs.org",
      "ip_address": "252.74.16.5"
    },
    {
      "id": 2,
      "first_name": "Kathrine",
      "last_name": "Aizkovitch",
      "email": "kaizkovitch1@bandcamp.com",
      "ip_address": "255.1.189.50"
    },
    {
      "id": 3,
      "first_name": "Gaven",
      "last_name": "Allanby",
      "email": "gallanby2@cloudflare.com",
      "ip_address": "33.223.203.230"
    },
...
```

He creado una base de datos vacía y he lanzado la misma petición y el resultado es:

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person | jq
{
  "data": []
}
```

En este caso el código de estatus devuelto es `200`; no se ha producido ningún error, por lo que el *slice* `people`, que se ha inicializado con 0 elementos antes de `rows.Next()`, está vacío.

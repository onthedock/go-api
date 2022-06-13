# Insertar y actualizar registros

En esta parte del tutorial vamos a insertar y actualizar registros en la base de datos.

Lanzaremos una petición `POST` que contiene un objeto JSON y lo procesaremos para generar un nuevo registro con la información contenida en él.

## Agregar una persona a la base de datos `AddPerson`

En el fichero `models/person.go`, añadimos la función `AddPerson`; a la que pasamos una `Person` y que devuelve `bool` (en función de si se ha añadido o no la persona a la base de datos) y un error.

Empezamos inicializando una transacción, con [`DB.Begin()`](https://pkg.go.dev/database/sql?utm_source=gopls#DB.Begin). A continuación, tras controlar un posible error, preparamos el comando para insertar el registro en la base de datos (excluimos `id` porque es autogenerado por la base de datos).

Lo insertamos mediante `stmt.Exec()` y finalizamos la transacción con [`tx.Commit()`](https://pkg.go.dev/database/sql?utm_source=gopls#Tx.Commit).

Si todo ha ido bien, devolvemos `true` (y ningún error, `nil`):

```go
// models/person.go
func AddPerson(newPerson Person) (bool, error) {

    tx, err := DB.Begin()
    if err != nil {
        log.Printf("Error opening the database: %s", err.Error())
        return false, err
    }
    stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email, ip_address) VALUES (?,?,?,?)")
    if err != nil {
        log.Printf("error preparing insert statament: %s", err.Error())
        return false, err
    }
    defer stmt.Close()

    _, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)
    if err != nil {
        log.Printf("error executing the query: %s", err.Error())
        return false, err
    }
    tx.Commit()
    return true, nil
}
```

## Llamar a `AddPerson`

Una vez creado el método para insertar los valores en la base de datos, volvemos a `main.go`. El objetivo es crear la función que recoja el objeto JSON enviado en la petición, convertirlo en una *struct* `Person` y pasarlo a la función que acabamos de crear para que lo inserte en la bbdd.

Creamos la variable `json` de tipo `models.Person`. Mediante `c.ShouldBindJSON(&json)`, que intenta asignar el contenido de la petición a la *struct* `json` (de tipo `Person`). Si hay un error, se devuelve para que se gestione en la aplicación (a diferencia de los métodos `Must...`, que fallan directamente con un 400 [Model binding and validation](https://pkg.go.dev/github.com/gin-gonic/gin@v1.8.1#readme-model-binding-and-validation)).

La función de *binding* permite validar campos requeridos, con valores comprendidos entre determinados límites, con un formato concreto, etc. Para ello, se deben decorar los campos en la *struct* (por ejemplo, con `binding:"required"` o `binding:"required,iso3166_1_alpha2"`)... Ver por ejemplo el artículo [Gin binding in Go: A tutorial with examples](https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/).

```go
// main.go
func addPerson(c *gin.Context) {
    var json models.Person
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    success, err := models.AddPerson(json)
    if success {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err})
    }
}
```

Compilamos la aplicación y la ejecutamos.

Podemos comprobar el funcionamiento de `AddPerson`; para ello creamos el fichero `user.json`:

```json
{
    "first_name": "Xavi",
    "last_name": "OnTheDock",
    "email": "onthedock@example.org",
    "ip_address": "127.0.0.1"
}
```

Usando `curl`:

```shell
$ curl --silent -X POST http://go.dev.vm:8080/api/v1/person -d @./user.json
{"message":"success"}
```

Como la versión actual de la aplicación sólo devuelve `{ "message": "success" }` o el error que se ha producido al intentar insertar un valor en la base de datos, no sabemos cuál es el `id` del nuevo registro (tampoco podemos buscar, por ahora).

La base de datos original del autor contiene 1000 registros, por lo que el registro insertado es el 1001:

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person/1001 | jq
{
  "data": {
    "id": 1001,
    "first_name": "Xavi",
    "last_name": "OnTheDock",
    "email": "onthedock@example.org",
    "ip_address": "127.0.0.1"
  }
}
```

## Actualizar un registro (`UpdatePerson`)

... # To be Done

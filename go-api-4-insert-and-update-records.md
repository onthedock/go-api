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

Para actualizar un registro, indicaremos qué registro pasando el `id`: `/api/v1/person/{id}` y usando el método HTTP `PUT`.

En `models/person.go`:

```go
// models/person.go
func UpdatePerson(ourPerson Person, id int) (bool, error) {
    tx, err := DB.Begin()
    if err != nil {
        fmt.Printf("Error connecting to db, %s", err.Error())
        return false, err
    }
    stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")
    if err != nil {
        fmt.Printf("Error preparing update: %s", err.Error())
        return false, err
    }
    defer stmt.Close()
    _, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, ourPerson.IpAddress, id)
    if err != nil {
        fmt.Printf("Error executing the query: %s", err.Error())
        return false, err
    }
    tx.Commit()
    return true, nil
}
```

Creamos la función para actualizar un registro en la base de datos.

Lo primero que me llama la atención es que se use `id int`, cuando en `GetPersonById(id string)` se usa un `string`. El autor no hace ningún comentario al respecto en el tutorial original, por lo que entiendo que se ha usado el razonamiento *por defecto* de que el `id` es un número, en contra de lo indicado anteriormente, de usar un `string` para evitar conversiones...

Por lo demás, se inicia una transacción con `DB.Begin()`, se prepara el *statement* para realizar el `UPDATE` en la base de datos y se ejecuta con `stmt.Exec(...)`, proporcinando los valores a insertar obtenidos desde el *struct* de tipo `Person` pasado en la función.

> En el original, faltan las comas entre los diferentes campos del comando `UPDATE`.
>
> Además, en `stmt.Exec()` se pasa `ourPerson.Id`, cuando en realidad se debe pasar `id`. Tal y como está el código en el tutorial, ourPerson.Id está **vacío**. Al realizar el `UPDATE` con `ourPerson.Id` se muestra el mensaje de éxito, pero si obtenemos de nuevo el registro con `GET`, veremos que **no se ha actualizado**.

Finalmente, si no se ha producido ningún error al ejecutar la consulta en la base de datos, se *compromete* (`tx.Commit()`) la transacción y se devuelve `true, nil`.

## Llamar a `updatePerson`

El método `updatePerson` es prácticamente igual a `addPerson`, solo que ahora tomamos el `id` de la URI:

```go
// main.go
func updatePerson(c *gin.Context) {
    var json models.Person
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    personId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    success, err:= models.UpdatePerson(json, personId)
    if success {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err})
    }
}
```

Obtenemos el objeto JSON enviado por el cliente como parte de la petición y lo *encajamos* en la variable `json` (que es una *struct* de tipo `Person`).

Si no hay ningún error, obtenemos también el `id` desde la petición y lo convertimos a `int`, ya que en `UpdatePerson` esperamos un `int`. Si se produce un error, el servidor no seguirá procesando la petición, por lo que, en mi opinión, es necesario incluir un `return` tras devolver el mensaje de error inválido.

Finalmente, si todo ha ido bien devolvemos el mensaje `success` o el error que se ha producido, en caso contrario.

```shell
 curl --silent -X PUT http://go.dev.vm:8080/api/v1/person/1001 -d @user.json
{"message":"success"}%
```

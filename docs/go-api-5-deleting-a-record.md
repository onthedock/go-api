# Borrando un registro

En esta última parte del tutorial, añadimos la capacidad de borrar registros de la base de datos a través de la API.

## Borrar un registro (`DeletePerson`)

En `models/person.go`

```go
// models/person.go
func DeletePerson(personId int) (bool, error) {
    tx, err := DB.Begin()
    if err != nil {
        fmt.Printf("Error connecting to database: %s", err.Error())
        return false, err
    }
    stmt, err := DB.Prepare("DELETE FROM people WHERE id = ?")
    if err != nil {
        fmt.Printf("Error preparing statement: %s", err.Error())
        return false, err
    }
    defer stmt.Close()
    _, err = stmt.Exec(personId)
    if err != nil {
        fmt.Printf("Error deleting record: %s", err.Error())
        return false, err
    }
    tx.Commit()
    return true, nil
}
```

> De nuevo el autor usa un `int` como identificador del registro a borrar en la definición de la función.

Iniciamos la transacción con `DB.Begin()`, preparamos la *query* y la ejecutamos. Si en cualquier momento se produce un error, devolvemos `false` y el error que se ha producido. Si no, finalizamos la transacción con `tx.Commit()` y devolvemos `true`.

## Llamar a `deletePerson`

Volvemos a `main.go` para actualizar la función `deletePerson`.

Obtenemos `id` del contexto de Gin (`gin.Context`), pero es un *string*, así que lo convertimos en un `int` usando `strconv.Atoi`. Una vez convertido en `int`, si no ha habido ningún error, lo usamos para llamar `models.DeletePerson(personId)`.

En función del resultado, devolvemos `{"message": "success"}` o el error que se ha producido:

```go
func deletePerson(c *gin.Context) {
    personId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        fmt.Printf("Error converting Id: %s", err.Error())
        return
    }
    success, err := models.DeletePerson(personId)
    if success {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err})
    }
}
```

Compilamos y ejecutamos la aplicación para validar que funciona correctamente; obtenemos un registro:

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person/1000     
{"data":{"id":1000,"first_name":"Christy","last_name":"Schankelborg","email":"cschankelborgrr@mayoclinic.com","ip_address":"242.79.255.50"}}
```

Lo borramos:

```shell
$ curl --silent -X DELETE http://go.dev.vm:8080/api/v1/person/1000
{"message":"success"}
```

Validamos que se ha eliminado:

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person/1000   
{"error":"No record found"}
```

## `OPTIONS`

El tutorial acaba implementando el [método `OPTIONS` de HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS). El objetivo de este método es que el cliente pueda obtener información de los métodos soportados por la API.

En `main.go`:

```go
func options(c *gin.Context) {
    ourOptions := "HTTP/1.1 200 OK\n" +
        "Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
        "Access-Control-Allow-Origin: http://go.dev.vm:8080\n" +
        "Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
        "Access-Control-Allow-Headers: Content-Type\n"

    c.String(http.StatusNoContent, ourOptions)
}
```

En el tutorial el autor devuelve un código de estado numérico (200), cuando para el resto de las peticiones ha usado las constantes definidas en el paquete `net/http`. En la documentación de la MDN (Mozilla Developer Network), el código devuelto es `204 - No content`. Sin embargo, si devuelvo `http.StatusNoContent`, no se muestra nada al hacer la petición a la API con el método `OPTIONS`, por lo que he dejado `http.StatusOK`:

```shell
$ curl -X OPTIONS  http://go.dev.vm:8080/api/v1/person
HTTP/1.1 200 OK
Allow: GET,POST,PUT,DELETE,OPTIONS
Access-Control-Allow-Origin: http://go.dev.vm:8080
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
Access-Control-Allow-Headers: Content-Type
```

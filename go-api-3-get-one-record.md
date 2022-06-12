# Obtener un registro

En esta parte del tutorial obtendremos un registro de la base de datos a partir de su `id`.

Creamos la función `GetPersonById` en `models/person.go`:

```go
func GetPersonById(id string) (Person, error) {

    stmt, err := DB.Prepare("SELECT id, first_name, last_name, email, ip_address from people WHERE id = ?")

    if err != nil {
        log.Printf("Error querying the database: %s", err.Error())
        return Person{}, err
    }

    person := Person{}

    sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email, &person.IpAddress)

    if sqlErr != nil {
        if sqlErr == sql.ErrNoRows {
            return Person{}, nil
        }
        return Person{}, sqlErr
    }
    return person, nil
}
```

Como queremos obtener un registro a partir del `id`, la función `GetPersonById` acepta el `id` como parámetro de entrada. El `id` es un `string`; podríamos usar un `int` y después convertirlo en `string` (para incluirlo en la *query*), pero de esta manera es más sencillo.

Si encontramos la persona en la base de datos, devolveremos `Person` (o un `error`).

El autor del tutorial indica que usar [DB.Prepare](https://pkg.go.dev/database/sql#DB.Prepare) permite evitar ataques de inyección SQL, lo que parece una buena idea. Además, usamos el nombre de los campos de forma explícita (en vez de `SELECT *`).

Ejecutamos la *query* pasando el `id` que queremos obtener de la base de datos y comprobamos si se ha producido algún error. [Query.Row](https://pkg.go.dev/database/sql#Stmt.QueryRow) devuelve `ErrNoRows` si no hay resultados, pero esto no es un error. Así que si se produce un error, comprobamos si se trata de `ErrNoRows`; en este caso, devolvemos `Person` (vacío) y `nil`; si no se trata de este error, devolvemos el err que se ha producido. Y si no se ha producido ningún error, devolvemos el resultado obtenido.

## Actualizar `GetPersonById`

Modificamos la función para llamar a `models.GetPersonById` y comprobamos si se ha producido un error con `checkErr(err)`.

```go
    // If person.Id is empty, no record is found
    if person.Id == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
        return
    } else {
        c.JSON(http.StatusOK, gin.H{"data": person})
    }
```

El autor del tutorial comprueba si el nombre en `person.FirstName` está vacío para determinar si se han encontrado resultados o no. En mi opinión, es más *genérico* comprobar por ejemplo, `person.id`; mi razonamiento es que el campo `id` es un campo más común, por lo que "la técnica" para validar si se han obtenido resultados de la consulta es más general y por tanto, aplicable en casi todos los casos. El campo `person.Id` es un `int`, por lo que verificamos que no es igual a `0`. Puede existir algún caso en que una base de datos permita `0` como un Id de registro válido, pero debería ser un caso extraño; del mismo modo que una persona no tenga "primer nombre" (y diría que éste caso puede ser más común).

Por otro lado, si hemos llegado hasta aquí, significa que no se ha producido ningún error, por lo que me parece más adecuado devolver el código de estado `http.StatusNotFound` (404) en vez de `http.StatusBadRequest` (400).

## Prueba: obtener un registro

Compilamos la aplicación y la ejecutamos.

Usando la base de datos vacía, intentamos obtener un registro (obtenemos un 404):

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person/1 | jq
{
  "error": "No record found"
}
```

En el log de Gin, vemos que obtenemos un código de estado 404:

```shell
[GIN] 2022/06/12 - 18:50:14 | 404 |      263.28µs |   192.168.1.139 | GET      "/api/v1/person/1"
```

Si cambiamos a la base de datos con valores:

```shell
$ curl --silent -X GET http://go.dev.vm:8080/api/v1/person/1 | jq
{
  "data": {
    "id": 1,
    "first_name": "Glyn",
    "last_name": "Quaife",
    "email": "gquaife0@edublogs.org",
    "ip_address": "252.74.16.5"
  }
}
```

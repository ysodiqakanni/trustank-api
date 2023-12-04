# Go RESTful API Starter Kit (Boilerplate)

The Project uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database access: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.


## Common Development Tasks

This section describes some common development tasks using this starter kit.

### Implementing a New Feature

Implementing a new feature typically involves the following steps:

1. Develop the service that implements the business logic supporting the feature. Please refer to `internal/album/service.go` as an example.
2. Develop the RESTful API exposing the service about the feature. Please refer to `internal/album/api.go` as an example.
3. Develop the repository that persists the data entities needed by the service. Please refer to `internal/album/repository.go` as an example.
4. Wire up the above components together by injecting their dependencies in the main function. Please refer to 
   the `album.RegisterHandlers()` call in `cmd/server/main.go`.

## Running the Project
You can run the project using several methods


To build the docker image, stay in the root folder and run
 `docker build -f cmd/server/dev.Dockerfile -t trustankbizapi:latest .`
### Running as a docker container
Run the image: `docker run -it -p 8080:8080 trustankbizapi:latest`

### To Tag and push the image to docker hub
$ docker tag trustankbizapi:latest ysodiqakanni/trustankbizapi:1.0.0
$ docker push ysodiqakanni/trustankbizapi:1.0.0

### Now that the image is in your dockerhub, you can now run the project using kubernetes
From the root folder, run `kubectl apply -f kubernetes/api.deployment.yml` 
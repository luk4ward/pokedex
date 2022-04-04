## TrueLayer, Pokedex
The service has been written using Go, version 1.16.

## Prerequisites
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/) 

## Build and Test

There are few options to build, test and run the service, all require some prerequisites listed above.

Default configuration can be found [here](https://github.com/luk4ward/pokedex/blob/master/config.yaml)

### Docker
```sh
docker build -t pokedex .
docker run --rm -ti -p 5050:5050 pokedex
```

### Docker Compose
```sh
docker-compose up
```
as per [compose](https://github.com/luk4ward/pokedex/blob/master/docker-compose.yaml) file service will be exposed on the same `5050` port

### Makefile

```sh
make run
```

for testing simply run
```sh
make test
```

provided [makefile](https://github.com/luk4ward/pokedex/blob/master/Makefile) contains few more commands that will make running or building process way easier.

## Endpoints

- `GET v1/pokemon/{name}`: fetches a Pokemon for a given name
- `GET v1/pokemon/translated/{name}`: fetches a Pokemon with a translated description for a given name
- `GET v1/_healthcheck`: handler for returning a 200 if service is alive.


## CI/CD

Two checks were added using Github Actions for making sure PRs are not breaking anything:
- `Docker Image CI` - for building docker image 
- `Go` - for buidling, formatting and running all tests 

## Future improvements
As stated in a comment using colored console writer is there only for visibility, it's inefficient and woudn't be used in production.

Also ue to some time limitation there are few areas of improvements and if I could spend more time on it I would definitely consider adding:
- [ ] Integration tests - even though everything is unit-tested, having integration tests that run against every PR as a part of CI/CD pipeline would be something much appreciated before going to production
- [ ] Caching - very important part bearing in mind that one of the third-party services is ratelimited
- [ ] Extended error handling - support of more error response codes
- [ ] Tracing & Metrics - better observability is super important when runing microservices

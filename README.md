# k8s-cluster-head-service
Head service for scaling up and down k8s cluster, and managing services and routine on the cluster. This project serves as a way to get myself familiar with golang.

## How to run locally
1. `go run src/server.go`, or after `air` is configured, run `air`.
1. Open browser at `http://0.0.0.0:3010`

## How to deploy
1. Install Heroku CLI, and run `heroku login`
1. `heroku create k8s-cluster-head-service --manifest` ([Heroku doc](https://devcenter.heroku.com/articles/build-docker-images-heroku-yml#creating-your-app-from-setup)), you'll get `https://k8s-cluster-head-service.herokuapp.com/ | https://git.heroku.com/k8s-cluster-head-service.git`
1. Commit git changes, including `heroku.yml`
1. Push to heroku remote `git push heroku master`
1. The app can be accessed at `https://k8s-cluster-head-service.herokuapp.com/`

# Reference

## Golang

### Exporting func or variable 

In order to export this function you need to capitalize it
https://tour.golang.org/basics/3

## Kickstart golang
- This github repo: https://github.com/rivernews/k8s-cluster-head-service

Getting started with a golang project
- Setup golang for local development: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-macos
- Golang project structure: https://github.com/golang-standards/project-layout
- What is `GOPATH` and `GOROOT`? https://stackoverflow.com/questions/7970390/what-should-be-the-values-of-gopath-and-goroot/7971481#7971481
- Golang tutorial for NodeJS developers
    - Concept comparison: https://deepu.tech/golang-for-javascript-developers-part-1/
- How to use `go mod init`: https://medium.com/@yunskilic/managing-dependencies-with-go-modules-4a6111d641cc
- Golang module v2: https://blog.golang.org/v2-go-modules
- More commands surrounding go modular: https://blog.golang.org/using-go-modules

Learning to write go code
- Basic control flow & playground: https://tour.golang.org/list

Deploying to Heroku
- Heroku doc: https://devcenter.heroku.com/articles/build-docker-images-heroku-yml#creating-your-app-from-setup
- Using docker container to deploy: https://dev.to/ilyakaznacheev/setup-build-automate-deploy-a-dockerized-app-to-heroku-fast-167

Web framework
- Line up stats - compare their stars!: https://github.com/mingrammer/go-web-framework-stars
- The most popular - gin: https://github.com/gin-gonic/gin
- Fiber: https://github.com/gofiber/fiber

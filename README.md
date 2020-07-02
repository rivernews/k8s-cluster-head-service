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

## Job queue
- To inspect job queues, run `workwebui -redis="REDIS_URL" -ns=my_app_namespace`, then navigate to brower `localhost:5040`.
  - If you don't have the CLI installed yet, [follow instruction](https://github.com/gocraft/work#run-the-web-ui) and run `go get github.com/gocraft/work/cmd/workwebui && go install github.com/gocraft/work/cmd/workwebui`.
  - ⚠️ You'll have to remove the username part from the redis URL, looks like it's not supported and will cause AUTh argument number error. Basically just [following this example](https://github.com/gocraft/work/issues/114#issuecomment-476822085).

## Testing

We're using [Testify](https://github.com/stretchr/testify) the test framework.
- `cd` into the directory where the test file resides
- Run `go test`

## Manual CI/CD API call

### Circle CI

#### Create a new pipeline

This will trigger a build for the specified branch.
POST against endpoint `/project/{project-slug}/pipeline`. For what `{project-slug}` is, see below.

Response:

```json
{"number":243,"state":"pending","id":"5c9ab317-3f41-4851-a3de-e5fb119da8e6","created_at":"2020-06-16T17:04:12.930Z"}
```


#### Get project meta info

Doesn't do much.

```sh
curl -Gv \
  --data-urlencode "circle-token=$CIRCLECI_TOKEN" \
  -H 'Accept: application/json' \
  https://circleci.com/api/v2/project/github%2Frivernews%2Firiversland2-kubernetes
```

#### Get pipelines of a project

Does't do much.
We only need the vcs, org and repo-name slug, which forms `{project-slug}`.

```sh
curl -Gv \
  --data-urlencode "circle-token=$CIRCLECI_TOKEN" \
  -H 'Accept: application/json' \
  https://circleci.com/api/v2/project/github%2Frivernews%2Firiversland2-kubernetes/pipeline/mine | python -mjson.tool
```

This endpoint is able to provide:
- All pipelines of this project
- Each pipeline contains
  - Pipeline id
  - Pipeline number
  - State

Response object:

```json
{ "items": [
    {
        "created_at": "2020-06-16T17:04:12.930Z",
        "errors": [],
        "id": "5c9ab317-3f41-4851-a3de-e5fb119da8e6",
        "number": 243,
        "project_slug": "gh/rivernews/iriversland2-kubernetes",
        "state": "created",
        "trigger": {
            "actor": {
                "avatar_url": "https://avatars1.githubusercontent.com/u/15918424?v=4",
                "login": "rivernews"
            },
            "received_at": "2020-06-16T17:04:12.904Z",
            "type": "api"
        },
        "updated_at": "2020-06-16T17:04:12.930Z",
        "vcs": {
            "branch": "destroy-release",
            "origin_repository_url": "https://github.com/rivernews/iriversland2-kubernetes",
            "provider_name": "GitHub",
            "revision": "a8dcdb66fdfd2413458d1b5b166c32e9ed1aa63f",
            "target_repository_url": "https://github.com/rivernews/iriversland2-kubernetes"
        }
    },
    ...
  ]
}
```

#### [Get a pipeline](https://circleci.com/docs/api/v2/#get-a-pipeline)

Doesn't do much.
Besides project slug, we need pipeline number. Get the number from `Get all pipelines` endpoint above, or **from the `POST pipeline` response**.

```sh
curl -Gv \
  --data-urlencode "circle-token=$CIRCLECI_TOKEN" \
  -H 'Accept: application/json' \
  https://circleci.com/api/v2/project/github%2Frivernews%2Firiversland2-kubernetes/pipeline/243 | python -mjson.tool
```

Response object:

```json
{
  "created_at": "2020-06-16T17:04:12.930Z",
  "errors": [],
  "id": "5c9ab317-3f41-4851-a3de-e5fb119da8e6",
  "number": 243,
  "project_slug": "gh/rivernews/iriversland2-kubernetes",
  "state": "created",
  "trigger": {
      "actor": {
          "avatar_url": "https://avatars1.githubusercontent.com/u/15918424?v=4",
          "login": "rivernews"
      },
      "received_at": "2020-06-16T17:04:12.904Z",
      "type": "api"
  },
  "updated_at": "2020-06-16T17:04:12.930Z",
  "vcs": {
      "branch": "destroy-release",
      "origin_repository_url": "https://github.com/rivernews/iriversland2-kubernetes",
      "provider_name": "GitHub",
      "revision": "a8dcdb66fdfd2413458d1b5b166c32e9ed1aa63f",
      "target_repository_url": "https://github.com/rivernews/iriversland2-kubernetes"
  }
}
```

#### Get all workflows of a pipeline

We won't use pipeline number - we use **pipeline id** instead, which can also be found in `POST pipeline` response.

The endpoint no longer needs project-slug; use `/pipeline/:id` instead. Mind that worflow is singular form, not plural.

We can conclude that pipeline number requires project slug, while pipeline uuid is used for locating pipeline object directly.

```sh
curl -Gv \
  --data-urlencode "circle-token=$CIRCLECI_TOKEN" \
  -H 'Accept: application/json' \
  https://circleci.com/api/v2/pipeline/5c9ab317-3f41-4851-a3de-e5fb119da8e6/workflow | python -mjson.tool
```

Response object:
```json
{
  "items": [
      {
          "created_at": "2020-06-16T17:04:13Z",
          "id": "c324f6b8-2a9a-401d-8656-32bf96f786b8",
          "name": "build-master",
          "pipeline_id": "5c9ab317-3f41-4851-a3de-e5fb119da8e6",
          "pipeline_number": 243,
          "project_slug": "gh/rivernews/iriversland2-kubernetes",
          "started_by": "f7b5bb29-fe45-4dfb-80d5-69064b0ea01f",
          "status": "success",
          "stopped_at": "2020-06-16T17:07:00Z"
      }
  ],
  "next_page_token": null
}
```

Mind that while we have two jobs (stages) for the pipeline, it's counted as within one workflow object. You can observe in `config.yml` that both jobs are under the single workflow `build-master`.

Here we can access the success / failure of a build. Can we get a `in progress` state? When we re-run the workflow and query again, we got response:

```json
{
  "items": [
      {
          "created_at": "2020-06-21T04:59:43Z",
          "id": "54e9a827-5893-4381-8453-880c0e104278",
          "name": "build-master",
          "pipeline_id": "5c9ab317-3f41-4851-a3de-e5fb119da8e6",
          "pipeline_number": 243,
          "project_slug": "gh/rivernews/iriversland2-kubernetes",
          "started_by": "f7b5bb29-fe45-4dfb-80d5-69064b0ea01f",
          "status": "running",
          "stopped_at": null
      },
      {
          "created_at": "2020-06-16T17:04:13Z",
          "id": "c324f6b8-2a9a-401d-8656-32bf96f786b8",
          "name": "build-master",
          "pipeline_id": "5c9ab317-3f41-4851-a3de-e5fb119da8e6",
          "pipeline_number": 243,
          "project_slug": "gh/rivernews/iriversland2-kubernetes",
          "started_by": "f7b5bb29-fe45-4dfb-80d5-69064b0ea01f",
          "status": "success",
          "stopped_at": "2020-06-16T17:07:00Z"
      }
  ],
  "next_page_token": null
}
```

You can see that you can just look at the first item - now status becomes `running`.

After we canceled the workflow, status becomes `canceled`.

A list of all possible values of status:

```
status	success
status	running
status	not_run
status	failed
status	error
status	failing
status	on_hold
status	canceled
```


### Travis CI

- Make request at `POST /repo/project-slug/requests/`
- Get request status at `GET /repo/project-slug/request/:id`
  - We should be able to get the build id here.
- Get build status at  `GET /build/:id`, of course, this `id` we'll be using build id, not the request id.
  - We get the `state` here. Possible values are
    - `passed`
    - or one of `:created, :received, :started, :passed, :failed, :errored, :canceled`, according to travis CI's [code base](https://github.com/travis-ci/travis-api/blob/master/lib/travis/model/build/states.rb#L25).

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
- Parsing http response
  - Parsing json into struct objects: https://stackoverflow.com/questions/33061117/in-golang-what-is-the-difference-between-json-encoding-and-marshalling
  - Using tag: https://blog.josephmisiti.com/parsing-json-responses-in-golang

## Working with Redigo - a redis golang client
- Tutorial [on parsing redis reply into struct](https://www.alexedwards.net/blog/working-with-redis).

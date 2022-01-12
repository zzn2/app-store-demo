# A demo API Server by golang

This is a demo API server created using [golang](https://go.dev/) and [gin-gonic](https://github.com/gin-gonic/gin).

It is an demo "app store", which stores metadata of apps and makes them queryable.
It does not use databases, only stores data in memory.

A sample app metadata looks like:

```yaml
title: App1
version: 0.0.1
maintainers:
- name: First Maintainer
  email: firstmaintainer@hotmail.com
- name: Second Maintainer
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
 ```


## Scenarios

### Create a new app

* Creates a new app.
```
POST /apps
```


### Get app metadata

* Get the app with specific title. If the app contains multiple versions, gets the latest version.
```
GET /apps/App1
```

* Get the app with specific title and version.
```
GET /apps/App1/versions/0.0.1
```

### List apps

* List all the apps.

> Paging/Field selecting functions are not supported in this demo.

```
GET /apps
```

* Filters could also be applied to search apps match the given rule set.

For example, the following query lists all the apps with title "App1" (possibly multiple versions can be listed.)
```
GET /apps?title=App1
```

Filters could be applied to multiple fields:
```
GET /apps?title=App1&versions=0.0.1
```

Filters can also be non-precise match by specifying [LHS Brackets](https://christiangiacomi.com/posts/rest-design-principles/). For example, the following query lists apps which title contains "App":
```
GET /apps?title[like]=App
```

Refer to [integration test scenarios](src/api_integration_test.go) for more use cases.


## Referred Links

A good tutorial for golang newbie
* https://go.dev/doc/tutorial/create-module

Modules and packages
* https://go.dev/blog/using-go-modules
* https://go.dev/blog/package-names
* https://rakyll.org/style-packages/

Error handling
* https://go.dev/blog/error-handling-and-go

Documentation
* https://go.dev/blog/godoc

Testing
* https://go.dev/doc/tutorial/add-a-test
* https://onsi.github.io/ginkgo/

Go source code, used as good examples for best practice
* https://go.dev/src/

Reflection and struct tags
* https://golangbot.com/reflection/
* https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go

Effective GO
* https://go.dev/doc/effective_go

Rest API design principals
* https://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api

Containerize
* https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e

K8s deployment & Helm
* https://docs.microsoft.com/en-us/learn/modules/aks-deploy-container-app/5-exercise-deploy-app
* https://www.studytonight.com/post/what-is-helm-in-the-kubernetes-world
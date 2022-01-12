# A demo API Server by golang

This is a demo API server created using [golang](https://go.dev/)uses it [gin-gonic](https://github.com/gin-gonic/gin) framework.

It is an demo "app store", which stores metadata of apps and makes them queryable.

It stores app metadata in memory only for demo usage.

A sample app meta looks like:

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

```
POST /apps
```

### Show a specific app (or optionally with a specific version)

```
GET /apps/App1
GET /apps/App1/versions/0.0.1
```

### List apps

List all the apps.

> Paging/Field selecting functions are not supported in this demo.

```
GET /apps
```

Filters could also be applied to search apps match the given rule set.

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

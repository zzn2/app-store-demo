# A demo API Server created by golang

This is a demo API server created by golang. It uses the [gin-gonic](https://github.com/gin-gonic/gin) web framework.

It is an demo "app store", which supports to store and search metadata of apps.

A sample app meta looks like:

```yaml
title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
- name: Firstname Lastname
  email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### blob of markdown
 More markdown
 ```


## Scenarios

### Create a new app

```
POST /apps
```

### Show a specific app (or optionally with a specific version)

```
GET /apps/App1
GET /apps/App1/version/0.0.1
```

### List apps

```
GET /apps
```

Also filters could be applied.
For example, the following query lists all the apps with title "App1" (possibly multiple versions will be listed.)
```
GET /apps?title=App1
```

Filters could be added to multiple fields:
```
GET /apps?title=App1&version=0.0.1
```

Filters can also be non-precise match by specifying with [LHS Brackets](https://christiangiacomi.com/posts/rest-design-principles/). For example, the following query will list apps which title contains "App":
```
GET /apps?title[like]=App
```

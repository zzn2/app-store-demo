The package is organized by responsibility
https://go.dev/blog/using-go-modules
https://rakyll.org/style-packages/

Package names convention
https://go.dev/blog/package-names


How to add method to a struct type?
https://www.geeksforgeeks.org/how-to-add-a-method-to-struct-type-in-golang/

A good tutorial for newbie
https://go.dev/doc/tutorial/create-module

Effective GO
https://go.dev/doc/effective_go

A name is exported if it begins with a capital letter. Otherwise, it is not exported.

Code executed as an application must be in a `main` package.

Structs explained
https://golangbot.com/structs/

How to add test
https://go.dev/doc/tutorial/add-a-test

Test framework
https://onsi.github.io/ginkgo/

Graceful shutdown:
https://github.com/gin-gonic/gin#graceful-shutdown-or-restart

Test REST APIs in VSCode
https://medium.com/refinitiv-developer-community/how-to-test-rest-api-with-visual-studio-code-rest-client-extensions-9f2e061d0299
https://github.com/Refinitiv-API-Samples/Article.RDP.VSCode.RESTClient

REST api design
https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/201

Fast iteration for gin
https://github.com/codegangsta/gin

QueryString design best practise
https://www.moesif.com/blog/technical/api-design/REST-API-Design-Best-Practices-for-Parameters-and-Query-String-Usage/

RESTAPI design principles
https://www.taniarascia.com/rest-api-sorting-filtering-pagination/#string-partial
https://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api

https://simonplend.com/how-to-build-filtering-into-your-rest-api/
https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/

GoLang enums
https://www.ribice.ba/golang-enums/

Godoc
https://go.dev/blog/godoc

Go source code
https://go.dev/src/

Error handling
https://go.dev/blog/error-handling-and-go

GO FAQ
https://go.dev/doc/faq#nil_error

Testing
https://pkg.go.dev/testing
https://gobyexample.com/testing

Reflection
https://golangbot.com/reflection/

Struct tags:
https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go

How to create a docker image
https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e

Why alpine
https://nickjanetakis.com/blog/the-3-biggest-wins-when-using-alpine-as-a-base-docker-image

K8s deployment 
https://docs.microsoft.com/en-us/learn/modules/aks-deploy-container-app/5-exercise-deploy-app

TODO:
* Add docs
* Add tests
* Error handling
* Add code coverage: https://github.com/gin-gonic/gin/pull/2998
* Add CI
* Add validator for post data
* Change gin to release mode: export GIN_MODE=release
* Handle null case of app store list
* Return the entity on POST call
* Provide Location header for POST response
* Add filter,search,sort,fields supports for list api
* Ensure to enable gzip
* Add support for paging
* Rate limiting
* Error code for error responses
* Document generation
* Try to calculate hash for existing app meta to avoid duplicate post.
* How to tell user errors vs system errors

* Do the exercise of go tour goroutine chapter

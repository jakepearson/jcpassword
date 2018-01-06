# Password hashing and encoding

Service to support requesting a password hash that can be retrieved on a separate `http` call.

# Getting Started

* Clone the [repository](https://github.com/jakepearson/jcpassword)
* Run the service with `go build main.go && main` or `go run main.go`
* To override the default port of `8080`, set the `PORT` environment variable
* All commits trigger a build job on [CircleCI](https://circleci.com/gh/jakepearson/jcpassword)
* If all tests pass, the source is pushed to [Heroku](https://dashboard.heroku.com/apps/jakepearson-password) for deployment to the [production instance](https://jakepearson-password.herokuapp.com/) (https://jakepearson-password.herokuapp.com/)
* Run the tests with `go test ./...`

# Routes

* `/hash` (POST): Request that a string is hashed with `SHA512` and `base64` encoded. To use, call to the endpoind with a `password` field defined in the query string. A `HashID` is returned that can be used to lookup the final result.
* `/hash/<HashID>` (GET): Request for the final result of the hash and encode operation. If the operation is complete, a status code of `200` will be returned with the hashed and encoded string in the body of the response. If the operation is not yet complete, a status code of `102` (Still Processing) will return with a message.
* `/shutdown` (GET): Request for the service to shutdown. The service will stop accepting requests, complete any open requests, and finally shutdown the process after `5` seconds.
* `/metrics` (GET): Request for simple metrics about the service. Returns a `json` string with fields for `TotalRequests` and `AverageResponseTime` in milliseconds.

# Future upgrades

* Change `/hash` endpoint to return the full `url` or `uri` to get the hash instead of just an `id`
* Replace in memory hash database with a DB that can survive restarts. We could then run multiple service instances and do rolling deploys
* Replace homegrown metrics system with something like [Prometheus](https://github.com/prometheus/prometheus)
* Switch from [Heroku](https://www.heroku.com/) logging to a different SaaS service with better performance and search capabilities
* Render output into `JUnit` format in [CircleCI](https://circleci.com/gh/jakepearson/jcpassword)
* Implement a [Swagger](https://swagger.io/) document and UI to make the service easier to discover and understand
# gomonit
![Go Version](https://img.shields.io/badge/go-1.8-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.9-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/gomonit)](https://goreportcard.com/report/github.com/Depado/gomonit)
[![Build Status](https://drone.depado.eu/api/badges/Depado/gomonit/status.svg)](https://drone.depado.eu/Depado/gomonit)
[![codecov](https://codecov.io/gh/Depado/gomonit/branch/master/graph/badge.svg)](https://codecov.io/gh/Depado/gomonit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/gomonit/blob/master/LICENSE)

Small soft to check if your services are running, providing a web interface.  

It also integrates with Github to fetch your commits, and with drone to fetch
the build list. Other versioning system and CI system will eventually be
supported in the future.

## Status and Goals

The implementation is simple. It simply checks if the root of the site answers with a 200 OK. It doesn't check anything else.

For now, this software uses a simple yaml configuration file, allowing users to define their hosts. This behavior is not documented for now. The final goal would be to have a more complex interface, allowing an admin user to add/remove/edit hosts using an admin interface and perhaps even check for some contents on the page.

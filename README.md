# WatchTask

[![GoDoc](https://godoc.org/github.com/Oliver-Fish/watchtask?status.png)](https://godoc.org/github.com/Oliver-Fish/watchtask)
[![Build Status](https://travis-ci.org/Oliver-Fish/watchtask.svg?branch=master)](https://travis-ci.org/Oliver-Fish/watchtask)
[![Go Report Card](https://goreportcard.com/badge/github.com/Oliver-Fish/watchtask)](https://goreportcard.com/report/github.com/Oliver-Fish/watchtask)

WatchTask is designed to be an ultra light task runner, it allows you to easily run defined commands on changes to defined paths.

## Install
```
go get github.com/Oliver-Fish/watchtask
```
or
```
go install github.com/Oliver-Fish/watchtask
```
## Usage
Run Command on changes in current directory
```
watchtask -c "go run main.go"
```
Run Multiple Commands on change
```
watchtask -c "go run main.go","echo TestTest"
```
Run Command on specified path 
```
watchtask -c "go run main.go" -p "~/webapp"
```

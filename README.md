# WatchTask
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
watchtask -c "go run main.go","echo TestTest
```
Run Command on specified path 
```
watchtask -c "go run main.go" -p "~/webapp"
```

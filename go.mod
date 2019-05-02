module github.com/bpoetzschke/bin.go

go 1.12

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.0.6

require (
	github.com/Sirupsen/logrus v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/nlopes/slack v0.5.0
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734 // indirect
	golang.org/x/sys v0.0.0-20190429190828-d89cdac9e872 // indirect
)

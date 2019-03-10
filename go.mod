module github.com/bryk-io/did-method

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/bryk-io/x v0.0.0-20190308202857-992b14f2aaf5
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.7.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kennygrant/sanitize v1.2.4
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	github.com/vmihailenco/msgpack v4.0.2+incompatible
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/crypto v0.0.0-20190131182504-b8fe1690c613
	golang.org/x/net v0.0.0-20190301231341-16b79f2e4e95
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
	google.golang.org/grpc v1.18.0
)

replace (
	github.com/dgraph-io/badger v1.5.5 => github.com/bryk-io/badger v1.5.5
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/bryk-io/go-grpc-middleware v1.0.1-0.20190202210917-0105da141832
)

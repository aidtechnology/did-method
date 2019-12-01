module github.com/bryk-io/did-method

require (
	github.com/bryk-io/x v0.0.0-20191201183948-176f35bef721
	github.com/gogo/googleapis v1.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/kennygrant/sanitize v1.2.4
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/crypto v0.0.0-20190927123631-a832865fa7ad
	google.golang.org/grpc v1.23.0
)

replace (
	github.com/dgraph-io/badger v1.5.5 => github.com/bryk-io/badger v1.5.5
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/bryk-io/go-grpc-middleware v1.0.1-0.20190615102816-71cc94c54bbc
)

go 1.13

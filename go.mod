module github.com/bryk-io/did-method

require (
	github.com/bryk-io/x v0.0.0-20190701133613-135e915a55a2
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1
	github.com/grpc-ecosystem/grpc-gateway v1.9.2
	github.com/kennygrant/sanitize v1.2.4
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
	google.golang.org/grpc v1.21.1
)

replace (
	github.com/dgraph-io/badger v1.5.5 => github.com/bryk-io/badger v1.5.5
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/bryk-io/go-grpc-middleware v1.0.1-0.20190615102816-71cc94c54bbc
)

go 1.13

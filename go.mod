module github.com/bryk-io/id

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/bryk-io/x v0.0.0-20190304144907-ed5b62ca5ffc
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kennygrant/sanitize v1.2.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	github.com/vmihailenco/msgpack v4.0.2+incompatible
	golang.org/x/crypto v0.0.0-20190131182504-b8fe1690c613
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
)

replace github.com/dgraph-io/badger v1.5.5 => github.com/bryk-io/badger v1.5.5

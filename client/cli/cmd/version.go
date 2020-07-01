package cmd

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/bryk-io/did-method/info"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		var components = map[string]string{
			"Version":    info.CoreVersion,
			"Build code": info.BuildCode,
			"OS/Arch":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			"Go version": runtime.Version(),
		}
		if info.BuildTimestamp != "" {
			st, err := strconv.ParseInt(info.BuildTimestamp, 10, 64)
			if err == nil {
				components["Release Date"] = time.Unix(st, 0).Format(time.RFC822)
			}
		}
		for k, v := range components {
			fmt.Printf("\033[21;37m%-13s:\033[0m %s\n", k, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

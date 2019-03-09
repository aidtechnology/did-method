package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Individual command parameters
type cParam struct {
	name      string
	usage     string
	flagKey   string
	byDefault interface{}
}

// Setup command parameters
func setupCommandParams(c *cobra.Command, params []cParam) (err error) {
	for _, p := range params {
		switch v := p.byDefault.(type) {
		case int:
			h := p.byDefault.(int)
			c.Flags().IntVar(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		case uint32:
			h := p.byDefault.(uint32)
			c.Flags().Uint32Var(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		case uint64:
			h := p.byDefault.(uint64)
			c.Flags().Uint64Var(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		case string:
			h := p.byDefault.(string)
			c.Flags().StringVar(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		case bool:
			h := p.byDefault.(bool)
			c.Flags().BoolVar(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		case []string:
			h := p.byDefault.([]string)
			c.Flags().StringSliceVar(&h, p.name, v, p.usage)
			if err = viper.BindPFlag(p.flagKey, c.Flags().Lookup(p.name)); err != nil {
				return err
			}
		}
	}
	return
}

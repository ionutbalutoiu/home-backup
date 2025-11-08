package config

import (
	"fmt"
	"strconv"
)

type DestResticParams struct {
	Repo     string `yaml:"repo"`
	KeepLast int    `yaml:"keep_last,omitempty"`
}

func (d *DestResticParams) ParseParams(params map[string]string) error {
	// repo
	if repo, ok := params["repo"]; ok {
		d.Repo = repo
	}
	// keep_last
	if keepLastStr, ok := params["keep_last"]; ok {
		keepLast, err := strconv.Atoi(keepLastStr)
		if err != nil {
			return fmt.Errorf("invalid int value for 'keep_last': %v", err)
		}
		d.KeepLast = keepLast
	}
	// validate params
	if err := d.validate(); err != nil {
		return err
	}
	return nil
}

func (d *DestResticParams) validate() error {
	if d.Repo == "" {
		return fmt.Errorf("restic destination 'repo' parameter is required")
	}
	if d.KeepLast < 0 {
		return fmt.Errorf("restic destination 'keep_last' parameter cannot be negative")
	} else if d.KeepLast == 0 {
		// set default value
		d.KeepLast = 10
	}
	return nil
}

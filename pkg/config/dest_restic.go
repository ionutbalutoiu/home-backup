package config

import (
	"fmt"
	"strconv"
)

const defaultKeepLast = 10
const defaultGroupBy = "host"

type DestResticParams struct {
	Repo     string `yaml:"repo"`
	KeepLast int    `yaml:"keep_last,omitempty"`
	GroupBy  string `yaml:"group_by,omitempty"`
}

func (d *DestResticParams) ParseParams(params map[string]string) error {
	if repo, ok := params["repo"]; ok {
		d.Repo = repo
	}
	if keepLastStr, ok := params["keep_last"]; ok {
		keepLast, err := strconv.Atoi(keepLastStr)
		if err != nil {
			return fmt.Errorf("invalid int value for 'keep_last': %w", err)
		}
		d.KeepLast = keepLast
	}
	if groupBy, ok := params["group_by"]; ok {
		d.GroupBy = groupBy
	}
	d.setDefaults()
	if err := d.validate(); err != nil {
		return err
	}
	return nil
}

func (d *DestResticParams) setDefaults() {
	if d.KeepLast == 0 {
		d.KeepLast = defaultKeepLast
	}
	if d.GroupBy == "" {
		d.GroupBy = defaultGroupBy
	}
}

func (d *DestResticParams) validate() error {
	if d.Repo == "" {
		return fmt.Errorf("restic destination 'repo' parameter is required")
	}
	if d.KeepLast < 0 {
		return fmt.Errorf("restic destination 'keep_last' parameter cannot be negative")
	}
	return nil
}

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type SrcDirectoryParams struct {
	Path string `yaml:"path"`
}

func (s *SrcDirectoryParams) ParseParams(params map[string]string) error {
	bytes, err := yaml.Marshal(params)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, s); err != nil {
		return err
	}
	if err := s.validate(); err != nil {
		return err
	}
	return nil
}

func (s *SrcDirectoryParams) validate() error {
	if s.Path == "" {
		return fmt.Errorf("directory source 'path' parameter is required")
	}
	stat, err := os.Stat(s.Path)
	if err != nil {
		return fmt.Errorf("directory source 'path' is not accessible: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("directory source 'path' is not a directory: %s", s.Path)
	}
	return nil
}

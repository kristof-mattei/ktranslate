package main

import (
	"context"

	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"

	"gopkg.in/yaml.v2"
)

// Loads config and uses this to override any flags.
func setConfig(cfg string) error {
	if cfg == "" { // Base case where not set.
		return nil
	}

	ms := kt.KTConf{}
	by, err := snmp_util.LoadFile(context.Background(), cfg)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(by, &ms)
	if err != nil {
		return err
	}

	return kt.ParseConfig(ms)
}

package plugin

import (
	"fmt"
	"log"

	"github.com/mojura/kiroku"
	s3 "github.com/mojura/sync-s3"
	"github.com/vroomy/vroomy"
)

var p Plugin

func init() {
	if err := vroomy.Register("mojura-source", &p); err != nil {
		log.Fatal(err)
	}
}

type Plugin struct {
	vroomy.BasePlugin

	source kiroku.Source
}

// Load will initialize the s3 client
func (p *Plugin) Load(env map[string]string) (err error) {
	var opts s3.Options
	opts.Key = env["s3-key"]
	opts.Secret = env["s3-secret"]
	opts.Bucket = env["s3-env"]
	opts.Region = env["s3-region"]

	if p.source, err = s3.New(opts); err != nil {
		err = fmt.Errorf("error s3 client: %v", err)
		return
	}

	return
}

// Backend exposes this plugin's data layer to other plugins
func (p *Plugin) Backend() interface{} {
	return p.source
}

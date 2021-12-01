package plugin

import (
	"fmt"
	"log"

	"github.com/gdbu/scribe"
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

	out *scribe.Scribe

	source kiroku.Source
}

// Load will initialize the s3 client
func (p *Plugin) Load(env map[string]string) (err error) {
	p.out = scribe.New("S3")
	switch env["mojura-sync-mode"] {
	case "development":
		p.out.Notification("Development mode enabled, disabling s3 DB syncing")
		p.source = &kiroku.NOOP{}
		return
	case "mirror":
		p.out.Notification("Mirror mode enabled")
	case "sync":
		p.out.Notification("Sync mode enabled, enabling s3 DB syncing")

	default:
		err = fmt.Errorf("invalid mode, <%s> is not supported, available modes are development, mirror, and sync", env["mode"])
		return
	}

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

package option

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

// Options contains start options
type Options struct {
	ConfigFile  string
	ConfigValue string
	Namespace   string
	Debug       bool
}

// NewOptions returns a new Options
func NewOptions() *Options {
	return &Options{}
}

// AddFlags add flags to app
func (opts *Options) AddFlags(app *cli.App) {

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Path to the kube config.",
			Destination: &opts.ConfigFile,
			Value:       filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		},
		cli.StringFlag{
			Name:        "configvalue",
			Usage:       "Kube config value in string. This can be set from environment value 'KUBECONFIG'",
			Destination: &opts.ConfigValue,
			EnvVar:      "KUBECONFIG",
		},
		cli.StringFlag{
			Name:        "namespace",
			Usage:       "Namespace to query.",
			Destination: &opts.Namespace,
			Value:       "default",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Run with debug mode",
			Destination: &opts.Debug,
		},
	}

	app.Flags = append(app.Flags, flags...)

}

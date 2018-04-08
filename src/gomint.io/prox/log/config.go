package log

import "flag"

var DebugEnabled = flag.Bool("debug", false, "--debug")

func init() {
	flag.Parse()
}
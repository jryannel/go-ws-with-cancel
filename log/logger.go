package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

var Debugf = log.Debugf
var Infof = log.Infof
var Warnf = log.Warnf
var Errorf = log.Errorf
var Fatalf = log.Fatalf
var Panicf = log.Panicf
var Debug = log.Debug
var Info = log.Info
var Warn = log.Warn
var Error = log.Error
var Fatal = log.Fatal
var Panic = log.Panic
var Debugln = log.Debugln
var Infoln = log.Infoln
var Warnln = log.Warnln
var Errorln = log.Errorln
var Fatalln = log.Fatalln
var Panicln = log.Panicln
var WithField = log.WithField
var WithFields = log.WithFields
var WithError = log.WithError

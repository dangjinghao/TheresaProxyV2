package Config

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"time"
)

func NewLogger(values logrus.Fields) *logrus.Entry {
	return logrus.WithFields(values)

}
func NewLoggerWithName(name string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"name": name,
	})

}
func init() {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		FieldsOrder:     []string{"name"},
	})
}

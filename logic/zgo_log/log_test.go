package zgo_log

import (
	"github.com/sirupsen/logrus"
	"golang-bootstrap/lib/log"
	"testing"
)

func TestT(t *testing.T) {
	//nl := T()
	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

}

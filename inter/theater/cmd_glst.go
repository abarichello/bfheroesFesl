package theater

import (
	"github.com/Synaxis/unstable/backend/inter/network"

	"github.com/sirupsen/logrus"
)

// GLST - CLIENT called to get a list of game servers? Irrelevant for heroes.
func (tm *Theater) GLST(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}
	logrus.Println("GLST was called")
}

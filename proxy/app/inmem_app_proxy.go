package app

import (
	"github.com/babbleio/babble/hashgraph"
	"github.com/sirupsen/logrus"
)

//InmemProxy is used for testing
type InmemAppProxy struct {
	submitCh    chan []byte
	commitedTxs [][]byte
	logger      *logrus.Logger
}

func NewInmemAppProxy(logger *logrus.Logger) *InmemAppProxy {
	if logger == nil {
		logger = logrus.New()
		logger.Level = logrus.DebugLevel
	}
	return &InmemAppProxy{
		submitCh:    make(chan []byte),
		commitedTxs: [][]byte{},
		logger:      logger,
	}
}

func (p *InmemAppProxy) SubmitCh() chan []byte {
	return p.submitCh
}

func (p *InmemAppProxy) CommitEvent(event hashgraph.Event) error {
	p.logger.WithField("event", event).Debug("InmemProxy CommitEvent")
	p.commitedTxs = append(p.commitedTxs, event.Body.Transactions...)
	return nil
}

//-------------------------------------------------------
//Implement AppProxy Interface

func (p *InmemAppProxy) SubmitTx(tx []byte) {
	p.submitCh <- tx
}

func (p *InmemAppProxy) GetCommittedTransactions() [][]byte {
	return p.commitedTxs
}

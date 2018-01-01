package proxy

import "github.com/babbleio/babble/hashgraph"

type AppProxy interface {
	SubmitCh() chan []byte
	CommitEvent(event hashgraph.Event) error
}

type BabbleProxy interface {
	CommitCh() chan hashgraph.Event
	SubmitTx(tx []byte) error
}

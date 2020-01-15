package mock

import "github.com/ElrondNetwork/elrond-go/p2p"

// NilAntifloodHandler is an empty implementation of P2PAntifloodHandler
// it does nothing
type NilAntifloodHandler struct {
}

func (nah *NilAntifloodHandler) ResetForTopic(topic string) {
}

func (nah *NilAntifloodHandler) SetMaxMessagesForTopic(topic string, maxNum uint32) {
}

// CanProcessMessage will always return nil, allowing messages to go to interceptors
func (nah *NilAntifloodHandler) CanProcessMessage(message p2p.MessageP2P, fromConnectedPeer p2p.PeerID) error {
	return nil
}

// CanProcessMessageOnTopic will always return nil, allowing messages to go to interceptors
func (nah *NilAntifloodHandler) CanProcessMessageOnTopic(message p2p.MessageP2P, fromConnectedPeer p2p.PeerID, topic string) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nah *NilAntifloodHandler) IsInterfaceNil() bool {
	return nah == nil
}

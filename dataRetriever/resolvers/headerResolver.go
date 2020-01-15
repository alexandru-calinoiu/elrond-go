package resolvers

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/typeConverters"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/logger"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/storage"
)

var log = logger.GetOrCreate("dataretriever/resolvers")

// HeaderResolver is a wrapper over Resolver that is specialized in resolving headers requests
type HeaderResolver struct {
	dataRetriever.TopicResolverSender
	headers          storage.Cacher
	hdrNonces        dataRetriever.Uint64SyncMapCacher
	hdrStorage       storage.Storer
	hdrNoncesStorage storage.Storer
	marshalizer      marshal.Marshalizer
	nonceConverter   typeConverters.Uint64ByteSliceConverter
	epochHandler     dataRetriever.EpochHandler
	shardCoordinator sharding.Coordinator
	antifloodHandler dataRetriever.P2PAntifloodHandler
}

// NewHeaderResolver creates a new header resolver
func NewHeaderResolver(
	senderResolver dataRetriever.TopicResolverSender,
	headers storage.Cacher,
	headersNonces dataRetriever.Uint64SyncMapCacher,
	hdrStorage storage.Storer,
	headersNoncesStorage storage.Storer,
	marshalizer marshal.Marshalizer,
	nonceConverter typeConverters.Uint64ByteSliceConverter,
	shardCoordinator sharding.Coordinator,
	antifloodHandler dataRetriever.P2PAntifloodHandler,
) (*HeaderResolver, error) {

	if check.IfNil(senderResolver) {
		return nil, dataRetriever.ErrNilResolverSender
	}
	if check.IfNil(headers) {
		return nil, dataRetriever.ErrNilHeadersDataPool
	}
	if check.IfNil(headersNonces) {
		return nil, dataRetriever.ErrNilHeadersNoncesDataPool
	}
	if check.IfNil(hdrStorage) {
		return nil, dataRetriever.ErrNilHeadersStorage
	}
	if check.IfNil(headersNoncesStorage) {
		return nil, dataRetriever.ErrNilHeadersNoncesStorage
	}
	if check.IfNil(marshalizer) {
		return nil, dataRetriever.ErrNilMarshalizer
	}
	if check.IfNil(nonceConverter) {
		return nil, dataRetriever.ErrNilUint64ByteSliceConverter
	}
	if check.IfNil(shardCoordinator) {
		return nil, dataRetriever.ErrNilShardCoordinator
	}
	if check.IfNil(antifloodHandler) {
		return nil, dataRetriever.ErrNilAntifloodHandler
	}

	hdrResolver := &HeaderResolver{
		TopicResolverSender: senderResolver,
		headers:             headers,
		hdrNonces:           headersNonces,
		hdrStorage:          hdrStorage,
		hdrNoncesStorage:    headersNoncesStorage,
		marshalizer:         marshalizer,
		nonceConverter:      nonceConverter,
		epochHandler:        &nilEpochHandler{},
		shardCoordinator:    shardCoordinator,
		antifloodHandler:    antifloodHandler,
	}

	return hdrResolver, nil
}

// SetEpochHandler sets the epoch handler for this component
func (hdrRes *HeaderResolver) SetEpochHandler(epochHandler dataRetriever.EpochHandler) error {
	if check.IfNil(epochHandler) {
		return dataRetriever.ErrNilEpochHandler
	}

	hdrRes.epochHandler = epochHandler
	return nil
}

func (hdrRes *HeaderResolver) topicNameForAntiflood() string {
	prefix := factory.ShardBlocksTopic
	if hdrRes.shardCoordinator.SelfId() > hdrRes.shardCoordinator.NumberOfShards() {
		prefix = factory.MetachainBlocksTopic
	}

	commIdentifier := hdrRes.shardCoordinator.CommunicationIdentifier(hdrRes.shardCoordinator.SelfId())
	suffix := hdrRes.TopicRequestSuffix()

	return prefix + commIdentifier + suffix
}

// ProcessReceivedMessage will be the callback func from the p2p.Messenger and will be called each time a new message was received
// (for the topic this validator was registered to, usually a request topic)
func (hdrRes *HeaderResolver) ProcessReceivedMessage(message p2p.MessageP2P, fromConnectedPeer p2p.PeerID) error {
	err := hdrRes.antifloodHandler.CanProcessMessage(message, fromConnectedPeer)
	if err != nil {
		log.Error("CanProcessMessage", "error", err.Error())
		return err
	}

	topic := hdrRes.topicNameForAntiflood()
	log.Error("topic on processed rcvd msg", "topic", topic)
	err = hdrRes.antifloodHandler.CanProcessMessageOnTopic(message, fromConnectedPeer, hdrRes.topicNameForAntiflood())
	if err != nil {
		return err
	}

	rd, err := hdrRes.parseReceivedMessage(message)
	if err != nil {
		return err
	}
	var buff []byte

	switch rd.Type {
	case dataRetriever.HashType:
		buff, err = hdrRes.resolveHeaderFromHash(rd.Value)
	case dataRetriever.NonceType:
		buff, err = hdrRes.resolveHeaderFromNonce(rd.Value)
	case dataRetriever.EpochType:
		buff, err = hdrRes.resolveHeaderFromEpoch(rd.Value)
	default:
		return dataRetriever.ErrResolveTypeUnknown
	}
	if err != nil {
		return err
	}
	if buff == nil {
		log.Trace("missing data",
			"data", rd)
		return nil
	}

	return hdrRes.Send(buff, message.Peer())
}

func (hdrRes *HeaderResolver) resolveHeaderFromNonce(key []byte) ([]byte, error) {
	// key is now an encoded nonce (uint64)

	// Search the nonce-key pair in cache-storage
	hash, err := hdrRes.hdrNoncesStorage.Get(key)
	if err != nil {
		log.Trace("hdrNoncesStorage.Get", "error", err.Error())
	}

	// Search the nonce-key pair in data pool
	if hash == nil {
		nonceBytes, err := hdrRes.nonceConverter.ToUint64(key)
		if err != nil {
			return nil, dataRetriever.ErrInvalidNonceByteSlice
		}

		value, ok := hdrRes.hdrNonces.Get(nonceBytes)
		if ok {
			value.Range(func(shardId uint32, existingHash []byte) bool {
				if shardId == hdrRes.TargetShardID() {
					hash = existingHash
					return false
				}

				return true
			})
		}

		if len(hash) == 0 {
			return nil, nil
		}
	}

	return hdrRes.resolveHeaderFromHash(hash)
}

// resolveHeaderFromHash resolves a header using its key (header hash)
func (hdrRes *HeaderResolver) resolveHeaderFromHash(key []byte) ([]byte, error) {
	value, ok := hdrRes.headers.Peek(key)
	if !ok {
		return hdrRes.hdrStorage.Get(key)
	}

	buff, err := hdrRes.marshalizer.Marshal(value)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

// resolveHeaderFromEpoch resolves a header using its key based on epoch
func (hdrRes *HeaderResolver) resolveHeaderFromEpoch(key []byte) ([]byte, error) {
	actualKey := key

	isUnknownEpoch, err := core.IsUnknownEpochIdentifier(key)
	if err != nil {
		return nil, err
	}

	if isUnknownEpoch {
		actualKey = []byte(core.EpochStartIdentifier(hdrRes.epochHandler.Epoch()))
	}

	return hdrRes.hdrStorage.Get(actualKey)
}

// parseReceivedMessage will transform the received p2p.Message in a RequestData object.
func (hdrRes *HeaderResolver) parseReceivedMessage(message p2p.MessageP2P) (*dataRetriever.RequestData, error) {
	rd := &dataRetriever.RequestData{}
	err := rd.Unmarshal(hdrRes.marshalizer, message)
	if err != nil {
		return nil, err
	}
	if rd.Value == nil {
		return nil, dataRetriever.ErrNilValue
	}

	return rd, nil
}

// RequestDataFromHash requests a header from other peers having input the hdr hash
func (hdrRes *HeaderResolver) RequestDataFromHash(hash []byte) error {
	return hdrRes.SendOnRequestTopic(&dataRetriever.RequestData{
		Type:  dataRetriever.HashType,
		Value: hash,
	})
}

// RequestDataFromNonce requests a header from other peers having input the hdr nonce
func (hdrRes *HeaderResolver) RequestDataFromNonce(nonce uint64) error {
	return hdrRes.SendOnRequestTopic(&dataRetriever.RequestData{
		Type:  dataRetriever.NonceType,
		Value: hdrRes.nonceConverter.ToByteSlice(nonce),
	})
}

// RequestDataFromEpoch requests a header from other peers having input the epoch
func (hdrRes *HeaderResolver) RequestDataFromEpoch(identifier []byte) error {
	return hdrRes.SendOnRequestTopic(&dataRetriever.RequestData{
		Type:  dataRetriever.EpochType,
		Value: identifier,
	})
}

// IsInterfaceNil returns true if there is no value under the interface
func (hdrRes *HeaderResolver) IsInterfaceNil() bool {
	return hdrRes == nil
}

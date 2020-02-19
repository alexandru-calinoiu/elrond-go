package resolverscontainer

import (
	"github.com/ElrondNetwork/elrond-go/core/random"
	triesFactory "github.com/ElrondNetwork/elrond-go/data/trie/factory"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/factory/containers"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/resolvers"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/resolvers/topicResolverSender"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

var _ dataRetriever.ResolversContainerFactory = (*shardResolversContainerFactory)(nil)

type shardResolversContainerFactory struct {
	*baseResolversContainerFactory
}

// NewShardResolversContainerFactory creates a new container filled with topic resolvers for shards
func NewShardResolversContainerFactory(
	args FactoryArgs,
) (*shardResolversContainerFactory, error) {
	if args.SizeCheckDelta > 0 {
		args.Marshalizer = marshal.NewSizeCheckUnmarshalizer(args.Marshalizer, args.SizeCheckDelta)
	}

	container := containers.NewResolversContainer()
	base := &baseResolversContainerFactory{
		container:                container,
		shardCoordinator:         args.ShardCoordinator,
		messenger:                args.Messenger,
		store:                    args.Store,
		marshalizer:              args.Marshalizer,
		dataPools:                args.DataPools,
		uint64ByteSliceConverter: args.Uint64ByteSliceConverter,
		intRandomizer:            &random.ConcurrentSafeIntRandomizer{},
		dataPacker:               args.DataPacker,
		triesContainer:           args.TriesContainer,
	}

	err := base.checkParams()
	if err != nil {
		return nil, err
	}

	return &shardResolversContainerFactory{
		baseResolversContainerFactory: base,
	}, nil
}

// Create returns a resolver container that will hold all resolvers in the system
func (srcf *shardResolversContainerFactory) Create() (dataRetriever.ResolversContainer, error) {
	err := srcf.generateTxResolvers(
		factory.TransactionTopic,
		dataRetriever.TransactionUnit,
		srcf.dataPools.Transactions(),
	)
	if err != nil {
		return nil, err
	}

	err = srcf.generateTxResolvers(
		factory.UnsignedTransactionTopic,
		dataRetriever.UnsignedTransactionUnit,
		srcf.dataPools.UnsignedTransactions(),
	)
	if err != nil {
		return nil, err
	}

	err = srcf.generateTxResolvers(
		factory.RewardsTransactionTopic,
		dataRetriever.RewardTransactionUnit,
		srcf.dataPools.RewardTransactions(),
	)
	if err != nil {
		return nil, err
	}

	err = srcf.generateHeaderResolvers()
	if err != nil {
		return nil, err
	}

	err = srcf.generateMiniBlocksResolvers()
	if err != nil {
		return nil, err
	}

	err = srcf.generatePeerChBlockBodyResolvers()
	if err != nil {
		return nil, err
	}

	err = srcf.generateMetablockHeaderResolvers()
	if err != nil {
		return nil, err
	}

	err = srcf.generateTrieNodesResolvers()
	if err != nil {
		return nil, err
	}

	return srcf.container, nil
}

//------- Hdr resolver

func (srcf *shardResolversContainerFactory) generateHeaderResolvers() error {
	shardC := srcf.shardCoordinator

	//only one shard header topic, for example: shardBlocks_0_META
	identifierHdr := factory.ShardBlocksTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(srcf.messenger, identifierHdr, emptyExcludePeersOnTopic)
	if err != nil {
		return err
	}

	hdrStorer := srcf.store.GetStorer(dataRetriever.BlockHeaderUnit)
	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		srcf.messenger,
		identifierHdr,
		peerListCreator,
		srcf.marshalizer,
		srcf.intRandomizer,
		shardC.SelfId(),
	)
	if err != nil {
		return err
	}

	hdrNonceHashDataUnit := dataRetriever.ShardHdrNonceHashDataUnit + dataRetriever.UnitType(shardC.SelfId())
	hdrNonceStore := srcf.store.GetStorer(hdrNonceHashDataUnit)
	resolver, err := resolvers.NewHeaderResolver(
		resolverSender,
		srcf.dataPools.Headers(),
		hdrStorer,
		hdrNonceStore,
		srcf.marshalizer,
		srcf.uint64ByteSliceConverter,
	)
	if err != nil {
		return err
	}
	//add on the request topic
	_, err = srcf.createTopicAndAssignHandler(
		identifierHdr+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierHdr, resolver)
}

//------- PeerChBlocks resolvers

func (srcf *shardResolversContainerFactory) generatePeerChBlockBodyResolvers() error {
	shardC := srcf.shardCoordinator

	//only one intrashard peer change blocks topic
	identifierPeerCh := factory.PeerChBodyTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	peerBlockBodyStorer := srcf.store.GetStorer(dataRetriever.PeerChangesUnit)

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(srcf.messenger, identifierPeerCh, emptyExcludePeersOnTopic)
	if err != nil {
		return err
	}

	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		srcf.messenger,
		identifierPeerCh,
		peerListCreator,
		srcf.marshalizer,
		srcf.intRandomizer,
		shardC.SelfId(),
	)
	if err != nil {
		return err
	}

	resolver, err := resolvers.NewGenericBlockBodyResolver(
		resolverSender,
		srcf.dataPools.MiniBlocks(),
		peerBlockBodyStorer,
		srcf.marshalizer,
	)
	if err != nil {
		return err
	}
	//add on the request topic
	_, err = srcf.createTopicAndAssignHandler(
		identifierPeerCh+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierPeerCh, resolver)
}

//------- MetaBlockHeaderResolvers

func (srcf *shardResolversContainerFactory) generateMetablockHeaderResolvers() error {
	shardC := srcf.shardCoordinator

	//only one metachain header block topic
	//this is: metachainBlocks
	identifierHdr := factory.MetachainBlocksTopic
	hdrStorer := srcf.store.GetStorer(dataRetriever.MetaBlockUnit)

	metaAndCrtShardTopic := factory.ShardBlocksTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	excludedPeersOnTopic := factory.TransactionTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(srcf.messenger, metaAndCrtShardTopic, excludedPeersOnTopic)
	if err != nil {
		return err
	}

	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		srcf.messenger,
		identifierHdr,
		peerListCreator,
		srcf.marshalizer,
		srcf.intRandomizer,
		sharding.MetachainShardId,
	)
	if err != nil {
		return err
	}

	hdrNonceStore := srcf.store.GetStorer(dataRetriever.MetaHdrNonceHashDataUnit)
	resolver, err := resolvers.NewHeaderResolver(
		resolverSender,
		srcf.dataPools.Headers(),
		hdrStorer,
		hdrNonceStore,
		srcf.marshalizer,
		srcf.uint64ByteSliceConverter,
	)
	if err != nil {
		return err
	}

	//add on the request topic
	_, err = srcf.createTopicAndAssignHandler(
		identifierHdr+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return err
	}

	return srcf.container.Add(identifierHdr, resolver)
}

func (srcf *shardResolversContainerFactory) generateTrieNodesResolvers() error {
	shardC := srcf.shardCoordinator

	keys := make([]string, 0)
	resolversSlice := make([]dataRetriever.Resolver, 0)

	identifierTrieNodes := factory.AccountTrieNodesTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	resolver, err := srcf.createTrieNodesResolver(identifierTrieNodes, triesFactory.UserAccountTrie)
	if err != nil {
		return err
	}

	resolversSlice = append(resolversSlice, resolver)
	keys = append(keys, identifierTrieNodes)

	identifierTrieNodes = factory.ValidatorTrieNodesTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	resolver, err = srcf.createTrieNodesResolver(identifierTrieNodes, triesFactory.PeerAccountTrie)
	if err != nil {
		return err
	}

	resolversSlice = append(resolversSlice, resolver)
	keys = append(keys, identifierTrieNodes)

	return srcf.container.AddMultiple(keys, resolversSlice)
}

// IsInterfaceNil returns true if there is no value under the interface
func (srcf *shardResolversContainerFactory) IsInterfaceNil() bool {
	return srcf == nil
}
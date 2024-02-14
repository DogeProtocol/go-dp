package proofofstake

import (
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"math/big"
	"time"
)

type PacketMap struct {
	round                 byte
	proposalDetailsMap    map[common.Address]*ProposalDetails
	proposalAckDetailsMap map[common.Address]*ProposalAckDetails
	precommitDetailsMap   map[common.Address]*PreCommitDetails
	commitDetailsMap      map[common.Address]*CommitDetails
}

func ParseConsensusPackets(parentHash common.Hash, consensusPackets *[]eth.ConsensusPacket, filteredValidatorDepositMap map[common.Address]*big.Int) (packetRoundMap map[byte]*PacketMap, err error) {
	packetRoundMap = make(map[byte]*PacketMap)

	packets := *consensusPackets
	for index, packet := range packets {
		if packet.ParentHash.IsEqualTo(parentHash) == false {
			return nil, errors.New("unexpected parenthash")
		}

		if packet.Signature == nil || packet.ConsensusData == nil || len(packet.Signature) == 0 || len(packet.ConsensusData) == 0 {
			return nil, errors.New("invalid consensus packet, nil data")
		}

		dataToVerify := append(packet.ParentHash.Bytes(), packet.ConsensusData...)
		digestHash := crypto.Keccak256(dataToVerify)
		pubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(digestHash, packet.Signature)
		if err != nil {
			return nil, err
		}
		if cryptobase.SigAlg.Verify(pubKey.PubData, digestHash, packet.Signature) == false {
			return nil, InvalidPacketErr
		}

		validator, err := cryptobase.SigAlg.PublicKeyToAddress(pubKey)
		if err != nil {
			log.Trace("invalid 3", "err", err)
			return nil, err
		}

		_, ok := filteredValidatorDepositMap[validator]
		if ok == false {
			return nil, errors.New("validator not part of block")
		}

		packetType := ConsensusPacketType(packet.ConsensusData[0])
		if packetType == CONSENSUS_PACKET_TYPE_PROPOSE_BLOCK {
			details := ProposalDetails{}

			err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
			if err != nil {
				return nil, err
			}

			if details.Round < byte(1) || details.Round > MAX_ROUND {
				return nil, errors.New("invalid round d")
			}

			blockProposer, err := getBlockProposer(parentHash, &filteredValidatorDepositMap, details.Round)
			if err != nil {
				return nil, err
			}
			if blockProposer.IsEqualTo(validator) == false {
				return nil, errors.New("invalid block proposer")
			}
			log.Trace("parseconsensuspackets propose", "details.Round", details.Round)
			_, ok := packetRoundMap[details.Round]
			if ok == false {
				packetRoundMap[details.Round] = &PacketMap{
					round:                 details.Round,
					proposalDetailsMap:    make(map[common.Address]*ProposalDetails),
					proposalAckDetailsMap: make(map[common.Address]*ProposalAckDetails),
					precommitDetailsMap:   make(map[common.Address]*PreCommitDetails),
					commitDetailsMap:      make(map[common.Address]*CommitDetails),
				}
			}
			packetMap := packetRoundMap[details.Round]
			pktTest, ok := packetMap.proposalDetailsMap[validator]
			if ok == true {
				log.Trace("duplicate proposal packet", "validator", validator, "details.Round", details.Round,
					"txn count", len(details.Txns), "pktTest.Round", pktTest.Round, "len(pktTest.Txns)", len(pktTest.Txns), "index", index, "len(*consensusPackets)", len(*consensusPackets))
				return nil, errors.New("duplicate proposal packet")
			} else {
				log.Trace("proposal packet", "validator", validator, "Round", details.Round, "count", len(details.Txns), "index", index)
			}
			proposalDetails := &ProposalDetails{
				Round: details.Round,
				Txns:  make([]common.Hash, len(details.Txns)),
			}
			for i, txn := range details.Txns {
				proposalDetails.Txns[i].CopyFrom(txn)
			}
			packetMap.proposalDetailsMap[validator] = proposalDetails
			packetRoundMap[details.Round] = packetMap
		} else if packetType == CONSENSUS_PACKET_TYPE_ACK_BLOCK_PROPOSAL {
			details := ProposalAckDetails{}

			err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
			if err != nil {
				return nil, err
			}

			if details.Round < byte(1) || details.Round > MAX_ROUND {
				return nil, errors.New("invalid round e")
			}

			_, ok := packetRoundMap[details.Round]
			if ok == false {
				packetRoundMap[details.Round] = &PacketMap{
					round:                 details.Round,
					proposalDetailsMap:    make(map[common.Address]*ProposalDetails),
					proposalAckDetailsMap: make(map[common.Address]*ProposalAckDetails),
					precommitDetailsMap:   make(map[common.Address]*PreCommitDetails),
					commitDetailsMap:      make(map[common.Address]*CommitDetails),
				}
			}

			packetMap := packetRoundMap[details.Round]
			_, ok = packetMap.proposalAckDetailsMap[validator]
			if ok == true {
				log.Warn("duplicate ack proposal packet", "validator", validator)
				return nil, errors.New("duplicate ack proposal packet")
			}
			proposalAckDetails := &ProposalAckDetails{
				Round:               details.Round,
				ProposalAckVoteType: details.ProposalAckVoteType,
			}
			proposalAckDetails.ProposalHash.CopyFrom(details.ProposalHash)
			if proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL && proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_OK {
				log.Trace("proposalAckDetails.ProposalAckVoteType", "ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType)
				return nil, errors.New("invalid vote type a")
			}

			if details.Round == MAX_ROUND && proposalAckDetails.ProposalAckVoteType != VOTE_TYPE_NIL {
				log.Trace("proposalAckDetails.ProposalAckVoteType", "ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType)
				return nil, errors.New("invalid vote type expecting nil")
			}

			packetMap.proposalAckDetailsMap[validator] = proposalAckDetails
			packetRoundMap[details.Round] = packetMap
		} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
			details := PreCommitDetails{}

			err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
			if err != nil {
				return nil, err
			}

			if details.Round < byte(1) || details.Round > MAX_ROUND {
				return nil, errors.New("invalid round c")
			}

			_, ok := packetRoundMap[details.Round]
			if ok == false {
				packetRoundMap[details.Round] = &PacketMap{
					round:                 details.Round,
					proposalDetailsMap:    make(map[common.Address]*ProposalDetails),
					proposalAckDetailsMap: make(map[common.Address]*ProposalAckDetails),
					precommitDetailsMap:   make(map[common.Address]*PreCommitDetails),
					commitDetailsMap:      make(map[common.Address]*CommitDetails),
				}
			}

			packetMap := packetRoundMap[details.Round]
			_, ok = packetMap.precommitDetailsMap[validator]
			if ok == true {
				log.Warn("duplicate precommit packet", "validator", validator)
				return nil, errors.New("duplicate precommit packet")
			}
			precommitDetails := &PreCommitDetails{
				Round: details.Round,
			}
			precommitDetails.PrecommitHash.CopyFrom(details.PrecommitHash)

			packetMap.precommitDetailsMap[validator] = precommitDetails
			packetRoundMap[details.Round] = packetMap
		} else if packetType == CONSENSUS_PACKET_TYPE_COMMIT_BLOCK {
			details := CommitDetails{}

			err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
			if err != nil {
				return nil, err
			}

			if details.Round < byte(1) || details.Round > MAX_ROUND {
				return nil, errors.New("invalid roun 4")
			}

			_, ok := packetRoundMap[details.Round]
			if ok == false {
				packetRoundMap[details.Round] = &PacketMap{
					round:                 details.Round,
					proposalDetailsMap:    make(map[common.Address]*ProposalDetails),
					proposalAckDetailsMap: make(map[common.Address]*ProposalAckDetails),
					precommitDetailsMap:   make(map[common.Address]*PreCommitDetails),
					commitDetailsMap:      make(map[common.Address]*CommitDetails),
				}
			}
			packetMap := packetRoundMap[details.Round]
			_, ok = packetMap.commitDetailsMap[validator]
			if ok == true {
				log.Warn("duplicate commit packet", "validator", validator)
				return nil, errors.New("duplicate commit packet")
			}
			commitDetails := &CommitDetails{
				Round: details.Round,
			}
			commitDetails.CommitHash.CopyFrom(details.CommitHash)

			packetMap.commitDetailsMap[validator] = commitDetails
			packetRoundMap[details.Round] = packetMap
		} else {
			return nil, errors.New("unknown packet type")
		}
	}

	return packetRoundMap, nil
}

func ValidatePackets(parentHash common.Hash, round byte, packetMap *PacketMap, voteType VoteType,
	filteredValidatorDepositMap *map[common.Address]*big.Int, totalBlockDepositValue *big.Int, minDepositRequired *big.Int, txns []common.Hash) error {
	valMap := *filteredValidatorDepositMap

	okVotesDepositValue := big.NewInt(0)
	nilVotesDepositValue := big.NewInt(0)

	var proposalHash common.Hash
	if voteType == VOTE_TYPE_OK {
		log.Trace("GetCombinedTxnHash a", "parentHash", parentHash, "round", round, "count", len(txns))
		proposalHash = GetCombinedTxnHash(parentHash, round, txns)
	} else {
		log.Trace("GetCombinedTxnHash b", "parentHash", parentHash, "round", round)
		proposalHash.CopyFrom(getNilVoteProposalHash(parentHash, round))
		if txns != nil && len(txns) > 0 {
			return errors.New("invalid transactions with nil vote")
		}
	}

	for v, proposalAckDetails := range packetMap.proposalAckDetailsMap {
		depositValue, ok := valMap[v]
		if ok == false {
			return errors.New("unrecognized validator")
		}

		if proposalAckDetails.Round != round {
			return errors.New("invalid round f")
		}
		log.Trace("val dep", "val", v, "depositValue", depositValue, "ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType, "ProposalHash", proposalAckDetails.ProposalHash)

		if proposalAckDetails.ProposalAckVoteType == VOTE_TYPE_NIL {
			if proposalAckDetails.ProposalHash.IsEqualTo(proposalHash) == false { //can be OK VOTE as well
				if voteType != VOTE_TYPE_OK { //can be ok VOTE as well
					log.Trace("proposal hash 2", "proposalHash", proposalHash, "proposalAckDetails.ProposalHash", proposalAckDetails.ProposalHash)
					return errors.New("invalid proposal hash")
				}
				continue
			}
			nilVotesDepositValue = common.SafeAddBigInt(nilVotesDepositValue, depositValue)
		} else if proposalAckDetails.ProposalAckVoteType == VOTE_TYPE_OK {
			if proposalAckDetails.ProposalHash.IsEqualTo(proposalHash) == false {
				if voteType != VOTE_TYPE_NIL { //can be NIL VOTE as well
					log.Trace("proposal hash 1", "proposalHash", proposalHash, "proposalAckDetails.ProposalHash", proposalAckDetails.ProposalHash, "voteType", voteType)
				}
				continue
			}
			okVotesDepositValue = common.SafeAddBigInt(okVotesDepositValue, depositValue)
		} else {
			return errors.New("invalid vote type b")
		}

		totalVotesDepositValue := common.SafeAddBigInt(nilVotesDepositValue, okVotesDepositValue)
		if totalVotesDepositValue.Cmp(totalBlockDepositValue) > 0 {
			return errors.New("invalid totalVotesDepositValue")
		}
	}

	log.Trace("ValidatePackets", "minDepositRequired", minDepositRequired, "okVotesDepositValue", okVotesDepositValue, "nilVotesDepositValue", nilVotesDepositValue,
		"voteType", voteType, "proposalAckDetails", len(packetMap.proposalAckDetailsMap), "txns", len(txns))
	var precommitHash common.Hash
	if voteType == VOTE_TYPE_NIL {
		if okVotesDepositValue.Cmp(minDepositRequired) >= 0 {
			return errors.New("VOTE_TYPE_NIL okVotesDepositValue error")
		}

		if nilVotesDepositValue.Cmp(minDepositRequired) < 0 {
			return errors.New("VOTE_TYPE_NIL nilVotesDepositValue error")
		}

		precommitHash = getNilVotePreCommitHash(parentHash, round)
	} else {
		if okVotesDepositValue.Cmp(minDepositRequired) < 0 {
			return errors.New("VOTE_TYPE_OK okVotesDepositValue error")
		}

		if nilVotesDepositValue.Cmp(minDepositRequired) >= 0 {
			return errors.New("VOTE_TYPE_OK nilVotesDepositValue error")
		}

		precommitHash = getOkVotePreCommitHash(parentHash, proposalHash, round)
	}
	commitHash := getCommitHash(precommitHash)
	precommitDepositValue := big.NewInt(0)
	for v, precommitDetails := range packetMap.precommitDetailsMap {
		depositValue, ok := valMap[v]
		if ok == false {
			return errors.New("unrecognized validator")
		}

		if precommitDetails.PrecommitHash.IsEqualTo(precommitHash) == false {
			errors.New("invalid precommithash")
		}

		precommitDepositValue = common.SafeAddBigInt(precommitDepositValue, depositValue)
	}

	if precommitDepositValue.Cmp(minDepositRequired) < 0 {
		return errors.New("precommit low deposit")
	}

	commitDepositValue := big.NewInt(0)
	for v, commitDetails := range packetMap.commitDetailsMap {
		depositValue, ok := valMap[v]
		if ok == false {
			return errors.New("unrecognized validator")
		}

		if commitDetails.CommitHash.IsEqualTo(commitHash) == false {
			errors.New("invalid commithash")
		}

		commitDepositValue = common.SafeAddBigInt(commitDepositValue, depositValue)
	}

	if commitDepositValue.Cmp(minDepositRequired) < 0 {
		return errors.New("precommit low deposit")
	}

	return nil
}

func ValidateBlockConsensusDataInner(txns []common.Hash, parentHash common.Hash, blockConsensusData *BlockConsensusData, blockAdditionalConsensusData *BlockAdditionalConsensusData,
	validatorDepositMap *map[common.Address]*big.Int) error {
	if blockConsensusData.Round < 1 {
		return errors.New("ValidateBlockConsensusData round min")
	}

	if blockConsensusData.Round >= MAX_ROUND && txns != nil && len(txns) > 0 { //todo: is this valid?
		return errors.New("ValidateBlockConsensusData round max")
	}

	if blockConsensusData.PrecommitHash.IsEqualTo(ZERO_HASH) {
		return errors.New("ValidateBlockConsensusData PrecommitHash zero_hash")
	}

	nilVotedProposers := make(map[common.Address]bool)
	if blockConsensusData.SlashedBlockProposers != nil {
		for _, proposer := range blockConsensusData.SlashedBlockProposers {
			nilVotedProposers[proposer] = true
			log.Trace("proposer slashed", "proposer", proposer)
		}
	}

	valMap := *validatorDepositMap
	filteredValidators, totalBlockDepositValue, minDepositRequired, err := filterValidators(parentHash, &valMap)
	if err != nil {
		return err
	}

	if MIN_BLOCK_DEPOSIT.Cmp(minDepositRequired) > 0 {
		return errors.New("min deposit required error")
	}

	if len(filteredValidators) < MIN_VALIDATORS {
		return errors.New("filteredValidators MIN_VALIDATORS")
	}

	if len(filteredValidators) > MAX_VALIDATORS {
		return errors.New("filteredValidators MAX_VALIDATORS")
	}

	var filteredValidatorDepositMap map[common.Address]*big.Int
	filteredValidatorDepositMap = make(map[common.Address]*big.Int)

	for v, _ := range filteredValidators {
		filteredValidatorDepositMap[v] = valMap[v]
	}

	roundBlockValidators := make(map[byte]common.Address)
	for r := byte(1); r <= blockConsensusData.Round; r++ {
		roundBlockValidators[r], err = getBlockProposer(parentHash, &filteredValidatorDepositMap, r)
		if err != nil {
			return err
		}
		log.Trace("roundBlockValidators[r]", "r", r, "roundBlockValidators[r]", roundBlockValidators[r])
	}

	if blockAdditionalConsensusData.ConsensusPackets == nil {
		return errors.New("nil ConsensusPackets")
	}

	packetRoundMap, err := ParseConsensusPackets(parentHash, &blockAdditionalConsensusData.ConsensusPackets, filteredValidatorDepositMap)
	if err != nil {
		return err
	}

	if blockConsensusData.VoteType == VOTE_TYPE_NIL {
		if len(txns) > 0 {
			return errors.New("txns in a NIL block")
		}
		if blockConsensusData.SelectedTransactions != nil && len(blockConsensusData.SelectedTransactions) > 0 {
			return errors.New("SelectedTransactions in a NIL vote")
		}
		if blockConsensusData.BlockProposer.IsEqualTo(ZERO_ADDRESS) == false {
			return errors.New("ValidateBlockConsensusData BlockProposer false")
		}

		if blockConsensusData.ProposalHash.IsEqualTo(getNilVoteProposalHash(parentHash, blockConsensusData.Round)) == false {
			return errors.New("proposal hash check failed")
		}

		precommitHash := getNilVotePreCommitHash(parentHash, blockConsensusData.Round)
		if blockConsensusData.PrecommitHash.IsEqualTo(precommitHash) == false {
			return errors.New("precommitHash hash check failed")
		}

		for r := byte(1); r <= blockConsensusData.Round; r++ {
			if r < MAX_ROUND {
				_, ok := nilVotedProposers[roundBlockValidators[r]]
				if ok == false {
					log.Trace("NilVotesProposer 1", "roundBlockValidators[r]", roundBlockValidators[r], "r", r, "parentHash", parentHash)
					return errors.New("nilVotedProposers 1")
				}
			}

			_, ok := packetRoundMap[r]
			if ok == false {
				log.Trace("could not find packetMap for round", "r", r)
				return errors.New("could not find packetMap for round")
			}
		}

		if len(nilVotedProposers) > int(blockConsensusData.Round) {
			return errors.New("unexpected number of nilVotedProposers")
		}

		packetMap := packetRoundMap[blockConsensusData.Round]
		err = ValidatePackets(parentHash, blockConsensusData.Round, packetMap, VOTE_TYPE_NIL, &filteredValidatorDepositMap, totalBlockDepositValue, minDepositRequired, blockConsensusData.SelectedTransactions)
		if err != nil {
			return err
		}

		//todo: deep validate block proposers
	} else if blockConsensusData.VoteType == VOTE_TYPE_OK {
		if blockConsensusData.BlockProposer.IsEqualTo(ZERO_ADDRESS) {
			return errors.New("ValidateBlockConsensusData BlockProposer true")
		}

		if blockConsensusData.BlockProposer.IsEqualTo(roundBlockValidators[blockConsensusData.Round]) == false {
			return errors.New("ValidateBlockConsensusData BlockProposer true")
		}

		if blockConsensusData.ProposalHash.IsEqualTo(ZERO_HASH) {
			log.Trace("ValidateBlockConsensusData ProposalHash zero_hash")
			return errors.New("ValidateBlockConsensusData ProposalHash zero_hash")
		}

		if blockConsensusData.SelectedTransactions == nil {
			if len(txns) > 0 {
				return errors.New("ValidateBlockConsensusData txns is non-empty but SelectedTransactions is nil")
			}
		} else {
			var selectedTxnsMap map[common.Hash]bool
			selectedTxnsMap = make(map[common.Hash]bool)
			for _, txn := range blockConsensusData.SelectedTransactions {
				selectedTxnsMap[txn] = true
			}
			if txns != nil {
				for _, txn := range txns {
					_, ok := selectedTxnsMap[txn]
					if ok == false {
						log.Trace("ValidateBlockConsensusData txn", "txn", txn)
						return errors.New("ValidateBlockConsensusData txns should be a subset of blockConsensusData.SelectedTransactions")
					}
				}
			}
		}

		if blockConsensusData.Round > 1 {
			for r := byte(1); r < blockConsensusData.Round; r++ {
				_, ok := nilVotedProposers[roundBlockValidators[r]]
				if ok == false {
					log.Trace("NilVotesProposer 2", "roundBlockValidators[r]", roundBlockValidators[r], "r", r, "parentHash", parentHash)
					return errors.New("nilVotedProposers 2")
				}
			}
			if len(blockConsensusData.SlashedBlockProposers) < int(blockConsensusData.Round-1) {
				log.Trace("SlashedBlockProposers", "len(nilVotedProposers)", len(nilVotedProposers), "int(blockConsensusData.Round)", int(blockConsensusData.Round))
				return errors.New("ValidateBlockConsensusData SlashedBlockProposers length")
			}
		}

		packetMap := packetRoundMap[blockConsensusData.Round]
		err = ValidatePackets(parentHash, blockConsensusData.Round, packetMap, VOTE_TYPE_OK, &filteredValidatorDepositMap, totalBlockDepositValue, minDepositRequired, blockConsensusData.SelectedTransactions)
		if err != nil {
			return err
		}

	} else {
		log.Trace("ValidateBlockConsensusData unexpected vote type", "vote type", blockConsensusData.VoteType)
		return errors.New("ValidateBlockConsensusData unexpected vote type")
	}

	return nil
}

// In this function, absolute time cannot be validated, since this function can get called at a different time, for example when new node is created and is reading old blocks
// Hence only basic checks are allowed
func ValidateBlockProposalTime(blockNumber uint64, proposedTime uint64) bool {
	if blockNumber == 1 || blockNumber%BLOCK_PERIOD_TIME_CHANGE == 0 {
		if proposedTime == 0 {
			return true
		}

		tm := time.Unix(int64(proposedTime), 0)
		if tm.Second() != 0 || tm.Nanosecond() != 0 { //No granularity at anything other than minute level allowed, to reduce ability to manipulate blockHash
			return false
		}
	} else {
		if proposedTime != 0 {
			return false
		}
	}

	return true
}

func ValidateBlockConsensusData(block *types.Block, validatorDepositMap *map[common.Address]*big.Int) error {
	header := block.Header()

	if header.ConsensusData == nil || header.UnhashedConsensusData == nil {
		return errors.New("ValidateBlockConsensusData nil")
	}

	blockConsensusData := &BlockConsensusData{}
	err := rlp.DecodeBytes(header.ConsensusData, blockConsensusData)
	if err != nil {
		return err
	}

	if blockConsensusData.SlashedBlockProposers == nil || blockConsensusData.SelectedTransactions == nil {
		return errors.New("ValidateBlockConsensusData SlashedBlockProposers or SelectedTransactions is nil")
	}

	blockAdditionalConsensusData := &BlockAdditionalConsensusData{}
	err = rlp.DecodeBytes(header.UnhashedConsensusData, blockAdditionalConsensusData)
	if err != nil {
		return err
	}

	if blockAdditionalConsensusData.ConsensusPackets == nil {
		return errors.New("ValidateBlockConsensusData ConsensusPackets is nil")
	}

	txns := block.Transactions()
	var txnList []common.Hash
	if txns != nil {
		txnList = make([]common.Hash, len(txns))
		for i, t := range txns {
			txnList[i].CopyFrom(t.Hash())
		}
	} else {
		txnList = make([]common.Hash, 0)
	}

	if ValidateBlockProposalTime(block.Number().Uint64(), blockConsensusData.BlockTime) == false {
		log.Warn("ValidateBlockProposalTime failed", "blockNumber", block.Number().Uint64(), "proposedTime", blockConsensusData.BlockTime)
		return errors.New("ValidateBlockProposalTime failed")
	}

	return ValidateBlockConsensusDataInner(txnList, header.ParentHash, blockConsensusData, blockAdditionalConsensusData, validatorDepositMap)
}

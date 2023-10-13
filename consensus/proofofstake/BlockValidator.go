package proofofstake

import (
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"math/big"
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
			fmt.Println("invalid 3", err)
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

			if details.Round < byte(1) || details.Round > MAX_ROUND_WITH_TXNS {
				return nil, errors.New("invalid round")
			}

			blockProposer, err := getBlockProposer(parentHash, &filteredValidatorDepositMap, details.Round)
			if err != nil {
				return nil, err
			}
			if blockProposer.IsEqualTo(validator) == false {
				return nil, errors.New("invalid block proposer")
			}
			//fmt.Println("parseconsensuspackets propose", details.Round)
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
				log.Warn("duplicate proposal packet", "validator", validator)
				fmt.Println("duplicate proposal packet", validator, details.Round, len(details.Txns), pktTest.Round, len(pktTest.Txns), index, len(*consensusPackets))
				return nil, errors.New("duplicate proposal packet")
			} else {
				//fmt.Println("proposal packet", validator, details.Round, len(details.Txns), index)
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

			if details.Round < byte(1) || details.Round > MAX_ROUND_WITH_TXNS {
				return nil, errors.New("invalid round")
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
			//fmt.Println("parseconsensuspackets ackProposal", details.Round)
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
				fmt.Println("proposalAckDetails.ProposalAckVoteType", proposalAckDetails.ProposalAckVoteType)
				return nil, errors.New("invalid vote type")
			}

			packetMap.proposalAckDetailsMap[validator] = proposalAckDetails
			packetRoundMap[details.Round] = packetMap
		} else if packetType == CONSENSUS_PACKET_TYPE_PRECOMMIT_BLOCK {
			details := PreCommitDetails{}

			err := rlp.DecodeBytes(packet.ConsensusData[1:], &details)
			if err != nil {
				return nil, err
			}

			if details.Round < byte(1) || details.Round > MAX_ROUND_WITH_TXNS {
				return nil, errors.New("invalid round")
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
			//fmt.Println("parseconsensuspackets precommit", details.Round)
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

			if details.Round < byte(1) || details.Round > MAX_ROUND_WITH_TXNS {
				return nil, errors.New("invalid round")
			}
			//fmt.Println("parseconsensuspackets commit", details.Round)
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

	//fmt.Println("parseconsensuspackets", len(*consensusPackets))

	return packetRoundMap, nil
}

func ValidatePackets(parentHash common.Hash, round byte, packetMap *PacketMap, voteType VoteType,
	filteredValidatorDepositMap *map[common.Address]*big.Int, totalBlockDepositValue *big.Int, minDepositRequired *big.Int, txns []common.Hash) error {
	valMap := *filteredValidatorDepositMap

	okVotesDepositValue := big.NewInt(0)
	nilVotesDepositValue := big.NewInt(0)

	var proposalHash common.Hash
	if voteType == VOTE_TYPE_OK {
		proposalHash = GetCombinedTxnHash(parentHash, round, txns)
	} else {
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
			return errors.New("invalid round")
		}

		if proposalAckDetails.ProposalAckVoteType == VOTE_TYPE_NIL {
			if proposalAckDetails.ProposalHash.IsEqualTo(ZERO_HASH) == false {
				return errors.New("invalid proposal hash, not zero hash when vote is nil")
			}
			nilVotesDepositValue = common.SafeAddBigInt(nilVotesDepositValue, depositValue)
		} else if proposalAckDetails.ProposalAckVoteType == VOTE_TYPE_OK {
			if proposalAckDetails.ProposalHash.IsEqualTo(proposalHash) == false {
				return errors.New("invalid proposal hash")
			}
			okVotesDepositValue = common.SafeAddBigInt(okVotesDepositValue, depositValue)
		} else {
			return errors.New("invalid vote type")
		}

		totalVotesDepositValue := common.SafeAddBigInt(nilVotesDepositValue, okVotesDepositValue)
		if totalVotesDepositValue.Cmp(totalBlockDepositValue) > 0 {
			return errors.New("invalid totalVotesDepositValue")
		}
	}

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

	if blockConsensusData.Round > MAX_ROUND_WITH_TXNS {
		return errors.New("ValidateBlockConsensusData round max")
	}

	if blockConsensusData.PrecommitHash.IsEqualTo(ZERO_HASH) {
		return errors.New("ValidateBlockConsensusData PrecommitHash zero_hash")
	}

	nilVotedProposers := make(map[common.Address]bool)
	if blockConsensusData.NilvotedBlockProposers != nil {
		for _, proposer := range blockConsensusData.NilvotedBlockProposers {
			nilVotedProposers[proposer] = true
			fmt.Println("proposer nilvoted", proposer)
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
		//fmt.Println("roundBlockValidators[r]", r, roundBlockValidators[r])
	}

	if blockAdditionalConsensusData.ConsensusPackets == nil {
		return errors.New("nil ConsensusPackets")
	}

	packetRoundMap, err := ParseConsensusPackets(parentHash, &blockAdditionalConsensusData.ConsensusPackets, filteredValidatorDepositMap)
	if err != nil {
		return err
	}
	//fmt.Println("packetRoundMap", len(packetRoundMap))

	if blockConsensusData.VoteType == VOTE_TYPE_NIL {
		if blockConsensusData.BlockProposer.IsEqualTo(ZERO_ADDRESS) == false {
			return errors.New("ValidateBlockConsensusData BlockProposer false")
		}

		if blockConsensusData.ProposalHash.IsEqualTo(ZERO_HASH) == false {
			return errors.New("proposal hash check failed")
		}

		precommitHash := crypto.Keccak256Hash(parentHash.Bytes(), ZERO_HASH.Bytes(), []byte{blockConsensusData.Round}, []byte{byte(VOTE_TYPE_NIL)})
		if blockConsensusData.PrecommitHash.IsEqualTo(precommitHash) == false {
			return errors.New("precommitHash hash check failed")
		}

		for r := byte(1); r <= blockConsensusData.Round; r++ {
			_, ok := nilVotedProposers[roundBlockValidators[r]]
			if ok == false {
				fmt.Println("NilVotesProposer 1", roundBlockValidators[r], r, parentHash)
				return errors.New("nilVotedProposers 1")
			}

			_, ok = packetRoundMap[r]
			if ok == false {
				fmt.Println("could not find packetMap for round", r)
				return errors.New("could not find packetMap for round")
			}
		}

		if len(nilVotedProposers) > int(blockConsensusData.Round) {
			return errors.New("unexpected number of nilVotedProposers")
		}

		packetMap := packetRoundMap[blockConsensusData.Round]
		err = ValidatePackets(parentHash, blockConsensusData.Round, packetMap, VOTE_TYPE_NIL, &filteredValidatorDepositMap, totalBlockDepositValue, minDepositRequired, txns)
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
			fmt.Println("ValidateBlockConsensusData ProposalHash zero_hash")
			return errors.New("ValidateBlockConsensusData ProposalHash zero_hash")
		}

		if blockConsensusData.Round > 1 {
			for r := byte(1); r < blockConsensusData.Round; r++ {
				_, ok := nilVotedProposers[roundBlockValidators[r]]
				if ok == false {
					fmt.Println("NilVotesProposer 2", roundBlockValidators[r], r, parentHash)
					return errors.New("nilVotedProposers 2")
				}
			}
			if len(blockConsensusData.NilvotedBlockProposers) < int(blockConsensusData.Round-1) {
				fmt.Println("NilvotedBlockProposers", len(nilVotedProposers), int(blockConsensusData.Round))
				return errors.New("ValidateBlockConsensusData NilvotedBlockProposers length")
			}
		}

		packetMap := packetRoundMap[blockConsensusData.Round]
		err = ValidatePackets(parentHash, blockConsensusData.Round, packetMap, VOTE_TYPE_OK, &filteredValidatorDepositMap, totalBlockDepositValue, minDepositRequired, txns)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("ValidateBlockConsensusData unexpected vote type")
		return errors.New("ValidateBlockConsensusData unexpected vote type")
	}

	return nil
}

func ValidateBlockConsensusData(block *types.Block, validatorDepositMap *map[common.Address]*big.Int) error {
	header := block.Header()

	if header.ConsensusData == nil || header.UnhashedConsensusData == nil {
		return errors.New("ValidateBlockConsensusData nil")
	}

	//fmt.Println("==================>ValidateBlockConsensusData", len(header.ConsensusData), len(header.UnhashedConsensusData))

	blockConsensusData := &BlockConsensusData{}
	err := rlp.DecodeBytes(header.ConsensusData, blockConsensusData)
	if err != nil {
		return err
	}

	blockAdditionalConsensusData := &BlockAdditionalConsensusData{}
	err = rlp.DecodeBytes(header.UnhashedConsensusData, blockAdditionalConsensusData)
	if err != nil {
		return err
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

	return ValidateBlockConsensusDataInner(txnList, header.ParentHash, blockConsensusData, blockAdditionalConsensusData, validatorDepositMap)
}

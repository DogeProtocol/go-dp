package backupmanager

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/rawdb"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/ethdb"
	"github.com/DogeProtocol/dp/log"
	"path/filepath"
	"sync"
)

type BackupManager struct {
	backupDir     string
	txBackupLock  sync.Mutex
	blkBackupLock sync.Mutex
	blockdb       *ethdb.Database
	txndb         *ethdb.Database
}

var singleInstance *BackupManager

func GetInstance() *BackupManager {
	return singleInstance
}

var instanceLock sync.Mutex

func NewBackupManager(backupDir string) (*BackupManager, error) {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	if singleInstance != nil {
		return singleInstance, nil
	}

	bm := &BackupManager{}

	err := bm.Initialize(backupDir)
	if err != nil {
		return nil, err
	}

	singleInstance = bm
	return bm, nil
}

func (b *BackupManager) Initialize(backupDir string) error {
	log.Debug("Initialize backup", "backupDir", backupDir)

	blockdbFilePath := filepath.Join(backupDir, "blockbackup.db")
	var blkdb ethdb.Database
	blkdb, err := rawdb.NewLevelDBDatabase(blockdbFilePath, 32, 0, "", false)
	if err != nil {
		return err
	}

	txndbFilePath := filepath.Join(backupDir, "txnbackup.db")
	var txndb ethdb.Database
	txndb, err = rawdb.NewLevelDBDatabase(txndbFilePath, 64, 0, "", false)
	if err != nil {
		return err
	}

	b.backupDir = backupDir
	b.blockdb = &blkdb
	b.txndb = &txndb

	return nil
}

func (b *BackupManager) BackupTransaction(tx *types.Transaction) error {
	b.txBackupLock.Lock()
	defer b.txBackupLock.Unlock()

	var buff bytes.Buffer
	buffWriter := bufio.NewWriter(&buff)

	err := tx.EncodeRLP(buffWriter)
	if err != nil {
		return err
	}
	err = buffWriter.Flush()
	if err != nil {
		return err
	}

	db := *b.txndb
	err = db.Put(tx.Hash().Bytes(), buff.Bytes())
	if err != nil {
		return err
	}

	log.Trace("BackupTransaction", "tx", tx.Hash())
	return nil
}

func (b *BackupManager) BackupBlock(blk *types.Block) error {
	b.blkBackupLock.Lock()
	defer b.blkBackupLock.Unlock()

	for _, tx := range blk.Transactions() {
		err := b.BackupTransaction(tx)
		if err != nil {
			return err
		}
	}

	var buff bytes.Buffer
	buffWriter := bufio.NewWriter(&buff)

	err := blk.EncodeRLP(buffWriter)
	if err != nil {
		return err
	}
	err = buffWriter.Flush()
	if err != nil {
		return err
	}

	db := *b.blockdb
	err = db.Put(blk.Hash().Bytes(), buff.Bytes())
	if err != nil {
		return err
	}

	//Mapping from block number to hash
	blkNumberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumberBytes, blk.NumberU64())
	blkNumberHash := crypto.Keccak256(blkNumberBytes)
	err = db.Put(blkNumberHash, blk.Hash().Bytes())
	if err != nil {
		return err
	}

	log.Trace("BackupBlock", "number", blk.Number(), "hash", blk.Hash())
	return nil
}

func (b *BackupManager) BlockExists(hash common.Hash) error {
	b.blkBackupLock.Lock()
	defer b.blkBackupLock.Unlock()

	db := *b.blockdb
	_, err := db.Get(hash.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (b *BackupManager) GetBlock(hash common.Hash) (*types.Block, error) {
	b.blkBackupLock.Lock()
	defer b.blkBackupLock.Unlock()

	db := *b.blockdb
	blockBytes, err := db.Get(hash.Bytes())
	if err != nil {
		return nil, err
	}

	return types.DecodeBlockFromRLP(blockBytes)
}

func (b *BackupManager) GetBlockHash(number uint64) (common.Hash, error) {
	b.blkBackupLock.Lock()
	defer b.blkBackupLock.Unlock()

	blkNumberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumberBytes, number)
	blkNumberHash := crypto.Keccak256(blkNumberBytes)

	db := *b.blockdb
	blockHashBytes, err := db.Get(blkNumberHash)
	if err != nil {
		return common.ZERO_HASH, err
	}

	if len(blockHashBytes) != len(common.ZERO_HASH.Bytes()) {
		return common.ZERO_HASH, errors.New("block hash length mismatch")
	}

	return common.BytesToHash(blockHashBytes), nil
}

func (b *BackupManager) TrsansactionExists(hash common.Hash) error {
	b.txBackupLock.Lock()
	defer b.txBackupLock.Unlock()

	db := *b.txndb
	_, err := db.Get(hash.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (b *BackupManager) Close() error {
	b.blkBackupLock.Lock()
	defer b.blkBackupLock.Unlock()

	b.txBackupLock.Lock()
	defer b.txBackupLock.Unlock()

	blkdb := *b.blockdb
	err := blkdb.Close()
	if err != nil {
		log.Debug("backup manager blockdb close error", "err", err)
		return err
	}

	txndb := *b.txndb
	err = txndb.Close()
	log.Debug("backup manager txndb close error", "err", err)
	if err != nil {
		return err
	}

	return nil
}

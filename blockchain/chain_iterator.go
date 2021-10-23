package blockchain

import "github.com/syndtr/goleveldb/leveldb"

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *leveldb.DB
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	data, err := iter.Database.Get(iter.CurrentHash,nil)
	Handle(err)
	block = block.Parse(data)
	iter.CurrentHash = block.PrevBlock
	return block
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}
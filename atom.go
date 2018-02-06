package atom

type Ledger string

const (
	LedgerBitcoin         Ledger = "bitcoin"
	LedgerBitcoinTestnet  Ledger = "bitcoin.testnet"
	LedgerEthereum        Ledger = "ethereum"
	LedgerEthereumRopsten Ledger = "ethereum.ropsten"
)

type LedgerAddress string

type Atom struct {
	ID         []byte
	Lock       []byte
	Fst        Ledger
	FstAddress LedgerAddress
	Snd        Ledger
	SndAddress LedgerAddress
}

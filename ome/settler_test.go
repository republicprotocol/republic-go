package ome_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/testutils"
)

const NumberOfNodes = 24

var _ = Describe("Settler", func() {
	var (
		storers   [NumberOfNodes]ComputationStorer
		smpcers   [NumberOfNodes]*testutils.Smpc
		contracts [NumberOfNodes]*omeBinder
		settles   [NumberOfNodes]Settler
	)

	Context("when a computation has been resolved to a match and been confirmed ", func() {
		BeforeEach(func() {
			for i := 0; i < NumberOfNodes; i++ {
				storer, err := leveldb.NewStore(fmt.Sprintf("./data-%v.out", i), time.Hour)
				Expect(err).ShouldNot(HaveOccurred())
				storers[i] = storer.SomerComputationStore()
				smpcers[i] = testutils.NewAlwaysMatchSmpc()
				contracts[i] = newOmeBinder()
				settles[i] = NewSettler(storers[i], smpcers[i], contracts[i], 0)
			}
		})

		AfterEach(func() {
			for i := 0; i < NumberOfNodes; i++ {
				os.RemoveAll(fmt.Sprintf("./data-%v.out", i))
			}
		})

		XIt("should be able to reconstruct the order and settle it.", func() {
			buyFragments, err := testutils.RandomBuyOrderFragments(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Expect(err).ShouldNot(HaveOccurred())
			sellFragments, err := testutils.RandomSellOrderFragments(int64(NumberOfNodes), int64(2*(NumberOfNodes+1)/3))
			Expect(err).ShouldNot(HaveOccurred())
			comp := NewComputation([32]byte{}, buyFragments[0], sellFragments[0], ComputationStateNil, true)
			for i := 0; i < NumberOfNodes; i++ {
				Expect(storers[i].PutComputation(comp)).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Expect(settles[i].Settle(comp)).ShouldNot(HaveOccurred())
			}

			for i := 0; i < NumberOfNodes; i++ {
				Expect(contracts[i].SettleCounts()).Should(Equal(1))
			}
		})
	})

	Context("volumeInEth", func() {
		buyPrice := uint64(40000000)
		buyVolume := uint64(80000000000)
		sellPrice := uint64(400000000)
		sellVolume := uint64(900000000000)

		It("calculates volume in BTC/ETH", func() {
			pair := testutils.TokensBTCETH

			expected := buyVolume

			buy := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, buyPrice, buyVolume, buyVolume, 0)
			sell := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, sellPrice, sellVolume, sellVolume, 1)

			Expect(VolumeInEth(buy, sell)).Should(Equal(expected))
		})

		It("calculates volume in ETH/*", func() {
			pair := testutils.TokensETHREN

			expected := uint64((buyVolume * ((buyPrice + sellPrice) / 2)) / 1e12)

			buy := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, buyPrice, buyVolume, buyVolume, 0)
			sell := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, sellPrice, sellVolume, sellVolume, 1)

			Expect(VolumeInEth(buy, sell)).Should(Equal(expected))
		})

		It("calculates volume in BTC/*", func() {
			pair := testutils.TokensBTCREN

			// We assume that BTC price is 100 times that of ETH (as an estimate)
			expected := uint64((buyVolume * ((buyPrice + sellPrice) / 2)) / 1e12 / 100)

			buy := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, buyPrice, buyVolume, buyVolume, 0)
			sell := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, sellPrice, sellVolume, sellVolume, 1)

			Expect(VolumeInEth(buy, sell)).Should(Equal(expected))
		})

		It("calculates volume in */*", func() {
			pair := order.Tokens((uint64(order.TokenREN) << 32) | uint64(order.TokenTUSD))

			buy := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, buyPrice, buyVolume, buyVolume, 0)
			sell := order.NewOrder(order.ParityBuy, order.TypeMidpoint, time.Now().Add(time.Hour), order.SettlementRenEx, pair, sellPrice, sellVolume, sellVolume, 1)

			_, err := VolumeInEth(buy, sell)
			Expect(err).Should(Equal(ErrUnableToCalculateVolumeInEth))
		})
	})
})

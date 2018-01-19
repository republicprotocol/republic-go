package network

import (
	"log"
	"time"

	identity "github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func SendOrderFragment(target identity.MultiAddress, orderFragment *compute.OrderFragment) (*identity.MultiAddress, error) {
	// Connect to the target.
	conn, err := Dial(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	// Create a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call the SendOrderFragment RPC on the target.
	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	to, err := target.Address()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	serializedOrderFragment.To = &rpc.Address{Address: to.String()}
	serializedOrderFragment.From = &rpc.Address{Address: ""}
	multi, err := client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}

	ret, err := identity.NewMultiAddressFromString(multi.Multi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func SerializeResultFragment(input *compute.ResultFragment) *rpc.ResultFragment {
	resultFragment := &rpc.ResultFragment{
		Id:                  []byte(input.ID),
		BuyOrderId:          []byte(input.BuyOrderID),
		SellOrderId:         []byte(input.SellOrderID),
		BuyOrderFragmentId:  []byte(input.BuyOrderFragmentID),
		SellOrderFragmentId: []byte(input.SellOrderFragmentID),
	}
	resultFragment.FstCodeShare = sss.ToBytes(input.FstCodeShare)
	resultFragment.SndCodeShare = sss.ToBytes(input.SndCodeShare)
	resultFragment.PriceShare = sss.ToBytes(input.PriceShare)
	resultFragment.MaxVolumeShare = sss.ToBytes(input.MaxVolumeShare)
	resultFragment.MinVolumeShare = sss.ToBytes(input.MinVolumeShare)
	return resultFragment
}

func DeserializeResultFragment(input *rpc.ResultFragment) (*compute.ResultFragment, error) {
	resultFragment := &compute.ResultFragment{
		ID:                  compute.ResultFragmentID(input.Id),
		BuyOrderID:          compute.OrderID(input.BuyOrderId),
		SellOrderID:         compute.OrderID(input.SellOrderId),
		BuyOrderFragmentID:  compute.OrderFragmentID(input.BuyOrderFragmentId),
		SellOrderFragmentID: compute.OrderFragmentID(input.SellOrderFragmentId),
	}

	var err error
	resultFragment.FstCodeShare, err = sss.FromBytes(input.FstCodeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.SndCodeShare, err = sss.FromBytes(input.SndCodeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.PriceShare, err = sss.FromBytes(input.PriceShare)
	if err != nil {
		return nil, err
	}
	resultFragment.MaxVolumeShare, err = sss.FromBytes(input.MaxVolumeShare)
	if err != nil {
		return nil, err
	}
	resultFragment.MinVolumeShare, err = sss.FromBytes(input.MinVolumeShare)
	if err != nil {
		return nil, err
	}
	return resultFragment, nil
}

func SerializeOrderFragment(input *compute.OrderFragment) *rpc.OrderFragment {
	orderFragment := &rpc.OrderFragment{
		Id:           []byte(input.ID),
		OrderId:      []byte(input.OrderID),
		OrderType:    int64(input.OrderType),
		OrderBuySell: int64(input.OrderBuySell),
	}
	orderFragment.FstCodeShare = sss.ToBytes(input.FstCodeShare)
	orderFragment.SndCodeShare = sss.ToBytes(input.SndCodeShare)
	orderFragment.PriceShare = sss.ToBytes(input.PriceShare)
	orderFragment.MaxVolumeShare = sss.ToBytes(input.MaxVolumeShare)
	orderFragment.MinVolumeShare = sss.ToBytes(input.MinVolumeShare)
	return orderFragment
}

func DeserializeOrderFragment(input *rpc.OrderFragment) (*compute.OrderFragment, error) {
	orderFragment := &compute.OrderFragment{
		ID:           compute.OrderFragmentID(input.Id),
		OrderID:      compute.OrderID(input.OrderId),
		OrderType:    compute.OrderType(input.OrderType),
		OrderBuySell: compute.OrderBuySell(input.OrderBuySell),
	}

	var err error
	orderFragment.FstCodeShare, err = sss.FromBytes(input.FstCodeShare)
	if err != nil {
		log.Println("FstCodeShare =", err)
		return nil, err
	}
	orderFragment.SndCodeShare, err = sss.FromBytes(input.SndCodeShare)
	if err != nil {
		log.Println("SndCodeShare =", err)
		return nil, err
	}
	orderFragment.PriceShare, err = sss.FromBytes(input.PriceShare)
	if err != nil {
		log.Println("PriceShare =", err)
		return nil, err
	}
	orderFragment.MaxVolumeShare, err = sss.FromBytes(input.MaxVolumeShare)
	if err != nil {
		log.Println("MaxVolumeShare =", err)
		return nil, err
	}
	orderFragment.MinVolumeShare, err = sss.FromBytes(input.MinVolumeShare)
	if err != nil {
		log.Println("MinVolumeShare =", err)
		return nil, err
	}
	return orderFragment, nil
}

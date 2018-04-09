package relay

type Relay struct {
}

func SendOrderToDarkOcean (order order.Order) error {

	fmt.Println(string(order.Type))

}
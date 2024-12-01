package main

import (
	"apishoppingcart/src/pb/shoppingcart"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	shoppingcart.ShoppingCartServiceServer
}

// lógica do servidor
func (s *server) AddItem(srv shoppingcart.ShoppingCartService_AddItemServer) error {
	var quantityItems int32 = 0
	var priceTotal float64 = 0.0
	for {
		newItem, err := srv.Recv()
		if err == io.EOF {
			return srv.Send(&shoppingcart.ShoppingCartTotal{
				QuantityItems: quantityItems,
				PriceTotal:    priceTotal,
			})
		}
		if err != nil {
			return fmt.Errorf("error on recv. error: %v", err)
		}

		quantityItems += newItem.GetQuantity()
		priceTotal += float64(newItem.GetPriceUnit() * float64(newItem.GetQuantity()))

		if err := srv.Send(&shoppingcart.ShoppingCartTotal{
			QuantityItems: quantityItems,
			PriceTotal:    priceTotal,
		}); err != nil {
			return fmt.Errorf("error on send. error: %v", err)
		}
	}
}

func main() {
	//subindo o servidor
	fmt.Println("starting grpc server")
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("error on get listener. error: ", err)
	}

	s := grpc.NewServer()
	shoppingcart.RegisterShoppingCartServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalln("error on server. error: ", err)
	}
}

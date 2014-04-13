package main

import (
	"flag"
	"fmt"
	"os"
	"pardal/net"
	"pardal/protocol"
)

var placa string

func init() {
	flag.StringVar(&placa, "placa", "", "Placa do veiculo")
}

func main() {
	flag.Parse()

	if len(placa) != 7 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	client := net.NewSinespClient()

	v, err := client.GetVehicleInfo(placa)
	if err != nil {
		if err == protocol.ErrVehicleNotFound {
			fmt.Println(placa, "Veiculo não encontrado")
		}
		return
	}

	fmt.Println(v)
}

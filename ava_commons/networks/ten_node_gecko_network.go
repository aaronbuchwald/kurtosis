package networks

import (
	"github.com/kurtosis-tech/kurtosis/ava_commons/services"
	"github.com/kurtosis-tech/kurtosis/commons/testnet"
	"github.com/palantir/stacktrace"
)

type TenNodeGeckoNetwork struct{
	geckoServices map[int]services.GeckoService
}
func (network TenNodeGeckoNetwork) GetGeckoService(i int) (services.GeckoService, error){
	if i < 0 || i >= len(network.geckoServices) {
		return services.GeckoService{}, stacktrace.NewError("Invalid Gecko service ID")
	}
	// TODO if we're just getting services back from the ServiceConfigBuilder, then how can we make assumptions here??
	service := network.geckoServices[i]
	return service, nil
}


type TenNodeGeckoNetworkLoader struct{}
func (loader TenNodeGeckoNetworkLoader) GetNetworkConfig(testImageName string) (*testnet.ServiceNetworkConfig, error) {
	factoryConfig := services.NewGeckoServiceFactoryConfig(
		testImageName,
		2,
		2,
		false,
		services.LOG_LEVEL_DEBUG)
	factory := testnet.NewServiceFactory(factoryConfig)

	builder := testnet.NewServiceNetworkConfigBuilder()
	config1 := builder.AddServiceConfiguration(*factory)
	bootNode0, err := builder.AddService(config1, make(map[int]bool))
	if err != nil {
		return nil, stacktrace.Propagate(err, "Could not add bootnode service")
	}
	bootNode1, err := builder.AddService(
		config1,
		map[int]bool{
			bootNode0: true,
		},
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Could not add dependent service")
	}
	bootNode2, err := builder.AddService(
		config1,
		map[int]bool{
			bootNode0: true,
			bootNode1: true,
		},
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Could not add dependent service")
	}
	bootNodeMap := map[int]bool{
		bootNode0: true,
		bootNode1: true,
		bootNode2: true,
	}
	for i:=3; i < 10; i++ {
		_, err := builder.AddService(
			config1,
			bootNodeMap,
		)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Could not add dependent service")
		}
	}

	return builder.Build(), nil
}
func (loader TenNodeGeckoNetworkLoader) LoadNetwork(ipAddrs map[int]string) (interface{}, error) {
	geckoServices := make(map[int]services.GeckoService)
	for serviceId, ipAddr := range ipAddrs {
		geckoServices[serviceId] = *services.NewGeckoService(ipAddr)
	}
	return TenNodeGeckoNetwork{
		geckoServices: geckoServices,
	}, nil
}

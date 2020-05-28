package services

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/kurtosis-tech/kurtosis/commons/testnet"
	"log"
	"strconv"
	"strings"
)

const (
	httpPort = 9650
	stakingPort = 9651
)

// ================= Gecko Service ==================================

type GeckoService struct {
	ipAddr string
}

func NewGeckoService(ipAddr string) *GeckoService {
	return &GeckoService{
		ipAddr:      ipAddr,
	}
}
func (g GeckoService) GetStakingSocket() testnet.ServiceSocket {
	stakingPort, err := nat.NewPort("tcp", strconv.Itoa(stakingPort))
	if err != nil {
		// Realllllly don't think we should deal with propagating this one.... it means the user mistyped an integer
		panic(err)
	}
	return *testnet.NewServiceSocket(g.ipAddr, stakingPort)
}

func (g GeckoService) GetJsonRpcSocket() testnet.ServiceSocket {
	httpPort, err := nat.NewPort("tcp", strconv.Itoa(httpPort))
	if err != nil {
		panic(err)
	}
	return *testnet.NewServiceSocket(g.ipAddr, httpPort)
}

// ================ Gecko Service Factory =============================
type geckoLogLevel string
const (
	LOG_LEVEL_VERBOSE geckoLogLevel = "verbo"
	LOG_LEVEL_DEBUG   geckoLogLevel = "debug"
	LOG_LEVEL_INFO    geckoLogLevel = "info"
)

type GeckoServiceFactoryConfig struct {
	dockerImage       string
	snowSampleSize    int
	snowQuorumSize    int
	stakingTlsEnabled bool
	logLevel          geckoLogLevel
}

func NewGeckoServiceFactoryConfig(dockerImage string,
	snowSampleSize int,
	snowQuorumSize int,
	stakingTlsEnabled bool,
	logLevel geckoLogLevel) *GeckoServiceFactoryConfig {
	return &GeckoServiceFactoryConfig{
		dockerImage:       dockerImage,
		snowSampleSize:    snowSampleSize,
		snowQuorumSize:    snowQuorumSize,
		stakingTlsEnabled: stakingTlsEnabled,
		logLevel:          logLevel,
	}
}

func (g GeckoServiceFactoryConfig) GetDockerImage() string {
	return g.dockerImage
}

func (g GeckoServiceFactoryConfig) GetUsedPorts() map[int]bool {
	return map[int]bool{
		httpPort:    true,
		stakingPort: true,
	}
}

func (g GeckoServiceFactoryConfig) GetStartCommand(publicIpAddr string, dependencies []testnet.Service) []string {
	publicIpFlag := fmt.Sprintf("--public-ip=%s", publicIpAddr)
	commandList := []string{
		"/gecko/build/ava",
		publicIpFlag,
		"--network-id=local",
		fmt.Sprintf("--http-port=%d", httpPort),
		fmt.Sprintf("--staking-port=%d", stakingPort),
		fmt.Sprintf("--log-level=%s", g.logLevel),
		fmt.Sprintf("--snow-sample-size=%d", g.snowSampleSize),
		fmt.Sprintf("--snow-quorum-size=%d", g.snowQuorumSize),
		fmt.Sprintf("--staking-tls-enabled=%v", g.stakingTlsEnabled),
	}

	// If bootstrap nodes are down then Gecko will wait until they are, so we don't actually need to busy-loop making
	// requests to the nodes
	if dependencies != nil && len(dependencies) > 0 {
		// TODO realllllllly wish Go had generics, so we didn't have to do this!
		avaDependencies := make([]AvaService, 0, len(dependencies))
		for _, service := range dependencies {
			avaDependencies = append(avaDependencies, service.(AvaService))
		}

		socketStrs := make([]string, 0, len(avaDependencies))
		for _, service := range avaDependencies {
			socket := service.GetStakingSocket()
			socketStrs = append(socketStrs, fmt.Sprintf("%s:%d", socket.GetIpAddr(), socket.GetPort().Int()))
		}
		joinedSockets := strings.Join(socketStrs, ",")
		commandList = append(commandList, "--bootstrap-ips=" + joinedSockets)
	}
	log.Printf("Command List: %+v", commandList)
	return commandList
}

func (g GeckoServiceFactoryConfig) GetServiceFromIp(ipAddr string) testnet.Service {
	return GeckoService{ipAddr: ipAddr}
}
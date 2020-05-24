/*

Contains types to represent nodes contained in Docker containers.

*/

package ava_commons

import (
	"fmt"
	"github.com/gmarchetti/kurtosis/commons"
	"log"
	"strings"
)

// Type representing a Gecko Node and which ports on the host machine it will use for HTTP and Staking.
type GeckoServiceConfig struct {
	geckoImageName string
	snowSampleSize int
	snowQuorumSize int
	stakingTlsEnabled bool
}

const (
	STAKING_PORT_ID commons.ServiceSpecificPort = 0

	HTTP_PORT = 9650
	STAKING_PORT = 9651
)

// TODO implement more Ava-specific params here, like snow quorum
func NewGeckoServiceConfig(
			dockerImage string,
			snowSampleSize int,
			snowQuorumSize int,
			stakingTlsEnabled bool) *GeckoServiceConfig {
	return &GeckoServiceConfig{
		geckoImageName:    dockerImage,
		snowSampleSize:    snowSampleSize,
		snowQuorumSize:    snowQuorumSize,
		stakingTlsEnabled: stakingTlsEnabled,
	}
}

func (g GeckoServiceConfig) GetDockerImage() string {
	return g.geckoImageName
}

func (g GeckoServiceConfig) GetJsonRpcPort() int {
	return HTTP_PORT
}

func (g GeckoServiceConfig) GetOtherPorts() map[commons.ServiceSpecificPort]int {
	result := make(map[commons.ServiceSpecificPort]int)
	result[STAKING_PORT_ID] = STAKING_PORT
	return result
}

// TODO The "ipAddrOffset" is a nasty janky hack because we have to give the public IP address! It will go away soon, when Gecko no longer needs this
// Argument will be a map of (IP,port) -> request to make to check if a node is up
func (g GeckoServiceConfig) GetContainerStartCommand(ipAddrOffset int, dependencyLivenessReqs map[commons.JsonRpcServiceSocket]commons.JsonRpcRequest) []string {
	commandList := []string{
		"/gecko/build/ava",
		// TODO this entire flag will go away soon!!
		fmt.Sprintf("--public-ip=127.0.0.1"),
		"--network-id=local",
		fmt.Sprintf("--http-port=%d", HTTP_PORT),
		fmt.Sprintf("--staking-port=%d", STAKING_PORT),
		"--log-level=verbo",
		fmt.Sprintf("--snow-sample-size=%d", g.snowSampleSize),
		fmt.Sprintf("--snow-quorum-size=%d", g.snowQuorumSize),
		fmt.Sprintf("--staking-tls-enabled=%v", g.stakingTlsEnabled),
	}

	// If bootstrap nodes are down then Gecko will wait until they are, so we don't actually need to busy-loop making
	// requests to the nodes

	if dependencyLivenessReqs != nil && len(dependencyLivenessReqs) > 0 {
		socketStrs := make([]string, 0, len(dependencyLivenessReqs))
		for socket, _ := range dependencyLivenessReqs {
			// TODO I hardcoded this to be the staking port rather than the RPC port, because there's currently no way to access the dependencys' staking port - fix!!!!!!!
			log.Printf("Socket structure: IPAddress: %s, Port: %v", socket.IPAddress, socket.Port)
			socketStrs = append(socketStrs, fmt.Sprintf("%s:%d", "127.0.0.1", 9651))
			// socketStrs = append(socketStrs, fmt.Sprintf("%s:%d", socket.IPAddress, socket.Port))
		}
		joinedSockets := strings.Join(socketStrs, ",")
		commandList = append(commandList, "--bootstrap-ips=" + joinedSockets)
	}

	/*
	// TODO DEbugging remove me!
	if len(dependencyLivenessReqs) > 0 {
		joined := strings.Join(commandList, " ")
		waitArgs := []string{
			"bash",
			"-c",
			"sleep 10; " + joined + "",
		}
		commandList = waitArgs
	}

	 */

	return commandList
}

func (g GeckoServiceConfig) GetLivenessRequest() commons.JsonRpcRequest {
	return commons.JsonRpcRequest{
		Endpoint: "/ext/P",
		Method: "platform.getCurrentValidators",
		RpcVersion: commons.RPC_VERSION_1_0,
		Params: make(map[string]string),
		ID: 1,   // Not really sure if we'd ever need to change this
	}
}


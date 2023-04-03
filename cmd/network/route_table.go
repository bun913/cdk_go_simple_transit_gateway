package network

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type routeToTransitGateway struct {
	scope       constructs.Construct
	sharedVpc   awsec2.Vpc
	workloadVpc awsec2.Vpc
	tgw         awsec2.CfnTransitGateway
}

func NewRouteToTransitGateway(scope constructs.Construct, sharedVpc awsec2.Vpc, workloadVpc awsec2.Vpc, tgw awsec2.CfnTransitGateway) routeToTransitGateway {
	return routeToTransitGateway{
		scope:       scope,
		sharedVpc:   sharedVpc,
		workloadVpc: workloadVpc,
		tgw:         tgw,
	}
}

// createRouteToTransitGateway creates a route to the Transit Gateway.
func (rttg routeToTransitGateway) CreateRouteToTransitGateway() {
	// Create a route table for the Transit Gateway.
	rttg.createRouteToTransitGateway("sharedToWorkload", rttg.sharedVpc)
	rttg.createRouteToTransitGateway("workloadToShared", rttg.workloadVpc)
}

func (rttg routeToTransitGateway) createRouteToTransitGateway(name string, fromVpc awsec2.Vpc) {
	subnets := fromVpc.SelectSubnets(&awsec2.SubnetSelection{
		SubnetGroupName: jsii.String("Private"),
	}).Subnets
	for i, subnet := range *subnets {
		routeName := fmt.Sprintf("%s%d", name, i)
		awsec2.NewCfnRoute(rttg.scope, jsii.String(routeName), &awsec2.CfnRouteProps{
			RouteTableId:         subnet.RouteTable().RouteTableId(),
			DestinationCidrBlock: jsii.String("0.0.0.0/0"),
			TransitGatewayId:     rttg.tgw.Ref(),
		})
	}
}

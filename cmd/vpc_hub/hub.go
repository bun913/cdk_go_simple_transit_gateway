package vpc_hub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Hub struct {
	scope constructs.Construct
}

func NewHub(scope constructs.Construct) Hub {
	return Hub{
		scope: scope,
	}
}

func (h Hub) CreateTransitGateway() awsec2.CfnTransitGateway {
	return awsec2.NewCfnTransitGateway(h.scope, jsii.String("TransitGateway"), &awsec2.CfnTransitGatewayProps{
		DefaultRouteTableAssociation: jsii.String("disable"),
		DefaultRouteTablePropagation: jsii.String("disable"),
	})
}

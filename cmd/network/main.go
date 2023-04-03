package network

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Network struct {
	scope          constructs.Construct
	vpcName        string
	cidr           string
	hasSSMEndpoint bool
}

func NewNetwork(scope constructs.Construct, vpcName string, cidr string, hasSSMEndpoint bool) Network {
	return Network{
		scope:          scope,
		vpcName:        vpcName,
		cidr:           cidr,
		hasSSMEndpoint: hasSSMEndpoint,
	}
}

func (nr Network) CreateNetworkResources() awsec2.Vpc {
	// VPC
	vpc := awsec2.NewVpc(nr.scope, &nr.vpcName, &awsec2.VpcProps{
		IpAddresses:        awsec2.IpAddresses_Cidr(jsii.String(nr.cidr)),
		MaxAzs:             jsii.Number(2),
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
		VpcName:            jsii.String("Shared"),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("Private"),
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(24),
			},
		},
	})
	// 指定した時のみVPCエンドポイントを追加
	if nr.hasSSMEndpoint {
		vpc.AddInterfaceEndpoint(jsii.String("SSM"), &awsec2.InterfaceVpcEndpointOptions{
			Service: awsec2.InterfaceVpcEndpointAwsService_SSM(),
		})
		vpc.AddInterfaceEndpoint(jsii.String("SSMMessage"), &awsec2.InterfaceVpcEndpointOptions{
			Service: awsec2.InterfaceVpcEndpointAwsService_SSM_MESSAGES(),
		})
		vpc.AddInterfaceEndpoint(jsii.String("EC2Messag"), &awsec2.InterfaceVpcEndpointOptions{
			Service: awsec2.InterfaceVpcEndpointAwsService_EC2_MESSAGES(),
		})
	}
	return vpc
}

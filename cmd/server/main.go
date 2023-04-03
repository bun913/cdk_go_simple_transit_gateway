package server

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Server struct {
	scope constructs.Construct
	name  string
	vpc   awsec2.Vpc
}

func NewServer(scope constructs.Construct, name string, vpc awsec2.Vpc) Server {
	return Server{
		scope: scope,
		name:  name,
		vpc:   vpc,
	}
}

func (sr Server) CreateServerResources() {
	awsec2.NewInstance(sr.scope, jsii.String(sr.name), &awsec2.InstanceProps{
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		MachineImage: awsec2.MachineImage_LatestAmazonLinux(&awsec2.AmazonLinuxImageProps{
			Generation: awsec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
		}),
		SsmSessionPermissions: jsii.Bool(true),
		Vpc:                   sr.vpc,
	})
}

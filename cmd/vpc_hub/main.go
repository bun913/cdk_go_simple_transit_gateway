package vpc_hub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
)

type HubParameters struct {
	scope       constructs.Construct
	sharedVpc   awsec2.Vpc
	workloadVpc awsec2.Vpc
}

func NewHubParameters(scope constructs.Construct, sharedVpc awsec2.Vpc, workloadVpc awsec2.Vpc) HubParameters {
	return HubParameters{
		scope:       scope,
		sharedVpc:   sharedVpc,
		workloadVpc: workloadVpc,
	}
}

type HubResult struct {
	Tgw awsec2.CfnTransitGateway
}

func (hp HubParameters) CreateHubResources() HubResult {
	hub := NewHub(hp.scope)
	// Transit Gatewayを作成
	tgw := hub.CreateTransitGateway()
	// Transit GatewayにVPCをアタッチ
	attachmentShared := NewVpcAttachment("SharedVpcAttachment", hp.sharedVpc, tgw, "TransitGateway")
	attachmentSharedVpc := attachmentShared.Attach()
	attchmentWorkload := NewVpcAttachment("WorkloadVpcAttachment", hp.workloadVpc, tgw, "TransitGateway")
	attachmentWorkloadVpc := attchmentWorkload.Attach()
	// VPC間の相互通信を許可するルートを作成
	rt := NewRouteTable("RouteTable", tgw)
	vpcInfo := NewVpcInfo(hp.sharedVpc, hp.workloadVpc, attachmentSharedVpc, attachmentWorkloadVpc)
	mutalRoute := NewMutalRoute(hp.scope, rt, vpcInfo)
	mutalRoute.CreateMutalRoute()
	return HubResult{
		Tgw: tgw,
	}
}

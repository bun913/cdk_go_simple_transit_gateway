package main

import (
	"transit_gw/cmd/hub_spoke"
	"transit_gw/cmd/network"
	"transit_gw/cmd/server"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcAndEc2StackProps struct {
	awscdk.StackProps
}

func NewVpcAndEc2Stack(scope constructs.Construct, id string, props *VpcAndEc2StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	// 共通タグを設定
	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("Sample"), nil)
	awscdk.Tags_Of(app).Add(jsii.String("Env"), jsii.String("Dev"), nil)

	stack := NewVpcAndEc2Stack(app, "Sample", &VpcAndEc2StackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	// 将来的にShard VPCになるものを配置するためのVPC
	sharedNetworkResource := network.NewNetwork(stack, "SharedVpc", "10.10.0.0/16", true)
	sharedVpc := sharedNetworkResource.CreateNetworkResources()
	severResource := server.NewServer(stack, "SharedVPCInstance", sharedVpc)
	severResource.CreateServerResources()

	// ワークロード用のVPCのようなもの
	workloadNetwork := network.NewNetwork(stack, "WorkloadVpc", "10.20.0.0/16", false)
	workloadVpc := workloadNetwork.CreateNetworkResources()
	workloadServer := server.NewServer(stack, "WorkloadVPCInstance", workloadVpc)
	workloadServer.CreateServerResources()

	// Transit GatewayによるVPC間の双方向通信を実現するためのリソースを作成
	hubParameters := hub_spoke.NewHubParameters(stack, sharedVpc, workloadVpc)
	hubResult := hubParameters.CreateHubResources()

	// EC2が属するサブネットのルートテーブルからTransit Gatewayへのルートを追加
	routeHubSubnetToTransit := network.NewRouteToTransitGateway(stack, "HubSubnetToTransitGW", sharedVpc, hubResult.Tgw, hubResult.HubAttachment)
	routeHubSubnetToTransit.CreateRouteToTransitGateway()
	routeSpokeSubnetToTransit := network.NewRouteToTransitGateway(stack, "SpokeSubnetToTransitGW", workloadVpc, hubResult.Tgw, hubResult.SpokeAttachment)
	routeSpokeSubnetToTransit.CreateRouteToTransitGateway()

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String("ap-northeast-1"),
	}
}

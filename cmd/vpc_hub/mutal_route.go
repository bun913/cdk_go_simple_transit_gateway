package vpc_hub

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MutalRoute struct {
	scope      constructs.Construct
	routeTable RouteTable
	vpcInfo    VpcInfo
}

func NewMutalRoute(scope constructs.Construct, routeTable RouteTable, vpcInfo VpcInfo) MutalRoute {
	return MutalRoute{
		scope:      scope,
		routeTable: routeTable,
		vpcInfo:    vpcInfo,
	}
}

func (mr MutalRoute) CreateMutalRoute() {
	transitRouteTable := mr.routeTable.Create()
	// Transit GatewayにVPCをアタッチ(アソシエーション・プロパゲーション)
	mr.createRouteBetweenVpc(transitRouteTable)
}

func (mr MutalRoute) createRouteBetweenVpc(routeTable awsec2.CfnTransitGatewayRouteTable) {
	// Transit GatewayにVPC関連付け
	awsec2.NewCfnTransitGatewayRouteTableAssociation(mr.scope, jsii.String("Vpc1AttachToRoute"), &awsec2.CfnTransitGatewayRouteTableAssociationProps{
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment1.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
	awsec2.NewCfnTransitGatewayRouteTableAssociation(mr.scope, jsii.String("Vpc2AttachToRoute"), &awsec2.CfnTransitGatewayRouteTableAssociationProps{
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment2.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
	// VPC1からVPC2へのルートを作成
	awsec2.NewCfnTransitGatewayRoute(mr.scope, jsii.String("Vpc1ToVpc2Route"), &awsec2.CfnTransitGatewayRouteProps{
		DestinationCidrBlock:       mr.vpcInfo.vpc2.VpcCidrBlock(),
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment2.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
	// VPC2からVPC1へのルートを作成
	awsec2.NewCfnTransitGatewayRoute(mr.scope, jsii.String("Vpc2ToVpc1Route"), &awsec2.CfnTransitGatewayRouteProps{
		DestinationCidrBlock:       mr.vpcInfo.vpc1.VpcCidrBlock(),
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment1.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
	// VPC1のプロパゲーションを追加
	awsec2.NewCfnTransitGatewayRouteTablePropagation(mr.scope, jsii.String("Vpc1Propagation"), &awsec2.CfnTransitGatewayRouteTablePropagationProps{
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment1.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
	// VPC2のプロパゲーションを追加
	awsec2.NewCfnTransitGatewayRouteTablePropagation(mr.scope, jsii.String("Vpc2Propagation"), &awsec2.CfnTransitGatewayRouteTablePropagationProps{
		TransitGatewayAttachmentId: mr.vpcInfo.vpcAttachment2.Ref(),
		TransitGatewayRouteTableId: routeTable.Ref(),
	})
}

type VpcInfo struct {
	vpc1           awsec2.Vpc
	vpc2           awsec2.Vpc
	vpcAttachment1 awsec2.CfnTransitGatewayAttachment
	vpcAttachment2 awsec2.CfnTransitGatewayAttachment
}

func NewVpcInfo(vpc1 awsec2.Vpc, vpc2 awsec2.Vpc, vpcAttachment1 awsec2.CfnTransitGatewayAttachment, vpcAttachment2 awsec2.CfnTransitGatewayAttachment) VpcInfo {
	return VpcInfo{
		vpc1:           vpc1,
		vpc2:           vpc2,
		vpcAttachment1: vpcAttachment1,
		vpcAttachment2: vpcAttachment2,
	}
}

type RouteTable struct {
	name string
	tgw  awsec2.CfnTransitGateway
}

func NewRouteTable(name string, tgw awsec2.CfnTransitGateway) RouteTable {
	return RouteTable{
		name: name,
		tgw:  tgw,
	}
}

func (ra RouteTable) Create() awsec2.CfnTransitGatewayRouteTable {
	return awsec2.NewCfnTransitGatewayRouteTable(ra.tgw, jsii.String("RouteTable"), &awsec2.CfnTransitGatewayRouteTableProps{
		TransitGatewayId: ra.tgw.Ref(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("Name"),
				Value: jsii.String(ra.name),
			},
		},
	})
}

package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	awsProvider "github.com/cdktf/cdktf-provider-aws-go/aws/v14/provider"
	awsSubnet "github.com/cdktf/cdktf-provider-aws-go/aws/v14/subnet"
	awsVpc "github.com/cdktf/cdktf-provider-aws-go/aws/v14/vpc"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

var regions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
}

var subnetsCidrs = []string{
	"10.0.1.0/28",
	"10.0.1.128/28",
}

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	provs := []awsProvider.AwsProvider{}

	for _, region := range regions {
		provs = append(provs, awsProvider.NewAwsProvider(stack, &region, &awsProvider.AwsProviderConfig{
			Alias:  &region,
			Region: &region,
		}))
	}

	for _, prov := range provs {
		defaultVpc := createDefaultVpc(stack, prov)
    createSubnetsWithCidrs(stack, prov, defaultVpc, subnetsCidrs)
	}

	return stack
}

func main() {
	app := cdktf.NewApp(&cdktf.AppConfig{
		Context: &map[string]interface{}{
			"excludeStackIdFromLogicalIds": "true",
			"allowSepCharsInLogicalIds":    "true",
		},
	})

	NewMyStack(app, "multiregion")

	app.Synth()
}

func createDefaultVpc(scope constructs.Construct, prov awsProvider.AwsProvider) awsVpc.Vpc {
	return awsVpc.NewVpc(scope, jsii.String(fmt.Sprintf("network-%s", *prov.Region())), &awsVpc.VpcConfig{
		Provider:  prov,
		CidrBlock: jsii.String("10.0.1.0/24"),
		Tags: &map[string]*string{
			"Name": jsii.String(fmt.Sprintf("test-vpc-%s", *prov.Region())),
		},
	})
}

func createSubnetsWithCidrs(scope constructs.Construct, prov awsProvider.AwsProvider, vpc awsVpc.Vpc, cidrs []string) []awsSubnet.Subnet {
	subnets := []awsSubnet.Subnet{}

	for i, cidr := range cidrs {
		subnets = append(
			subnets,
			awsSubnet.NewSubnet(
				scope,
				jsii.String(fmt.Sprintf("subnet-%s-%d", *prov.Region(), i)),
				&awsSubnet.SubnetConfig{
					Provider:  prov,
					VpcId:     vpc.Id(),
					CidrBlock: &cidr,
					Tags: &map[string]*string{
						"Name": jsii.String(fmt.Sprintf("mySubnet-%d", i)),
					},
				},
			),
		)
	}
	return subnets
}

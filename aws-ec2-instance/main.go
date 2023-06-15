package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/dataawsami"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/instance"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/keypair"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/securitygroup"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/vpcsecuritygroupegressrule"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v14/vpcsecuritygroupingressrule"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/mrWinston/cdktf-playground/aws-ec2-instance/generated/vpc"
)

const (
	AMZ_LINUX_AMI_FILTER = "al2023-ami-2023.*-x86_64"
	INSTANCE_TYPE        = "t3a.micro"
	PUBKEY               = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC8JDT86qc0SOML82+V0UMG4NRq+lxDGexl1MENrecQHCufYAhpDZmjbf3gLWRtLSvNcZ8NjwInHOKWBbI8R6nlJ7ac99ZH2ATqAOBmFkmvFAjpmPiVRJjemM26Tkf4o9hd/KxgWMg3bTnunsPgjBhnw8xdSMXyD0ygKO2pOCVW6/XQ5LTRFO4z/OkpY4tJiiLwos1WDJfBKVs8amKi4q+gn86R58hj9f9F3SwLSlG2uYmVLYouFfqymtpkfFEJek2XrCnBLfDK5GjWd4fDrE2o6iBR6xgaUzc/BdBMrx97qVqCLvJEjeDtoxbHJWi1X5HJdBTSrwLjdlXDFG+QxsrJEuTO/wQxbqJ/SVbOnn0VkTZym5zvEkka8SZNvV+hCjgCkpbUlY6NKRSvPOc98DuEAWAW3rPRirRj5BaqKaCuDa5NydSEO9bBxYMoL3l153i3HqNp0VBOXvxpYT7DSYiFYwgGTH3jw0WRveMWdNtA2jh22B6EFw8j5DvmdDGxjyLUxCtFY3hgO8fK8dYArGJEkg9/KZmhX6vBDFlZ6w1X6su+C6le9dxP2V5IKn2j6qbknYUDtc6vxnmaMMkwIJIBbfeI7Un164mNzqlgFNqq/yOsIHhFOXsE3LDFoCuTRM8oD+2vBtxjGWkwDZ8fOG6xj/RPLDSxzdYLHIWGoNJl3Q== maschulz@redhat.com"
  NAME = "teststack"
)

var azs = []*string{
	jsii.String("use2-az2"),
	jsii.String("use2-az3"),
}

var privateSnCidrs = []*string{
	jsii.String("10.0.1.0/24"),
	jsii.String("10.0.2.0/24"),
}
var publicSnCidrs = []*string{
	jsii.String("10.0.101.0/24"),
	jsii.String("10.0.102.0/24"),
}

func createAwsAMIResource(scope constructs.Construct, id string) dataawsami.DataAwsAmi {
	return dataawsami.NewDataAwsAmi(scope, &id, &dataawsami.DataAwsAmiConfig{
		Filter: []dataawsami.DataAwsAmiFilter{{
			Name:   jsii.String("name"),
			Values: &[]*string{jsii.String(AMZ_LINUX_AMI_FILTER)},
		}},
		MostRecent: jsii.Bool(true),
		Owners:     &[]*string{jsii.String("amazon")},
	})
}

func securityGroupWithPorts(scope constructs.Construct, name string, ingressPorts []int, vpcId *string) securitygroup.SecurityGroup{
  sg := securitygroup.NewSecurityGroup(scope, jsii.String(name), &securitygroup.SecurityGroupConfig{
  	Description:         jsii.String("Security Group " + name),
  	NamePrefix:          &name,
  	Tags:                &map[string]*string{},
  	VpcId:               vpcId,
  })

  for _, port := range ingressPorts {
    vpcsecuritygroupingressrule.NewVpcSecurityGroupIngressRule(scope, jsii.String(fmt.Sprintf("%s_allow_%d", name, port)), &vpcsecuritygroupingressrule.VpcSecurityGroupIngressRuleConfig{
    	IpProtocol:                jsii.String("tcp"),
    	CidrIpv4:                  jsii.String("0.0.0.0/0"),
    	Description:               jsii.String(fmt.Sprintf("Allow ingress on port %d", port)),
    	FromPort:                  jsii.Number(port),
    	SecurityGroupId:           sg.Id(),
    	ToPort:                    jsii.Number(port),
    })
  }

  vpcsecuritygroupegressrule.NewVpcSecurityGroupEgressRule(scope, jsii.String(fmt.Sprintf("%s_allow_egress", name)), &vpcsecuritygroupegressrule.VpcSecurityGroupEgressRuleConfig{
  	IpProtocol:                jsii.String("-1"),
  	CidrIpv4:                  jsii.String("0.0.0.0/0"),
  	Description:               jsii.String("Allow all outbound traffic"),
  	SecurityGroupId:           sg.Id(),
  })

  return sg
}

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	provider.NewAwsProvider(stack, jsii.String("aws"), &provider.AwsProviderConfig{
		Region: jsii.String("us-east-2"),
	})

	amzAmi := createAwsAMIResource(stack, "amz_ami")

	vpcModule := vpc.NewVpc(stack, jsii.String("vpc"), &vpc.VpcConfig{
		Name:           jsii.String(NAME),
		Cidr:           jsii.String("10.0.0.0/16"),
		Azs:            &azs,
		PublicSubnets:  &publicSnCidrs,
		PrivateSubnets: &privateSnCidrs,

		EnableNatGateway:   jsii.Bool(true),
		SingleNatGateway:   jsii.Bool(true),
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
	})

	kp := keypair.NewKeyPair(stack, jsii.String("key"), &keypair.KeyPairConfig{
		PublicKey:     jsii.String(PUBKEY),
		KeyNamePrefix: jsii.String("sample-key"),
	})

  sg := securityGroupWithPorts(stack, "teststack", []int{22}, vpcModule.VpcIdOutput())

	instance.NewInstance(stack, jsii.String("inst"), &instance.InstanceConfig{
		Ami:            amzAmi.Id(),
		KeyName:        kp.KeyName(),
		InstanceType:   jsii.String(INSTANCE_TYPE),
		VpcSecurityGroupIds: &[]*string{sg.Id()},
		SubnetId:       jsii.String(cdktf.Fn_Element(cdktf.Token_AsAny(vpcModule.PublicSubnetsOutput()), jsii.Number(0)).(string)),
		Tags:           &map[string]*string{},
	})

	cdktf.NewTerraformOutput(stack, jsii.String("cluster-private-subnet"), &cdktf.TerraformOutputConfig{
		Value: vpcModule.PrivateSubnetsOutput(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("cluster-public-subnet"), &cdktf.TerraformOutputConfig{
		Value: vpcModule.PublicSubnetsOutput(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("static-output"), &cdktf.TerraformOutputConfig{
		Value: jsii.String("test"),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(&cdktf.AppConfig{
		Context: &map[string]interface{}{
			"excludeStackIdFromLogicalIds": "true",
			"allowSepCharsInLogicalIds":    "true",
		},
	})

	NewMyStack(app, "prod")

	app.Synth()
}

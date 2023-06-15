package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
  nullProvider "github.com/cdktf/cdktf-provider-null-go/null/v6/provider"
  "github.com/aws/jsii-runtime-go"
  nullResource "github.com/cdktf/cdktf-provider-null-go/null/v6/resource"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

  nullProvider.NewNullProvider(stack, jsii.String("null"), &nullProvider.NullProviderConfig{})

  res := nullResource.NewResource(stack, jsii.String("test"), &nullResource.ResourceConfig{})

  cdktf.NewTerraformOutput(stack, jsii.String("testoutput"), &cdktf.TerraformOutputConfig{
  	Value: res.Id(),
  	Description:  jsii.String("something"),
  })
  cdktf.NewTerraformVariable(stack, jsii.String("testvar"), &cdktf.TerraformVariableConfig{
  	Description: jsii.String("test"),
    Default: jsii.String("yeah"),
  })

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "prod")
	NewMyStack(app, "stage")

	app.Synth()
}

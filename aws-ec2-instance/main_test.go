package main

import (
	"fmt"
	"testing"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

// The tests below are example tests, you can find more information at
// https://cdk.tf/testing


var stack = NewMyStack(cdktf.Testing_App(nil), "stack")
var synth = cdktf.Testing_Synth(stack, jsii.Bool(false))



//func TestShouldContainContainer(t *testing.T){
//	assertion := cdktf.Testing_ToHaveResource(synth, docker.Container_TfResourceType())
//
//	if !*assertion  {
//		t.Error(assertion.Message())
//	}
//}

func TestMain(t *testing.T){
	assertion := cdktf.Testing_ToBeValidTerraform(cdktf.Testing_FullSynth(stack))
  fmt.Println(cdktf.Testing_RenderConstructTree(stack))

	if !*assertion  {
		t.Error()
	}
}


package containerpatterns_test

import (
	//	"fmt"
	containerpatterns "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"log"
	"testing"
)

var (
	ecsProjectApp   core.App   = nil
	ecsProjectStack core.Stack = nil
	ecsProjectProps containerpatterns.EcsProjectProps
)

func setupEcsProject() {
	ecsProjectApp = core.NewApp(&core.AppProps{
		AnalyticsReporting: jsii.Bool(false),
	})
	ecsProjectStack = core.NewStack(ecsProjectApp, jsii.String("NonLoadBalancedEc2ServiceStack"), &core.StackProps{
		Env: &core.Environment{
			Account: jsii.String("123456789012"),
			Region:  jsii.String("us-east-1"),
		},
	})
	ecsProjectProps = containerpatterns.EcsProjectProps{
		NetworkProps: brznetwork.VpcProps{
			Id:        "vpc-12345",
			IsDefault: true,
		},
		ComputeProps: containerpatterns.EcsComputeProps{
			VpcId: "vpc-12345",
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: true,
			},
			AsgCapacityProviders: nil,
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:        "test-env-bucket",
				IsVersioned: false,
			},
			LoadBalancer:      containerpatterns.LoadBalancerOptions{},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{},
		},
		NonLoadBalancedEc2ServicesProps: nil,
		LoadBalancedEc2ServicesProps:    nil,
	}

}

func teardown() {
	ecsProjectApp = nil
	ecsProjectStack = nil
	ecsProjectProps = containerpatterns.EcsProjectProps{}
	jsii.Close()
	log.Println("Resetting EcsProject resources")
}

func TestEcsProject(t *testing.T) {
	setupEcsProject()
	containerpatterns.NewEcsProject(ecsProjectStack, jsii.String("EcsProjectWithComputeOnly"), &ecsProjectProps)
	template := assertions.Template_FromStack(ecsProjectStack, &assertions.TemplateParsingOptions{})
	template.ResourceCountIs(jsii.String("AWS::CloudFormation::Stack"), jsii.Number(1))
    t.Cleanup(teardown)
}

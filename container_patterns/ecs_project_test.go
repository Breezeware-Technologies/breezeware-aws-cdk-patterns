package containerpatterns_test

import (
	//	"fmt"
	containerpatterns "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"
)

var (
	app             core.App   = nil
	stack           core.Stack = nil
	ecsProjectProps containerpatterns.EcsProjectProps
)

func setupEcsProject() {
	app = core.NewApp(&core.AppProps{
		AnalyticsReporting: jsii.Bool(false),
	})
	stack = core.NewStack(app, jsii.String("NonLoadBalancedEc2ServiceStack"), &core.StackProps{
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
	app = nil
	stack = nil
	ecsProjectProps = containerpatterns.EcsProjectProps{}
	jsii.Close()
	log.Println("Resetting EcsProject resources")
}

func TestEcsProject(t *testing.T) {
	setupEcsProject()
	containerpatterns.NewEcsProject(stack, jsii.String("EcsProjectWithComputeOnly"), &ecsProjectProps)
	template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
	template.ResourceCountIs(jsii.String("AWS::CloudFormation::Stack"), jsii.Number(1))
}

type mockEcsProject struct {
	mock.Mock
	constructs.Construct
	network                         ec2.IVpc
	compute                         ecs.Cluster
	applicationLoadBalancer         elb2.ApplicationLoadBalancer
	ec2ContainerApplicationServices []ecs.Ec2Service
}

type MockEcsProject interface {
	//	mock.Mock
	// Network returns an IVpc as the network-working component of the EcsProject.
	Network() ec2.IVpc
	// Compute returns Cluster as the container compute environment.
	Compute() ecs.Cluster
	// ApplicationLoadBalancer returns the load-balancer associated for the EcsProject as a part of managing load & traffic.
	ApplicationLoadBalancer() elb2.ApplicationLoadBalancer
	// Ec2ContainerApplicationServices returns a slice of ECS EC2 based services that are currently in the EcsProject
	Ec2ContainerApplicationServices() []ecs.Ec2Service
}

func (m *mockEcsProject) Network() ec2.IVpc {
	m.Called()
	return m.network
}

func TestEcsProject_Compute(t *testing.T) {
	var m mockEcsProject
	m.On("Network", ecsProjectProps).Return(ec2.NewVpc(nil, nil, nil))
}

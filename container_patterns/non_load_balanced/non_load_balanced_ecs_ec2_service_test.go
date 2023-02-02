package containerpatterns_test

import (
	//	"fmt"
	brzNlbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/non_load_balanced"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/go-cmp/cmp"
	"log"
	"testing"
)

var (
	app   core.App   = nil
	stack core.Stack = nil
)

func setup() {
	app = core.NewApp(&core.AppProps{
		AnalyticsReporting: jsii.Bool(false),
	})
	stack = core.NewStack(app, jsii.String("NonLoadBalancedEc2ServiceStack"), &core.StackProps{
		Env: &core.Environment{
			Account: jsii.String("123456789012"),
			Region:  jsii.String("us-east-1"),
		},
	})
}

func TestNewNonLoadBalancedEc2Service_ServiceDiscovery(t *testing.T) {
	t.Run("With Service Discovery Enabled", func(t *testing.T) {
		setup()
		brzNlbEc2Service.NewNonLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
			Cluster: brzNlbEc2Service.ClusterProps{
				Vpc: brznetwork.VpcProps{
					IsDefault: true,
				},
			},
			LogGroupName: "test-log-group",
			TaskDefinition: brzNlbEc2Service.TaskDefinition{
				NetworkMode: brzNlbEc2Service.TaskDefintionNetworkModeAwsVpc,
				EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
					BucketName: "test-bucket",
				},
				RequiresVolume: false,
				ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
					{
						ContainerName: "test-service",
						Image:         "test-service",
						RegistryType:  brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
						ImageTag:      "latest",
						IsEssential:   true,
						Cpu:           512,
						Memory:        1024,
						PortMappings: []ecs.PortMapping{
							{
								ContainerPort: jsii.Number(8443),
								Protocol:      ecs.Protocol_TCP,
							},
						},
						EnvironmentFileObjectKey: "test-service/service.env",
					},
				},
			},
			IsTracingEnabled: false,
			DesiredTaskCount: 1,
			CapacityProviders: []string{
				"testT2Small",
			},
			IsServiceDiscoveryEnabled: true,
			ServiceDiscovery: brzNlbEc2Service.ServiceDiscoveryProps{
				ServiceName: "test-service",
				ServicePort: 8443,
				CloudMapNamespace: brzNlbEc2Service.CloudMapNamespaceProps{
					NamespaceName: "golang-cdk.test",
					NamespaceId:   "ns-golangtest",
					NamespaceArn:  "arn:aws:servicediscovery:us-east-1:123456789012:namespace/ns-golangtest",
				},
			},
		})
		// then
		template := assertions.Template_FromStack(stack, nil)
		template.ResourceCountIs(jsii.String("AWS::ServiceDiscovery::Service"), jsii.Number(1))
		template.HasResourceProperties(jsii.String("AWS::ServiceDiscovery::Service"), map[string]interface{}{
			"DnsConfig": map[string]interface{}{
				"DnsRecords": []map[string]interface{}{
					{
						"TTL":  60,
						"Type": "A",
					},
				},
				"NamespaceId":   "ns-golangtest",
				"RoutingPolicy": "MULTIVALUE",
			},
			"HealthCheckCustomConfig": map[string]interface{}{
				"FailureThreshold": 1,
			},
			"Name":        "test-service",
			"NamespaceId": "ns-golangtest",
		})
	})
	t.Run("With Service Discovery Disabled", func(t *testing.T) {
		setup()
		brzNlbEc2Service.NewNonLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
			Cluster: brzNlbEc2Service.ClusterProps{
				Vpc: brznetwork.VpcProps{
					IsDefault: true,
				},
			},
			LogGroupName: "test-log-group",
			TaskDefinition: brzNlbEc2Service.TaskDefinition{
				NetworkMode: brzNlbEc2Service.TaskDefintionNetworkModeAwsVpc,
				EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
					BucketName: "test-bucket",
				},
				RequiresVolume: false,
				ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
					{
						ContainerName: "test-service",
						Image:         "test-service",
						RegistryType:  brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
						ImageTag:      "latest",
						IsEssential:   true,
						Cpu:           512,
						Memory:        1024,
						PortMappings: []ecs.PortMapping{
							{
								ContainerPort: jsii.Number(8443),
								Protocol:      ecs.Protocol_TCP,
							},
						},
						EnvironmentFileObjectKey: "test-service/service.env",
					},
				},
			},
			IsTracingEnabled: false,
			DesiredTaskCount: 1,
			CapacityProviders: []string{
				"testT2Small",
			},
			IsServiceDiscoveryEnabled: false,
			ServiceDiscovery: brzNlbEc2Service.ServiceDiscoveryProps{
				ServiceName: "test-service",
				ServicePort: 8443,
				CloudMapNamespace: brzNlbEc2Service.CloudMapNamespaceProps{
					NamespaceName: "golang-cdk.test",
					NamespaceId:   "ns-golangtest",
					NamespaceArn:  "arn:aws:servicediscovery:us-east-1:123456789012:namespace/ns-golangtest",
				},
			},
		})
		// then
		template := assertions.Template_FromStack(stack, nil)
		template.ResourceCountIs(jsii.String("AWS::ServiceDiscovery::Service"), jsii.Number(0))
	})
	t.Cleanup(teardown)

}

func teardown() {
	app = nil
	stack = nil
	log.Println("Tearing down tested resources")
}

func TestNewNonLoadBalancedEc2Service_LogGroupCreation(t *testing.T) {

	setup()
	logGroupName := "testLogGroup"

	brzNlbEc2Service.NewNonLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
		Cluster: brzNlbEc2Service.ClusterProps{
			Vpc: brznetwork.VpcProps{
				IsDefault: true,
			},
		},
		LogGroupName: logGroupName,
		TaskDefinition: brzNlbEc2Service.TaskDefinition{
			//			FamilyName:  "test-service",
			NetworkMode: brzNlbEc2Service.TaskDefintionNetworkModeBridge,
			EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
				BucketName: "test-bucket",
			},
			RequiresVolume: false,
			ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
				{
					ContainerName: "test-service",
					Image:         "test-service",
					RegistryType:  brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
					ImageTag:      "latest",
					IsEssential:   true,
					Cpu:           512,
					Memory:        1458,
					PortMappings: []ecs.PortMapping{
						{
							ContainerPort: jsii.Number(8443),
							Protocol:      ecs.Protocol_TCP,
						},
					},
					EnvironmentFileObjectKey: "test-service/service.env",
				},
			},
		},
		IsTracingEnabled: false,
		DesiredTaskCount: 1,
		CapacityProviders: []string{
			"testT2Small",
		},
		IsServiceDiscoveryEnabled: false,
		ServiceDiscovery:          brzNlbEc2Service.ServiceDiscoveryProps{},
	})
	template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
	template.ResourceCountIs(jsii.String("AWS::Logs::LogGroup"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::Logs::LogGroup"), map[string]interface{}{
		"RetentionInDays": 14,
	})
	template.HasResource(jsii.String("AWS::Logs::LogGroup"), map[string]interface{}{
		"Type": "AWS::Logs::LogGroup",
		"Properties": map[string]interface{}{
			"LogGroupName":    logGroupName,
			"RetentionInDays": 14,
		},
		"UpdateReplacePolicy": "Delete",
		"DeletionPolicy":      "Delete",
	})
	t.Cleanup(teardown)

}

func TestNewNonLoadBalancedEc2Service_TracingEnabled(t *testing.T) {

	setup()
	brzNlbEc2Service.NewNonLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
		Cluster: brzNlbEc2Service.ClusterProps{
			Vpc: brznetwork.VpcProps{
				IsDefault: true,
			},
		},
		LogGroupName: "test-log-group",
		TaskDefinition: brzNlbEc2Service.TaskDefinition{
			NetworkMode: brzNlbEc2Service.TaskDefintionNetworkModeBridge,
			EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
				BucketName: "test-bucket",
			},
			RequiresVolume: false,
			ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
				{
					ContainerName: "test-service",
					Image:         "test-service",
					RegistryType:  brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
					ImageTag:      "latest",
					IsEssential:   true,
					Cpu:           512,
					Memory:        1024,
					PortMappings: []ecs.PortMapping{
						{
							ContainerPort: jsii.Number(8443),
							Protocol:      ecs.Protocol_TCP,
						},
					},
					EnvironmentFileObjectKey: "test-service/service.env",
				},
			},
		},
		IsTracingEnabled: true,
		DesiredTaskCount: 1,
		CapacityProviders: []string{
			"testT2Small",
		},
		IsServiceDiscoveryEnabled: false,
		ServiceDiscovery:          brzNlbEc2Service.ServiceDiscoveryProps{},
	})
	template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
	template.ResourceCountIs(jsii.String("AWS::ECS::TaskDefinition"), jsii.Number(1))
	logConfigurationCapture := assertions.NewCapture(nil)
	template.HasResourceProperties(jsii.String("AWS::ECS::TaskDefinition"), map[string]interface{}{
		"ContainerDefinitions": []map[string]interface{}{
			{
				"Command": []string{
					"--config=/etc/ecs/ecs-default-config.yaml",
				},
				"Cpu":              256,
				"Essential":        true,
				"Image":            "amazon/aws-otel-collector:v0.25.0",
				"LogConfiguration": logConfigurationCapture,
				"Memory":           256,
				"Name":             "otel-xray",
				"PortMappings": []map[string]interface{}{
					{
						"ContainerPort": 2000,
						"HostPort":      0,
						"Protocol":      "udp",
					},
					{
						"ContainerPort": 4317,
						"HostPort":      0,
						"Protocol":      "tcp",
					},
					{
						"ContainerPort": 8125,
						"HostPort":      0,
						"Protocol":      "udp",
					},
				},
			},
			{
				"Cpu": 512,
				"DependsOn": []map[string]interface{}{
					{
						"Condition":     "START",
						"ContainerName": "otel-xray",
					},
				},
				"EnvironmentFiles": []interface{}{
					map[string]interface{}{
						"Type": "s3",
						"Value": map[string]interface{}{
							"Fn::Join": []interface{}{
								"",
								[]interface{}{
									"arn:",
									map[string]interface{}{
										"Ref": "AWS::Partition",
									},
									":s3:::test-bucket/test-service/service.env",
								},
							},
						},
					},
				},
				"Essential": true,
				"Image": map[string]interface{}{
					"Fn::Join": []interface{}{
						"",
						[]interface{}{
							"123456789012.dkr.ecr.us-east-1.",
							map[string]interface{}{
								"Ref": "AWS::URLSuffix",
							},
							"/test-service:latest",
						},
					},
				},
				"Links": []string{
					"otel-xray:otel-xray",
				},
				"LogConfiguration": map[string]interface{}{
					"LogDriver": "awslogs",
					"Options": map[string]interface{}{
						"awslogs-group": map[string]interface{}{
							"Ref": "LoadBalancedEc2ServiceLogGroup73C29131",
						},
						"awslogs-stream-prefix": "test-service",
						"awslogs-region":        "us-east-1",
					},
				},
				"Memory": 1024,
				"Name":   "test-service",
				"PortMappings": []map[string]interface{}{
					{
						"ContainerPort": 8443,
						"HostPort":      0,
						"Protocol":      "tcp",
					},
				},
			},
		},
	})
	template.ResourceCountIs(jsii.String("AWS::ECS::TaskDefinition"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::ECS::TaskDefinition"), map[string]interface{}{
		"NetworkMode": "bridge",
	})
	expectedLogConfiguration := &map[string]interface{}{
		"LogDriver": "awslogs",
		"Options": map[string]interface{}{
			"awslogs-group": map[string]interface{}{
				"Ref": "LoadBalancedEc2ServiceLogGroup73C29131",
			},
			"awslogs-stream-prefix": "otel",
			"awslogs-region":        "us-east-1",
		},
	}
	if !cmp.Equal(logConfigurationCapture.AsObject(), expectedLogConfiguration) {
		t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedLogConfiguration, logConfigurationCapture.AsObject())
	}
	t.Cleanup(teardown)

}

func TestNewNonLoadBalancedEc2Service_WithoutLoadBalancerListenerAndTargetGroup(t *testing.T) {
	setup()
	brzNlbEc2Service.NewNonLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
		Cluster: brzNlbEc2Service.ClusterProps{
			Vpc: brznetwork.VpcProps{
				IsDefault: true,
			},
		},
		LogGroupName: "test-log-group",
		TaskDefinition: brzNlbEc2Service.TaskDefinition{
			NetworkMode: brzNlbEc2Service.TaskDefintionNetworkModeBridge,
			EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
				BucketName: "test-bucket",
			},
			RequiresVolume: false,
			ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
				{
					ContainerName: "test-service",
					Image:         "test-service",
					RegistryType:  brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
					ImageTag:      "latest",
					IsEssential:   true,
					Cpu:           512,
					Memory:        1024,
					PortMappings: []ecs.PortMapping{
						{
							ContainerPort: jsii.Number(8443),
							Protocol:      ecs.Protocol_TCP,
						},
					},
					EnvironmentFileObjectKey: "test-service/service.env",
				},
			},
		},
		IsTracingEnabled: true,
		DesiredTaskCount: 1,
		CapacityProviders: []string{
			"testT2Small",
		},
		IsServiceDiscoveryEnabled: false,
		ServiceDiscovery:          brzNlbEc2Service.ServiceDiscoveryProps{},
	})
	template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
	template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::ListenerRule"), jsii.Number(0))
	template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::TargetGroup"), jsii.Number(0))
}

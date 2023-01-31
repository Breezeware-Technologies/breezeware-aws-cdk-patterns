package containerpatterns_test

import (
	//	"fmt"
	brzLbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/load_balanced"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
	"log"
	"testing"
)

var loadBalancedEc2ServiceProps brzLbEc2Service.LoadBalancedEc2ServiceProps

func TestNewLoadBalancedEc2Service(t *testing.T) {
	setup()
	defer teardown()
	// given
	app := core.NewApp(nil)

	// when
	stack := core.NewStack(app, jsii.String("LoadBalancedEc2ServiceStack"), &core.StackProps{
		Env: &core.Environment{
			Account: jsii.String("305251478828"),
			Region:  jsii.String("us-east-1"),
		},
	})
	//	fmt.Println("LoadBalancedEC2ServiceProps: ", loadBalancedEc2ServiceProps)
	brzLbEc2Service.NewLoadBalancedEc2Service(stack, jsii.String("LoadBalancedEc2Service"), &loadBalancedEc2ServiceProps)

	// then
	template := assertions.Template_FromStack(stack, nil)
	template.ResourceCountIs(jsii.String("AWS::Logs::LogGroup"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::ListenerRule"), map[string]interface{}{
		"Priority": 1,
		"Conditions": []map[string]interface{}{
			{
				"Field": "host-header",
				"HostHeaderConfig": map[string]interface{}{
					"Values": []string{
						"nginx.dynamostack.com",
					},
				},
			},
			{
				"Field": "path-pattern",
				"PathPatternConfig": map[string]interface{}{
					"Values": []string{
						"/api/*",
					},
				},
			},
		},
	})
	template.HasResource(jsii.String("AWS::ElasticLoadBalancingV2::TargetGroup"), map[string]interface{}{
		//		"Type": "AWS::ElasticLoadBalancingV2::TargetGroup",
		"Properties": map[string]interface{}{
			"HealthCheckEnabled":         true,
			"HealthCheckIntervalSeconds": 30,
			"HealthCheckPath":            "/api/health-status",
			"Matcher": map[string]interface{}{
				"HttpCode": "200",
			},
			"Port":     80,
			"Protocol": "HTTP",
			"TargetGroupAttributes": []map[string]interface{}{
				{
					"Key":   "stickiness.enabled",
					"Value": "false",
				},
			},
			"TargetType": "instance",
			//			"VpcId":      "vpc-535bd136",
		},
		//		"Metadata": map[string]interface{}{
		//			"aws:cdk:path": "EcsProject/EcsProject/ComputeStack/LoadBalancedEc2ContainerApplicationServicesStack/LoadBalancedService0/ApplicationTargetGroup/Resource",
		//		},
	})

	//    json := template.ToJSON()
	//    log.Println("Synthesized Template: ", json)
}

func setup() {
	// log.Println("Entering setup()")
	loadBalancedEc2ServiceProps = brzLbEc2Service.LoadBalancedEc2ServiceProps{
		Cluster: brzLbEc2Service.ClusterProps{
			ClusterName: "<dummy>",
			Vpc: brznetwork.VpcProps{
				Id:        "<dummy>",
				IsDefault: true,
			},
		},
		LogGroupName: "GolangCdkDemoService",
		TaskDefinition: brzLbEc2Service.TaskDefinition{
			FamilyName:  "rpc-service",
			NetworkMode: brzLbEc2Service.TaskDefintionNetworkModeBridge,
			EnvironmentFile: brzLbEc2Service.EnvironmentFile{
				BucketName: "test-bucket",
				BucketArn:  "<dummy>",
			},
			RequiresVolume: false,
			ApplicationContainers: []brzLbEc2Service.ContainerDefinition{
				{
					ContainerName:            "rpc-service",
					Image:                    "rpc-service",
					RegistryType:             brzLbEc2Service.ContainerDefinitionRegistryAwsEcr,
					ImageTag:                 "latest",
					IsEssential:              true,
					Cpu:                      512,
					Memory:                   1458,
					EnvironmentFileObjectKey: "rpc-service/prod/app.env",
					PortMappings: []ecs.PortMapping{
						{
							ContainerPort: jsii.Number(8443),
							Protocol:      ecs.Protocol_TCP,
						},
					},
				},
			},
		},
		IsTracingEnabled: true,
		DesiredTaskCount: 1,
		CapacityProviders: []string{
			"GolangCdkDemoT2Small",
		},
		IsServiceDiscoveryEnabled: false,
		ServiceDiscovery:          brzLbEc2Service.ServiceDiscoveryProps{},
//		IsLoadBalancerEnabled:     true,
		LoadBalancer: brzLbEc2Service.LoadBalancerProps{
			ListenerArn:           "<dummy>",
			SecurityGroupId:       "<dummy>",
			TargetHealthCheckPath: "/api/health-status",
			ListenerRuleProps: brzLbEc2Service.ListenerRuleProps{
				Priority:      1,
				PathCondition: "/api/*",
				HostCondition: "nginx.dynamostack.com",
			},
		},
		LoadBalancerTargetOptions: ecs.LoadBalancerTargetOptions{
			ContainerName: jsii.String("rpc-service"),
			ContainerPort: jsii.Number(8443),
			Protocol:      ecs.Protocol_TCP,
		},
	}
	// log.Printf("Initialized LoadBalancedEc2ServiceProps: %v", loadBalancedEc2ServiceProps)
	log.Println("Initialized LoadBalancedEc2ServiceProps")
	// log.Println("Leaving setup()")
}

func teardown() {
	log.Println("Tearing down tested resources")
}

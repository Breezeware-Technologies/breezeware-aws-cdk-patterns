package containerpatterns_test

import (
	"fmt"
	containerpatterns "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns"
	brzLbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/load_balanced"
	brzNlbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/non_load_balanced"
	"github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
)

//var ecsProjectProps containerpatterns.EcsProjectProps

func Example() {

	ecsProject := containerpatterns.NewEcsProject(awscdk.NewStack(awscdk.NewApp(nil), jsii.String("LoadBalancedEc2ServiceStack"), &awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: jsii.String("305251478828"),
			Region:  jsii.String("us-east-1"),
		},
	}), jsii.String("EcsProject"), &containerpatterns.EcsProjectProps{
		NetworkProps: network.VpcProps{
			IsDefault: true,
		},
		ComputeProps: containerpatterns.EcsComputeProps{
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "GolangCdkDemo",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: false,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GolangCdkDemoT2Micro",
						InstanceClass: ec2.InstanceClass_BURSTABLE2,
						InstanceSize:  ec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GolangCdkDemoT2Micro",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GolangCdkDemoT2Small",
						InstanceClass: ec2.InstanceClass_BURSTABLE2,
						InstanceSize:  ec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GolangCdkDemoT2Small",
					},
				},
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:        "golang-cdk-demo-" + *awscdk.Aws_REGION(),
				IsVersioned: false,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "GolangCdkDemo",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "golang-cdk.demo",
				Description: "Golang CDK demo service discovery",
			},
		},
		NonLoadBalancedEc2ServicesProps: []brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
			{
				Cluster: brzNlbEc2Service.ClusterProps{
					Vpc: network.VpcProps{
						IsDefault: true,
					},
				},
				LogGroupName: "GolangCdkDemoDb",
				TaskDefinition: brzNlbEc2Service.TaskDefinition{
					FamilyName:     "rpc-service-db",
					NetworkMode:    brzNlbEc2Service.TaskDefintionNetworkModeAwsVpc,
					RequiresVolume: true,
					Volumes: []brzNlbEc2Service.Volume{
						{
							Name: "rpc-service-db",
							Size: "10",
						},
					},
					ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
						{
							ContainerName:            "rpc-service-db",
							Image:                    "rpc-service-db",
							RegistryType:             brzNlbEc2Service.ContainerDefinitionRegistryAwsEcr,
							ImageTag:                 "latest",
							IsEssential:              true,
							Cpu:                      512,
							Memory:                   896,
							EnvironmentFileObjectKey: "rpc-service-db/prod/db.env",
							VolumeMountPoint: []ecs.MountPoint{
								{
									ContainerPath: jsii.String("/var/lib/postgresql/data"),
									ReadOnly:      jsii.Bool(false),
									SourceVolume:  jsii.String("rpc-service-db"),
								},
							},
							PortMappings: []ecs.PortMapping{
								{
									ContainerPort: jsii.Number(5432),
									Protocol:      ecs.Protocol_TCP,
								},
							},
						},
					},
				},
				IsTracingEnabled: false,
				DesiredTaskCount: 1,
				CapacityProviders: []string{
					"GolangCdkDemoT2Micro",
				},
				IsServiceDiscoveryEnabled: true,
				ServiceDiscovery: brzNlbEc2Service.ServiceDiscoveryProps{
					ServiceName: "rpc-service-db",
					ServicePort: 5432,
				},
			},
		},
		LoadBalancedEc2ServicesProps: []brzLbEc2Service.LoadBalancedEc2ServiceProps{
			{
				Cluster: brzLbEc2Service.ClusterProps{
					ClusterName: "<dummy>",
					Vpc: network.VpcProps{
						Id:        "<dummy>",
						IsDefault: true,
					},
				},
				LogGroupName: "GolangCdkDemoService",
				TaskDefinition: brzLbEc2Service.TaskDefinition{
					FamilyName:     "rpc-service",
					NetworkMode:    brzLbEc2Service.TaskDefintionNetworkModeBridge,
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
				//				IsLoadBalancerEnabled:     true,
				LoadBalancer: brzLbEc2Service.LoadBalancerProps{
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
			},
		},
	})
	fmt.Println(len(ecsProject.Ec2ContainerApplicationServices()))
	fmt.Println(len(*ecsProject.ApplicationLoadBalancer().Listeners()))
	// Output:
	// 2
	// 0
}

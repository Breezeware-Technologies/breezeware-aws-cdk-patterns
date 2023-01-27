package containerpatterns

import (
	brzLbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/load_balanced"
	brzNlbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/non_load_balanced"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	//	brzLbEc2Service "breezeware-aws-cdk-patterns-samples/container_patterns/load_balanced"
	//	brzNlbEc2Service "breezeware-aws-cdk-patterns-samples/container_patterns/non_load_balanced"
	//	brznetwork "breezeware-aws-cdk-patterns-samples/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elb2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	//	"log"

	//	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	//	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"strconv"
)

// EcsProjectProps properties for creating a EcsProject
type EcsProjectProps struct {
	NetworkProps                    brznetwork.VpcProps
	ComputeProps                    EcsComputeProps
	NonLoadBalancedEc2ServicesProps []brzNlbEc2Service.NonLoadBalancedEc2ServiceProps
	LoadBalancedEc2ServicesProps    []brzLbEc2Service.LoadBalancedEc2ServiceProps
}

// EcsProject containing network, compute and relavant ec2ContainerApplicationServices in AWS ECS
type ecsProject struct {
	constructs.Construct
	network                         ec2.IVpc
	compute                         ecs.Cluster
	applicationLoadBalancer         elb2.ApplicationLoadBalancer
	ec2ContainerApplicationServices []ecs.Ec2Service
}

type EcsProject interface {
	Network() ec2.IVpc
	Compute() ecs.Cluster
	ApplicationLoadBalancer() elb2.ApplicationLoadBalancer
	Ec2ContainerApplicationServices() []ecs.Ec2Service
}

func (prj *ecsProject) Network() ec2.IVpc {
	return prj.network
}

func (prj *ecsProject) Compute() ecs.Cluster {
	return prj.compute
}

func (prj *ecsProject) ApplicationLoadBalancer() elb2.ApplicationLoadBalancer {
	return prj.applicationLoadBalancer
}

func (prj *ecsProject) Ec2ContainerApplicationServices() []ecs.Ec2Service {
	return prj.ec2ContainerApplicationServices
}

type ComputeStackProps struct {
	core.NestedStackProps
	Vpc ec2.IVpc
}

// Example:
//
//	brzecs.NewEcsProject(stack, jsii.String("EcsProject"), &brzecs.EcsProjectProps{
//	   NetworkProps: brznetwork.VpcProps{
//	       IsDefault: true,
//	       },
//	       ComputeProps: brzecs.EcsComputeProps{
//	       VpcId: "<dummy>",
//	       Cluster: brzecs.ClusterOptions{
//	           Name:                             "GolangCdkDemo",
//	           ContainerInsights:                false,
//	           IsAsgCapacityProviderEnabled:     true,
//	           IsFargateCapacityProviderEnabled: false,
//	           },
//	           AsgCapacityProviders: []brzecs.AsgCapacityProviders{
//	           {
//	               AutoScalingGroup: brzecs.AsgProps{
//	                   Name:          "GolangCdkDemoT2Micro",
//	                   InstanceClass: ec2.InstanceClass_BURSTABLE2,
//	                   InstanceSize:  ec2.InstanceSize_MICRO,
//	                   MinCapacity:   0,
//	                   MaxCapacity:   2,
//	                   SshKeyName:    "breezethru-demo-key-pair",
//	                   },
//	                   CapacityProvider: brzecs.AsgCapacityProviderProps{
//	                   Name: "GolangCdkDemoT2Micro",
//	                   },
//	                   },
//	                   {
//	               AutoScalingGroup: brzecs.AsgProps{
//	                   Name:          "GolangCdkDemoT2Small",
//	                   InstanceClass: ec2.InstanceClass_BURSTABLE2,
//	                   InstanceSize:  ec2.InstanceSize_SMALL,
//	                   MinCapacity:   0,
//	                   MaxCapacity:   2,
//	                   SshKeyName:    "breezethru-demo-key-pair",
//	                   },
//	                   CapacityProvider: brzecs.AsgCapacityProviderProps{
//	                   Name: "GolangCdkDemoT2Small",
//	                   },
//	                   },
//	                   },
//	                   EnvironmentFileBucket: brzecs.BucketOptions{
//	           Name:        "golang-cdk-demo-" + *awscdk.Aws_REGION(),
//	           IsVersioned: false,
//	           },
//	           LoadBalancer: brzecs.LoadBalancerOptions{
//	           Name:                   "GolangCdkDemo",
//	           ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
//	           },
//	           CloudmapNamespace: brzecs.CloudmapNamespaceProps{
//	           Name:        "golang-cdk.demo",
//	           Description: "Golang CDK demo service discovery",
//	           },
//	           },
//	           NonLoadBalancedEc2ServicesProps: []brzNlbEc2Service.NonLoadBalancedEc2ServiceProps{
//	       {
//	           Cluster: brzNlbEc2Service.ClusterProps{
//	               ClusterName: "<dummy>",
//	               Vpc: brznetwork.VpcProps{
//	                   Id:        "<dummy>",
//	                   IsDefault: true,
//	                   },
//	                   },
//	                   LogGroupName: "GolangCdkDemoDb",
//	                   TaskDefinition: brzNlbEc2Service.TaskDefinition{
//	               FamilyName:  "rpc-service-db",
//	               NetworkMode: brzNlbEc2Service.TASK_DEFINTION_NETWORK_MODE_AWS_VPC,
//	               EnvironmentFile: brzNlbEc2Service.EnvironmentFile{
//	                   BucketName: "<dummy>",
//	                   BucketArn:  "<dummy>",
//	                   },
//	                   RequiresVolume: true,
//	                   Volumes: []brzNlbEc2Service.Volume{
//	                   {
//	                       Name: "rpc-service-db",
//	                       Size: "10",
//	                       },
//	                       },
//	                       ApplicationContainers: []brzNlbEc2Service.ContainerDefinition{
//	                   {
//	                       ContainerName:            "rpc-service-db",
//	                       Image:                    "rpc-service-db",
//	                       RegistryType:             brzNlbEc2Service.CONTAINER_DEFINITION_REGISTRY_AWS_ECR,
//	                       ImageTag:                 "latest",
//	                       IsEssential:              true,
//	                       Cpu:                      512,
//	                       Memory:                   896,
//	                       EnvironmentFileObjectKey: "rpc-service-db/prod/db.env",
//	                       VolumeMountPoint: []ecs.MountPoint{
//	                           {
//	                               ContainerPath: jsii.String("/var/lib/postgresql/data"),
//	                               ReadOnly:      jsii.Bool(false),
//	                               SourceVolume:  jsii.String("rpc-service-db"),
//	                               },
//	                               },
//	                               PortMappings: []ecs.PortMapping{
//	                           {
//	                               ContainerPort: jsii.Number(5432),
//	                               Protocol:      ecs.Protocol_TCP,
//	                               },
//	                               },
//	                               },
//	                               },
//	                               },
//	                               IsTracingEnabled: false,
//	                               DesiredTaskCount: 1,
//	                               CapacityProviders: []string{
//	               "GolangCdkDemoT2Micro",
//	               },
//	               IsServiceDiscoveryEnabled: true,
//	               ServiceDiscovery: brzNlbEc2Service.ServiceDiscoveryProps{
//	               ServiceName: "rpc-service-db",
//	               ServicePort: 5432,
//	               CloudMapNamespace: brzNlbEc2Service.CloudMapNamespaceProps{
//	                   NamespaceName: "<dummy>",
//	                   NamespaceId:   "<dummy>",
//	                   NamespaceArn:  "<dummy>",
//	                   },
//	                   },
//	                   },
//	                   },
//	                   LoadBalancedEc2ServicesProps: []brzLbEc2Service.LoadBalancedEc2ServiceProps{
//	       {
//	           Cluster: brzLbEc2Service.ClusterProps{
//	               ClusterName: "<dummy>",
//	               Vpc: brznetwork.VpcProps{
//	                   Id:        "<dummy>",
//	                   IsDefault: true,
//	                   },
//	                   },
//	                   LogGroupName: "GolangCdkDemoService",
//	                   TaskDefinition: brzLbEc2Service.TaskDefinition{
//	               FamilyName:  "rpc-service",
//	               NetworkMode: brzLbEc2Service.TASK_DEFINTION_NETWORK_MODE_BRIDGE,
//	               EnvironmentFile: brzLbEc2Service.EnvironmentFile{
//	                   BucketName: "<dummy>",
//	                   BucketArn:  "<dummy>",
//	                   },
//	                   RequiresVolume: false,
//	                   ApplicationContainers: []brzLbEc2Service.ContainerDefinition{
//	                   {
//	                       ContainerName:            "rpc-service",
//	                       Image:                    "rpc-service",
//	                       RegistryType:             brzLbEc2Service.CONTAINER_DEFINITION_REGISTRY_AWS_ECR,
//	                       ImageTag:                 "latest",
//	                       IsEssential:              true,
//	                       Cpu:                      512,
//	                       Memory:                   1458,
//	                       EnvironmentFileObjectKey: "rpc-service/prod/app.env",
//	                       PortMappings: []ecs.PortMapping{
//	                           {
//	                               ContainerPort: jsii.Number(8443),
//	                               Protocol:      ecs.Protocol_TCP,
//	                               },
//	                               },
//	                               },
//	                               },
//	                               },
//	                               IsTracingEnabled: true,
//	                               DesiredTaskCount: 1,
//	                               CapacityProviders: []string{
//	               "GolangCdkDemoT2Small",
//	               },
//	               IsServiceDiscoveryEnabled: false,
//	               ServiceDiscovery:          brzLbEc2Service.ServiceDiscoveryProps{},
//	               IsLoadBalancerEnabled:     true,
//	               LoadBalancer: brzLbEc2Service.LoadBalancerProps{
//	               ListenerArn:           "<dummy>",
//	               SecurityGroupId:       "<dummy>",
//	               TargetHealthCheckPath: "/api/health-status",
//	               ListenerRuleProps: brzLbEc2Service.ListenerRuleProps{
//	                   Priority:      1,
//	                   PathCondition: "/api/*",
//	                   HostCondition: "nginx.dynamostack.com",
//	                   },
//	                   },
//	                   LoadBalancerTargetOptions: ecs.LoadBalancerTargetOptions{
//	               ContainerName: jsii.String("rpc-service"),
//	               ContainerPort: jsii.Number(8443),
//	               Protocol:      ecs.Protocol_TCP,
//	               },
//	               },
//	               },
//	               })
func NewEcsProject(scope constructs.Construct, id *string, props *EcsProjectProps) EcsProject {
	this := constructs.NewConstruct(scope, id)

	// network stack configuration
	var vpcLookupOptions ec2.VpcLookupOptions = ec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	}
	if props.NetworkProps != (brznetwork.VpcProps{}) {
		var networkProps brznetwork.VpcProps = props.NetworkProps
		if !networkProps.IsDefault {
			vpcLookupOptions = ec2.VpcLookupOptions{
				IsDefault: jsii.Bool(false),
				Region:    core.Aws_REGION(),
				VpcId:     jsii.String(networkProps.Id),
			}
		}
	}
	vpc := ec2.Vpc_FromLookup(this, jsii.String("Network"), &vpcLookupOptions)

	// storage configuration
	//	storageStack := core.NewNestedStack(this, jsii.String("StorageStack"), &core.NestedStackProps{})
	//	storageProps := props.ComputeProps.EnvironmentFileBucket

	// compute stack configuration
	computeStack := core.NewNestedStack(this, jsii.String("ComputeStack"), &core.NestedStackProps{
		Description: jsii.String("Cloudformation stack for handling Container based compute resources"),
	})

	computeProps := props.ComputeProps
	computeProps.VpcId = *vpc.VpcId()
	ecsContainerCompute := NewContainerCompute(computeStack, jsii.String("EcsCompute"), &computeProps)

	// nested stack references
	cluster := ecsContainerCompute.Cluster()
	clusterSecurityGroups := ecsContainerCompute.ClusterSecurityGroups()
	applicationLoadBalancer := ecsContainerCompute.LoadBalancer()
	cloudMapNamespace := ecsContainerCompute.CloudMapNamespace()
	asgCapacityProviders := ecsContainerCompute.AsgCapacityProviders()
	loadBalancerSecurityGroupId := ecsContainerCompute.LoadBalancerSecurityGroup().SecurityGroupId()
	httpsListenerArn := ecsContainerCompute.HttpsListener().ListenerArn()
	environmentFileBucket := ecsContainerCompute.EnvironmentFileBucket()

	// applications stack configuration
	var capacityProviders []string
	//	var lbEc2ServiceCapacityProviderStrategies []*ecs.CapacityProviderStrategy
	for _, asgCapacityProvider := range asgCapacityProviders {
		capacityProviders = append(capacityProviders, *asgCapacityProvider.CapacityProviderName())
		// creates CapacityProviderStrategy from the CapacityProviderName
		//		lbEc2ServiceCapacityProviderStrategies = append(lbEc2ServiceCapacityProviderStrategies, &ecs.CapacityProviderStrategy{
		//			CapacityProvider: asgCapacityProvider.CapacityProviderName(),
		//			Base:             jsii.Number(0),
		//			Weight:           jsii.Number(1),
		//		})
	}

	var ec2Services []ecs.Ec2Service
	var nonLoadBalancedServicesStack core.NestedStack
	if len(props.NonLoadBalancedEc2ServicesProps) > 0 {
		// non load-balanced applications stack configuration
		nonLoadBalancedServicesStack = core.NewNestedStack(computeStack, jsii.String("NonLoadBalancedEc2ContainerApplicationServicesStack"), &core.NestedStackProps{})
		for index, nlbEc2ServiceProps := range props.NonLoadBalancedEc2ServicesProps {

			nlbEc2ServiceProps.Cluster = brzNlbEc2Service.ClusterProps{
				ClusterName: *cluster.ClusterName(),
				Vpc: brznetwork.VpcProps{
					Id:        computeProps.VpcId,
					IsDefault: nlbEc2ServiceProps.Cluster.Vpc.IsDefault,
				},
				SecurityGroups: clusterSecurityGroups,
			}
			//			var nlbEc2ServiceCapacityProviders []string
			//			for _, provider1 := range nlbEc2ServiceProps.CapacityProviders {
			//				log.Println("Checking for capacity provider with name: ", provider1)
			//				for _, provider2 := range capacityProviders {
			//					log.Println("Available capacity provider name: ", provider2)
			//					if provider1 == provider2 {
			//						log.Printf("User provided capacity provider %v matched with provider %v ", provider1, provider2)
			//						nlbEc2ServiceCapacityProviders = append(nlbEc2ServiceCapacityProviders, provider2)
			//						break
			//					}
			//				}
			//			}
			//			log.Printf("Capacity Providers for Non-Load Balanced EC2 Services: %v", nlbEc2ServiceCapacityProviders)
			//			nlbEc2ServiceProps.CapacityProviders = nlbEc2ServiceCapacityProviders
			//			log.Println("Capacity Providers: ", nlbEc2ServiceProps.CapacityProviders)
			if nlbEc2ServiceProps.IsServiceDiscoveryEnabled {
				nlbEc2ServiceProps.ServiceDiscovery = brzNlbEc2Service.ServiceDiscoveryProps{
					ServiceName: nlbEc2ServiceProps.ServiceDiscovery.ServiceName,
					ServicePort: nlbEc2ServiceProps.ServiceDiscovery.ServicePort,
					CloudMapNamespace: brzNlbEc2Service.CloudMapNamespaceProps{
						NamespaceName: *cloudMapNamespace.NamespaceName(),
						NamespaceId:   *cloudMapNamespace.NamespaceId(),
						NamespaceArn:  *cloudMapNamespace.NamespaceArn(),
					},
				}
			}
			nlbEc2ServiceProps.TaskDefinition.EnvironmentFile = brzNlbEc2Service.EnvironmentFile{
				BucketName: *environmentFileBucket.BucketName(),
				BucketArn:  *environmentFileBucket.BucketArn(),
			}
			nonLoadBalancedEc2Service := brzNlbEc2Service.NewNonLoadBalancedEc2Service(nonLoadBalancedServicesStack, jsii.String("LoadBalancedService"+strconv.FormatInt(int64(index), 10)), &nlbEc2ServiceProps)
			ec2Services = append(ec2Services, nonLoadBalancedEc2Service.Service())
		}
	}

	if len(props.LoadBalancedEc2ServicesProps) > 0 {
		// load-balanced applications stack configuration
		loadBalancedServicesStack := core.NewNestedStack(computeStack, jsii.String("LoadBalancedEc2ContainerApplicationServicesStack"), &core.NestedStackProps{
			Description: jsii.String("Cloudformation stack for handling Load Balanced ECS EC2Service(s) (applications)"),
		})
		// waits for non load-balanced applications stack creation to complete
		loadBalancedServicesStack.AddDependency(nonLoadBalancedServicesStack, jsii.String("Wait for non load-balanced applications to start before starting the load balanced applications"))

		for index, lbEc2ServiceProps := range props.LoadBalancedEc2ServicesProps {

			lbEc2ServiceProps.Cluster = brzLbEc2Service.ClusterProps{
				ClusterName: *cluster.ClusterName(),
				Vpc: brznetwork.VpcProps{
					Id:        computeProps.VpcId,
					IsDefault: lbEc2ServiceProps.Cluster.Vpc.IsDefault,
				},
				SecurityGroups: clusterSecurityGroups,
			}
			//			var lbEc2ServiceCapacityProviders []string
			//			for _, provider1 := range lbEc2ServiceProps.CapacityProviders {
			//				for _, provider2 := range capacityProviders {
			//					if provider1 == provider2 {
			//						lbEc2ServiceCapacityProviders = append(lbEc2ServiceCapacityProviders, provider2)
			//						break
			//					}
			//				}
			//			}
			//			log.Printf("Capacity Providers for Load Balanced EC2 Services: %v", lbEc2ServiceCapacityProviders)
			//			lbEc2ServiceProps.CapacityProviders = lbEc2ServiceCapacityProviders
			//			log.Println("Capacity Providers: ", lbEc2ServiceProps.CapacityProviders)
			lbEc2ServiceProps.ServiceDiscovery = brzLbEc2Service.ServiceDiscoveryProps{
				ServiceName: lbEc2ServiceProps.ServiceDiscovery.ServiceName,
				ServicePort: lbEc2ServiceProps.ServiceDiscovery.ServicePort,
				CloudMapNamespace: brzLbEc2Service.CloudMapNamespaceProps{
					NamespaceName: *cloudMapNamespace.NamespaceName(),
					NamespaceId:   *cloudMapNamespace.NamespaceId(),
					NamespaceArn:  *cloudMapNamespace.NamespaceArn(),
				},
			}
			lbEc2ServiceProps.IsLoadBalancerEnabled = true
			lbEc2ServiceProps.LoadBalancer = brzLbEc2Service.LoadBalancerProps{
				SecurityGroupId:       *loadBalancerSecurityGroupId,
				TargetHealthCheckPath: lbEc2ServiceProps.LoadBalancer.TargetHealthCheckPath,
				ListenerArn:           *httpsListenerArn,
				ListenerRuleProps: brzLbEc2Service.ListenerRuleProps{
					Priority:      lbEc2ServiceProps.LoadBalancer.ListenerRuleProps.Priority,
					PathCondition: lbEc2ServiceProps.LoadBalancer.ListenerRuleProps.PathCondition,
					HostCondition: lbEc2ServiceProps.LoadBalancer.ListenerRuleProps.HostCondition,
				},
			}
			lbEc2ServiceProps.TaskDefinition.EnvironmentFile = brzLbEc2Service.EnvironmentFile{
				BucketName: *environmentFileBucket.BucketName(),
				BucketArn:  *environmentFileBucket.BucketArn(),
			}
			loadBalancedEc2Service := brzLbEc2Service.NewLoadBalancedEc2Service(loadBalancedServicesStack, jsii.String("LoadBalancedService"+strconv.FormatInt(int64(index), 10)), &lbEc2ServiceProps)
			ec2Services = append(ec2Services, loadBalancedEc2Service.Service())
		}
	}

	return &ecsProject{this, vpc, cluster, applicationLoadBalancer, ec2Services}
}

// Package containerpatterns provides level 3 CDK constructs for container based patterns.
//
// Provides patterns like:
//   - load-balanced ECS service based on EC2.
//   - non load-balanced ECS service based on EC2.
package containerpatterns

import (
	brzLbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/load_balanced"
	brzNlbEc2Service "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns/non_load_balanced"
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elb2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"strconv"
)

// A EcsProjectProps represents properties for creating an EcsProject construct.
type EcsProjectProps struct {
	NetworkProps                    brznetwork.VpcProps                               // network-related properties of the ECS project
	ComputeProps                    EcsComputeProps                                   // compute configuration properties of the ECS project
	NonLoadBalancedEc2ServicesProps []brzNlbEc2Service.NonLoadBalancedEc2ServiceProps // list of Load-Balanced ECS EC2-based Service properties
	LoadBalancedEc2ServicesProps    []brzLbEc2Service.LoadBalancedEc2ServiceProps     // list of non Load-Balanced ECS EC2-based Service properties
}

// ecsProject represents the container-based ECS project pattern/construct.
type ecsProject struct {
	constructs.Construct
	network                         ec2.IVpc
	compute                         ecs.Cluster
	applicationLoadBalancer         elb2.ApplicationLoadBalancer
	ec2ContainerApplicationServices []ecs.Ec2Service
}

// EcsProject provides implementations for the ecsProject.
type EcsProject interface {
	// Network returns an IVpc as the network-working component of the EcsProject.
	Network() ec2.IVpc
	// Compute returns Cluster as the container compute environment.
	Compute() ecs.Cluster
	// ApplicationLoadBalancer returns the load-balancer associated for the EcsProject as a part of managing load & traffic.
	ApplicationLoadBalancer() elb2.ApplicationLoadBalancer
	// Ec2ContainerApplicationServices returns a slice of ECS EC2 based services that are currently in the EcsProject
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

// NewEcsProject creates a new ECS based construct with cloudformation stack(s)
// for managing the compute and nested stack(s) for the applications inside the compute stack from EcsProjectProps
func NewEcsProject(scope constructs.Construct, id *string, props *EcsProjectProps) EcsProject {
	this := constructs.NewConstruct(scope, id)

	// network stack creation (parent stack)
	// VpcLookupOptions with default value for lookingup the default VPC of the account
	var vpcLookupOptions ec2.VpcLookupOptions = ec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	}
	if props.NetworkProps != (brznetwork.VpcProps{}) {
		var networkProps brznetwork.VpcProps = props.NetworkProps
		// Overrides the default values configured from th actual props
		if !networkProps.IsDefault {
			vpcLookupOptions = ec2.VpcLookupOptions{
				IsDefault: jsii.Bool(false),
				Region:    core.Aws_REGION(),
				VpcId:     jsii.String(networkProps.Id),
			}
		}
	}
	// Looks-up VPC using CDK CLI
	vpc := ec2.Vpc_FromLookup(this, jsii.String("Network"), &vpcLookupOptions)

	// compute nested stack creation
	computeStack := core.NewNestedStack(this, jsii.String("ComputeStack"), &core.NestedStackProps{
		Description: jsii.String("Cloudformation stack for handling Container based compute resources"),
	})

	computeProps := props.ComputeProps
	// compute nested stack references
	computeProps.VpcId = *vpc.VpcId()
	ecsContainerCompute := NewContainerCompute(computeStack, jsii.String("EcsCompute"), &computeProps)

	// applications nested stack references
	cluster := ecsContainerCompute.Cluster()
	clusterSecurityGroups := ecsContainerCompute.ClusterSecurityGroups()
	applicationLoadBalancer := ecsContainerCompute.ApplicationLoadBalancer()
	cloudMapNamespace := ecsContainerCompute.CloudMapNamespace()
	asgCapacityProviders := ecsContainerCompute.AsgCapacityProviders()
	loadBalancerSecurityGroupId := ecsContainerCompute.AlbSecurityGroup().SecurityGroupId()
	httpsListenerArn := ecsContainerCompute.AlbHttpsListener().ListenerArn()
	environmentFileBucket := ecsContainerCompute.EnvironmentFileBucket()

	// applications nested stack with load-balanced & non load-balanced applications
	var capacityProviders []string
	for _, asgCapacityProvider := range asgCapacityProviders {
		capacityProviders = append(capacityProviders, *asgCapacityProvider.CapacityProviderName())
	}

	var ec2Services []ecs.Ec2Service
	var nonLoadBalancedServicesStack core.NestedStack
	// creates a non load-balanced applications nested stack only if the length of the props.NonLoadBalancedEc2ServicesProps is greater than zero or not empty
	if len(props.NonLoadBalancedEc2ServicesProps) > 0 {
		// non load-balanced applications nested stack creation
		nonLoadBalancedServicesStack = core.NewNestedStack(computeStack, jsii.String("NonLoadBalancedEc2ContainerApplicationServicesStack"), &core.NestedStackProps{})
		// Overrides the values from the actual props.NonLoadBalancedEc2ServicesProps with nested stack references
		for index, nlbEc2ServiceProps := range props.NonLoadBalancedEc2ServicesProps {

			nlbEc2ServiceProps.Cluster = brzNlbEc2Service.ClusterProps{
				ClusterName: *cluster.ClusterName(),
				Vpc: brznetwork.VpcProps{
					Id:        computeProps.VpcId,
					IsDefault: nlbEc2ServiceProps.Cluster.Vpc.IsDefault,
				},
				SecurityGroups: clusterSecurityGroups,
			}
			// Overrides ServiceDiscovery properties if service discovery is enabled for the service
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

	// creates a load-balanced applications nested stack only if the length of the props.LoadBalancedEc2ServicesProps is greater than zero or not empty
	if len(props.LoadBalancedEc2ServicesProps) > 0 {
		// non load-balanced applications nested stack creation
		loadBalancedServicesStack := core.NewNestedStack(computeStack, jsii.String("LoadBalancedEc2ContainerApplicationServicesStack"), &core.NestedStackProps{
			Description: jsii.String("Cloudformation stack for handling Load Balanced ECS EC2Service(s) (applications)"),
		})

		if len(props.NonLoadBalancedEc2ServicesProps) > 0 {
			// waits for non load-balanced applications stack creation to complete. Specifically used for cases like Database service to start before application service creation
			loadBalancedServicesStack.AddDependency(nonLoadBalancedServicesStack, jsii.String("Wait for non load-balanced applications to start before starting the load balanced applications"))
		}

		// Overrides the values from the actual props.LoadBalancedEc2ServicesProps with nested stack references
		for index, lbEc2ServiceProps := range props.LoadBalancedEc2ServicesProps {

			lbEc2ServiceProps.Cluster = brzLbEc2Service.ClusterProps{
				ClusterName: *cluster.ClusterName(),
				Vpc: brznetwork.VpcProps{
					Id:        computeProps.VpcId,
					IsDefault: lbEc2ServiceProps.Cluster.Vpc.IsDefault,
				},
				SecurityGroups: clusterSecurityGroups,
			}
			lbEc2ServiceProps.ServiceDiscovery = brzLbEc2Service.ServiceDiscoveryProps{
				ServiceName: lbEc2ServiceProps.ServiceDiscovery.ServiceName,
				ServicePort: lbEc2ServiceProps.ServiceDiscovery.ServicePort,
				CloudMapNamespace: brzLbEc2Service.CloudMapNamespaceProps{
					NamespaceName: *cloudMapNamespace.NamespaceName(),
					NamespaceId:   *cloudMapNamespace.NamespaceId(),
					NamespaceArn:  *cloudMapNamespace.NamespaceArn(),
				},
			}
//			lbEc2ServiceProps.IsLoadBalancerEnabled = true
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

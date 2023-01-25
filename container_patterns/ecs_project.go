package containerpatterns

import (
	//	brznetwork "breezeware-aws-cdk-patterns-samples/network"
	//	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"strconv"
)

type EcsProjectProps struct {
	NetworkProps                    brznetwork.VpcProps
	ComputeProps                    EcsComputeProps
	NonLoadBalancedEc2ServicesProps []ecs.Ec2ServiceProps
	LoadBalancedEc2ServicesProps    []LoadBalancedEc2ServiceProps
}

// ECS Load Balanced Container applications ECS
type ecsProject struct {
	constructs.Construct
	network      ec2.IVpc
	compute      ecs.Cluster
	applications []ecs.Ec2Service
}

type EcsProject interface {
	Network() ec2.IVpc
	Compute() ecs.Cluster
	Applications() []ecs.Ec2Service
}

func (prj *ecsProject) Network() ec2.IVpc {
	return prj.network
}

func (prj *ecsProject) Compute() ecs.Cluster {
	return prj.compute
}

func (prj *ecsProject) Applications() []ecs.Ec2Service {
	return prj.applications
}

type ComputeStackProps struct {
	core.NestedStackProps
	Vpc ec2.IVpc
}

func NewEcsProject(scope constructs.Construct, id *string, props *EcsProjectProps) EcsProject {
	this := constructs.NewConstruct(scope, id)

	// network configuration
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
				//				VpcName:   jsii.String(networkProps.Name),
			}
		}
	}
	vpcFromLookup := ec2.Vpc_FromLookup(this, jsii.String("Network"), &vpcLookupOptions)

	// storage configuration
	//	storageStack := core.NewNestedStack(this, jsii.String("StorageStack"), &core.NestedStackProps{})
	//	storageProps := props.ComputeProps.EnvironmentFileBucket

	// compute configuration
	computeStack := core.NewNestedStack(this, jsii.String("ComputeStack"), &core.NestedStackProps{})

	computeProps := props.ComputeProps
	computeProps.VpcId = *vpcFromLookup.VpcId()
	ecsContainerCompute := NewContainerCompute(computeStack, jsii.String("EcsCompute"), &computeProps)

	cluster := ecsContainerCompute.Cluster()
	clusterSecurityGroups := ecsContainerCompute.ClusterSecurityGroups()
	cloudMapNamespace := ecsContainerCompute.CloudMapNamespace()
	asgCapacityProviders := ecsContainerCompute.AsgCapacityProviders()
	loadBalancerSecurityGroupId := ecsContainerCompute.LoadBalancerSecurityGroup().SecurityGroupId()
	httpsListenerArn := ecsContainerCompute.HttpsListener().ListenerArn()
	environmentFileBucket := ecsContainerCompute.EnvironmentFileBucket()

	// applications configuration
	var services []ecs.Ec2Service
	var capacityProviderNames []string
	var capacityProviderStrategies []*ecs.CapacityProviderStrategy
	for _, asgCapacityProvider := range asgCapacityProviders {
		capacityProviderNames = append(capacityProviderNames, *asgCapacityProvider.CapacityProviderName())
		capacityProviderStrategies = append(capacityProviderStrategies, &ecs.CapacityProviderStrategy{
			CapacityProvider: asgCapacityProvider.CapacityProviderName(),
			Base:             jsii.Number(0),
			Weight:           jsii.Number(1),
		})
	}
	// non load-balanced applications
	//	nonLoadBalancedServicesStack := core.NewNestedStack(computeStack, jsii.String("NonLoadBalancedServicesStack"), &core.NestedStackProps{})
	//	for index, nlbServiceProps := range props.NonLoadBalancedEc2ServicesProps {
	//		nonLoadBalancedEc2ServiceProps := nlbServiceProps
	//		nonLoadBalancedEc2ServiceProps.Cluster = cluster
	//		nonLoadBalancedEc2ServiceProps.CloudMapOptions = &ecs.CloudMapOptions{
	//			CloudMapNamespace: cloudMapNamespace,
	//			Container:         nlbServiceProps.CloudMapOptions.Container,
	//			ContainerPort:     nlbServiceProps.CloudMapOptions.ContainerPort,
	//			DnsRecordType:     nlbServiceProps.CloudMapOptions.DnsRecordType,
	//			DnsTtl:            nlbServiceProps.CloudMapOptions.DnsTtl,
	//			FailureThreshold:  nlbServiceProps.CloudMapOptions.FailureThreshold,
	//			Name:              nlbServiceProps.CloudMapOptions.Name,
	//		}
	//		nonLoadBalancedEc2ServiceProps.CapacityProviderStrategies = &capacityProviderStrategies
	//		for _, containerDef := range *nonLoadBalancedEc2ServiceProps.TaskDefinition.Containers() {
	//			for _, contEnvFile := range *containerDef.EnvironmentFiles() {
	//				envFile := *contEnvFile
	//				envFile.S3Location = awss3.Location{
	//					BucketName:    containerCompute.EnvironmentFileBucket(),
	//					ObjectKey:     nil,
	//					ObjectVersion: nil,
	//				}
	//			}
	//			//			containerDef.EnvironmentFiles() = []ecs.EnvironmentFile{
	//			//				ecs.AssetEnvironmentFile_FromBucket(jsii.String(""), jsii.String(containerDef.EnvironmentFiles()), nil),
	//			//			}
	//		}
	//
	//		nonLoadBalancedService := ecs.NewEc2Service(nonLoadBalancedServicesStack, jsii.String("NonLoadBalancedService"+strconv.FormatInt(int64(index), 10)), &nonLoadBalancedEc2ServiceProps)
	//		services = append(services, nonLoadBalancedService)
	//
	//	}

	// load-balanced applications configuration
	loadBalancedServicesStack := core.NewNestedStack(computeStack, jsii.String("LoadBalancedServicesStack"), &core.NestedStackProps{})
	for index, lbServiceProps := range props.LoadBalancedEc2ServicesProps {

		loadBalancedEc2ServiceProps := lbServiceProps
		loadBalancedEc2ServiceProps.Cluster = ClusterProps{
			ClusterName: *cluster.ClusterName(),
			Vpc: brznetwork.VpcProps{
				Id:        computeProps.VpcId,
				IsDefault: true,
			},
			SecurityGroups: clusterSecurityGroups,
		}
		loadBalancedEc2ServiceProps.CapacityProviderStrategies = capacityProviderNames
		loadBalancedEc2ServiceProps.ServiceDiscovery = ServiceDiscoveryProps{
			ServiceName: lbServiceProps.ServiceDiscovery.ServiceName,
			ServicePort: lbServiceProps.ServiceDiscovery.ServicePort,
			CloudMapNamespace: CloudMapNamespaceProps{
				NamespaceName: *cloudMapNamespace.NamespaceName(),
				NamespaceId:   *cloudMapNamespace.NamespaceId(),
				NamespaceArn:  *cloudMapNamespace.NamespaceArn(),
			},
		}
		loadBalancedEc2ServiceProps.LoadBalancer = LoadBalancerProps{
			LoadBalancerSecurityGroupId: *loadBalancerSecurityGroupId,
			TargetHealthCheckPath:       lbServiceProps.LoadBalancer.TargetHealthCheckPath,
			ListenerArn:                 *httpsListenerArn,
			ListenerRuleProps: ListenerRuleProps{
				Priority:      lbServiceProps.LoadBalancer.ListenerRuleProps.Priority,
				PathCondition: lbServiceProps.LoadBalancer.ListenerRuleProps.PathCondition,
				HostCondition: lbServiceProps.LoadBalancer.ListenerRuleProps.HostCondition,
			},
		}
		loadBalancedEc2ServiceProps.TaskDefinition.EnvironmentFile = EnvironmentFile{
			BucketName: *environmentFileBucket.BucketName(),
			BucketArn:  *environmentFileBucket.BucketArn(),
		}
		loadBalancedEc2Service := NewLoadBalancedEc2Service(loadBalancedServicesStack, jsii.String("LoadBalancedService"+strconv.FormatInt(int64(index), 10)), &loadBalancedEc2ServiceProps)
		services = append(services, loadBalancedEc2Service.Service())
	}
	return &ecsProject{this, vpcFromLookup, cluster, services}
}

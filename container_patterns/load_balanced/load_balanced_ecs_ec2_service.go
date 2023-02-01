package containerpatterns

import (
	"math/rand"
	"strconv"

	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elb2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	cloudwatchlogs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	servicediscovery "github.com/aws/aws-cdk-go/awscdk/v2/awsservicediscovery"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// Custom types for enum handling
type (
	// Networkmode enum type
	Networkmode string
	// RegistryType enum type
	RegistryType string
	// LoadBalancerTargetProtocol enum type
	LoadBalancerTargetProtocol string
)

// constants
const (
	DefaultLogRetention       cloudwatchlogs.RetentionDays = cloudwatchlogs.RetentionDays_TWO_WEEKS // Default CloudWatch Log retention with value TWO_WEEKS
	DefaultDockerVolumeDriver string                       = "rexray/ebs"                           // Default Docker volume driver with value 'rexay/ebs'
	DefaultDockerVolumeType   string                       = "gp2"                                  // Default Docker volume type with value 'gp2'
	OtelContainerImage        string                       = "amazon/aws-otel-collector:v0.25.0"    // OTEL conatiner image with tag
)

const (
	TaskDefintionNetworkModeBridge   Networkmode     = "BRIDGE"               // Bridge mode Netwok mode of the ECS task definition
	TaskDefintionNetworkModeAwsVpc   Networkmode     = "AWS_VPC"              // AWS_VPC mode Netwok mode of the ECS task definition
	DefaultTaskDefinitionNetworkMode ecs.NetworkMode = ecs.NetworkMode_BRIDGE // Default Network mode of the ECS task definition with value ecs.NetworkMode_BRIDGE
)

const (
	ContainerDefinitionRegistryAwsEcr  RegistryType = "ECR"                             // ECR registry type for pulling image
	ContainerDefinitionRegistryOthers  RegistryType = "OTHERS"                          // Other registry types for pulling image
	DefaultContainerDefinitionRegistry RegistryType = ContainerDefinitionRegistryAwsEcr // Default registry type for pulling image with value ContainerDefinitionRegistryAwsEcr
)

const (
	loadBalancerTargetProtocolTcp     LoadBalancerTargetProtocol = "TCP"            // TCP target protocol for the Load-Balancer
	loadBalancerTargetProtocolUdp     LoadBalancerTargetProtocol = "UDP"            // UDP target protocol for the Load-Balancer
	DefaultLoadBalancerTargetProtocol ecs.Protocol               = ecs.Protocol_TCP // Default target protocol for the Load-Balancer with value ecs.Protocol_TCP
)

// LoadBalancedEc2ServiceProps represents the properties that are needed to create a Application Load-balanced EC2 Service inside ECS
type LoadBalancedEc2ServiceProps struct {
	Cluster                   ClusterProps                  `field:"required"` // cluster properties for the EC2 based Service
	LogGroupName              string                        `field:"required"` // name of the Log Group that will be created for the Load-Balanced EC2 based Service
	TaskDefinition            TaskDefinition                `field:"required"` // task-definition for the EC2 based service
	IsTracingEnabled          bool                          `field:"required"` // flag representing whether service level tracing is enabled or not
	DesiredTaskCount          float64                       `field:"required"` // number of task(s) that is desired to run at all time inside the EC2 based Service
	CapacityProviders         []string                      `field:"required"` // capacity providers to provision the EC2 instance infrastructure needed for the service's task(s) to run
	IsServiceDiscoveryEnabled bool                          `field:"required"` // flag representing whether service discovery is enabled or not for internal service identification/discovery inside the ECS Cluster
	ServiceDiscovery          ServiceDiscoveryProps         `field:"optional"` // service discovery properties if IsServiceDiscoveryEnabled flag is true
	LoadBalancer              LoadBalancerProps             `field:"required"` // application load-balancer properties for the EC2 based Service
	LoadBalancerTargetOptions ecs.LoadBalancerTargetOptions `field:"required"` // application load-balancer target configuration for the EC2 based Service
}

// ClusterProps represents the properties for retrieving a Cluster
type ClusterProps struct {
	ClusterName    string               `field:"required"` // name of the cluster
	Vpc            brznetwork.VpcProps  `field:"required"` // vpc properties
	SecurityGroups []ec2.ISecurityGroup `field:"required"` // security groups of the Auto-Scaling Groups associated with the Cluster if Auto-Scaling Group based Capacity providers are configured
}

// TaskDefinition represents the properties for creating a ecs.TaskDefinition
type TaskDefinition struct {
	FamilyName            string                `field:"optional"` // family name of the task-definition, used for grouping
	NetworkMode           Networkmode           `field:"required"` // network mode for the task-definition, can be either TaskDefintionNetworkModeBridge or TaskDefintionNetworkModeAwsVpc. Default value configured to DefaultTaskDefinitionNetworkMode
	EnvironmentFile       EnvironmentFile       `field:"required"` // environment file propeties
	TaskPolicy            iam.PolicyDocument    `field:"optional"` // task policy of the task-definition. Gives the container(s) specified in the TaskDefinition access to the AWS service(s)
	ApplicationContainers []ContainerDefinition `field:"required"` // application container definitions. Container specification like registry, image, tag, logging, etc
	RequiresVolume        bool                  `field:"required"` // flag representing whether the task-definition requires volume to persist data
	Volumes               []Volume              `field:"optional"` // volumes for the task-definition. Only used if RequiresVolume falg is true, else ommited even if configured
}

// EnvironmentFile represents the S3 Bucket options for environment file(s) in the TaskDefinition
type EnvironmentFile struct {
	BucketName string `field:"required"` // bucket name
	BucketArn  string `field:"optional"` // bucket arn
}

// ContainerDefinition represents container-definition for the TaskDefinition
type ContainerDefinition struct {
	ContainerName            string            `field:"required"` // name of the container when instantiated
	Image                    string            `field:"required"` // container image without tag
	RegistryType             RegistryType      `field:"required"` // type of registry for pullin image. Default registry type configured to use DefaultContainerDefinitionRegistry
	ImageTag                 string            `field:"required"` // container image's tag
	IsEssential              bool              `field:"required"` // falg representing whether the container should be considered as essential for the TaskDefinition
	Commands                 []string          `field:"optional"` // shell commands to be supplied when instantiating the container
	EntryPointCommands       []string          `field:"optional"` // entrypoint commands to be supplied when instantiating the container
	Cpu                      float64           `field:"required"` // cpu allocation for the container
	Memory                   float64           `field:"required"` // memory allocation for the container.
	PortMappings             []ecs.PortMapping `field:"required"` // port mapping of the container. Dynamic port mapping is used, if Host port is not configured. Useful for auto-scaling of tasks under load-balancer
	EnvironmentFileObjectKey string            `field:"required"` // object key of the environment file present in the S3 bucket
	VolumeMountPoint         []ecs.MountPoint  `field:"optional"` // volume mounts incase of data persistence needed by the container
}

// Volume represents properties for creating a EBS Volume for the TaskDefinition
type Volume struct {
	Name string `field:"required"` // name of the EBS volume to be created for the EC2 service
	Size string `field:"required"` // size of the EBS volume to be crated for the EC2 service
}

// ServiceDiscoveryProps represents properties for service discovery using CLoudMap for the EC2 based Service
type ServiceDiscoveryProps struct {
	ServiceName       string                 `field:"required"` // cloudmap namespace service name of the EC2 based ECS Service for service discovery
	ServicePort       float64                `field:"required"` // port configuration of the EC2 based ECS Service for service discovery
	CloudMapNamespace CloudMapNamespaceProps `field:"required"` // cloudmap namespace properties
}

// CloudMapNamespaceProps represents properties for retrieving the CloudMapNamespace
type CloudMapNamespaceProps struct {
	NamespaceName string `field:"required"` // name of the namespace
	NamespaceId   string `field:"required"` // id of the namespace
	NamespaceArn  string `field:"required"` // arn of the namespace
}

// LoadBalancerProps represents the properties for associating the Application Load-Balancer with the EC2 based ECS Service
type LoadBalancerProps struct {
	SecurityGroupId       string            `field:"required"` // id of the load-balancer's security group
	TargetHealthCheckPath string            `field:"required"` // health check path of the target to validate target health.
	ListenerArn           string            `field:"required"` // arn of the HTTPS listener associated with the load balancer
	ListenerRuleProps     ListenerRuleProps `field:"required"` // listener rule properties for handling traffic to multiple targets inside the load-balancer. Service(s) will be registered as individual targets inside a single load balancer
}

// ListenerRuleProps represents the Application Load-Balancer listener rule properties
type ListenerRuleProps struct {
	Priority      float64 `field:"required"` // rule priority
	PathCondition string  `field:"optional"` // path for path based routing, like '/api/*' will be routed to a target
	HostCondition string  `field:"optional"` // host for host based routing, like 'app.example.com' will be routed to a target
}

// loadBalancedEc2Service construct type
type loadBalancedEc2Service struct {
	constructs.Construct
	logGroup   cloudwatchlogs.LogGroup
	ec2Service ecs.Ec2Service
}

// LoadBalancedEc2Service provides implementation for the loadBalancedEc2Service
type LoadBalancedEc2Service interface {
	// LogGroup returns the Log Group created for the Load-Balanced EC2 service
	LogGroup() cloudwatchlogs.LogGroup
	// Service returns the EC2 service created
	Service() ecs.Ec2Service
}

func (s *loadBalancedEc2Service) Service() ecs.Ec2Service {
	return s.ec2Service
}

func (s *loadBalancedEc2Service) LogGroup() cloudwatchlogs.LogGroup {
	return s.logGroup
}

// NewLoadBalancedEc2Service creates a new Load-Balanced ECS EC2 based service.
//
// Internally creates a Log Group for each service created
// and attaches a policy statement to the task role if present or creates a new task role with AWS XRay access if tracing is enabled
func NewLoadBalancedEc2Service(scope constructs.Construct, id *string, props *LoadBalancedEc2ServiceProps) LoadBalancedEc2Service {
	this := constructs.NewConstruct(scope, id)

	var taskPolicyDocument iam.PolicyDocument = nil
	if props.TaskDefinition.TaskPolicy != nil {

		taskPolicyDocument = props.TaskDefinition.TaskPolicy
		if props.IsTracingEnabled {
			taskPolicyDocument.AddStatements(
				createTaskContainerDefaultXrayPolciyStatement(),
			)
		}
	} else {
		if props.IsTracingEnabled {
			taskPolicyDocument = iam.NewPolicyDocument(&iam.PolicyDocumentProps{
				AssignSids: jsii.Bool(true),
				Statements: &[]iam.PolicyStatement{
					createTaskContainerDefaultXrayPolciyStatement(),
				},
			})
		}
	}

	vpc := lookupVpc(this, id, &props.Cluster.Vpc)

	var networkMode ecs.NetworkMode = DefaultTaskDefinitionNetworkMode
	var loadBalancedServiceTargetType elb2.TargetType = elb2.TargetType_IP
	var serviceSecurityGroups []ec2.ISecurityGroup = nil
	if props.TaskDefinition.NetworkMode == TaskDefintionNetworkModeAwsVpc {
		// task definition network mode configuration
		networkMode = ecs.NetworkMode_AWS_VPC
		// load balancer target type configuration
		loadBalancedServiceTargetType = elb2.TargetType_IP
		// service security group creation & configuration
		sg := ec2.NewSecurityGroup(this, jsii.String("ServiceSecurityGroup"), &ec2.SecurityGroupProps{
			Vpc:              vpc,
			AllowAllOutbound: jsii.Bool(true),
		})
		sg.AddIngressRule(ec2.Peer_AnyIpv4(), ec2.Port_Tcp(jsii.Number(props.ServiceDiscovery.ServicePort)), nil, nil)
		serviceSecurityGroups = append(serviceSecurityGroups, sg)

	} else if props.TaskDefinition.NetworkMode == TaskDefintionNetworkModeBridge {
		networkMode = ecs.NetworkMode_BRIDGE
		loadBalancedServiceTargetType = elb2.TargetType_INSTANCE
	}

	var taskRole iam.Role = nil
	if taskPolicyDocument != nil {
		taskRole = iam.NewRole(this, jsii.String("TaskRole"), &iam.RoleProps{
			AssumedBy: iam.NewServicePrincipal(jsii.String("ecs-tasks."+*core.Aws_URL_SUFFIX()), &iam.ServicePrincipalOpts{}),
			InlinePolicies: &map[string]iam.PolicyDocument{
				*jsii.String("DefaultPolicy"): taskPolicyDocument,
			},
		})
	}

	taskDef := ecs.NewEc2TaskDefinition(this, jsii.String("Ec2TaskDefinition"), &ecs.Ec2TaskDefinitionProps{
		Family:      jsii.String(props.TaskDefinition.FamilyName),
		NetworkMode: networkMode,
		ExecutionRole: iam.NewRole(this, jsii.String("ExecutionRole"), &iam.RoleProps{
			AssumedBy: iam.NewServicePrincipal(jsii.String("ecs-tasks."+*core.Aws_URL_SUFFIX()), &iam.ServicePrincipalOpts{}),
			InlinePolicies: &map[string]iam.PolicyDocument{
				*jsii.String("DefaultPolicy"): iam.NewPolicyDocument(
					&iam.PolicyDocumentProps{
						AssignSids: jsii.Bool(true),
						Statements: &[]iam.PolicyStatement{
							iam.NewPolicyStatement(
								&iam.PolicyStatementProps{
									Actions: &[]*string{
										jsii.String("s3:GetBucketLocation"),
									},
									Effect: iam.Effect_ALLOW,
									Resources: &[]*string{
										&props.TaskDefinition.EnvironmentFile.BucketArn,
									},
								},
							),
						},
					},
				),
			},
		}),
		TaskRole: taskRole,
	})

	if props.TaskDefinition.RequiresVolume {
		for _, volume := range props.TaskDefinition.Volumes {
			var vol ecs.Volume = ecs.Volume{
				Name: jsii.String(volume.Name),
				DockerVolumeConfiguration: &ecs.DockerVolumeConfiguration{
					Driver:        jsii.String(DefaultDockerVolumeDriver),
					Scope:         ecs.Scope_SHARED,
					Autoprovision: jsii.Bool(true),
					DriverOpts: &map[string]*string{
						"volumetype": jsii.String(DefaultDockerVolumeType),
						"size":       jsii.String(volume.Size),
					},
				},
			}
			taskDef.AddVolume(&vol)
		}
	}

	// Creates a CloudWatch Log Group for each service with removal policy set to DESTROY
	logGroup := cloudwatchlogs.NewLogGroup(this, jsii.String("LogGroup"), &cloudwatchlogs.LogGroupProps{
		LogGroupName: jsii.String(props.LogGroupName),
		Retention:    DefaultLogRetention,
	})
	logGroup.ApplyRemovalPolicy(core.RemovalPolicy_DESTROY)

	// adds otel container-defintion to task-defintion if tracing is enabled
	var otelContainerDef ecs.ContainerDefinition = nil
	if props.IsTracingEnabled {
		otelContainerDef = ecs.NewContainerDefinition(scope, jsii.String(RandomString(4)+"OtelContainerDefinition"), &ecs.ContainerDefinitionProps{
			TaskDefinition: taskDef,
			ContainerName:  jsii.String("otel-xray"),
			Image:          ecs.ContainerImage_FromRegistry(jsii.String(OtelContainerImage), &ecs.RepositoryImageProps{}),
			Cpu:            jsii.Number(256),
			MemoryLimitMiB: jsii.Number(256),
			Logging:        createAwsLogDriverForContainer(logGroup, "otel"),
			Command: &[]*string{
				jsii.String("--config=/etc/ecs/ecs-default-config.yaml"),
			},
			// Dynamic port mapping is used
			PortMappings: &[]*ecs.PortMapping{
				{
					ContainerPort: jsii.Number(2000),
					Protocol:      ecs.Protocol_UDP,
				},
				{
					ContainerPort: jsii.Number(4317),
					Protocol:      ecs.Protocol_TCP,
				},
				{
					ContainerPort: jsii.Number(8125),
					Protocol:      ecs.Protocol_UDP,
				},
			},
		})
	}

	for index, containerDef := range props.TaskDefinition.ApplicationContainers {
		// updates task definition with statements providing container access to specific environment file in th S3 bucket
		taskDef.AddToExecutionRolePolicy(
			createEnvironmentFileObjectReadOnlyAccessPolicyStatement(
				props.TaskDefinition.EnvironmentFile.BucketArn,
				containerDef.EnvironmentFileObjectKey),
		)
		// creates a container definition and associates it with the task definition
		cd := createContainerDefinition(
			this,
			"Container"+strconv.FormatInt(int64(index), 10),
			containerDef,
			taskDef,
			s3.Bucket_FromBucketName(
				this,
				jsii.String("EnvironmentFileBucket"),
				jsii.String(props.TaskDefinition.EnvironmentFile.BucketName),
			),
			logGroup,
			props.IsTracingEnabled,
		)
		if props.TaskDefinition.RequiresVolume {
			// addes volume mounts for the container definition
			cd.AddMountPoints(convertContainerVolumeMountPoints(containerDef.VolumeMountPoint)...)
		}

		// creates a container link between the actual container and the otel xray container if tracing is enabled
		// Only works in NetworkMode_BRIDGE mode
		if props.IsTracingEnabled && otelContainerDef != nil {
			cd.AddLink(otelContainerDef, jsii.String("otel-xray"))
			cd.AddContainerDependencies(&ecs.ContainerDependency{
				Condition: ecs.ContainerDependencyCondition_START,
				Container: otelContainerDef,
			})
		}
	}

	var capacityProviderStrategies []*ecs.CapacityProviderStrategy
	for _, cps := range props.CapacityProviders {
		capacityProviderStrategy := createCapacityProviderStrategy(cps)
		capacityProviderStrategies = append(capacityProviderStrategies, &capacityProviderStrategy)
	}

	// creates cloudmap options for EC2 based ECS Service if service discovery is enabled
	var cmOpts *ecs.CloudMapOptions = nil
	if props.IsServiceDiscoveryEnabled {
		cmOpts = &ecs.CloudMapOptions{
			DnsTtl:            core.Duration_Minutes(jsii.Number(1)),
			DnsRecordType:     servicediscovery.DnsRecordType_A,
			ContainerPort:     jsii.Number(props.ServiceDiscovery.ServicePort),
			Name:              jsii.String(props.ServiceDiscovery.ServiceName),
			CloudMapNamespace: retrieveCloudMapNamespaceService(this, props.ServiceDiscovery),
		}
	}
	// builds a Ec2ServiceProps
	ec2ServiceProps := ecs.Ec2ServiceProps{
		Cluster: ecs.Cluster_FromClusterAttributes(this, jsii.String("Cluster"), &ecs.ClusterAttributes{
			ClusterName:    jsii.String(props.Cluster.ClusterName),
			Vpc:            vpc,
			SecurityGroups: &props.Cluster.SecurityGroups,
		}),
		CapacityProviderStrategies: &capacityProviderStrategies,
		TaskDefinition:             taskDef,
		DesiredCount:               &props.DesiredTaskCount,
		CircuitBreaker: &ecs.DeploymentCircuitBreaker{
			Rollback: jsii.Bool(true),
		},
		// placement strategy default to memory
		PlacementStrategies: &[]ecs.PlacementStrategy{
			ecs.PlacementStrategy_PackedByMemory(),
		},
		CloudMapOptions: cmOpts,
		// tag(s) propagation
		PropagateTags:        ecs.PropagatedTagSource_SERVICE,
		EnableECSManagedTags: jsii.Bool(true),
	}
	// configures security groups for the EC2Service if task-definition's network mode is TaskDefintionNetworkModeAwsVpc
	if props.TaskDefinition.NetworkMode == TaskDefintionNetworkModeAwsVpc {
		ec2ServiceProps.SecurityGroups = &serviceSecurityGroups
	}

	// creates an EC2Service based on the Ec2ServiceProps
	ec2Service := ecs.NewEc2Service(this, jsii.String("Ec2Service"), &ec2ServiceProps)

	// creates a application load-balancer target group for the service
	ecsServiceTargetGroup := elb2.NewApplicationTargetGroup(this, jsii.String("ApplicationTargetGroup"), &elb2.ApplicationTargetGroupProps{
		HealthCheck: &elb2.HealthCheck{
			Enabled:          jsii.Bool(true),
			HealthyHttpCodes: jsii.String("200"),
			Path:             jsii.String(props.LoadBalancer.TargetHealthCheckPath),
			Interval:         core.Duration_Seconds(jsii.Number(30)),
		},
		TargetType: loadBalancedServiceTargetType,
		Vpc:        vpc,
		Protocol:   elb2.ApplicationProtocol_HTTP,
		Targets: &[]elb2.IApplicationLoadBalancerTarget{
			ec2Service.LoadBalancerTarget(&props.LoadBalancerTargetOptions),
		},
	})

	var elbListenerConditions []elb2.ListenerCondition
	pathCondition := props.LoadBalancer.ListenerRuleProps.PathCondition
	hostCondition := props.LoadBalancer.ListenerRuleProps.HostCondition
	if hostCondition != "" && pathCondition != "" {
		elbListenerConditions = append(elbListenerConditions, elb2.ListenerCondition_HostHeaders(jsii.Strings(props.LoadBalancer.ListenerRuleProps.HostCondition)))
		elbListenerConditions = append(elbListenerConditions, elb2.ListenerCondition_PathPatterns(jsii.Strings(props.LoadBalancer.ListenerRuleProps.PathCondition)))
	} else if hostCondition != "" {
		elbListenerConditions = append(elbListenerConditions, elb2.ListenerCondition_HostHeaders(jsii.Strings(props.LoadBalancer.ListenerRuleProps.HostCondition)))
	} else if pathCondition != "" {
		elbListenerConditions = append(elbListenerConditions, elb2.ListenerCondition_PathPatterns(jsii.Strings(props.LoadBalancer.ListenerRuleProps.PathCondition)))
	}

	// creates a listener rule for the application load-balancer for routing traffic to the EC2Service
	elb2.NewApplicationListenerRule(this, jsii.String("ALBListenerRule"), &elb2.ApplicationListenerRuleProps{
		Priority:   jsii.Number(props.LoadBalancer.ListenerRuleProps.Priority),
		Action:     elb2.ListenerAction_Forward(&[]elb2.IApplicationTargetGroup{ecsServiceTargetGroup}, &elb2.ForwardOptions{}),
		Conditions: &elbListenerConditions,
		Listener: elb2.ApplicationListener_FromApplicationListenerAttributes(this, jsii.String("ALBListener"), &elb2.ApplicationListenerAttributes{
			ListenerArn:   jsii.String(props.LoadBalancer.ListenerArn),
			SecurityGroup: ec2.SecurityGroup_FromSecurityGroupId(this, jsii.String("ALBSecurityGroup"), jsii.String(props.LoadBalancer.SecurityGroupId), &ec2.SecurityGroupImportOptions{}),
		}),
	})

	return &loadBalancedEc2Service{this, logGroup, ec2Service}
}

// createContainerDefinition creates a container-definition and associates it to the taskDef
func createContainerDefinition(scope constructs.Construct, id string, containerDef ContainerDefinition, taskDef ecs.TaskDefinition, taskDefEnvFileBucket s3.IBucket, logGroup cloudwatchlogs.ILogGroup, tracingEnabled bool) ecs.ContainerDefinition {
	cd := ecs.NewContainerDefinition(scope, jsii.String(id), &ecs.ContainerDefinitionProps{
		TaskDefinition: taskDef,
		ContainerName:  &containerDef.ContainerName,
		//		Command:        convertContainerCommands(containerDef.Commands),
		//		EntryPoint:     convertContainerEntryPointCommands(containerDef.EntryPointCommands),
		Essential:      jsii.Bool(containerDef.IsEssential),
		Image:          configureContainerImage(scope, containerDef.RegistryType, containerDef.Image, containerDef.ImageTag),
		Cpu:            jsii.Number(containerDef.Cpu),
		MemoryLimitMiB: jsii.Number(containerDef.Memory),
		EnvironmentFiles: &[]ecs.EnvironmentFile{
			ecs.AssetEnvironmentFile_FromBucket(taskDefEnvFileBucket, jsii.String(containerDef.EnvironmentFileObjectKey), nil),
		},
		Logging:      createAwsLogDriverForContainer(logGroup, containerDef.ContainerName),
		PortMappings: convertContainerPortMappings(containerDef.PortMappings),
	})

	return cd
}

//func convertContainerCommands(cmds []string) *[]*string {
//	var commands []*string
//	for _, cmd := range cmds {
//		commands = append(commands, jsii.String(cmd))
//	}
//	return &commands
//}

//func convertContainerEntryPointCommands(cmds []string) *[]*string {
//	var entryPointCmds []*string
//	for _, cmd := range cmds {
//		entryPointCmds = append(entryPointCmds, jsii.String(cmd))
//	}
//	return &entryPointCmds
//}

// convertContainerPortMappings converts PortMapping(s) to pointers
func convertContainerPortMappings(pm []ecs.PortMapping) *[]*ecs.PortMapping {
	var portMapping []*ecs.PortMapping
	for _, mapping := range pm {
		portMapping = append(portMapping, &mapping)
	}
	return &portMapping

}

// convertContainerVolumeMountPoints converts MountPoint(s) to pointers
func convertContainerVolumeMountPoints(pm []ecs.MountPoint) []*ecs.MountPoint {
	var mountPoints []*ecs.MountPoint
	for _, mount := range pm {
		mountPoints = append(mountPoints, &mount)
	}
	return mountPoints
}

// configureContainerImage creates a ContainerImage based on the registryType
func configureContainerImage(scope constructs.Construct, registryType RegistryType, image string, tag string) ecs.ContainerImage {
	if registryType == ContainerDefinitionRegistryAwsEcr {
		return ecs.ContainerImage_FromEcrRepository(ecr.Repository_FromRepositoryName(scope, jsii.String("EcrRepository"), jsii.String(image)), jsii.String(tag))
	} else {
		return ecs.ContainerImage_FromRegistry(jsii.String(image+":"+tag), &ecs.RepositoryImageProps{})
	}
}

// createAwsLogDriverForContainer creates an AwsLogDriver for the container to handle logs with AWS CloudWatch service
func createAwsLogDriverForContainer(logGroup cloudwatchlogs.ILogGroup, prefix string) ecs.LogDriver {
	logDriver := ecs.AwsLogDriver_AwsLogs(&ecs.AwsLogDriverProps{
		LogGroup:     logGroup,
		StreamPrefix: jsii.String(prefix),
	})
	return logDriver
}

// lookupVpc looks-up for the vpc using th VpcProps
func lookupVpc(scope constructs.Construct, id *string, props *brznetwork.VpcProps) ec2.IVpc {
	vpc := ec2.Vpc_FromLookup(scope, jsii.String("Vpc"), &ec2.VpcLookupOptions{
		VpcId:     jsii.String(props.Id),
		IsDefault: jsii.Bool(props.IsDefault),
	})
	return vpc
}

// createCapacityProviderStrategy creates a CapacityProviderStrategy for the EC2Service with default weight configured to 1
func createCapacityProviderStrategy(name string) ecs.CapacityProviderStrategy {
	capacityProviderStrategy := ecs.CapacityProviderStrategy{
		CapacityProvider: jsii.String(name),
		Weight:           jsii.Number(1),
	}

	return capacityProviderStrategy
}

// createEnvironmentFileObjectReadOnlyAccessPolicyStatement creates a policy statement that will be attached to the EC2Service's execution role policy
//
// adds access only to the particular s3 object mentioned as environment file in the container-definition
func createEnvironmentFileObjectReadOnlyAccessPolicyStatement(bucket string, key string) iam.PolicyStatement {

	policy := iam.NewPolicyStatement(
		&iam.PolicyStatementProps{
			Effect: iam.Effect_ALLOW,
			Actions: &[]*string{
				jsii.String("s3:GetObject"),
			},
			Resources: &[]*string{
				jsii.String(bucket + "/" + key),
			},
		},
	)

	return policy
}

// createTaskContainerDefaultXrayPolciyStatement creates a policy statement that will be attached to the EC2Service's task role policy
//
// adds access to all reources in the AWS Xray
func createTaskContainerDefaultXrayPolciyStatement() iam.PolicyStatement {
	policy := iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("xray:GetSamplingRules"),
			jsii.String("xray:GetSamplingStatisticSummaries"),
			jsii.String("xray:GetSamplingTargets"),
			jsii.String("xray:PutTelemetryRecords"),
			jsii.String("xray:PutTraceSegments"),
		},
		Effect: iam.Effect_ALLOW,
		// TODO: update resource section for OTEL policy
		Resources: &[]*string{jsii.String("*")},
	})

	return policy
}

// retrieveCloudMapNamespaceService retrieves the CLoudMap Namespace from ServiceDiscoveryProps
func retrieveCloudMapNamespaceService(scope constructs.Construct, sd ServiceDiscoveryProps) servicediscovery.IPrivateDnsNamespace {
	privateNamespace := servicediscovery.PrivateDnsNamespace_FromPrivateDnsNamespaceAttributes(
		scope, jsii.String("CloudMapNamespace"), &servicediscovery.PrivateDnsNamespaceAttributes{
			NamespaceArn:  jsii.String(sd.CloudMapNamespace.NamespaceArn),
			NamespaceId:   jsii.String(sd.CloudMapNamespace.NamespaceId),
			NamespaceName: jsii.String(sd.CloudMapNamespace.NamespaceName),
		},
	)
	return privateNamespace
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

package containerpatterns

import (
	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	autoscaling "github.com/aws/aws-cdk-go/awscdk/v2/awsautoscaling"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsservicediscovery"
	servicediscovery "github.com/aws/aws-cdk-go/awscdk/v2/awsservicediscovery"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

var (
	// vpc for the Compute construct
	vpc ec2.IVpc
	// clusterSecurityGroups are the Security Group(s) associated with the Auto-Scaling group based Capacity Providers of the Cluster
	clusterSecurityGroups []ec2.ISecurityGroup
	// loadBalancerSecurityGroup is the Security Group associated with the Application Load-Balancer for Services inside the Cluster
	loadBalancerSecurityGroup ec2.ISecurityGroup
	// loadBalancer is the application loadbalancer associated with the Cluster
	loadBalancer elbv2.ApplicationLoadBalancer
	// httpsListener is the application listener associated with the Application Load-Balancer for Services inside the Cluster
	httpsListener elbv2.IApplicationListener
	// cloudmapNamespace is the cloudmap privateDnsNamespace associated with the vpc created for the Cluster
	cloudmapNamespace awsservicediscovery.IPrivateDnsNamespace
)

// ecsCompute represents the compute pattern/construct based on ECS
type ecsCompute struct {
	constructs.Construct
	cluster                   ecs.Cluster
	clusterSecurityGroups     []ec2.ISecurityGroup
	environmentFileBucket     s3.Bucket
	loadbalancer              elbv2.ApplicationLoadBalancer
	loadBalancerSecurityGroup ec2.ISecurityGroup
	asgCapacityProviders      []ecs.AsgCapacityProvider
	cloudmapNamespace         servicediscovery.IPrivateDnsNamespace
	httpsListener             elbv2.IApplicationListener
}

// EcsCompute provides implementations for the ecsCompute
type EcsCompute interface {
	// Cluster returns an ECS Cluster as the compute component.
	Cluster() ecs.Cluster
	// ClusterSecurityGroups returns the Security Groups associated with the ASG Capacity Provider of the ECS Cluster compute constuct
	ClusterSecurityGroups() []ec2.ISecurityGroup
	// EnvironmentFileBucket returns the S3 Bucket that is created for the ECS Cluster compute construct for handling the environment file(s) of the services
	EnvironmentFileBucket() s3.Bucket
	// ApplicationLoadBalancer returns the Application Load-Balancer that is created for the ECS Cluster compute construct for routing the traffic between the services
	ApplicationLoadBalancer() elbv2.ApplicationLoadBalancer
	// AlbHttpsListener returns the :443 listener of the Application Load-Balancer for attaching the path-based listener rules for the underlying services if load-balanced
	AlbHttpsListener() elbv2.IApplicationListener
	// AlbSecurityGroup returns the Security Group of the Application Load-Balancer that is created for the ECS Cluster compute construct
	AlbSecurityGroup() ec2.ISecurityGroup
	// AsgCapacityProvider returns the ASG based Capacity Providers of the Cluster compute construct
	AsgCapacityProviders() []ecs.AsgCapacityProvider
	// CloudMapNamespace returns the CloudMapNamespace created for the Cluster compute construct for handling internal routing i.e private DNS
	CloudMapNamespace() servicediscovery.IPrivateDnsNamespace
}

func (cc *ecsCompute) Cluster() ecs.Cluster {
	return cc.cluster
}

func (cc *ecsCompute) ClusterSecurityGroups() []ec2.ISecurityGroup {
	return cc.clusterSecurityGroups
}

func (cc *ecsCompute) EnvironmentFileBucket() s3.Bucket {
	return cc.environmentFileBucket
}

func (cc *ecsCompute) ApplicationLoadBalancer() elbv2.ApplicationLoadBalancer {
	return cc.loadbalancer
}

func (cc *ecsCompute) AlbSecurityGroup() ec2.ISecurityGroup {
	return cc.loadBalancerSecurityGroup
}

func (cc *ecsCompute) AsgCapacityProviders() []ecs.AsgCapacityProvider {
	return cc.asgCapacityProviders
}

func (cc *ecsCompute) CloudMapNamespace() servicediscovery.IPrivateDnsNamespace {
	return cc.cloudmapNamespace
}

func (cc *ecsCompute) AlbHttpsListener() elbv2.IApplicationListener {
	return cc.httpsListener
}

// LoadBalancerOptions represents the options for configuring the Load-Balancer.
//
//   - vpc is not exported and will be configured internally from the package level variable ecsCompute.vpc
//   - SSL certificate should be created separately and the arn should be passed to configure SSL for the Load-Balancer's :443 listener
type LoadBalancerOptions struct {
	Name                   string   // Name of the Application Load-Balancer
	ListenerCertificateArn string   // ListenerCertificateArn is the ARN of the ACM Certificate for the :443 listener
	vpc                    ec2.IVpc // vpc in which the Application Load-Balancer will be created
}

// CloudmapNamespaceProps represents properties for creating the CloudMapNamespace.
//
//   - vpc is not exported and will be configured internally from the package level variable ecsCompute.vpc
type CloudmapNamespaceProps struct {
	Name        string   // Name of the CloudMapNamespace
	Description string   // Description of the CloudMapNamespace
	vpc         ec2.IVpc // vpc in which the CloudMapNamespace will be created
}

type securityGroupProps struct {
	Name        string
	Description string
	vpc         ec2.IVpc
}

// AsgCapacityProvider represents an ASG Capacity Provider for the Cluster compute construct
type AsgCapacityProvider struct {
	AutoScalingGroup AsgProps
	CapacityProvider AsgCapacityProviderProps
}

// AsgProps represent the properties of the EC2 Auto-Scaling Group tp be created for the Cluster compute
type AsgProps struct {
	Name        string  // Name of the Auto-Scaling Group
	MinCapacity float64 // Minimum capacity of the Auto-Scaling Group, i.e the minimum number of EC2 instance(s) to be present inside the ASG
	MaxCapacity float64 // Maximum capacity of the Auto-Scaling Group, i.e the maximum number of EC2 instance(s) that can be present inside the ASG
	//	DesiredCapacity float64
	SshKeyName    string            // SshKeyName is the name of the SSH key for the EC2 instance(s) inside the ASG
	InstanceClass ec2.InstanceClass // InstanceClass is the type of EC2 instance such as T2, T3, etc.
	InstanceSize  ec2.InstanceSize  // InstanceSize is the size of EC2 instance such as small, micro, large, etc.
	vpc           ec2.IVpc          // vpc in which the EC2 instance(s) will be created from the ASG
}

// AsgCapacityProviderProps represents the properties of the Auto-Scaling Group based Capacity Provider in the Cluster compute construct
type AsgCapacityProviderProps struct {
	Name string // Name of the ASG Capacity provider
}

// A EcsComputeProps represents properties for creating an EcsCompute construct.
type EcsComputeProps struct {
	VpcId                 string                 // VpcId is the id of the vpc in which the Compute will be created
	Cluster               ClusterOptions         // Cluster options of the Compute construct
	AsgCapacityProviders  []AsgCapacityProvider  // AsgCapacityProviders is the Auto-Scaling Group based Capacity Providers for the Cluster compute
	EnvironmentFileBucket BucketOptions          // EnvironmentFileBucket provides options for creating a S3 bucket for handling environment files for the service(s) inside the compute construct
	LoadBalancer          LoadBalancerOptions    // LoadBalancer provides options for creating an Application Load-Balancer for handling the traffic inside the Cluster between the service(s)
	CloudmapNamespace     CloudmapNamespaceProps // CloudMapNamespace represents properties for creating the CloudMapNamespace
}

// ClusterOptions represents options for creating a Cluster inside the compute construct
type ClusterOptions struct {
	Name                             string   // Name of the ECS Cluster
	ContainerInsights                bool     // ContainerInsights flag represents whether ContainerInsights option is enabled for implementing observability
	IsAsgCapacityProviderEnabled     bool     // IsAsgCapacityProviderEnabled flag represents whether Auto-Scaling Group based Capacity Provider(s) is enabled for the Cluster
	IsFargateCapacityProviderEnabled bool     // IsAsgCapacityProviderEnabled flag represents whether Fargate Capacity Provider(s) is enabled for the Cluster
	vpc                              ec2.IVpc // vpc in which the Cluster will be created
}

// BucketOptions represents options for creating a S3 Bucket for handling environment file for the services inside the Cluster compute construct
type BucketOptions struct {
	Name              string // Name of the S3 Bucket
	IsVersioned       bool   // IsVersioned flag represents whether object(s) should be versioned or not inside the S3 Bucket
	AutoDeleteEnabled bool   // IsVersioned flag represents whether object(s) should be versioned or not inside the S3 Bucket
}

// NewEcsCompute creates a new ECS based compute constructfrom EcsComputeProps
func NewEcsCompute(scope constructs.Construct, id *string, props *EcsComputeProps) EcsCompute {

	this := constructs.NewConstruct(scope, id)

	vpc = LookupVpc(scope, jsii.String("LookUpVpc"), &brznetwork.VpcProps{Id: props.VpcId})

	cluster :=
		createCluster(this, jsii.String("EcsCluster"), &props.Cluster)

	var capacityProviders []ecs.AsgCapacityProvider
	if props.Cluster.IsAsgCapacityProviderEnabled {
		for _, asgCapacityProvider := range props.AsgCapacityProviders {

			autoScalingGroup := createAutoScalingGroup(this,
				jsii.String(asgCapacityProvider.AutoScalingGroup.Name+"AutoscalingGroup"),
				&asgCapacityProvider.AutoScalingGroup, *cluster.ClusterName())

			capacityProvider := createCapacityProvider(this,
				jsii.String(asgCapacityProvider.CapacityProvider.Name+"AsgCapacityProvider"),
				&asgCapacityProvider.CapacityProvider, autoScalingGroup)
			capacityProviders = append(capacityProviders, capacityProvider)

			cluster.AddAsgCapacityProvider(capacityProvider, &ecs.AddAutoScalingGroupCapacityOptions{})
		}
	}
	envFileBucket := s3.NewBucket(this, jsii.String("EnvironmentFileBucket"), &s3.BucketProps{
		BucketName:        jsii.String(props.EnvironmentFileBucket.Name),
		Versioned:         jsii.Bool(props.EnvironmentFileBucket.IsVersioned),
		AutoDeleteObjects: jsii.Bool(props.EnvironmentFileBucket.AutoDeleteEnabled),
		RemovalPolicy:     core.RemovalPolicy_DESTROY,
	})

	if props.LoadBalancer != (LoadBalancerOptions{}) {

		loadBalancer = createLoadBalancer(this, jsii.String("LoadBalanerSetup"), &props.LoadBalancer)

		httpsListener = createHttpsListener(this, jsii.String("AlbHttpsListener"), &props.LoadBalancer, loadBalancer)

		createHttpListener(this, jsii.String("HttpListener"), loadBalancer)
	}

	if props.CloudmapNamespace != (CloudmapNamespaceProps{}) {
		cloudmapNamespace = createCloudMapNamespace(this, jsii.String("CloudMapNamespace"), &props.CloudmapNamespace)
	}

	return &ecsCompute{this, cluster, clusterSecurityGroups, envFileBucket, loadBalancer, loadBalancerSecurityGroup, capacityProviders, cloudmapNamespace, httpsListener}
}

// LookupVpc looks-up for the Vpc using the VpcProps an returns a IVpc
func LookupVpc(scope constructs.Construct, id *string, props *brznetwork.VpcProps) ec2.IVpc {
	vpc := ec2.Vpc_FromLookup(scope, id, &ec2.VpcLookupOptions{
		VpcId: jsii.String(props.Id),
	})
	return vpc
}

// createCluster creates an ECS Cluster from the ClusterOptions
func createCluster(scope constructs.Construct, id *string, props *ClusterOptions) ecs.Cluster {
	if props.IsFargateCapacityProviderEnabled {
		cluster := ecs.NewCluster(scope, id, &ecs.ClusterProps{
			ClusterName:                    jsii.String(props.Name),
			ContainerInsights:              jsii.Bool(props.ContainerInsights),
			EnableFargateCapacityProviders: jsii.Bool(true),
			Vpc:                            vpc,
		})
		return cluster
	} else {
		cluster := ecs.NewCluster(scope, id, &ecs.ClusterProps{
			ClusterName:                    jsii.String(props.Name),
			ContainerInsights:              jsii.Bool(props.ContainerInsights),
			EnableFargateCapacityProviders: jsii.Bool(false),
			Vpc:                            vpc,
		})
		return cluster
	}
}

// createLbSecurityGroup creates a Security Group for the Application Load-Balancer inside the Vpc with default inbound/ingress rules with access for ports 80(HTTP) and 443(HTTPS)
func createLbSecurityGroup(scope constructs.Construct, id *string, props *securityGroupProps, vpc ec2.IVpc) ec2.ISecurityGroup {
	lbSecurityGroup := ec2.NewSecurityGroup(scope, id, &ec2.SecurityGroupProps{
		AllowAllOutbound:  jsii.Bool(true),
		Vpc:               vpc,
		SecurityGroupName: &props.Name,
		Description:       &props.Description,
	})

	lbSecurityGroup.AddIngressRule(
		ec2.Peer_AnyIpv4(),
		ec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Default HTTPS Port"),
		jsii.Bool(false),
	)

	lbSecurityGroup.AddIngressRule(
		ec2.Peer_AnyIpv4(),
		ec2.Port_Tcp(jsii.Number(80)),
		jsii.String("Default HTTP Port"),
		jsii.Bool(false),
	)

	return lbSecurityGroup
}

// createLoadBalancer creates an Application Load-Balancer for handling the service routing and traffic with Security Group from ecsCompute.lbSecurityGroup
func createLoadBalancer(scope constructs.Construct, id *string, props *LoadBalancerOptions) elbv2.ApplicationLoadBalancer {

	loadBalancerSecurityGroup = createLbSecurityGroup(scope, jsii.String(props.Name+"SecurityGroup"), &securityGroupProps{
		Name:        props.Name + "SecurityGroup",
		Description: "Security group for " + props.Name,
	},
		vpc,
	)

	lb := elbv2.NewApplicationLoadBalancer(scope, id, &elbv2.ApplicationLoadBalancerProps{
		LoadBalancerName: jsii.String(props.Name),
		Vpc:              vpc,
		InternetFacing:   jsii.Bool(true),
		VpcSubnets:       &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PUBLIC},
		IdleTimeout:      core.Duration_Seconds(jsii.Number(120)),
		IpAddressType:    elbv2.IpAddressType_IPV4,
		SecurityGroup:    loadBalancerSecurityGroup,
	})
	return lb
}

// createHttpsListener creates a HTTPS Listener in the Application Load-Balancer for the port :443 with default target group from LoadBalancerOptions
func createHttpsListener(scope constructs.Construct, id *string, props *LoadBalancerOptions, lb elbv2.IApplicationLoadBalancer) elbv2.IApplicationListener {
	httpsListener := elbv2.NewApplicationListener(scope, jsii.String("LoadbalancerHttpsListener"), &elbv2.ApplicationListenerProps{
		LoadBalancer: lb,
		Certificates: &[]elbv2.IListenerCertificate{
			elbv2.ListenerCertificate_FromArn(jsii.String(props.ListenerCertificateArn))},
		Protocol: elbv2.ApplicationProtocol_HTTPS,
		Port:     jsii.Number(443),
		DefaultTargetGroups: &[]elbv2.IApplicationTargetGroup{
			elbv2.NewApplicationTargetGroup(
				scope,
				jsii.String("DefaultTargetGroup"),
				&elbv2.ApplicationTargetGroupProps{
					TargetGroupName: jsii.String(props.Name + "DefaultTargetGroup"),
					TargetType:      elbv2.TargetType_INSTANCE,
					Vpc:             vpc,
					Protocol:        elbv2.ApplicationProtocol_HTTP,
					Port:            jsii.Number(8080),
				},
			),
		},
	})
	return httpsListener
}

// createHttpListener creates a HTTP Listener in the Application Load-Balancer for the port :80 with default action forwarding to the HTTPS listener inside the ALB from LoadBalancerOptions
func createHttpListener(scope constructs.Construct, id *string, lb elbv2.IApplicationLoadBalancer) {

	elbv2.NewApplicationListener(scope, jsii.String("LoadbalancerHttpListener"), &elbv2.ApplicationListenerProps{
		Port:         jsii.Number(80),
		LoadBalancer: lb,
		DefaultAction: elbv2.ListenerAction_Redirect(
			&elbv2.RedirectOptions{
				Host:      jsii.String("#{host}"),
				Protocol:  jsii.String("HTTPS"),
				Port:      jsii.String("443"),
				Path:      jsii.String("/#{path}"),
				Query:     jsii.String("#{query}"),
				Permanent: jsii.Bool(true),
			}),
	})
}

// createCloudMapNamespace creates a CloudMapNamespace from CloudmapNamespaceProps
func createCloudMapNamespace(scope constructs.Construct, id *string, props *CloudmapNamespaceProps) servicediscovery.IPrivateDnsNamespace {
	cloudmapNamespace := servicediscovery.NewPrivateDnsNamespace(scope, id, &servicediscovery.PrivateDnsNamespaceProps{
		Name:        jsii.String(props.Name),
		Description: jsii.String(props.Description),
		Vpc:         vpc,
	})
	return cloudmapNamespace
}

// createAsgSecurityGroup creates a Security Group for the Austo-Scaling Group with default outbound/internet access from securityGroupProps
func createAsgSecurityGroup(scope constructs.Construct, id *string, props *securityGroupProps) ec2.SecurityGroup {
	asgSecurityGroup := ec2.NewSecurityGroup(scope, id, &ec2.SecurityGroupProps{
		AllowAllOutbound:  jsii.Bool(true),
		Vpc:               vpc,
		SecurityGroupName: &props.Name,
		Description:       &props.Description,
	})
	return asgSecurityGroup
}

// createAsgPolicyDocument creates an IAM Policy for the Auto-Scaling Group Instance Profile Role for handling EBS volume attachment
func createAsgPolicyDocument() iam.PolicyDocument {
	pd := iam.NewPolicyDocument(&iam.PolicyDocumentProps{
		Statements: &[]iam.PolicyStatement{iam.NewPolicyStatement(&iam.PolicyStatementProps{Effect: iam.Effect_ALLOW,
			Actions: &[]*string{
				jsii.String("ec2:AttachVolume"),
				jsii.String("ec2:CreateVolume"),
				jsii.String("ec2:DeleteVolume"),
				jsii.String("ec2:DescribeAvailabilityZones"),
				jsii.String("ec2:DescribeInstances"),
				jsii.String("ec2:DescribeVolumes"),
				jsii.String("ec2:DescribeVolumeAttribute"),
				jsii.String("ec2:DetachVolume"),
				jsii.String("ec2:DescribeVolumeStatus"),
				jsii.String("ec2:ModifyVolumeAttribute"),
				jsii.String("ec2:DescribeTags"),
				jsii.String("ec2:CreateTags"),
			},
			Resources: &[]*string{jsii.String("*")}})},
	})
	return pd
}

// createAsgRole creates an IAM Instance Profile Role for the Auto-Scaling Group
func createAsgRole(scope constructs.Construct, id *string, props *AsgProps, policyDocument iam.PolicyDocument) iam.IRole {
	role := iam.NewRole(scope, id, &iam.RoleProps{
		Description:    jsii.String("Iam role for autoscaling group " + props.Name),
		InlinePolicies: &map[string]iam.PolicyDocument{"Ec2VolumeAccess": policyDocument},
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
		},
		RoleName:  jsii.String(props.Name + "InstanceProfileRole"),
		AssumedBy: iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &iam.ServicePrincipalOpts{}),
	})
	return role
}

// createAutoScalingGroup creates a Auto-Scaling Group for the cluster
func createAutoScalingGroup(scope constructs.Construct, id *string, props *AsgProps, clusterName string) autoscaling.AutoScalingGroup {
	asgPolicyDocument := createAsgPolicyDocument()

	role := createAsgRole(scope, jsii.String("IamRole"+props.Name), props, asgPolicyDocument)

	asgSecurityGroup := createAsgSecurityGroup(scope, jsii.String(props.Name+"SecurityGroup"), &securityGroupProps{
		Name:        props.Name + "SecurityGroup",
		Description: "SecurityGroup for " + props.Name,
		vpc:         vpc,
	})
	clusterSecurityGroups = append(clusterSecurityGroups, asgSecurityGroup)

	asg := autoscaling.NewAutoScalingGroup(scope, id, &autoscaling.AutoScalingGroupProps{
		AutoScalingGroupName: jsii.String(props.Name),
		MinCapacity:          jsii.Number(props.MinCapacity),
		MaxCapacity:          jsii.Number(props.MaxCapacity),
		InstanceType:         ec2.InstanceType_Of(props.InstanceClass, props.InstanceSize),
		MachineImage:         createMachineImage(),
		SecurityGroup:        asgSecurityGroup,
		UserData:             ec2.UserData_ForLinux(&ec2.LinuxUserDataOptions{Shebang: jsii.String("#!/bin/bash")}),
		VpcSubnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PUBLIC},
		Vpc:                  vpc,
		KeyName:              jsii.String(props.SshKeyName),
		Role:                 role,
	})

	asg.UserData().AddCommands(
		jsii.String("sudo yum -y update"),
		jsii.String("sudo yum -y install wget"),
		jsii.String("sudo touch /etc/ecs/ecs.config"),
		jsii.String("sudo amazon-linux-extras disable docker"),
		jsii.String("sudo amazon-linux-extras install -y ecs"),
		jsii.String("echo \"ECS_CLUSTER="+clusterName+"\" >>  /etc/ecs/ecs.config"),
		jsii.String("echo \"ECS_AWSVPC_BLOCK_IMDS=true\" >> /etc/ecs/ecs.config"),
		jsii.String("sudo systemctl enable --now --no-block ecs.service"),
		jsii.String("sudo systemctl enable amazon-ssm-agent"),
		jsii.String("sudo systemctl start amazon-ssm-agent"),
		jsii.String("docker plugin install rexray/ebs REXRAY_PREEMPT=true EBS_REGION="+*core.Aws_REGION()+" --grant-all-permissions"),
	)
	return asg
}

// createMachineImage creates a MachineImage for the Auto-Scaling Group instance provisioning
func createMachineImage() ec2.IMachineImage {
	image := ec2.NewAmazonLinuxImage(&ec2.AmazonLinuxImageProps{
		CpuType:        ec2.AmazonLinuxCpuType_X86_64,
		Edition:        ec2.AmazonLinuxEdition_STANDARD,
		Generation:     ec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
		Virtualization: ec2.AmazonLinuxVirt_HVM,
		Kernel:         ec2.AmazonLinuxKernel_KERNEL5_X,
	})
	return image
}

// createCapacityProvider creates an ECS Cluster Auto-Scaling Group
// based Capacity Provider from AsgCapacityProviderProps with default configurations like:
//   - asjkdh
//   - enabled managed scaling,
//   - disabled managed termination protection,
//   - 100 percent taget capacity &
//   - access to Instance Role from the Cluster's service containers/tasks
func createCapacityProvider(scope constructs.Construct, id *string, props *AsgCapacityProviderProps, asg autoscaling.IAutoScalingGroup) ecs.AsgCapacityProvider {
	asgCapacityProvider := ecs.NewAsgCapacityProvider(scope, id, &ecs.AsgCapacityProviderProps{
		AutoScalingGroup:                   asg,
		EnableManagedScaling:               jsii.Bool(true),
		EnableManagedTerminationProtection: jsii.Bool(false),
		TargetCapacityPercent:              jsii.Number(100),
		CapacityProviderName:               jsii.String(props.Name),
		CanContainersAccessInstanceRole:    jsii.Bool(true),
	})
	return asgCapacityProvider
}

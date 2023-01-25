package containerpatterns

import (
	//	brznetwork "breezeware-aws-cdk-patterns-samples/network"

	brznetwork "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/network"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	autoscaling "github.com/aws/aws-cdk-go/awscdk/v2/awsautoscaling"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	servicediscovery "github.com/aws/aws-cdk-go/awscdk/v2/awsservicediscovery"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

var (
	vpc                       ec2.IVpc
	clusterSecurityGroups     []ec2.ISecurityGroup
	loadBalancerSecurityGroup ec2.ISecurityGroup
)

type containerCompute struct {
	constructs.Construct
	cluster                   ecs.Cluster
	clusterSecurityGroups     []ec2.ISecurityGroup
	environmentFileBucket     s3.Bucket
	loadbalancer              elbv2.IApplicationLoadBalancer
	loadBalancerSecurityGroup ec2.ISecurityGroup
	asgCapacityProviders      []ecs.AsgCapacityProvider
	cloudmapNamespace         servicediscovery.IPrivateDnsNamespace
	httpsListener             elbv2.IApplicationListener
}

type ContainerCompute interface {
	Cluster() ecs.Cluster
	ClusterSecurityGroups() []ec2.ISecurityGroup
	EnvironmentFileBucket() s3.Bucket
	LoadBalancer() elbv2.IApplicationLoadBalancer
	LoadBalancerSecurityGroup() ec2.ISecurityGroup
	AsgCapacityProviders() []ecs.AsgCapacityProvider
	CloudMapNamespace() servicediscovery.IPrivateDnsNamespace
	HttpsListener() elbv2.IApplicationListener
}

func (cc *containerCompute) Cluster() ecs.Cluster {
	return cc.cluster
}

func (cc *containerCompute) ClusterSecurityGroups() []ec2.ISecurityGroup {
	return cc.clusterSecurityGroups
}

func (cc *containerCompute) EnvironmentFileBucket() s3.Bucket {
	return cc.environmentFileBucket
}

func (cc *containerCompute) LoadBalancer() elbv2.IApplicationLoadBalancer {
	return cc.loadbalancer
}

func (cc *containerCompute) LoadBalancerSecurityGroup() ec2.ISecurityGroup {
	return cc.loadBalancerSecurityGroup
}

func (cc *containerCompute) AsgCapacityProviders() []ecs.AsgCapacityProvider {
	return cc.asgCapacityProviders
}

func (cc *containerCompute) CloudMapNamespace() servicediscovery.IPrivateDnsNamespace {
	return cc.cloudmapNamespace
}

func (cc *containerCompute) HttpsListener() elbv2.IApplicationListener {
	return cc.httpsListener
}

type LoadBalancerOptions struct {
	Name                   string
	ListenerCertificateArn string
	vpc                    ec2.IVpc
}

type CloudmapNamespaceProps struct {
	Name        string
	Description string
	vpc         ec2.IVpc
}

type securityGroupProps struct {
	Name        string
	Description string
	vpc         ec2.IVpc
}

type AsgCapacityProviders struct {
	AutoScalingGroup AsgProps
	CapacityProvider AsgCapacityProviderProps
}

type AsgProps struct {
	Name            string
	MinCapacity     float64
	MaxCapacity     float64
	DesiredCapacity float64
	SshKeyName      string
	InstanceClass   ec2.InstanceClass
	InstanceSize    ec2.InstanceSize
	vpc             ec2.IVpc
}

type AsgCapacityProviderProps struct {
	Name string
}

type EcsComputeProps struct {
	VpcId                 string
	Cluster               ClusterOptions
	AsgCapacityProviders  []AsgCapacityProviders
	EnvironmentFileBucket BucketOptions
	LoadBalancer          LoadBalancerOptions
	CloudmapNamespace     CloudmapNamespaceProps
}

type ClusterOptions struct {
	Name                             string
	ContainerInsights                bool
	IsAsgCapacityProviderEnabled     bool
	IsFargateCapacityProviderEnabled bool
	vpc                              ec2.IVpc
}

type BucketOptions struct {
	Name        string
	IsVersioned bool
}

func NewContainerCompute(scope constructs.Construct, id *string, props *EcsComputeProps) ContainerCompute {

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
		BucketName: jsii.String(props.EnvironmentFileBucket.Name),
		Versioned:  jsii.Bool(props.EnvironmentFileBucket.IsVersioned),
	})

	loadBalancer := createLoadBalancer(this, jsii.String("LoadBalanerSetup"), &props.LoadBalancer)

	httpsListener := createHttpsListener(this, jsii.String("HttpsListener"), &props.LoadBalancer, loadBalancer)

	createHttpListener(this, jsii.String("HttpListener"), loadBalancer)

	cloudmapNamespace := createCloudMapNamespace(this, jsii.String("CloudMapNamespace"), &props.CloudmapNamespace)

	return &containerCompute{this, cluster, clusterSecurityGroups, envFileBucket, loadBalancer, loadBalancerSecurityGroup, capacityProviders, cloudmapNamespace, httpsListener}
}

func LookupVpc(scope constructs.Construct, id *string, props *brznetwork.VpcProps) ec2.IVpc {
	vpc := ec2.Vpc_FromLookup(scope, id, &ec2.VpcLookupOptions{
		VpcId: jsii.String(props.Id),
	})
	return vpc
}

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

func createLoadBalancer(scope constructs.Construct, id *string, props *LoadBalancerOptions) elbv2.IApplicationLoadBalancer {

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

func createCloudMapNamespace(scope constructs.Construct, id *string, props *CloudmapNamespaceProps) servicediscovery.IPrivateDnsNamespace {
	cloudmapNamespace := servicediscovery.NewPrivateDnsNamespace(scope, id, &servicediscovery.PrivateDnsNamespaceProps{
		Name:        jsii.String(props.Name),
		Description: jsii.String(props.Description),
		Vpc:         vpc,
	})
	return cloudmapNamespace
}

func createAsgSecurityGroup(scope constructs.Construct, id *string, props *securityGroupProps) ec2.SecurityGroup {
	asgSecurityGroup := ec2.NewSecurityGroup(scope, id, &ec2.SecurityGroupProps{
		AllowAllOutbound:  jsii.Bool(true),
		Vpc:               vpc,
		SecurityGroupName: &props.Name,
		Description:       &props.Description,
	})
	return asgSecurityGroup
}

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

func createAsgRole(scope constructs.Construct, id *string, props *AsgProps, policyDocument iam.PolicyDocument) iam.IRole {
	role := iam.NewRole(scope, id, &iam.RoleProps{
		Description:    jsii.String("Iam role for autoscaling group " + props.Name),
		InlinePolicies: &map[string]iam.PolicyDocument{"Ec2VolumeAccess": policyDocument},
		RoleName:       jsii.String(props.Name + "InstanceProfileRole"),
		AssumedBy:      iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &iam.ServicePrincipalOpts{}),
	})
	return role
}

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
		jsii.String("docker plugin install rexray/ebs REXRAY_PREEMPT=true EBS_REGION="+*core.Aws_REGION()+" --grant-all-permissions"),
	)
	return asg
}

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

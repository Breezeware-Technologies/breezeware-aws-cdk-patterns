package containerpatterns_test

import (
	"testing"

	containerpatterns "github.com/Breezeware-Technologies/breezeware-aws-cdk-patterns/container_patterns"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/go-cmp/cmp"
)

var (
	app   core.App   = nil
	stack core.Stack = nil
)

func setup() {
	app = core.NewApp(&core.AppProps{
		AnalyticsReporting: jsii.Bool(false),
	})
	stack = core.NewStack(app, jsii.String("NewStacksetup"), &core.StackProps{
		Env: &core.Environment{
			Account: jsii.String("123456789012"),
			Region:  jsii.String("us-east-1"),
		},
	})
}
func TestEcsComputeClusterResource(t *testing.T) {

	//Test Case 1
	t.Run("cluster with container insights enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                true,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: true,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ECS::Cluster"), jsii.Number(1))

		template.HasResourceProperties(jsii.String("AWS::ECS::Cluster"), map[string]interface{}{
			"ClusterName": "test-cluster",
			"ClusterSettings": []map[string]interface{}{{
				"Name":  "containerInsights",
				"Value": "enabled",
			},
			}})
	})

	//Test Case 2
	t.Run("cluster with container insights disabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: true,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ECS::Cluster"), jsii.Number(1))

		template.HasResourceProperties(jsii.String("AWS::ECS::Cluster"), map[string]interface{}{
			"ClusterName": "test-cluster",
			"ClusterSettings": []map[string]interface{}{{
				"Name":  "containerInsights",
				"Value": "disabled",
			},
			}})
	})

	//Test Case 3
	t.Run("cluster with fargate and asg capacity provider enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                true,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: true,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		capacityProviderAssociationsClusterCapture := assertions.NewCapture(nil)
		capacityProviderAssociationsListCpsCapture := assertions.NewCapture(nil)
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), jsii.Number(1))
		template.HasResourceProperties(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), &map[string]interface{}{
			"CapacityProviders":               capacityProviderAssociationsListCpsCapture,
			"Cluster":                         capacityProviderAssociationsClusterCapture,
			"DefaultCapacityProviderStrategy": []string{},
		})
		expectedCluster := &map[string]interface{}{
			"Ref": "EcsComputeEcsClusterFF2AB253",
		}
		expectedCapacityProviders := &[]interface{}{
			"FARGATE",
			"FARGATE_SPOT",
			map[string]interface{}{
				"Ref": "EcsComputeGoLangMicroAsgCapacityProviderAsgCapacityProviderB28553B0",
			},
			map[string]interface{}{
				"Ref": "EcsComputeGoLangSmallAsgCapacityProviderAsgCapacityProvider86BBE533",
			}}

		if !cmp.Equal(capacityProviderAssociationsClusterCapture.AsObject(), expectedCluster) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCluster, capacityProviderAssociationsClusterCapture.AsObject())
		}
		if !cmp.Equal(capacityProviderAssociationsListCpsCapture.AsArray(), expectedCapacityProviders) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCapacityProviders, capacityProviderAssociationsListCpsCapture.AsArray())
		}
	})

	//Test Case 4
	t.Run("cluster with fargate capacity provider enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                true,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: true,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		capacityProviderAssociationsClusterCapture := assertions.NewCapture(nil)
		capacityProviderAssociationsListCpsCapture := assertions.NewCapture(nil)
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), jsii.Number(1))
		template.HasResourceProperties(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), &map[string]interface{}{
			"CapacityProviders":               capacityProviderAssociationsListCpsCapture,
			"Cluster":                         capacityProviderAssociationsClusterCapture,
			"DefaultCapacityProviderStrategy": []string{},
		})
		expectedCluster := &map[string]interface{}{
			"Ref": "EcsComputeEcsClusterFF2AB253",
		}
		expectedCapacityProviders := &[]interface{}{
			"FARGATE",
			"FARGATE_SPOT",
		}

		if !cmp.Equal(capacityProviderAssociationsClusterCapture.AsObject(), expectedCluster) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCluster, capacityProviderAssociationsClusterCapture.AsObject())
		}
		if !cmp.Equal(capacityProviderAssociationsListCpsCapture.AsArray(), expectedCapacityProviders) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCapacityProviders, capacityProviderAssociationsListCpsCapture.AsArray())
		}
	})

	//Test Case 5
	t.Run("cluster with asg capacity provider enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                true,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: false,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		capacityProviderAssociationsClusterCapture := assertions.NewCapture(nil)
		capacityProviderAssociationsListCpsCapture := assertions.NewCapture(nil)
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), jsii.Number(1))
		template.HasResourceProperties(jsii.String("AWS::ECS::ClusterCapacityProviderAssociations"), &map[string]interface{}{
			"CapacityProviders":               capacityProviderAssociationsListCpsCapture,
			"Cluster":                         capacityProviderAssociationsClusterCapture,
			"DefaultCapacityProviderStrategy": []string{},
		})
		expectedCluster := &map[string]interface{}{
			"Ref": "EcsComputeEcsClusterFF2AB253",
		}
		expectedCapacityProviders := &[]interface{}{
			map[string]interface{}{
				"Ref": "EcsComputeGoLangMicroAsgCapacityProviderAsgCapacityProviderB28553B0",
			},
			map[string]interface{}{
				"Ref": "EcsComputeGoLangSmallAsgCapacityProviderAsgCapacityProvider86BBE533",
			}}

		if !cmp.Equal(capacityProviderAssociationsClusterCapture.AsObject(), expectedCluster) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCluster, capacityProviderAssociationsClusterCapture.AsObject())
		}
		if !cmp.Equal(capacityProviderAssociationsListCpsCapture.AsArray(), expectedCapacityProviders) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedCapacityProviders, capacityProviderAssociationsListCpsCapture.AsArray())
		}

		template.ResourceCountIs(jsii.String("AWS::ECS::CapacityProvider"), jsii.Number(2))

		capacityProviderAsgArnCapture := assertions.NewCapture(nil)
		expectedAsgArn := &map[string]interface{}{
			"Ref": "EcsComputeGoLangMicroAsgAutoscalingGroupASG58A27314",
		}
		template.HasResourceProperties(jsii.String("AWS::ECS::CapacityProvider"), &map[string]interface{}{
			"AutoScalingGroupProvider": &map[string]interface{}{
				"AutoScalingGroupArn": capacityProviderAsgArnCapture,
				"ManagedScaling": &map[string]interface{}{
					"Status":         "ENABLED",
					"TargetCapacity": 100,
				},
			},
			"Name": "GoLangMicroAsgCapacityProvider",
		})
		if !cmp.Equal(capacityProviderAsgArnCapture.AsObject(), expectedAsgArn) {
			t.Errorf("\nExpected value: %v\nand\nActual value: %v\nNOT EQUAL", expectedAsgArn, capacityProviderAsgArnCapture.AsObject())
		}
	})

	//Test Case 6
	t.Run("verify autoscaling group resources", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: "vpc-3456789",
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                true,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: false,
			},
			AsgCapacityProviders: []containerpatterns.AsgCapacityProvider{
				{
					AutoScalingGroup: containerpatterns.AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "demo-pair",
					},
					CapacityProvider: containerpatterns.AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})

		template.ResourceCountIs(jsii.String("AWS::AutoScaling::AutoScalingGroup"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("AWS::AutoScaling::LaunchConfiguration"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("AWS::IAM::Role"), jsii.Number(3))
		template.ResourceCountIs(jsii.String("AWS::IAM::InstanceProfile"), jsii.Number(1))
		template.HasResourceProperties(jsii.String("AWS::AutoScaling::AutoScalingGroup"), &map[string]interface{}{
			"AutoScalingGroupName": "GoLangMicroAsg",
			"MaxSize":              "2",
			"MinSize":              "0",
		})
		template.HasResourceProperties(jsii.String("AWS::AutoScaling::LaunchConfiguration"), &map[string]interface{}{
			"InstanceType": "t2.micro",
			"KeyName":      "demo-pair",
		})
		template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), &map[string]interface{}{
			"GroupName": "GoLangMicroAsgSecurityGroup",
		})
		template.HasResourceProperties(jsii.String("AWS::IAM::Role"), &map[string]interface{}{
			"RoleName": "GoLangMicroAsgInstanceProfileRole",
			"Policies": []interface{}{
				map[string]interface{}{
					"PolicyDocument": map[string]interface{}{
						"Version": "2012-10-17",
						"Statement": []interface{}{
							map[string]interface{}{
								"Effect":   "Allow",
								"Resource": "*",
								"Action": []interface{}{
									"ec2:AttachVolume",
									"ec2:CreateVolume",
									"ec2:DeleteVolume",
									"ec2:DescribeAvailabilityZones",
									"ec2:DescribeInstances",
									"ec2:DescribeVolumes",
									"ec2:DescribeVolumeAttribute",
									"ec2:DetachVolume",
									"ec2:DescribeVolumeStatus",
									"ec2:ModifyVolumeAttribute",
									"ec2:DescribeTags",
									"ec2:CreateTags",
								},
							},
						},
					},
					"PolicyName": "Ec2VolumeAccess",
				},
			},
			"ManagedPolicyArns": []interface{}{
				map[string]interface{}{
					"Fn::Join": []interface{}{
						"",
						[]interface{}{
							"arn:",
							map[string]interface{}{
								"Ref": "AWS::Partition",
							},
							":iam::aws:policy/AmazonSSMManagedInstanceCore",
						},
					},
				},
			},
		})

	})
}

func TestEcsComputeLoadBalancerResource(t *testing.T) {
	//Test Case 1
	t.Run("cluster with loadbalancer resource enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::LoadBalancer"), jsii.Number(0))
	})

	//Test Case 2
	t.Run("cluster with loadbalancer resource enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			LoadBalancer: containerpatterns.LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::LoadBalancer"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::TargetGroup"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::Listener"), jsii.Number(2))
		template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(1))

		template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::LoadBalancer"), &map[string]interface{}{
			"Name":   "ClusterAlb",
			"Scheme": "internet-facing",
			"Type":   "application",
		})

		template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::Listener"), &map[string]interface{}{
			"Port":     443,
			"Protocol": "HTTPS",
			"Certificates": []interface{}{
				&map[string]interface{}{
					"CertificateArn": "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
				},
			},
		})

		template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::Listener"), &map[string]interface{}{
			"DefaultActions": []interface{}{
				&map[string]interface{}{
					"RedirectConfig": &map[string]interface{}{
						"Host":       "#{host}",
						"Path":       "/#{path}",
						"Port":       "443",
						"Protocol":   "HTTPS",
						"Query":      "#{query}",
						"StatusCode": "HTTP_301",
					},
				},
			},
		})

		template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::TargetGroup"), &map[string]interface{}{
			"Name":     "ClusterAlbDefaultTargetGroup",
			"Port":     8080,
			"Protocol": "HTTP",
			"TargetGroupAttributes": []interface{}{
				&map[string]interface{}{
					"Key":   "stickiness.enabled",
					"Value": "false",
				},
			},
			"TargetType": "instance",
		})

		template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), &map[string]interface{}{
			"GroupName": "ClusterAlbSecurityGroup",
			"SecurityGroupIngress": []interface{}{
				&map[string]interface{}{
					"CidrIp":      "0.0.0.0/0",
					"Description": "Default HTTPS Port",
					"FromPort":    443,
					"IpProtocol":  "tcp",
					"ToPort":      443,
				},
				&map[string]interface{}{
					"CidrIp":      "0.0.0.0/0",
					"Description": "Default HTTP Port",
					"FromPort":    80,
					"IpProtocol":  "tcp",
					"ToPort":      80,
				},
			},
		})
	})
}

func TestEcsComputeCloudMapResource(t *testing.T) {
	//Test Case 1
	t.Run("cluster with cloud map resource enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
			CloudmapNamespace: containerpatterns.CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ServiceDiscovery::PrivateDnsNamespace"), jsii.Number(1))

		template.HasResourceProperties(jsii.String("AWS::ServiceDiscovery::PrivateDnsNamespace"), &map[string]interface{}{
			"Name": "brz.demo",
		})
	})

	//Test Case 2
	t.Run("cluster with cloud map resource disabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: true,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::ServiceDiscovery::PrivateDnsNamespace"), jsii.Number(0))
	})
}

func TestEcsComputeS3Resource(t *testing.T) {
	//Test Case 1
	t.Run("cluster's env s3 bucket with versioning and auto delete disabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       false,
				AutoDeleteEnabled: false,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::S3::Bucket"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("Custom::S3AutoDeleteObjects"), jsii.Number(0))

		template.HasResource(jsii.String("AWS::S3::Bucket"), &map[string]interface{}{
			"Properties": &map[string]interface{}{
				"BucketName": "demp-s3-bucket-brz",
			},
			"UpdateReplacePolicy": "Delete",
			"DeletionPolicy":      "Delete",
		})
	})

	//Test Case 2
	t.Run("cluster's env s3 bucket with versioning enabled and auto delete disabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       true,
				AutoDeleteEnabled: false,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::S3::Bucket"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("Custom::S3AutoDeleteObjects"), jsii.Number(0))

		template.HasResource(jsii.String("AWS::S3::Bucket"), &map[string]interface{}{
			"Properties": &map[string]interface{}{
				"BucketName": "demp-s3-bucket-brz",
				"VersioningConfiguration": &map[string]interface{}{
					"Status": "Enabled",
				},
			},
			"UpdateReplacePolicy": "Delete",
			"DeletionPolicy":      "Delete",
		})
	})

	//Test Case 3
	t.Run("cluster's env s3 with versioning disabled and auto delete enabled", func(t *testing.T) {
		setup()
		containerpatterns.NewContainerCompute(stack, jsii.String("EcsCompute"), &containerpatterns.EcsComputeProps{
			VpcId: *jsii.String("vpc-535bd136"),
			Cluster: containerpatterns.ClusterOptions{
				Name:                             "test-cluster",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     false,
				IsFargateCapacityProviderEnabled: false,
			},
			EnvironmentFileBucket: containerpatterns.BucketOptions{
				Name:              "demp-s3-bucket-brz",
				IsVersioned:       false,
				AutoDeleteEnabled: true,
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})
		template.ResourceCountIs(jsii.String("AWS::S3::Bucket"), jsii.Number(1))
		template.ResourceCountIs(jsii.String("Custom::S3AutoDeleteObjects"), jsii.Number(1))

		template.HasResource(jsii.String("AWS::S3::Bucket"), &map[string]interface{}{
			"Properties": &map[string]interface{}{
				"BucketName": "demp-s3-bucket-brz",
			},
			"UpdateReplacePolicy": "Delete",
			"DeletionPolicy":      "Delete",
		})
	})
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("1234567890"),
		Region:  jsii.String("us-east-1"),
	}
}

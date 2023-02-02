package containerpatterns

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	core "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

func TestSavibenefitsIacStack(t *testing.T) {

	t.Run("ClusterResourcePresentTest", func(t *testing.T) {
		// GIVEN
		app := awscdk.NewApp(nil)
		// WHEN
		stack := core.NewStack(app, jsii.String("LoadBalancedEc2ServiceStack"), &core.StackProps{
			Env: env(),
		})

		NewContainerCompute(stack, jsii.String("EcsCompute"), &EcsComputeProps{
			VpcId: *jsii.String(""),
			Cluster: ClusterOptions{
				Name:                             "ClusterGoLang",
				ContainerInsights:                false,
				IsAsgCapacityProviderEnabled:     true,
				IsFargateCapacityProviderEnabled: true,
			},
			LoadBalancer: LoadBalancerOptions{
				Name:                   "ClusterAlb",
				ListenerCertificateArn: "arn:aws:acm:us-east-1:305251478828:certificate/3f5f3c4f-5e6c-40de-a588-41cca514bbeb",
			},
			CloudmapNamespace: CloudmapNamespaceProps{
				Name:        "brz.demo",
				Description: "service discovery namespace",
			},
			AsgCapacityProviders: []AsgCapacityProvider{
				{
					AutoScalingGroup: AsgProps{
						Name:          "GoLangMicroAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_MICRO,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: AsgCapacityProviderProps{
						Name: "GoLangMicroAsgCapacityProvider",
					},
				},
				{
					AutoScalingGroup: AsgProps{
						Name:          "GoLangSmallAsg",
						InstanceClass: awsec2.InstanceClass_BURSTABLE2,
						InstanceSize:  awsec2.InstanceSize_SMALL,
						MinCapacity:   0,
						MaxCapacity:   2,
						SshKeyName:    "breezethru-demo-key-pair",
					},
					CapacityProvider: AsgCapacityProviderProps{
						Name: "GoLangSmallAsgCapacityProvider",
					},
				},
			},
		})

		// THEN
		template := assertions.Template_FromStack(stack, &assertions.TemplateParsingOptions{})

		template.HasResourceProperties(jsii.String("AWS::ECS::Cluster"), map[string]interface{}{
			"ClusterName": "ClusterGoLang",
		})

		template.ResourceCountIs(jsii.String("AWS::ECS::Cluster"), jsii.Number(1))
	})

}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("123456789"),
		Region:  jsii.String("us-east-1"),
	}
}

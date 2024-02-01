package awsec2

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Instance struct {
	InstanceID string
	AmiID      string
}

func (i *Instance) String() string {
	return fmt.Sprintf(" -> InstanceID: %s, AmiID: %s\n", i.InstanceID, i.AmiID)
}

type ASGSummary struct {
	AutoScalingGroupName string
	ASGAmiID             string
	Instances            []Instance
}

func (a *ASGSummary) String() string {
	return fmt.Sprintf("AutoScalingGroupName: %s, ASGAmiID: %s, Instances: \n%v", a.AutoScalingGroupName, a.ASGAmiID, a.Instances)
}

// ListAutoScalingGroupsWithSubstring returns a list of AutoScalingGroups that contain the given substring
func ListAutoScalingGroupsWithSubstring(substring string) ([]types.AutoScalingGroup, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := autoscaling.NewFromConfig(cfg)

	var result []types.AutoScalingGroup

	// Use paginator to handle paginated results
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, &autoscaling.DescribeAutoScalingGroupsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, group := range page.AutoScalingGroups {
			if strings.Contains(*group.AutoScalingGroupName, substring) {
				result = append(result, group)
			}
		}
	}

	return result, nil
}

// GetAmiIDForInstance returns the AMI ID for the given instance
func GetAmiIDForInstance(instanceID string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return "", err
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("EC2 instance with ID %s not found", instanceID)
	}

	amiID := *result.Reservations[0].Instances[0].ImageId
	return amiID, nil
}

// ASGSummaries returns a list of ASGSummary objects for the given list of AutoScalingGroups
func ASGSummaries(ASGs []types.AutoScalingGroup) (summaries []ASGSummary, err error) {
	for _, asg := range ASGs {
		var instances []Instance
		for _, instance := range asg.Instances {
			amiID, err := GetAmiIDForInstance(*instance.InstanceId)
			if err != nil {
				return summaries, fmt.Errorf("error getting AMI ID for instance %s: %w", *instance.InstanceId, err)
			}
			instances = append(instances, Instance{
				InstanceID: *instance.InstanceId,
				AmiID:      amiID,
			})
		}
		summaries = append(summaries, ASGSummary{
			AutoScalingGroupName: *asg.AutoScalingGroupName,
			ASGAmiID:             *asg.LaunchTemplate.LaunchTemplateId,
			Instances:            instances,
		})
	}
	return summaries, err
}
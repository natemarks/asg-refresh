package awsec2

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// TrimASGName trims the given string to remove the "asg-" suffix
func TrimASGName(input string) string {
	index := strings.LastIndex(input, "asg-")

	if index != -1 {
		return input[:index]
	}

	return input
}

// Instance represents an EC2 instance
type Instance struct {
	InstanceID string
	AmiID      string
}

// String returns a string representation of the instance
func (i *Instance) String() string {
	return fmt.Sprintf(" -> InstanceID: %s, AmiID: %s\n", i.InstanceID, i.AmiID)
}

// ASGSummary represents a summary of an AutoScalingGroup
type ASGSummary struct {
	AutoScalingGroupName string
	ASGAmiID             string
	Instances            []Instance
}

// String returns a string representation of the ASGSummary
func (a *ASGSummary) String() (result string) {
	for _, instance := range a.Instances {
		if instance.AmiID != a.ASGAmiID {
			result = result + fmt.Sprintf("MISMATCH: %s [ %s ] %s -> %s", TrimASGName(a.AutoScalingGroupName), instance.InstanceID, instance.AmiID, a.ASGAmiID) + "\n"
		} else {
			result = result + fmt.Sprintf("MATCH: %s [ %s ] %s = %s", TrimASGName(a.AutoScalingGroupName), instance.InstanceID, instance.AmiID, a.ASGAmiID) + "\n"
		}
	}
	return result
}

// ListAutoScalingGroupsWithSubstring returns a list of AutoScalingGroups that contain the given substring
func ListAutoScalingGroupsWithSubstring(substring string, log *zerolog.Logger) ([]types.AutoScalingGroup, error) {
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
				log.Debug().Msgf("Found ASG: %s", *group.AutoScalingGroupName)
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

// GetAMIByLaunchConfigurationName returns the AMI ID for the given launch configuration name
func GetAMIByLaunchConfigurationName(launchConfigurationName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := autoscaling.NewFromConfig(cfg)

	input := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []string{launchConfigurationName},
	}

	resp, err := client.DescribeLaunchConfigurations(context.TODO(), input)
	if err != nil {
		return "", err
	}

	if len(resp.LaunchConfigurations) == 0 {
		return "", fmt.Errorf("launch configuration not found: %s", launchConfigurationName)
	}

	amiID := resp.LaunchConfigurations[0].ImageId
	return *amiID, nil
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
		ASGamiID, err := GetAMIByLaunchConfigurationName(*asg.LaunchConfigurationName)
		if err != nil {
			return summaries, fmt.Errorf("error getting AMI ID for launch configuration %s: %w", *asg.LaunchConfigurationName, err)
		}
		summaries = append(summaries, ASGSummary{
			AutoScalingGroupName: *asg.AutoScalingGroupName,
			ASGAmiID:             ASGamiID,
			Instances:            instances,
		})
	}
	return summaries, err
}

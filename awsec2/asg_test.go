package awsec2

import (
	"os"
	"testing"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

func TestListAutoScalingGroupsWithSubstring(t *testing.T) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	type args struct {
		substring string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test ListAutoScalingGroupsWithSubstring",
			args: args{
				substring: "-sandbox-",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ListAutoScalingGroupsWithSubstring(tt.args.substring, &log)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAutoScalingGroupsWithSubstring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestASGSummaries(t *testing.T) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ASGs, err := ListAutoScalingGroupsWithSubstring("-sandbox-", &log)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ASGs []types.AutoScalingGroup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test ASGSummaries",
			args: args{
				ASGs: ASGs,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ASGSummaries(tt.args.ASGs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ASGSummaries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

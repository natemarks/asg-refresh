package awsec2

import (
	"testing"
)

func TestListAutoScalingGroupsWithSubstring(t *testing.T) {
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
			_, err := ListAutoScalingGroupsWithSubstring(tt.args.substring)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAutoScalingGroupsWithSubstring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAmiIDForInstance(t *testing.T) {
	type args struct {
		instanceID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test GetAmiIDForInstance",
			args: args{
				instanceID: "i-0b8a3bd2749a9c238",
			},
			want:    "ami-0b789cb897e3bc975",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAmiIDForInstance(tt.args.instanceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAmiIDForInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAmiIDForInstance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

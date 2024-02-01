package main

import (
	"fmt"

	"github.com/natemarks/asg-refresh/awsec2"
)

func main() {
	//cfg, err := GetConfig()
	//if err != nil {
	//	panic(err)
	//}
	cfg := Config{
		update: true,
		asgID:  "-sandbox-",
		debug:  true,
	}
	log := cfg.GetLogger()
	log.Info().Msgf("config: %+v", cfg)
	ASGs, err := awsec2.ListAutoScalingGroupsWithSubstring(cfg.asgID)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting AutoScalingGroups")
	}
	summaries, err := awsec2.ASGSummaries(ASGs)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting ASG summaries")
	}
	for _, summary := range summaries {
		fmt.Println(summary.String())
	}
}

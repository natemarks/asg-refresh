# asg-refresh
Check the EC2 instances in autoscaling groups to determing if teh AMIID in the instances match their ASG launch configuration.

The command below will report on all autoscaling groups that have the string '-sandbox-' in their name.

```bash
asg-refresh -asg=-sandbox-

# OR
asg-refresh -asg="-sandbox-"

```


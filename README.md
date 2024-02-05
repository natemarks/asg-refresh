# asg-refresh


## Usage: Report Mode

Report on all matching autoscaling groups. Provide the AMI ID configured for the autoscaling group AND for each of the launched instances.

The command below will report on all autoscaling groups that have teh string '-sandbox-' in their name.

```bash
asg-refresh -asgID=-sandbox-
```

## Usage: Update Mode

Run an instance refresh on a single matching autoscaling group and wait for teh refresh to complete. This mode will fail if it matches more than one autoscaling group.

The command below will update the autoscaling group that has the string '-sandbox-' in its name.

```bash
asg-refresh -asgID=-sandbox- -update
```

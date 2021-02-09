# aws-config-compliance-prometheus-exporter
Prometheus Exporter for AWS Config Compliance

## How to run

### Local

```
$ go run main.go
```

### Binary

Get the binary file from [Releases](https://github.com/chaspy/aws-config-compliance-prometheus-exporter/releases) and run it.

### Docker

```
$ docker run chaspy/aws-config-compliance-prometheus-exporter:v0.1.0
```

## Metrics

```
$ curl -s localhost:8080/metrics | grep aws_custom_config_compliance
# HELP aws_custom_config_compliance Number of compliance
# TYPE aws_custom_config_compliance gauge
aws_custom_config_compliance{cap_exceeded="false",compliance="COMPLIANT",config_rule_name="securityhub-efs-encrypted-check-bd414301"} 0
aws_custom_config_compliance{cap_exceeded="false",compliance="INSUFFICIENT_DATA",config_rule_name="securityhub-dms-replication-not-public-1f6729b8"} 0
aws_custom_config_compliance{cap_exceeded="false",compliance="INSUFFICIENT_DATA",config_rule_name="securityhub-ec2-managedinstance-patch-compliance-440fg71a"} 0
aws_custom_config_compliance{cap_exceeded="false",compliance="NON_COMPLIANT",config_rule_name="eip-attached"} 2
aws_custom_config_compliance{cap_exceeded="false",compliance="NON_COMPLIANT",config_rule_name="s3-bukcet-logging-enabled"} 23
```

## IAM Role

The following policy must be attached to the AWS role to be executed.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "configservice:DescribeComplianceByConfigRule",
            ],
            "Resource": "*"
        }
    ]
}
```

## Datadog Autodiscovery

If you use Datadog, you can use [Kubernetes Integration Autodiscovery](https://docs.datadoghq.com/agent/kubernetes/integrations/?tab=kubernetes) feature.

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/configservice"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Compliance struct {
	ConfigRuleName string
	Compliance            string
	CapExceeded bool
	CappedCount int64
}

var (
	//nolint:gochecknoglobals
	compliance = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "aws_custom",
		Subsystem: "config",
		Name:      "compliance",
		Help:      "Number of compliance",
	},
		[]string{"config_rule_name", "compliance", "capped_count", "cap_exceeded"},
	)
)

func main() {
	interval, err := getInterval()
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(compliance)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		// register metrics as background
		for range ticker.C {
			err := snapshot()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func snapshot() error {
	compliance.Reset()

	Compliances, err := getcompliances()
	if err != nil {
		return fmt.Errorf("failed to get compliances: %w", err)
	}

	for _, Compliance := range Compliances {

		fmt.Printf("compliance %v\n",Compliance)

		labels := prometheus.Labels{
			"config_rule_name": Compliance.ConfigRuleName,
			"compliance":       Compliance.Compliance,
			"capped_count":    "hoge",
			"cap_exceeded":       "hoge",
		}
		compliance.With(labels).Set(1)
	}

	return nil
}

func getInterval() (int, error) {
	const defaultAWSAPIIntervalSecond = 300
	AWSAPIInterval := os.Getenv("AWS_API_INTERVAL")
	if len(AWSAPIInterval) == 0 {
		return defaultAWSAPIIntervalSecond, nil
	}

	integerAWSAPIInterval, err := strconv.Atoi(AWSAPIInterval)
	if err != nil {
		return 0, fmt.Errorf("failed to read Datadog Config: %w", err)
	}

	return integerAWSAPIInterval, nil
}

func getcompliances() ([]Compliance, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := configservice.New(sess)
	input := &configservice.DescribeComplianceByConfigRuleInput{}

	result, err := svc.DescribeComplianceByConfigRule(input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe compliance: %w", err)
	}

	Compliances := make([]Compliance, len(result.ComplianceByConfigRules))
	for i, comp := range result.ComplianceByConfigRules {
		Compliances[i] = Compliance{
			ConfigRuleName: *comp.ConfigRuleName,
			Compliance:            *comp.Compliance.ComplianceType,
			//CapExceeded:     *comp.Compliance.ComplianceContributorCount.CapExceeded,
			//CappedCount: *comp.Compliance.ComplianceContributorCount.CappedCount,
		}
	}

	return Compliances, nil
}

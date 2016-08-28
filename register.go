package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	ip := os.Getenv("REGISTER_IP")
	if ip == "" {
		log.Fatal("You must set the REGISTER_IP environment variable")
	}
	domain := os.Getenv("REGISTER_DOMAIN")
	if domain == "" {
		log.Fatal("You must set the REGISTER_DOMAIN environment variable")
	}
	dnsType := os.Getenv("DNS_TYPE")
	if dnsType == "" {
		dnsType = "A"
	}

	// normelize domain
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Global services, region doesn't really matter
	r53 := route53.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})

	// Get hosted zone that matches the largest portion of the input dns For
	// example, if you have a hosted zone for example.com. and foo.example.com.,
	// if the input domain is bob.foo.example.com then it should match the latter
	// We ultimately want the zone id to use in future calls to update this zone
	var zid string
	longestDomainYet := ""
	err := r53.ListHostedZonesPages(&route53.ListHostedZonesInput{}, func(zones *route53.ListHostedZonesOutput, lastPage bool) bool {
		for _, zone := range zones.HostedZones {
			if strings.HasSuffix(domain, *zone.Name) {
				if len(*zone.Name) > len(longestDomainYet) {
					longestDomainYet = *zone.Name
					zid = *zone.Id
				}
			}
		}
		return true
	})
	if err != nil && zid == "" {
		log.Fatalf("Could not get domain: %v", err)
	}
	if err != nil {
		log.Printf("Error while listing domains, but we got a zone id, trying it anyways: %v\n", err)
	}

	_, err = r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zid,
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String(fmt.Sprintf("Update %s by r53-registerer", time.Now())),
			Changes: []*route53.Change{{
				Action: aws.String("UPSERT"),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name: &domain,
					Type: &dnsType,
					TTL:  aws.Int64(0),
					ResourceRecords: []*route53.ResourceRecord{{
						Value: &ip,
					}},
				},
			}},
		},
	})

	if err != nil {
		log.Fatalf("Error updating resource record set: %v\n", err)
	}
	log.Printf("Updated %v to %v", domain, ip)
}

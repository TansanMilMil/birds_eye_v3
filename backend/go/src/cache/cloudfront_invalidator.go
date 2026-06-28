package cache

import (
	"log"

	"github.com/birdseyeapi/birds_eye_v3/go/src/aws"
	"github.com/birdseyeapi/birds_eye_v3/go/src/env"
)

type CloudFrontInvalidator struct{}

func (c *CloudFrontInvalidator) Invalidate() bool {
	err := aws.CreateInvalidation(
		env.GetEnv("BIRDSEYE_AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID", ""),
		[]string{
			"/news/*",
		})
	if err != nil {
		// CloudFront invalidation is optional — log and continue
		log.Printf("CloudFront invalidation skipped (non-fatal): %v", err)
		return false
	}
	log.Printf("CloudFront invalidation created successfully")
	return true
}

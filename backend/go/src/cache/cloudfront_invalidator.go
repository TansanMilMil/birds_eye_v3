package cache

import (
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
		println("Error creating CloudFront invalidation:", err.Error())
	} else {
		println("CloudFront invalidation created successfully")
	}

	return err == nil
}

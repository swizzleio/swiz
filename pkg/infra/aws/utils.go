package aws

import "github.com/aws/aws-sdk-go-v2/service/ec2/types"

// strOrEmpty returns a string or empty string if it's null
func strOrEmpty(str *string) string {
	if nil == str {
		return ""
	}

	return *str
}

// getTagValue returns a tag value or an empty string if it's not found
func getTagValue(key string, tags []types.Tag) string {
	for _, t := range tags {
		if strOrEmpty(t.Key) == key {
			return strOrEmpty(t.Value)
		}
	}

	return ""
}

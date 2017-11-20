package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSGameliftBuild_basic(t *testing.T) {
	var conf gamelift.Build

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGameliftBuildDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayRestAPIConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSGameliftBuildExists("aws_gamelift_build.test", &conf),
					// TODO
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "name", "bar"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "description", ""),
					resource.TestCheckResourceAttrSet("aws_gamelift_build.test", "created_date"),
					resource.TestCheckNoResourceAttr("aws_gamelift_build.test", "binary_media_types"),
				),
			},
			{
				Config: testAccAWSAPIGatewayRestAPIUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSGameliftBuildExists("aws_gamelift_build.test", &conf),
					// TODO
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "name", "test"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "description", "test"),
					resource.TestCheckResourceAttrSet("aws_gamelift_build.test", "created_date"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "binary_media_types.#", "1"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "binary_media_types.0", "application/octet-stream"),
				),
			},
		},
	})
}

func testAccCheckAWSGameliftBuildExists(n string, res *gamelift.Build) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Gamelift Build ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).gameliftconn

		req := &gamelift.DescribeBuildInput{
			BuildId: aws.String(rs.Primary.ID),
		}
		describe, err := conn.DescribeBuild(req)
		if err != nil {
			return err
		}

		if *describe.Id != rs.Primary.ID {
			return fmt.Errorf("Gamelift Build not found")
		}

		*res = *describe.Build

		return nil
	}
}

func testAccCheckAWSGameliftBuildDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).gameliftconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_gamelift_build" {
			continue
		}

		req := gamelift.DescribeBuildInput{
			BuildId: aws.String(d.Id()),
		}
		describe, err := conn.DescribeBuild(&req)
		if err == nil {
			if len(describe.Items) != 0 &&
				*describe.Items[0].Id == rs.Primary.ID {
				return fmt.Errorf("Gamelift Build still exists")
			}
		}

		return err
	}

	return nil
}

const testAccAWSAPIGatewayRestAPIConfig = `
resource "aws_gamelift_build" "test" {
  // TODO
}
`

const testAccAWSAPIGatewayRestAPIUpdateConfig = `
resource "aws_gamelift_build" "test" {
  // TODO
}
`

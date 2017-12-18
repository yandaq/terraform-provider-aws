package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSGameliftBuild_basic(t *testing.T) {
	var conf gamelift.Build

	buildName := fmt.Sprintf("tf_acc_build_%s", acctest.RandString(8))
	bucketName := fmt.Sprintf("tf_acc_bucket_%s", acctest.RandString(8))
	roleName := fmt.Sprintf("tf_acc_role_%s", acctest.RandString(8))
	roleArnRe := regexp.MustCompile(":" + roleName + "$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGameliftBuildDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGameliftBuildBasicConfig(buildName, bucketName, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSGameliftBuildExists("aws_gamelift_build.test", &conf),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "name", buildName),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "operating_system", "AMAZON_LINUX"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.#", "1"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.0.bucket", bucketName),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.0.key", "tf-acc-test"),
					resource.TestMatchResourceAttr("aws_gamelift_build.test", "storage_location.0.role_arn", roleArnRe),
				),
			},
			{
				Config: testAccAWSGameliftBuildBasicUpdateConfig(buildName, bucketName, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSGameliftBuildExists("aws_gamelift_build.test", &conf),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "name", buildName),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "operating_system", "AMAZON_LINUX"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.#", "1"),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.0.bucket", bucketName),
					resource.TestCheckResourceAttr("aws_gamelift_build.test", "storage_location.0.key", "tf-acc-test"),
					resource.TestMatchResourceAttr("aws_gamelift_build.test", "storage_location.0.role_arn", roleArnRe),
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
		out, err := conn.DescribeBuild(req)
		if err != nil {
			return err
		}

		b := out.Build

		if *b.BuildId != rs.Primary.ID {
			return fmt.Errorf("Gamelift Build not found")
		}

		*res = *b

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
			BuildId: aws.String(rs.Primary.ID),
		}
		out, err := conn.DescribeBuild(&req)
		if err == nil {
			if *out.Build.BuildId == rs.Primary.ID {
				return fmt.Errorf("Gamelift Build still exists")
			}
		}

		return err
	}

	return nil
}

func testAccAWSGameliftBuildBasicConfig(buildName, bucketName, roleName string) string {
	return fmt.Sprintf(`
resource "aws_gamelift_build" "test" {
  name = "%s"
  operating_system = "AMAZON_LINUX"
  storage_location {
    bucket = "${aws_s3_bucket.test.bucket}"
    key = "tf-acc-test" // TODO?
    role_arn = "${aws_iam_role.test.arn}"
  }
}

resource "aws_s3_bucket" "test" {
  bucket = "%s"
}

resource "aws_iam_role" "test" {
  name = "%s"
  path = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "gamelift.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
`, buildName, bucketName, roleName)
}

func testAccAWSGameliftBuildBasicUpdateConfig(buildName, bucketName, roleName string) string {
	return fmt.Sprintf(`
resource "aws_gamelift_build" "test" {
  name = "%s"
  operating_system = "AMAZON_LINUX"
  storage_location {
    bucket = "${aws_s3_bucket.test.bucket}"
    key = "tf-acc-test" // TODO?
    role_arn = "${aws_iam_role.test.arn}"
  }
  version = "??" // TODO
}

resource "aws_s3_bucket" "test" {
  bucket = "%s"
}

resource "aws_iam_role" "test" {
  name = "%s"
  path = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
`, buildName, bucketName, roleName)
}

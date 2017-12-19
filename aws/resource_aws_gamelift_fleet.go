package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsGameliftFleet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsGameliftFleetCreate,
		Read:   resourceAwsGameliftFleetRead,
		Update: resourceAwsGameliftFleetUpdate,
		Delete: resourceAwsGameliftFleetDelete,

		Schema: map[string]*schema.Schema{
			"routing_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"build_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ec2_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ec2_inbound_permissions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_paths": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_groups": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"new_game_session_protection_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_vpc_aws_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_creation_limit_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime_configuration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_launch_parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_launch_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsGameliftFleetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).gameliftconn

	rs := expandGameliftIpPermissions(d.Get("routing_strategy").([]interface{}))
	input := gamelift.CreateFleetInput{
		BuildId:         aws.String(d.Get("build_id").(string)),
		EC2InstanceType: aws.String(d.Get("ec2_instance_type").(string)),
		Name:            aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}
	if v, ok := d.GetOk("ec2_inbound_permissions"); ok {
		input.EC2InboundPermissions = expandGameliftIpPermissions(v.([]interface{}))
	}
	if v, ok := d.GetOk("log_paths"); ok {
		input.LogPaths = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("metric_groups"); ok {
		input.MetricGroups = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("new_game_session_protection_policy"); ok {
		input.NewGameSessionProtectionPolicy = aws.String(v.(string))
	}
	if v, ok := d.GetOk("peer_vpc_aws_account_id"); ok {
		input.PeerVpcAwsAccountId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("peer_vpc_id"); ok {
		input.PeerVpcId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("resource_creation_limit_policy"); ok {
		input.ResourceCreationLimitPolicy = expandResourceCreationLimitPolicy(v.([]interface{}))
	}
	if v, ok := d.GetOk("runtime_configuration"); ok {
		input.RuntimeConfiguration = expandRuntimeConfiguration(v.([]interface{}))
	}
	if v, ok := d.GetOk("server_launch_parameters"); ok {
		input.ServerLaunchParameters = aws.String(v.(string))
	}
	if v, ok := d.GetOk("server_launch_path"); ok {
		input.ServerLaunchPath = aws.String(v.(string))
	}

	log.Printf("[INFO] Creating Gamelift Fleet: %s", input)
	out, err := conn.CreateFleet(&input)
	if err != nil {
		return err
	}

	// TODO: waiter

	d.SetId(*out.FleetAttributes.FleetId)

	return resourceAwsGameliftFleetRead(d, meta)
}

func resourceAwsGameliftFleetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).gameliftconn

	log.Printf("[INFO] Describing Gamelift Fleet: %s", d.Id())
	out, err := conn.DescribeFleetAttributes(&gamelift.DescribeFleetAttributesInput{
		FleetIds: aws.StringSlice([]string{d.Id()}),
	})
	if err != nil {
		return err
	}
	attributes := out.FleetAttributes
	if len(attributes) < 1 {
		d.SetId("")
		log.Printf("[WARN] TODO")
		return nil
	}
	if len(attributes) != 1 {
		return fmt.Errorf("Expected exactly 1 Gamelift fleet, found %d under %q",
			len(attributes), d.Id())
	}
	fleet := attributes[0]

	d.Set("build_id", fleet.BuildId)
	d.Set("description", fleet.Description)
	d.Set("fleet_arn", fleet.FleetArn)
	d.Set("log_paths", flattenStringList(fleet.LogPaths))
	d.Set("metric_groups", flattenStringList(fleet.MetricGroups))
	d.Set("name", fleet.Name)
	d.Set("new_game_session_protection_policy", fleet.NewGameSessionProtectionPolicy)
	d.Set("operating_system", fleet.OperatingSystem)
	d.Set("resource_creation_limit_policy", flattenGameliftResourceCreationLimitPolicy(fleet.ResourceCreationLimitPolicy))
	d.Set("server_launch_parameters", fleet.ServerLaunchParameters)
	d.Set("server_launch_path", fleet.ServerLaunchPath)

	return nil
}

func resourceAwsGameliftFleetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).gameliftconn

	log.Printf("[INFO] Updating Gamelift Fleet: %s", d.Id())
	_, err := conn.UpdateFleetAttributes(&gamelift.UpdateFleetAttributesInput{
	Description: aws.String(d.Get("Description").(string)),
	FleetId: aws.String(d.Get("FleetId").(string)),
	MetricGroups: expandStringList(v.([]interface{})),
	Name: aws.String(d.Get("Name").(string)),
	NewGameSessionProtectionPolicy:aws.String(d.Get("NewGameSessionProtectionPolicy").(string)),
	ResourceCreationLimitPolicy: // TODO
		})
	if err != nil {
		return err
	}

	return resourceAwsGameliftFleetRead(d, meta)
}

func resourceAwsGameliftFleetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).gameliftconn

	log.Printf("[INFO] Deleting Gamelift Fleet: %s", d.Id())
	_, err := conn.DeleteFleet(&gamelift.DeleteFleetInput{
		AliasId: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}

	return nil
}

func expandGameliftIpPermissions(cfg []interface{}) []*gamelift.IpPermission {
	strategy := cfg[0].(map[string]interface{})
	return &gamelift.IpPermission{
		FleetId: aws.String(strategy["fleet_id"].(string)),
		Message: aws.String(strategy["message"].(string)),
		Type:    aws.String(strategy["type"].(string)),
	}
}

func flattenGameliftIpPermissions(rs []*gamelift.IpPermission) []interface{} {
	m := make(map[string]interface{}, 0)
	m["fleet_id"] = *rs.FleetId
	m["message"] = *rs.Message
	m["type"] = *rs.Type

	return []interface{}{m}
}

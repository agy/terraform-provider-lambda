package invoke

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

const (
	invocation = "RequestResponse"
	success    = 200
)

func dataSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"fn_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"fn_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"payload": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

			"response": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {
	fnName := d.Get("fn_name").(string)
	fnVersion := d.Get("fn_version").(string)
	payload := d.Get("payload").(map[string]interface{})

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	input := &lambda.InvokeInput{
		// TODO(agy): Implement ClientContext
		// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#InvokeInput
		FunctionName:   aws.String(fnName),
		InvocationType: aws.String(invocation),
		Payload:        jsonPayload,
	}

	if fnVersion != "" {
		input.Qualifier = aws.String(fnVersion)
	}

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	ctx := context.Background()
	svc := lambda.New(sess)

	output, err := svc.InvokeWithContext(ctx, input)
	if err != nil {
		return err
	}

	statusCode := int(*output.StatusCode)
	if statusCode != success {
		return fmt.Errorf("HTTP request error. Response code: %d", output.StatusCode)
	}

	unescaped, err := unescape(output.Payload)
	if err != nil {
		return errors.Wrap(err, "cannot unescape lambda response")
	}

	response := map[string]interface{}{}

	if err := json.Unmarshal([]byte(unescaped), &response); err != nil {
		return errors.Wrap(err, "cannot unmarshal json from lambda response")
	}

	d.Set("response", response)
	d.SetId(time.Now().UTC().String())

	return nil
}

// unescape accepts a slice of bytes and returns the string value with escaped
// characters removed.
func unescape(b []byte) (string, error) {
	unescaped, err := strconv.Unquote(string(b))
	if err != nil {
		return "", err
	}

	return unescaped, nil
}

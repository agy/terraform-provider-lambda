# terraform-provider-lambda

A simple [Terraform provider](https://www.terraform.io/docs/providers/index.html) to invoke a syncronous [AWS Lambda](https://aws.amazon.com/documentation/lambda/) function and
provide the output.

Similar to the [HTTP provider](https://www.terraform.io/docs/providers/http/index.html).

## Example Usage

_Example Lambda Function: hello_

```
import json


# Return the JSON event that we receive
def lambda_handler(event, context):
    e = json.dumps(event, indent=2)
    return e
```

_Example Terraform config:_

```
data "lambda_function_invoke" "hello" {
  fn_name = "hello"

  payload = {
    foo = "bar"
  }
}

output "lambda_output" {
  value = "${data.lambda_function_invoke.hello.*.response}"
}
```

_Example terminal session:_

```
$ terraform apply
data.lambda_function_invoke.hello: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

lambda_output = [
    {
        foo = bar
    }
]
```

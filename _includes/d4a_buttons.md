{% capture aws_button_latest %}
<a href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/stable/Docker.tmpl" data-rel="Stable-2" target="blank" class="aws-deploy">![Docker for AWS](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)</a>
{% endcapture %}
{% capture aws_blue_latest %}
<a class="button darkblue-btn aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/stable/Docker.tmpl" data-rel="Stable-2" target="blank">Deploy Docker for AWS (stable)</a>
{% endcapture %}
{% capture aws_blue_edge %}
<a class="button darkblue-btn aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl" data-rel="Beta-18" target="blank">Deploy Docker for AWS (beta)</a>
{% endcapture %}
{% capture aws_blue_vpc_latest %}
<a class="button darkblue-btn aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/stable/Docker-no-vpc.tmpl" data-rel="Stable-1" target="blank">Deploy Docker for AWS (stable)<br/><small>uses your existing VPC</small></a>
{% endcapture %}
{% capture aws_blue_vpc_edge %}
<a class="button darkblue-btn aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker-no-vpc.tmpl" data-rel="Beta-18" target="blank">Deploy Docker for AWS (beta)<br/><small>uses your existing VPC</small></a>
{% endcapture %}

{% capture azure_blue_latest %}
<a class="button darkblue-btn azure-deploy" href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fstable%2FDocker.tmpl" data-rel="Stable-2" target="blank">Deploy Docker for Azure (stable)</a>
{% endcapture %}
{% capture azure_blue_edge %}
<a class="button darkblue-btn azure-deploy" href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fedge%2FDocker.tmpl" data-rel="Beta-18" target="blank">Deploy Docker for Azure (beta)</a>
{% endcapture %}
{% capture azure_button_latest %}
<a href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fstable%2FDocker.tmpl" data-rel="Stable-2" target="_blank" class="azure-deploy">![Docker for Azure](http://azuredeploy.net/deploybutton.png)</a>
{% endcapture %}


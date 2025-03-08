{% assign d4a_stable = "CE-Stable-1" %}
{% assign d4a_edge = "CE-Edge-1" %}
{% assign d4a_test = "CE-Test-1" %}


{% capture aws_blue_latest %}
<a class="button outline-btn min-hgt aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/stable/Docker.tmpl" data-rel="{{ d4a_stable }}" target="blank">Deploy Docker Community Edition (CE) for AWS (stable)</a>
{% endcapture %}

{% capture aws_blue_vpc_latest %}
<a class="button outline-btn min-hgt aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/stable/Docker-no-vpc.tmpl" data-rel="{{ d4a_stable }}" target="blank">Deploy Docker Community Edition (CE) for AWS (stable)<br/><small>uses your existing VPC</small></a>
{% endcapture %}

{% capture aws_blue_edge %}
<a class="button outline-btn min-hgt aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl" data-rel="{{ d4a_edge }}" target="blank">Deploy Docker Community Edition (CE) for AWS (edge)</a>
{% endcapture %}

{% capture aws_blue_vpc_edge %}
<a class="button outline-btn min-hgt aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker-no-vpc.tmpl" data-rel="{{ d4a_edge }}" target="blank">Deploy Docker Community Edition (CE) for AWS (edge)<br/><small>uses your existing VPC</small></a>
{% endcapture %}

{% capture aws_blue_test %}
<a class="button outline-btn aws-deploy" href="https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=Docker&templateURL=https://editions-us-east-1.s3.amazonaws.com/aws/test/Docker.tmpl" data-rel="{{ d4a_test }}" target="blank">Deploy Docker Community Edition (CE) for AWS (test)</a>
{% endcapture %}

{% capture azure_blue_latest %}
<a class="button outline-btn azure-deploy" href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fstable%2FDocker.tmpl" data-rel="{{ d4a_stable }}" target="blank">Deploy Docker Community Edition (CE) for Azure (stable)</a>
{% endcapture %}

{% capture azure_blue_edge %}
<a class="button outline-btn azure-deploy" href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fedge%2FDocker.tmpl" data-rel="{{ d4a_edge }}" target="blank">Deploy Docker Community Edition (CE) for Azure (edge)</a>
{% endcapture %}

{% capture azure_button_latest %}
<a href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fdownload.docker.com%2Fazure%2Fstable%2FDocker.tmpl" data-rel="Stable-2" target="_blank" class="azure-deploy">![Docker for Azure](http://azuredeploy.net/deploybutton.png)</a>
{% endcapture %}

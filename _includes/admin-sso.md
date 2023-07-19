{% if include.product == "admin" %}
  {% assign product_name="Docker Admin" %}
  {% if include.layer == "company" %}
    {% assign sso_config_link = "[configuring SSO](/admin/company/settings/sso-configuration/)" %}
  {% else %}
    {% assign sso_config_link = "[configuring SSO](/admin/organization/security-settings/sso-configuration/)" %}
  {% endif %}
{% else %}
  {% assign product_name="Docker Hub" %}
  {% assign sso_config_link = "[configuring SSO](/single-sign-on/configure/)" %}
{% endif %}

SSO allows users to authenticate using their identity providers (IdPs) to access Docker. SSO is available for a whole company, and all associated organizations, or an individual organization that has a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

## How it works

When SSO is enabled, users are redirected to your IdP's authentication page to sign in. They cannot authenticate using their Docker login credentials (Docker ID and password). Docker currently supports Service Provider Initiated SSO flow. Your users must sign in to Docker Hub or Docker Desktop to initiate the SSO authentication process.

The following diagram shows how SSO operates and is managed in Docker Hub and Docker Desktop. In addition, it provides information on how to authenticate between your IdP.

![SSO architecture](/single-sign-on/images/sso-architecture.png)

## How to set it up

Before enabling SSO in Docker, administrators must first configure their IdP to work with Docker. Docker provides the Assertion Consumer Service (ACS) URL and the Entity ID. Administrators use this information to establish a connection between their IdP server and Docker Hub.

After establishing the connection between the IdP server and Docker, administrators sign in to {{ product_name }} and complete the SSO enablement process.

When you enable SSO for your company, a first-time user can sign in to Docker Hub using their company's domain email address. They're then added to your company and assigned to the company team in the organization.

Administrators can then choose to enforce SSO login and effortlessly manage SSO connections for their individual company.

## Prerequisites

* You must first notify your company about the new SSO login procedures.
* Verify that your members have Docker Desktop version 4.4.2, or later, installed on their machines.
* If your organization uses the Docker Hub CLI, new org members must [create a Personal Access Token (PAT)](/docker-hub/access-tokens/) to sign in to the CLI. There is a grace period for existing users, which will expire in the near future. Before the grace period ends, your users can sign in from Docker Desktop CLI using their previous credentials until PATs are mandatory.
In addition, you should add all email addresses to your IdP.
* Confirm that all CI/CD pipelines have replaced their passwords with PATs.
* For your service accounts, add your additional domains or enable it in your IdP.

## What's next?

- Start {{ sso_config_link }}
- Explore the [FAQs](/single-sign-on/faqs/)

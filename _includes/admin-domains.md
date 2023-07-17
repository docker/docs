{% if include.product == "admin" %}
  {% assign product_link="[Docker Admin](https://admin.docker.com)" %}
  {% if include.layer == "company" %}
    {% assign domain_navigation="Select your company in the left navigation drop-down menu, and then select **Domain management**." %}
  {% else" %}
    {% assign domain_navigation="Select your organization in the left navigation drop-down menu, and then select **Domain management**." %}
  {% endif %}
{% else %}
  {% assign product_link="[Docker Hub](https://hub.docker.com)" %}
  {% assign domain_navigation="Navigate to the domain settings page for your organization or company.
    - Organization: Select **Organizations**, your organization, **Settings**, and then **Security**.
    - Company: Select **Organizations**, your company, and then **Settings**." %}
{% endif %}



1. Sign in to {{ product_link }}{: target="_blank" rel="noopener" class="_"}.
2. {{ domain_navigation }}
3. Select **Add a domain** and continue with the on-screen instructions to add the TXT Record Value to your domain name system (DNS).

    >**Note**
    >
    > Format your domains without protocol or www information, for example, `yourcompany.example`. This should include all email domains and subdomains users will use to access Docker, for example `yourcompany.example` and `us.yourcompany.example`. Public domains such as `gmail.com`, `outlook.com`, etc. aren’t permitted. Also, the email domain should be set as the primary email.

4. Once you have waited 72 hours for the TXT Record verification, you can then select **Verify** next to the domain you've added, and follow the on-screen instructions.
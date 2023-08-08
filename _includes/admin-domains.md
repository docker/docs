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
3. Select **Add a domain**.
4. Continue with the on-screen instructions to get a verification code for your domain as a **TXT Record Value**.

    >**Note**
    >
    > Format your domains without protocol or www information, for example, `yourcompany.example`. This should include all email domains and subdomains users will use to access Docker, for example `yourcompany.example` and `us.yourcompany.example`. Public domains such as `gmail.com`, `outlook.com`, etc. aren’t permitted.

5. Add your domain verification code as a new TXT record to your Domain Name System (DNS). The steps to do so may vary depending on your DNS provider.

   >**Note**
   >
   > Make sure that the TXT record name that you create on your DNS matches the domain you registered on Docker in Step 4. For example, if you registered the subdomain `us.yourcompany.example`, you need to create a TXT record within the same name/zone `us`. A root domain such as `yourcompany.example` needs a TXT record on the root zone, which is typically denoted with the `@` name for the record.

6. Once you have waited 72 hours for the TXT record verification, you can then select **Verify** next to the domain you've added, and follow the on-screen instructions.
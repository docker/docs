Production Readiness: Docker Hub Front-End (hub-web-v2)
================================

Testing
-------

  * **What is the max traffic load that your service has been tested with?**
  Hub UI has not been load tested.

  * **How has the service been soak-tested?**
  Hub UI has not been soak tested.

  Monitoring
  ----------

  * **How do you monitor?**
  New Relic for server monitoring, BugSnag for JavaScript errors and PagerDuty for alerting.

  * **What’s the link(s) to the dashboard(s)?**
  New Relic: https://rpm.newrelic.com/accounts/532547/applications/8853774
  BugSnag: https://bugsnag.com/docker/hub-prod/errors
  PagerDuty: https://docker.pagerduty.com/services/PKZG21B

  * **Do you use an exception tracking service (e.g. Bugsnag, Sentry, New Relic)?**
  Yes, BugSnag and New Relic.

  * **What’s the health check endpoint? And what checks does that endpoint perform?**
  https://hub.docker.com/_health/

  * **What external services do you depend on? How do you monitor them and handle failures?**
  Hub API Gateway and all downstream Docker Cloud services.
  Google Tag Manager (gtm.js)
  Recurly (recurly.js)
    
  

  * **What’s the link to view the logs?**

  Alerting
  --------

  * **How do you know if your service is down?**
  PagerDuty alerts
  Prometheus alerts

  * **What are the metrics that you alert on?**
  500's from Front End containers

  * **Have you tested tripping on an alert to page somebody?**
  Not manually tested.  But production systems are paging properly.

  * **What’s the link to your on-call schedule?**
  https://docker.pagerduty.com/schedules#P88XAI9

  * **Where is your on-call run-book?**
  https://docker.atlassian.net/wiki/display/DE/Hub+UI+Runbook

  Disaster
  --------

  * **What’s the plan if your persistence layer blows up?**
  Front-end is stateless so this shouldn't be required, but restart Container if unsure.

  * **What’s the plan if any of your external service dependencies blows up?**
  Hub API Gateway or downstream service - find service owner, escalate/alert through PagerDuty, contact service team via Slack.
  Google Tag Manager - disable Google Tag Manager from https://tagmanager.google.com - single signon with docker.com account
  Recurly problem - check status.recurly.com, escalate/alert Billing team through PagerDuty, contact service team via Slack.
  Update status.io describing impact to UI if any.


  Security
  --------

  * **Is the service exposed on the public internet? Does it require TLS?**
  https://hub.docker.com/

  * **How do you store production secrets?**
  Front-End does not store secrets.  JWT is stored in user's browser cookie.

  * **What is your authentication model (both user authentication and service-to-service authentication)?**
  oauth

  * **Do you store any sensitive user data (emails, phone numbers, intellectual property)?**
  JWT in cookie.

  Release process
  ---------------

  * **What’s the link to your docs on how to do a release?**
  https://docker.atlassian.net/wiki/display/DH/Hub+frontend+Deployment+Process

  * **How long does it take to release a code fix to production?**
  4-8 hours

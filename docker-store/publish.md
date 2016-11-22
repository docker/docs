---
description: Submit a product for the Docker Store
keywords: Docker, docker, store, purchase images
title: Submit a product to Docker Store
---

To submit an image to the Docker Store, you must first [apply to
join](https://store.docker.com/publisher/signup) our publisher program. You'll
receive a welcome email when you're accepted into the program.

Once you've been accepted, click the link in your acceptance email, or go to the
[Docker Store](https://store.docker.com) and click **Publish a Product**.

-----------------

## Before you begin

Before you start, there are a few things you should know.

**The Docker Store moderation process**

The Docker Store team validates submitted products to ensure quality and
security, and to make sure your product information is complete and helpful for
potential customers.

To do this, you'll provide your product binaries and some information about the
product (the "product manifest") which will be reviewed by a moderator. If
changes are needed, the moderation team will notify you by email. At that point,
you can make changes and resubmit the product.

**Use private repositories**

The source for your product must be in a **private** repository in either Docker
Cloud or Docker Hub. This allows us to provide feedback to help you produce
excellent products _before_ you make your items available to the public.

**Product tiers**

You can create several different tiers for a single product. For example, you
might have Free, Basic, and Enterprise versions of a single product, each with
their own features, support levels, and subscription pricing.

At minimum, each product tier must be represented by a unique tag within a
repository, however you can also select tags for each tier from several
different repositories or namespaces.


**Save and continue**

We'll ask for a lot of information to display on your product page, and we know
that you may not be have all of it available right away. As you fill out your
product information, you can always save your work and come back to work on it
later, before you submit it.

Spot a typo? You can always edit and resubmit your product information.
Resubmitted product information goes through the same moderation process, but
small changes should take less time to validate.

## Select repositories

Start with private repositories on Docker Cloud or Docker Hub.

Select at least one repository by choosing a user or organization (the
namespace), then select a repository from that account, and then a tag.

Optionally, click **Add another repository** and repeat this process for any
product tiers you plan to offer on the Docker Store. For example you might have
a Free tier, a Basic tier, and an Enterprise tier, each represented by a
different namespace/repo/tag combination.


Make sure you have read the required **Vendor agreement**, and check the box to
indicate your agreement.

Click **Save and Continue**.

## Add company information

Fill out your Publisher Details. If you've already done this, for example if
you've already submitted a product, skip to the next section.

Your company name comes from the [initial sign up
form](https://store.docker.com/publisher/signup) you filled out, however you can
change your details on this screen if needed.

Provide a URL to the logo that represents your company or organization. This
logo must be at least 512x512 pixels.

Provide the URL of your company website.

## Add product information

Next, fill out the Product Details.

Provide a tagline: a short description of your product in 140 characters or
less. This appears in Store search results along with the product icon, so make
it useful.

Provide a URL to the image that will represent the overall product, again at
least 512x512 pixels. Remember that in the Docker Store, this product icon
displays for every product tier, so you may need to make it general.

Select any categories that apply to your image. These categories help customers
find your image when they search the Docker Store.

Add a longer product description. If the tagline is your elevator pitch to get
the customer's attention, the long description is your chance to highlight what
makes your software great. Don't neglect it.

Provide the URL for the product's support pages. This can be as simple as a
troubleshooting section in your product's README file, or a link to your
company's Support site or knowledge base.

Finally, add some screenshots. These should be 1920x1200 pixels or larger, and
should show your product in use.

Click **Save and Continue** to save your changes and go on to the next screen.

## Product tier offerings

For each repository you selected in the first step, you'll be prompted to create
a Product Tier.

**Default tiers**

The "Default" product tier is the one that is selected on your product's Docker
Store listing page until the customer switches to another tier. You can use the
default option to highlight a specific product tier, or to help your customers
when you expect most of them will want a specific tier rather than another one.

**One month trials**

Docker Store allows you to offer a one-month free trial for any of your paid
subscription products. When you select this option, Docker Store begins the
subscription right away, but does not charge the user the monthly fee until the
beginning of their second month. The user can cancel at any time during the
first month trial period and not be charged.

**Free product tiers**

To create a free product tier, enter a monthly price of $0. Free subscriptions
are treated exactly the same as paid subscriptions, except they do not produce
monthly charges or invoices.

### Create product tiers

For each binary you selected in step one, you'll see a section for product tier
information. Choose one tier to make the Default tier. Then fill out the
information for each individual tier.

For each tier, add a tier name and monthly subscription price. For example, you
might have a tier called "Free" for $0, and a tier called and "Enterprise" for
$10 per month. Optionally, you can choose to offer a free one-month trial.

Select the source repository for each tier.

> **Note**: At this time, you can only select one source for each pricing tier. Support for bundled products is coming at a later date.

For each tier, enter a description. This description tells the customer what's
included or different about this product tier.

Paste the link to your software's license agreement in the next field. This
allows the customer to read and review your license agreement before purchasing.

Finally, provide installation instructions for this product tier.

Repeat this process for each tier.

## What's next?

When you submitted the repository information for your product in the first
step, we began the Docker Security Scan process. You'll receive notification of
your scan results in a few days. During that time, we'll also review the product
information you submitted to make sure it meets our quality guidelines. If any
changes are needed, or if security vulnerabilities are discovered, you'll get an
email explaining what needs to be changed.

When your product's image is secure and the product information meets our
quality guidelines, you'll receive an email notification that the product is
ready to publish to the Docker Store.

Once you receive this email you can go to the Docker Store and click "Publish"
to make your product available.

> **Tip**: Docker does not automatically make the approved product available. This means you can time the product's release on the Docker Store with announcements or marketing activity.

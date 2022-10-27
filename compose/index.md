---
description: Introduction and Overview of Compose
keywords: documentation, docs, docker, compose, orchestration, containers
title: Overview
redirect_from:
 - /compose/cli-command/
 - /compose/networking/swarm/
 - /compose/overview/
 - /compose/swarm/
 - /compose/completion/
---

Compose is a tool for defining and running multi-container Docker applications.
With Compose, you use a YAML file to configure your application's services.
Then, with a single command, you create and start all the services
from your configuration.

Compose works in all environments: production, staging, development, testing, as
well as CI workflows. It also has commands for managing the whole lifecycle of your application:

 * Start, stop, and rebuild services
 * View the status of running services
 * Stream the log output of running services
 * Run a one-off command on a service

The key features of Compose that make it effective are:

* [Have multiple isolated environments on a single host](features-uses.md#have-multiple-isolated-environments-on-a-single-host)
* [Preserves volume data when containers are created](features-uses.md#preserves-volume-data-when-containers-are-created)
* [Only recreate containers that have changed](features-uses.md#only-recreate-containers-that-have-changed)
* [Supports variables and moving a composition between environments](features-uses.md#supports-variables-and-moving-a-composition-between-environments)

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/compose/install/"><img src="/assets/images/download.svg" alt="Download and install" width="70" height="70"></a>
             </div>
                 <h2 id="docker-compose"><a href="/compose/install/">Install Compose </a></h2>
                <p>Follow the instructions on how to install Docker Compose.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/compose/gettingstarted/"><img src="/assets/images/explore.svg" alt="Docker Compose" width="70" height="70"></a>
            </div>
                <h2 id="docker-compose"><a href="/compose/gettingstarted/">Try Compose</a></h2>
                <p>Learn the key concepts of Docker Compose whilst building a simple Python web application.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/compose/release-notes/"><img src="/assets/images/note-add.svg" alt="Release notes" width="70" height="70"></a>
            </div>
                <h2 id="docker-compose"><a href="/compose/release-notes/">View the release notes</a></h2>
                <p>Find out about the latest enhancements and bug fixes.</p>
        </div>
    </div>
    </div>
        <!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/compose/features-uses/"><img src="/assets/images/help.svg" alt="FAQs" width="70" height="70"></a>
            </div>
                <h2 id="docker-compose"><a href="/compose/features-uses/">Understand key features of Compose</a></h2>
                <p>Understand its key features and explore common use cases.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
                 <a href="/compose/compose-file/"><img src="/assets/images/all-inbox.svg" alt="Additional resources" width="70" height="70"></a>
          </div>
                <h2 id="docker-compose"><a href="/compose/compose-file/">Explore the Compose file reference</a></h2>
                <p>Find information on defining services, networks, and volumes for a Docker application.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/compose/faq/"><img src="/assets/images/sms.svg" alt="Give feedback" width="70" height="70"></a>
            </div>
                <h2 id="docker-compose"><a href="/compose/faq/">Browse common FAQs</a></h2>
                <p>Explore general FAQs and find out how to give feedback.</p>
        </div>
     </div>
    </div>
</div>




---
title: Overview
description: Overall index for Docker Extensions SDK documentation
keywords: Docker, Extensions, sdk
redirect_from:
  - /desktop/extensions-sdk/dev/overview/
---

Use the resources in this section to create your own Docker Extension.

Extensions are packaged as specially formatted Docker images, which our CLI tool helps to build. At the root of the image filesystem is a `metadata.json` file which describes the content of the extension. It is a fundamental element of a Docker extension.

An extension can contain a UI part and backend parts that run either on the host or in the Desktop virtual machine.
For further details, see the [architecture page](architecture/index.md).

Extensions are distributed through the Docker Hub.
Development of extensions can be done locally without the need to push the extension to Docker Hub. See [Extensions distribution](extensions/DISTRIBUTION.md) for further details.

{% include extensions-form.md %}


<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/desktop/extensions-sdk/process/"><img src="/assets/images/process.svg" alt="Process" width="70" height="70"></a>
                 </div>
                 <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/process/">The build and publish process</a></h2>
                <p> Understand the process for building and publishing an extension.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/extensions-sdk/quickstart/"><img src="/assets/images/explore.svg" alt="Quickstart" width="70" height="70"></a>
            </div>
                <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/quickstart/">Quickstart guide</a></h2>
                <p>Follow the quickstart guide to build a basic Docker Extension quickly.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/desktop/extensions-sdk/design/design-guidelines/"><img src="/assets/images/design.svg" alt="Design quidelines" width="70" height="70"></a>
            </div>
                <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/design/design-guidelines/">View the design guidelines</a></h2>
                <p>Ensure your extension aligns to Docker's design guidelines and principles</p>
        </div>
    </div>
    </div>
    <!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/extensions-sdk/extensions/"><img src="/assets/images/publish.svg" alt="Publish" width="70" height="70"></a>
            </div>
                <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/extensions/">Publish your extension</a></h2>
                <p>Understand how to publish your extension to the Marketplace.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
                 <a href="/desktop/extensions-sdk/dev/kubernetes/"><img src="/assets/images/sync.svg" alt="Kubernetes" width="70" height="70"></a>
          </div>
                <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/dev/kubernetes/">Interacting with Kubernetes</a></h2>
                <p>Find information on how to interact indirectly with a Kubernetes cluster from your Docker extension.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/desktop/extensions-sdk/extensions/multi-arch/"><img src="/assets/images/build-multi-platform.svg" alt="Multi-arch" width="70" height="70"></a>
            </div>
                <h2 id="docker-extensions"><a href="/desktop/extensions-sdk/extensions/multi-arch/">Multi-arch extensions</a></h2>
                <p>Build your extension for multiple architectures.</p>
        </div>
     </div>
    </div>
</div>

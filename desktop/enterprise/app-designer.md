---
title: Application Designer
description: Docker Desktop Enterprise Application Designer
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise, templates, designer
redirect_from:
- /ee/desktop/app-designer/
---

## Overview

The Application Designer helps Docker developers quickly create new
Docker apps using a library of templates. To start the Application
Designer, select the **Design new application** menu entry.

![The Application Designer lets you choose an existing template or create a custom application.](./images/app-design-start.png "Application Designer")

The list of available templates is provided:

![You can tab through the available application templates. A description of each template is provided.](./images/app-design-choose.png "Available templates for application creation")

After selecting a template, you can then customize your application, For
example, if you select **Flask / NGINX / MySQL**, you can then

- select a different version of python or mysql; and

- choose different external ports:

![You can customize your application, which includes specifying database, proxy, and other details.](./images/app-design-custom.png "Customizing your application")

You can then name your application and customize the disk location:

![You can also customize the name and location of your application.](./images/app-design-custom2.png "Naming and specifying a location for your application")

When you select **Assemble**, your application is created.

![When you assemble your application, a status screen is displayed.](./images/app-design-test.png "Assembling your application")

Once assembled, the following screen allows you to run the application. Select **Run application** to pull the images and start the containers:

![When you run your application, the terminal displays output from the application.](./images/app-design-run.png "Running your application")

Use the corresponding buttons to start and stop your application. Select **Open in Finder** on Mac or **Open in Explorer** on Windows to
view application files on disk. Select **Open in Visual Studio Code** to open files with an editor. Note that debug logs from the application are displayed in the lower part of the Application Designer
window.

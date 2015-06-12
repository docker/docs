<!--[metadata]>
+++
title = "Organize our issues"
description = "Organize our issues"
keywords = ["governance, board, members, profiles"]
[menu.main]
parent="smn_registry_ref"
+++
<![end-metadata]-->

# Organize our issues

The Docker projects use GitHub issues to record issues and feature requests that
come in from contributors.  Help us organize our work by triaging. Triage is the
process of reviewing incoming issue tickets, gathering more information about
the issue, and verifying whether or not the issue is valid.

You should triage if you want to discover which Docker features other contributors
think are important. Triage is a great choice if you have an interest
or experience in software product management or project management.


# What kind of issues can I triage?

Docker users and contributors create new issues if they want to:

* report a problem they had with Docker software
* request a new feature 
* ask a question

# How do I triage?

Follow these steps:

1. Sign up for a <a href="https://github.com" target="_blank">Github account</a>.

2. Visit a Docker repository and press the **Watch** button.

This tells GitHub to notify you of new issues. Depending on your settings,
notification go to your GitHub or email inbox. Some of repositories you can watch are:
<table class="tg" >
    <col width="20%">
    <col width="80%">
    <tr>
        <td class="tg-031e"><a href="https://github.com/docker/docker" target="_blank">docker/docker</a></td>
        <td class="tg-031e">Docker the open-source application container engine</td>
    </tr>
    <tr>
        <td class="tg-031e"><a href="https://github.com/docker/machine" target="_blank">docker/machine</a></td>
        <td class="tg-031e">Software for the easy and quick creation of Docker hosts on your computer, on cloud providers, and inside your own data center.</td>
    </tr>
    <tr>
        <td class="tg-031e"><a href="https://github.com/kitematic/kitematic" target="_blank">kitematic/kitematic</a></td>
        <td class="tg-031e">Kitematic is a simple application for managing Docker containers on Mac OS X.</td>
    </tr>
    </tr>
    <tr>
        <td class="tg-031e"><a href="https://github.com/docker/swarm" target="_blank">docker/swarm</a></td>
        <td class="tg-031e">Native clustering for Docker; manage several Docker hosts as a single, virtual host.</td>
    </tr>
    <tr>
        <td class="tg-031e"><a href="https://github.com/docker/compose" target="_blank">docker/compose</a></td>
        <td class="tg-031e">Define and run complex applications using one or many interlinked containers.</td>
    </tr>
</table>


See <a href="https://github.com/docker" target="_blank">the complete list of
Docker repositories</a> on GitHub.
	
3. Choose an issue from the <a
href="https://github.com/docker/docker/issues?q=is%3Aopen+is%3Aissue+-label%
3Akind%2Fproposal+-label%3Akind%2Fenhancement+-label%3Akind%2Fbug+-label%3Akind%
2Fcleanup+-label%3Akind%2Fgraphics+-label%3Akind%2Fwriting+-label%3Akind%
2Fsecurity+-label%3Akind%2Fquestion+-label%3Akind%2Fimprovement+-label%3Akind%
2Ffeature" target="_blank">list of untriaged issues</a>.

4. Follow the <a
href="https://github.com/docker/docker/blob/master/project/ISSUE-TRIAGE.md"
target="_blank">the triage process</a> to triage the issue. 

The triage process asks you to add both a `kind/` and a `exp/` label to each
issue. Because you are not a Docker maintainer, you add these through comments.
Simply add a `+label` keyword to an issue comment: 

![Example](/images/triage-label.png)

For example, the `+exp/beginner` and `+kind/writing` labels would triage an issue as
beginner writing task. For descriptions of valid labels, see the <a
href="https://github.com/docker/docker/blob/master/project/ISSUE-TRIAGE.md">the triage process</a> 

5. Triage another issue.


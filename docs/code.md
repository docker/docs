<!--[metadata]>
+++
title = "Contribute code"
description = "Contribute code"
keywords = ["governance, board, members, profiles"]
[menu.main]
parent="mn_opensource"
weight=4
+++
<![end-metadata]-->

# Contribute code

If you'd like to improve the code of any of Docker's projects, we would love to
have your contributions. All of our projects' code repositories are on GitHub:

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
</td>
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

If you want to contribute to the `docker/docker` repository you should be
familiar with or a invested in learning Go or the Markdown language.  If you
know other languages, investigate our
other repositories&mdash;not all of them run on Go.

# Code contribution workflow

Below is the general workflow for contributing Docker code or documentation.
If you are an experienced open source contributor you may be familiar with this
workflow. If you are new or just need reminders, the steps below link to more
detailed documentation in Docker's project contributors guide.

1. <a href="http://docs.docker.com/project/software-required/"
target="_blank">Get the software</a> you need.

	This explains how to install a couple of tools used in our development
	environment.  What you need (or don't need) might surprise you.

2. <a href="http://docs.docker.com/project/set-up-git/"
target="_blank">Configure Git and fork the repo</a>.

	Your Git configuration can make it easier for you to contribute. 
	Configuration is especially key if are new to contributing or to Docker.

3. <a href="http://docs.docker.com/project/set-up-dev-env/"
target="_blank">Learn to work with the Docker development container</a>.
	
	Docker developers run `docker` in `docker`.  If you are a geek,
	this is a pretty cool experience.
4. <a href="http://docs.docker.com/project/find-an-issue/"
target="_blank">Claim an issue</a> to work on.

	We created a filter listing <a href="http://goo.gl/Hsp2mk" target="_blank">all open
	and unclaimed issues</a> for Docker. 

5. <a
href="http://docs.docker.com/project/work-issue/" target="_blank">Work on the
issue</a>.

	If you change or add code or docs to a project, you should test your changes
	as you work. This page explains <a
	href="http://docs.docker.com/project/test-and-docs/" target="_blank">how to
	test in our development environment</a>.  
	
	Also, remember to always **sign your commits** as you work! To sign your
	commits, include the `-s` flag in your commit like this:
	
		$ git commit -s -m "Add commit with signature example"
	
	If you don't sign <a href="https://twitter.com/gordontheturtle"
	target="_blank">Gordon</a> will get you!

6. <a href="http://docs.docker.com/project/create-pr/" target="_blank">Create a
pull request</a>.

	If you make a change to fix an issue, add reference to the issue in the pull
	request. Here is an example of a perfect pull request with a good description,
	issue reference, and signature in the commit:
	
	![Sign commits and issues](/images/bonus.png)
	
	We have also have checklist that describes [what each pull request
	needs](#what-is-the-pre-pull-request-checklist).
	

7. <a href="http://docs.docker.com/project/review-pr/"
target="_blank">Participate in the pull request review</a> till a successful
merge.


## FAQ and troubleshooting tips for coders

This section contains some frequently asked questions and tips for
troubleshooting problems in your code contribution.

- [How do I set my signature?](#how-do-i-set-my-signature:cb7f612e17aad7eb26c06709ef92a867)
- [How do I track changes from the docker repo upstream?](#how-do-i-track-changes-from-the-docker-repo-upstream:cb7f612e17aad7eb26c06709ef92a867)
- [How do I format my Go code?](#how-do-i-format-my-go-code:cb7f612e17aad7eb26c06709ef92a867)
- [What is the pre-pull request checklist?](#what-is-the-pre-pull-request-checklist:cb7f612e17aad7eb26c06709ef92a867)
- [How should I comment my code?](#how-should-i-comment-my-code:cb7f612e17aad7eb26c06709ef92a867)
- [How do I rebase my feature branch?](#how-do-i-rebase-my-feature-branch:cb7f612e17aad7eb26c06709ef92a867)

### How do I set my signature {#how-do-i-set-my-signature}

1. Change to the root of your `docker-fork` repository.

        $ cd docker-fork

2. Set your `user.name` for the repository.

        $ git config --local user.name "FirstName LastName"

3. Set your `user.email` for the repository.

        $ git config --local user.email "emailname@mycompany.com"
        
### How do I track changes from the docker repo upstream

Set your local repo to track changes upstream, on the `docker` repository. 

1. Change to the root of your Docker repository.

		$ cd docker-fork

2. Add a remote called `upstream` that points to `docker/docker`

 		$ git remote add upstream https://github.com/docker/docker.git



### How do I format my Go code

Run `gofmt -s -w filename.go` on each changed file before committing your changes.
Most editors have plug-ins that do the formatting automatically.

### What is the pre-pull request checklist

* Sync and cleanly rebase on top of Docker's `master`; do not mix multiple branches
  in the pull request.

* Squash your commits into logical units of work using
  `git rebase -i` and `git push -f`. 

* If your code requires a change to tests or documentation, include code,test,
and documentation changes in the same commit as your code; this ensures a
revert would remove all traces of the feature or fix.

* Reference each issue in your pull request description (`#XXXX`). 

### How should I comment my code?

The Go blog wrote about code comments, it is <a href="http://goo.gl/fXCRu"
target="_blank">a single page explanation</a>. A summary follows:

- Comments begin with two forward `//` slashes.
- To document a type, variable, constant, function, or even a package, write a
regular comment directly preceding the elements declaration, with no intervening blank
line. 
- Comments on package declarations should provide general package documentation. 
- For packages that need large amounts of introductory documentation: the
package comment is placed in its own file.
- Subsequent lines of text are considered part of the same paragraph; you must
leave a blank line to separate paragraphs.
-  Indent pre-formatted text relative to the surrounding comment text (see gob's doc.go for an example).
- URLs are converted to HTML links; no special markup is necessary.

### How do I rebase my feature branch?

Always rebase and squash your commits before making a pull request. 

1. Fetch any of the last minute changes from `docker/docker`.

        $ git fetch upstream master

3. Start an interactive rebase.

        $ git rebase -i upstream/master

4. Rebase opens an editor with a list of commits.

			pick 1a79f55 Tweak some of images
			pick 3ce07bb Add a new line 

	If you run into trouble, `git --rebase abort` removes any changes and gets you
back to where you started. 

4. Squash the `pick` keyword with `squash` on all but the first commit.

		pick 1a79f55 Tweak some of images
		squash 3ce07bb Add a new line 

	After closing the file, `git` opens your editor again to edit the commit
	message. 

5. Edit and save your commit message.

		$ git commit -s

 	Make sure your message includes your signature.

8. Push any changes to your fork on GitHub.

        $ git push origin my-keen-feature


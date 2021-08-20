# Having issues with Docker Desktop and Windows Home?

I did too. It took a lot of troubleshooting, and there wasn't anything on the Docker repo mentioning it as far as I found. So I've created a quick little document to summarize.

WSL
-----
First, make sure you have WSL enabled and that you've restarted after that - it requires a restart to actually take place. If you do have it enabled and the problem persists, 
you'll likely have the same problem as me. 

Hyper-V
-----
Docker Destop utilizes something called Hyper-V containers, which is something only available for Windows 10 Pro and above. There is supposed to be a version of Docker Desktop for
Home that circumvents this, but it isn't working for me. There is a workaround, and it is basically installing containers and telling your system that you actually have the right requirments.
https://poweruser.blog/docker-on-windows-10-without-hyper-v-a529897ed1cc

These instructions worked flawlessly for me (although you have to do them every startup). They were difficult to find, that's for sure!

All set!
-----
At this point, my issue got resolved. Hopefully this helps people out there!

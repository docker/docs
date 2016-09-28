from fabric.api import run

def start_project(email="none", user="none", auth="none", beta_password="maybejustnotsomeemptyspaceyea?", sha="latest", new_relic_key="", new_relic_app_name="hub-stage-node"):
        run('docker rm $(docker ps -a -q) > /dev/null 2>&1 || :')
        run('docker rmi $(docker images -q) > /dev/null 2>&1 || :')
	run('cd /home/')
	run('docker login -e %s -u %s -p %s' % (email, user, auth))
	run('docker pull bagel/hub-prod:%s' % sha)
        run('docker pull bagel/haproxy_beta:latest')
	run("docker ps | awk '{if($1 != \"CONTAINER\"){print $1}}' | xargs -r docker kill")
        # We should tag the image with the git commit and deploy that instead of "latest"
	run('docker run -dp 7001:3000 -e ENV=production --restart=on-failure:5 -e HUB_API_BASE_URL=https://hub-beta-stage.docker.com -e REGISTRY_API_BASE_URL=https://hub-beta-stage.docker.com -e NEW_RELIC_LICENSE_KEY=%s -e NEW_RELIC_APP_NAME=%s bagel/hub-prod:%s' % (new_relic_key, new_relic_app_name, sha))
        # HAProxy doesn't change a lot. We should check the image names before killing/rebooting
        run('docker run -dp 80:80 -p 443:443 -e BETA_PASSWORD=%s --restart=on-failure:5 -v /opt/haproxy.pem:/haproxy/keys/hub-beta.docker.com/hub-beta.docker.pem bagel/haproxy_beta:latest' % beta_password)

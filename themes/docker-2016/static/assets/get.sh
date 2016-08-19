
#SITE="http://moby:mobysecure@moby-docker.pantheonsite.io"
SITE="https://docker.com"



wget -O index.html ${SITE}/
wget -O css/app2.css ${SITE}/sites/all/themes/docker/assets/css/app2.css
wget -O css/allcss.css ${SITE}/sites/all/themes/docker/assets/css/allcss.css
wget -O css/p2p.css ${SITE}/sites/all/themes/docker/assets/css/p2p.css
wget -O css/mobile_responsive.css ${SITE}/sites/all/themes/docker/assets/css/mobile_responsive.css
wget -O css/responsive.css ${SITE}/sites/all/themes/docker/assets/css/responsive.css


wget -O js/app.js ${SITE}/sites/all/themes/docker/assets/js/app.js
wget -O js/alljs.js ${SITE}/sites/all/themes/docker/assets/js/alljs.js

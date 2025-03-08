FROM ubuntu:trusty

RUN apt-get update && \
	DEBIAN_FRONTEND=noninteractive apt-get install -yq ruby ruby-dev build-essential git awscli && \
	gem install --no-ri --no-rdoc bundler
WORKDIR /app
ADD Gemfile /app/Gemfile
ADD Gemfile.lock /app/Gemfile.lock
RUN bundle install
ADD . /app

ENV PATH=$PATH:/app AWS_ACCESS_KEY_ID=**setme** AWS_SECRET_ACCESS_KEY=**setme** BUCKET=**setme**

EXPOSE 4567
CMD ["bundle", "exec", "middleman", "server"]

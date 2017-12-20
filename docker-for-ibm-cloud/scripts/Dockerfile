#################################################
# Watson Conversation Demo (Containerized)
#################################################
# Example Dockerfile for an IBM Watson Conversation service deployed to a Docker for IBM Cloud swarm.

# Standard Ubuntu Container
FROM ubuntu:16.04

# Updates and Package Installs
RUN apt-get update
RUN apt-get install --no-install-recommends -y git vim net-tools ca-certificates curl bzip2 jq

# Install Node
RUN curl -sL https://deb.nodesource.com/setup_6.x | bash -
RUN apt-get install --no-install-recommends -y nodejs

# Install and configure
RUN git clone https://github.com/watson-developer-cloud/conversation-simple.git
RUN cd conversation-simple && npm install

RUN echo "export CONVERSATION_USERNAME=\$(cat /run/secrets/watson_conversation | jq \".username\" | sed 's/\"//g')" > /conversation-simple/run.sh && \
    echo "export CONVERSATION_PASSWORD=\$(cat /run/secrets/watson_conversation | jq \".password\" | sed 's/\"//g')" >> /conversation-simple/run.sh && \
    echo "npm start" >> /conversation-simple/run.sh && \
    chmod +x /conversation-simple/run.sh

# Runtime
EXPOSE 3000
WORKDIR /conversation-simple

CMD /conversation-simple/run.sh

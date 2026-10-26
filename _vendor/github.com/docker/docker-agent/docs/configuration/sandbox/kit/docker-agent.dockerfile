FROM docker/docker-agent-sbx-templates:latest
# Cloud uploaded-kit injection does not derive named credential sentinels.
ENV OPENAI_API_KEY=proxy-managed
USER agent
WORKDIR /home/agent/workspace
ENTRYPOINT ["docker-agent", "run", "--sandbox=false", "default"]
CMD []

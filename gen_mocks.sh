#!/bin/sh

mockgen -destination pkg/jobrunner-framework/mock_Closer.go -package jobrunner-framework io Closer
mockgen -destination pkg/jobrunner-framework/worker/mock_Executor.go -package worker -self_package github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker Executor
mockgen -destination pkg/jobrunner-framework/schema/mock_JobrunnerManager.go -package schema -self_package github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema JobrunnerManager

# fix imports https://github.com/golang/mock/issues/30 remove this hack when this issue is closed
ag 'github.com/docker/dhe-deploy/vendor/' --files-with-matches jobrunner-framework | xargs sed -i 's|github.com/docker/dhe-deploy/vendor/||'
gofmt -w -s pkg/jobrunner-framework/schema/mock_JobrunnerManager.go
gofmt -w -s pkg/jobrunner-framework/worker/mock_Executor.go

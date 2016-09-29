using System;

namespace Docker.Core.backend
{
    public interface IBackendInstallFeatureExceptionReplier
    {
        void SetQuitAction(Action action);
        void AskToUser(BackendInstallFeatureException exception, ITaskQueue taskQueue);
    }
}

using System;

namespace Docker.Core.backend
{
    public interface IBackendExceptionReplier
    {
        void SetQuitAction(Action action);
        void Warn(string message);
        void QuitWithMessage(string message, ITaskQueue taskQueue);
    }
}

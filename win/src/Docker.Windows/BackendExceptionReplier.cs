using System;
using Docker.Core;
using Docker.Core.backend;
using Docker.WPF;

namespace Docker
{
    internal class BackendExceptionReplier : IBackendExceptionReplier
    {
        private Action _quitAction;

        public void SetQuitAction(Action action)
        {
            _quitAction = action;
        }

        public void Warn(string message)
        {
            WarnMessageBox.ShowConfirm(message);
        }

        public void QuitWithMessage(string message, ITaskQueue taskQueue)
        {
            if (message.Length > 0)
            {
                QuitMessageBox.ShowOk(message);
            }
            taskQueue.Queue(_quitAction);
        }
    }
}
using System;
using System.Collections.Concurrent;
using System.Threading;
using Docker.Core.backend;

namespace Docker.Core
{
    public interface ITaskQueue
    {
        event TaskQueue.VoidDelegate Started;
        event TaskQueue.BoolDelegate Ended;

        void Queue(Action task);
        void QueueWithWaitMessage(string message, Action task);
        void Shutdown(string message, Action finalTask);
        void Shutdown();
    }

    public class TaskQueue : ITaskQueue
    {
        private readonly INotifications _notifications;
        private readonly IBackendExceptionReplier _backendExceptionReplier;
        private readonly IBackendInstallFeatureExceptionReplier _backendAskExceptionReplier;
        private readonly ISupplier<BlockingCollection<Action>> _taskQueue;
        private readonly Action _poisonPill = () => { };

        private Thread _jobsThread;
        private bool _shuttingDown;
        private long _runningJobs;
        private long LeftJobsCount => Interlocked.Read(ref _runningJobs) + _taskQueue.Get().Count;

        public delegate void VoidDelegate();
        public delegate void BoolDelegate(bool exceptionOccured);

        public event VoidDelegate Started;
        public event BoolDelegate Ended;

        public TaskQueue(INotifications notifications, IBackendExceptionReplier backendExceptionReplier, IBackendInstallFeatureExceptionReplier backendAskExceptionReplier)
        {
            _notifications = notifications;
            _backendExceptionReplier = backendExceptionReplier;
            _backendAskExceptionReplier = backendAskExceptionReplier;

            _taskQueue = Suppliers<BlockingCollection<Action>>.Memoize(() =>
            {
                var taskQueue = new BlockingCollection<Action>();
                _jobsThread = new Thread(() =>
                {
                    foreach (var task in taskQueue.GetConsumingEnumerable())
                    {

                        if (task == _poisonPill) return;

                        try
                        {
                            Started?.Invoke();
                        }
                        catch
                        {
                            // Ignore
                        }

                        Interlocked.Increment(ref _runningJobs);
                        var exceptionOccured = false;
                        try
                        {
                            task();
                        }
                        catch (BackendWarnException ex)
                        {
                            _backendExceptionReplier.Warn(ex.Message);
                        }
                        catch (BackendQuitException ex)
                        {
                            _backendExceptionReplier.QuitWithMessage(ex.Message, this);
                        }
                        catch (BackendInstallFeatureException ex)
                        {
                            _backendAskExceptionReplier.AskToUser(ex, this);
                        }
                        catch (Exception ex)
                        {
                            if (!_shuttingDown)
                            {
                                exceptionOccured = true;
                                _notifications.NotifyError(ex);
                            }
                        }
                        finally
                        {
                            Interlocked.Decrement(ref _runningJobs);
                        }

                        try
                        {
                            Ended?.Invoke(exceptionOccured);
                        }
                        catch
                        {
                            // Ignore
                        }
                    }
                })
                { Name = "Docker Tasks Thread" };
                _jobsThread.SetApartmentState(ApartmentState.STA);
                _jobsThread.Start();

                return taskQueue;
            });

        }

        private bool IsBusy()
        {
            return LeftJobsCount > 0;
        }

        public void Queue(Action task)
        {
            _taskQueue.Get().Add(task);
        }

        public void QueueWithWaitMessage(string message, Action task)
        {
            if (IsBusy())
            {
                _notifications.Notify(message, "Wait for running jobs to complete", null);
            }

            Queue(task);
        }

        public void Shutdown(string message, Action finalTask)
        {
            QueueWithWaitMessage(message, () =>
            {
                _shuttingDown = true;
                finalTask.Invoke();
            });
            Shutdown();
        }

        public void Shutdown()
        {
            Queue(_poisonPill);
            _jobsThread.Join();
        }
    }
}
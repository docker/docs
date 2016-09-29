using System;
using System.Management;

namespace Docker.Core
{
    public class SystemChangesNotifier
    {
        private readonly Logger _logger;

        public delegate void VoidDelegate();
        public event VoidDelegate OnDiskMount;
        public event VoidDelegate OnDiskUnmount;

        public SystemChangesNotifier()
        {
            _logger = new Logger(GetType());
            RegisterDiskEvents();
        }

        private void RegisterWMIQuery(string targetInstance, string eventClassName, ManagementScope scope, EventArrivedEventHandler eventHandler)
        {
            var eventWatcher = new ManagementEventWatcher(scope, new WqlEventQuery
            {
                EventClassName = eventClassName,
                WithinInterval = new TimeSpan(0, 0, 1),
                Condition = $"TargetInstance ISA '{targetInstance}'"
            });

            try
            {
                eventWatcher.EventArrived += eventHandler;
                eventWatcher.Start();
            }
            catch (Exception)
            {
                eventWatcher.Stop();
            }
        }

        private void RegisterDiskEvents()
        {
            try
            {
                var scope = new ManagementScope("root\\CIMV2");
                RegisterWMIQuery("Win32_LogicalDisk", "__InstanceCreationEvent", scope, (_, e) =>
                {
                    if (MustTriggerDiskEvent(e)) OnDiskMount?.Invoke();
                });
                RegisterWMIQuery("Win32_LogicalDisk", "__InstanceDeletionEvent", scope, (_, e) =>
                {
                    if (MustTriggerDiskEvent(e)) OnDiskUnmount?.Invoke();
                });
            }
            catch (Exception ex)
            {
                _logger.Error($"failed to init: {ex.Message}");
            }
        }

        private bool MustTriggerDiskEvent(EventArrivedEventArgs e)
        {
            var baseObject = e.NewEvent.Properties["TargetInstance"];

            return (baseObject.Value as ManagementBaseObject)?.Properties["DeviceID"].Value != null;
        }
    }
}

using System;
using Docker.Core;

namespace Docker.WPF
{
    public interface IActions
    {
        void Start();
        void StopVm();
        void ResetToDefault();
        void ShutdownVm();
        void MigrateVolume(string name, string volumePath);
        void RestartVm(Action<Settings> changes = null);
        void SwitchDaemon();
    }
}
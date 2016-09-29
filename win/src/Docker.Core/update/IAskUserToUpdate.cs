using System;

namespace Docker.Core.Update
{
    public interface IInstallUpdateWindow
    {
        void Open(AvailableUpdate latestUpdate, Action startingUpdateAction);
    }
}
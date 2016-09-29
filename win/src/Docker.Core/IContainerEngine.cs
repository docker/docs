namespace Docker.Core
{
    public enum ContainerEngineStatus
    {
        Starting,
        Started,
        FailedToStart,
        Stopping,
        Stopped,
        FailedToStop,
        Destroying,
        Destroyed,
        FailedToDestroy
    }

    public interface IContainerEngine
    {
        void Start(Settings settings);
        void Stop();
        void Destroy(bool keepVolume);
    }
}

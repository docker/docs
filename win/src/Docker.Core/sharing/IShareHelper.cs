namespace Docker.Core.sharing
{
    public interface IShareHelper
    {
        void UpdateMounts(Settings settings);
        void ResetCredential();
        bool Mount(string drive, Settings settings);
        void Unmount(string drive);
    }
}

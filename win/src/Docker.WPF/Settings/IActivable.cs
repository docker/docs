using Docker.Core;

namespace Docker.WPF
{
    internal interface IActivable
    {
        void Refresh(Settings settings);
    }

    internal interface IAlwaysEnabled
    {
    }
}

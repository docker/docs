namespace Docker.Core
{
    public class DevMode
    {
        public bool On { get; set; }

        public DevMode()
        {
#if DEBUG
            On = true;
#endif
        }
    }
}

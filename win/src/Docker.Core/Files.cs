using System.IO;

namespace Docker.Core
{
    public interface IFiles
    {
        void ForceMove(string source, string destination);
    }

    public class Files : IFiles
    {
        public void ForceMove(string source, string destination)
        {
            File.Delete(destination);
            File.Move(source, destination);
        }
    }
}
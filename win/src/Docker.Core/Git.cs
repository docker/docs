using System.IO;

namespace Docker.Core
{
    public class Git
    {
        private readonly string _sha1File;
        private string _sha1;

        public Git() : this(Paths.InResourcesDir("sha1"))
        {
        }

        internal Git(string sha1File)
        {
            _sha1File = sha1File;
        }

        public string Sha1()
        {
            return _sha1 ?? (_sha1 = Load());
        }

        public string ShortSha1()
        {
            var sha1 = Sha1();
            if (sha1.Length >= 7)
            {
                return sha1.Substring(0, 7);
            }
            return sha1;
        }

        private string Load()
        {
            try
            {
                return File.ReadAllText(_sha1File);
            }
            catch
            {
                return "";
            }
        }
    }
}

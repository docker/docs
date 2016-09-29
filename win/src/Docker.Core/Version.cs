using System;
using System.Reflection;
using System.Text.RegularExpressions;
using static System.StringComparison;

namespace Docker.Core
{
    public interface IVersion
    {
        string ToHumanString();
        string ToHumanStringWithBuildNumber();
        string ToAppVersionName();
        string ToCrossPlatformName();
        int Major();
        int Minor();
        int Build();
    }

    public class Version : IVersion
    {
        private readonly System.Version _version;
        private readonly string _humanVersion;

        public Version()
        {
            _version = Assembly.GetExecutingAssembly().GetName().Version;
            _humanVersion = LoadHumanVersion();
        }

        internal Version(System.Version version, string humanVersion)
        {
            _version = version;
            _humanVersion = humanVersion;
        }

        internal Version(int major, int minor, int build, int revision, string humanVersion)
        {
            _version = new System.Version(major,minor,build,revision);
            _humanVersion = humanVersion;
        }

        public string ToHumanString()
        {
            return _humanVersion;
        }

        public string ToHumanStringWithBuildNumber()
        {
            var buildNumber = Revision() == 0 ? "local" : $"{Revision()}";
            return $"{_humanVersion} (build: {buildNumber})";
        }

        public string ToAppVersionName()
        {
            return ParseAppVersionName(_humanVersion);
        }

        public string ToCrossPlatformName()
        {
            var index = _humanVersion.LastIndexOf("-", Ordinal);
            if (index != -1)
            {
                return _humanVersion.Substring(0, index);
            }
            return $"{ _version.Major}.{ _version.Minor}.{ _version.Build}-{ToAppVersionName()}";
        }

        public int Major()
        {
            return _version.Major;
        }

        public int Minor()
        {
            return _version.Minor;
        }

        public int Build()
        {
            return _version.Build;
        }

        public int Revision()
        {
            return _version.Revision;
        }

        public override string ToString()
        {
            return $"{ _version.Major}.{ _version.Minor}.{ _version.Build}.{ _version.Revision}";
        }

        private string LoadHumanVersion()
        {
            var attribute = Attribute.GetCustomAttribute(Assembly.GetExecutingAssembly(), typeof(AssemblyInformationalVersionAttribute)) as AssemblyInformationalVersionAttribute;
            return attribute == null ? ToString() : attribute.InformationalVersion;
        }

        public static string ParseAppVersionName(string humanString)
        {
            var scheme = new Regex("[0-9]+\\.[0-9]+\\.[0-9]+(?:-rc[0-9]+)?-(beta[0-9]+)(-[a-z]+)?");

            return scheme.Match(humanString).Groups[1].Value;
        }
    }
}

using System.Globalization;
using Microsoft.VisualBasic.Devices;
using Microsoft.Win32;

namespace Docker.Core
{
    public static class Env
    {
        public static int MaxMemory => (int)(new ComputerInfo().TotalPhysicalMemory/(1024*1024));

        public static class Os
        {
            private static string GetStringInWinNtVersionRegistry(string keyName)
            {
                try
                {
                    using (var key = Registry.LocalMachine.OpenSubKey(@"SOFTWARE\Microsoft\Windows NT\CurrentVersion"))
                    {
                        var o = key?.GetValue(keyName);
                        if (o != null)
                        {
                            return o as string;
                        }
                    }
                }
                catch
                {
                    // ignored
                }
                return "";
            }

            public static string Name => GetStringInWinNtVersionRegistry("ProductName");
            public static string Edition => GetStringInWinNtVersionRegistry("EditionID");
            public static string ReleaseId => GetStringInWinNtVersionRegistry("ReleaseId");
            public static string BuildNumber => GetStringInWinNtVersionRegistry("CurrentBuild");
            public static string BuildLabName => GetStringInWinNtVersionRegistry("BuildLabEx");
            public static string Language => CultureInfo.InstalledUICulture.TwoLetterISOLanguageName;
        }
    }
}
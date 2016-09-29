using System;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Threading;

namespace Docker
{
    public static class SingleInstance
    {
        private static class NativeMethods
        {
            internal const int HwndBroadcast = 0xffff;
            internal static readonly int WmShowfirstinstance = RegisterWindowMessage($"WM_SHOWFIRSTINSTANCE|{Assembly.GetExecutingAssembly().GetType().GUID}");

            [DllImport("user32", CharSet = CharSet.Unicode)]
            private static extern int RegisterWindowMessage(string message);

            [DllImport("user32")]
            internal static extern bool PostMessage(IntPtr hwnd, int msg, IntPtr wparam, IntPtr lparam);
        }

        private static Mutex _mutex;

        public static bool Start()
        {
            bool onlyInstance;
            string mutexName = $"Local\\{Assembly.GetExecutingAssembly().GetType().GUID}";

            // if you want your app to be limited to a single instance
            // across ALL SESSIONS (multiple users & terminal services), then use the following line instead:
            // string mutexName = String.Format("Global\\{0}", ProgramInfo.AssemblyGuid);

            _mutex = new Mutex(true, mutexName, out onlyInstance);
            GC.KeepAlive(_mutex);

            if (!onlyInstance)
            {
                ShowFirstInstance();
            }

            return onlyInstance;
        }

        private static void ShowFirstInstance()
        {
            NativeMethods.PostMessage(
                (IntPtr) NativeMethods.HwndBroadcast,
                NativeMethods.WmShowfirstinstance,
                IntPtr.Zero,
                IntPtr.Zero);
        }

        public static void Release()
        {
            try
            {
                _mutex?.ReleaseMutex();
            }
            catch
            {
                // ignored
            }
        }
    }
}
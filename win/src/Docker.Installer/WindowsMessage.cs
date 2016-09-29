using System;
using System.Runtime.InteropServices;

namespace Docker.Installer
{
    class WindowsMessage
    {
        //https://msdn.microsoft.com/en-us/library/windows/desktop/ms644952(v=vs.85).aspx

        static class NativeMethods
        {
            [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
            public static extern IntPtr SendMessageTimeout(
                IntPtr hWnd,
                uint msg,
                UIntPtr wParam,
                IntPtr lParam,
                uint fuFlags,
                uint uTimeout,
                out UIntPtr lpdwResult
                );
        }

        public static void BroadcastSettingsChange()
        {
            UIntPtr result;
            var setting = Marshal.StringToHGlobalUni("Environment");

            NativeMethods.SendMessageTimeout(
                (IntPtr) 0xFFFF, //HWND_BROADCAST
                0x001A, //WM_SETTINGCHANGE
                (UIntPtr) 0,
                setting,
                0x0002, // SMTO_ABORTIFHUNG
                5000,
                out result
                );

            Marshal.FreeHGlobal(setting);
        }
    }
}
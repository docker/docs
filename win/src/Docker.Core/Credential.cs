using System;
using System.Runtime.InteropServices;
using System.Security;

namespace Docker.Core
{
    [Serializable]
    public class Credential
    {
        public Credential(string user, string password)
        {
            User = user;
            Password = password;
        }

        public bool EmptyPassword => (Password == null) || (Password.Length == 0);

        public string User { get; }

        public string Password { get; }

        private bool Equals(Credential other)
        {
            return string.Equals(User, other.User) && string.Equals(Password, other.Password);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != GetType()) return false;
            return Equals((Credential) obj);
        }

        public override int GetHashCode()
        {
            unchecked
            {
                return ((User?.GetHashCode() ?? 0)*397) ^ (Password?.GetHashCode() ?? 0);
            }
        }

        public override string ToString()
        {
            return $"{User}:**********";
        }

        public static Credential FromSecureCredential(string username, SecureString password)
        {
            return new Credential(username, ConvertToUnsecureString(password));
        }

        private static string ConvertToUnsecureString(SecureString pwd)
        {
            var clearTextPwd = IntPtr.Zero;
            try
            {
                clearTextPwd = Marshal.SecureStringToGlobalAllocUnicode(pwd);
                return Marshal.PtrToStringUni(clearTextPwd);
            }
            finally
            {
                Marshal.ZeroFreeGlobalAllocUnicode(clearTextPwd);
            }
        }
    }
}
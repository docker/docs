using Docker.Core;

namespace Docker.WPF.Credentials
{
    public interface ICredentialAsker
    {
        bool AskCredential();
        void Reset();
        Credential GetStoredCredential();
    }

    public class FixedCredentialAsker : ICredentialAsker
    {
        private Credential Credential { get; }

        public FixedCredentialAsker(string username, string password)
        {
            Credential = new Credential(username, password);
        }

        public bool AskCredential()
        {
            return true;
        }

        public void Reset()
        {
        }

        public Credential GetStoredCredential()
        {
            return Credential;
        }
    }
}
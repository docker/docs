using Docker.Core;

namespace Docker.WPF
{
    public class HttpProxySettingsViewModel : BindableBase
    {
        private readonly IActions _actions;

        private string _exclusionList;
        private string _httpAddress;
        private string _httpsAddress;
        private bool _useHttpProxy;
        private bool _useSameAddressForHttpAndHttps;

        private HttpProxySettingsViewModel(IActions actions)
        {
            _actions = actions;
        }

        public bool IsHttpsTextboxEnabled => UseHttpProxy && !UseSameAddressForHttpAndHttps;

        public string ExclusionList
        {
            get { return _exclusionList; }
            set
            {
                _exclusionList = value;
                OnPropertyChanged();
            }
        }

        public string HttpAddress
        {
            get { return _httpAddress; }
            set
            {
                _httpAddress = value;
                OnPropertyChanged();
                if (UseSameAddressForHttpAndHttps)
                {
                    HttpsAddress = value;
                }
            }
        }

        public string HttpsAddress
        {
            get { return _httpsAddress; }
            set
            {
                _httpsAddress = value;
                OnPropertyChanged();
            }
        }

        public bool UseHttpProxy
        {
            get { return _useHttpProxy; }
            set
            {
                _useHttpProxy = value;
                OnPropertyChanged();
                OnPropertyChanged(nameof(IsHttpsTextboxEnabled));
            }
        }

        public bool UseSameAddressForHttpAndHttps
        {
            get { return _useSameAddressForHttpAndHttps; }
            set
            {
                _useSameAddressForHttpAndHttps = value;
                OnPropertyChanged();
                OnPropertyChanged(nameof(IsHttpsTextboxEnabled));
                if (value)
                {
                    HttpsAddress = HttpAddress;
                }
            }
        }

        internal void Apply()
        {
            _actions.RestartVm(_ =>
            {
                _.UseHttpProxy = UseHttpProxy;
                _.ProxyHttp = HttpAddress ?? string.Empty;
                _.ProxyHttps = UseSameAddressForHttpAndHttps ? _.ProxyHttp : HttpsAddress ?? string.Empty;
                _.ProxyExclude = ExclusionList;
            });
        }

        public static HttpProxySettingsViewModel Load(Settings settings, IActions actions, ISettingsLoader settingsLoader)
        {
            var result = new HttpProxySettingsViewModel(actions) {UseHttpProxy = settings.UseHttpProxy};
            if (result.UseHttpProxy)
            {
                result.HttpAddress = settings.ProxyHttp;
                if (settings.ProxyHttp == settings.ProxyHttps)
                {
                    result.UseSameAddressForHttpAndHttps = true;
                }
                else
                {
                    result.HttpsAddress = settings.ProxyHttps;
                }
                result.ExclusionList = settings.ProxyExclude;
            }
            else
            {
                result.UseSameAddressForHttpAndHttps = true;
            }

            return result;
        }
    }
}
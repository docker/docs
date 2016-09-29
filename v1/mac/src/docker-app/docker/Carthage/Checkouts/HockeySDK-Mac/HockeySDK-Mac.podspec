Pod::Spec.new do |s|
  s.name              = 'HockeySDK-Mac'
  s.version           = '4.0.3'

  s.summary           = 'Collect live crash reports, get feedback from your users, distribute your betas, and get usage data.'
  s.description       = <<-DESC
                        HockeyApp is a service to distribute beta apps, collect crash reports as well as usage data and
                        communicate with your app's users.
                        
                        It improves the testing process dramatically and can be used for both beta
                        and App Store builds.
                        DESC

  s.homepage          = 'http://hockeyapp.net/'
  s.documentation_url = "http://hockeyapp.net/help/sdk/mac/#{s.version}/"

  s.license           = { :type => 'MIT', :file => 'HockeySDK-Mac/LICENSE.txt' }
  s.author            = { 'Microsoft' => 'support@hockeyapp.net' }
  s.source = { :http => "https://github.com/bitstadium/HockeySDK-Mac/releases/download/#{s.version}/HockeySDK-Mac-#{s.version}.zip" }

  s.platform              = :osx, '10.7'
  s.osx.deployment_target = 10.7
  s.requires_arc          = false
  
  s.vendored_frameworks   = "HockeySDK-Mac/HockeySDK.framework"
  s.resource              = "HockeySDK-Mac/HockeySDK.framework"
  s.pod_target_xcconfig   = { 'LD_RUNPATH_SEARCH_PATHS' => '@executable_path/../Frameworks' }

end

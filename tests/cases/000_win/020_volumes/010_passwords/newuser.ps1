$testUsername = $Args[0]
$testUser = Get-LocalUser $testUsername -ea SilentlyContinue

if (!$testUser) {
    $password = ConvertTo-SecureString "Docker4theTeam!" -AsPlainText -Force
    New-LocalUser $testUsername -Password $password -FullName "Docker Test User" -Description "Local user for 2e2 tests." -PasswordNeverExpires
}
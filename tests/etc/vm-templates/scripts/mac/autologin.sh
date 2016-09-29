if [ "$AUTOLOGIN" != "true" ] && [ "$AUTOLOGIN" != "1" ]; then
  exit
fi

echo "Enabling automatic GUI login for the '$USERNAME' user.."

python /private/tmp/set_kcpassword.py "$PASSWORD"

/usr/bin/defaults write /Library/Preferences/com.apple.loginwindow autoLoginUser "$USERNAME"

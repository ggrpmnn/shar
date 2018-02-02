# shar

shar is a simple CLI tool designed to help you keep an eye on who's talking to your Raspberry Pi. It simply parses debian's `/var/log/auth.log` and looks for failed login attempts, then displays them in an easy-to-digest format.

Because it needs access to your `auth.log` file, shar will need to be run as a user with `root` privileges.

Currently, shar only supports looking at SSH login attempts.

# shar

shar is a simple CLI tool designed to help you keep an eye on who's talking to your Raspberry Pi. It simply parses debian's `/var/log/auth.log` file and looks for failed SSH login attempts, then displays them in an easy-to-digest format.

Because it needs access to your `auth.log` file, shar will need to be run as a user with `root` privileges.

Currently, shar only supports looking at SSH login attempts.

### Installation

The easiest way to get this tool up and running is to install `git` on your RaspberryPi (or debian-based Linux machine), then clone this repository. Once the code has been pulled, it can be installed by running `go install`.

### Usage

Running `shar` is as simple as running `sudo shar` in your terminal (sudo is required to grant the app access to the `auth.log` file). Options for output and filtering results can be found by using the help (`-h`) flag.

### TODOs

[ ] Use batch requesting to avoid dealing with ip-api rate limiting

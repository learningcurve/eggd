eggd
====

eggd is a git and foreman autodeploy daemon for Procfile apps on Amazon EC2
instances. It's written in Go and uses an interface to the Linux filesystem's
inotify API to restart foreman when a git push is detected.

(Experimental, not for general use yet.)

Installation
------------

You will need Go installed. Running

    sudo apt-get install golang

suffices. To build and install eggd to <code>/usr/local/bin</code>:

    git clone https://github.com/learningcurve/eggd.git
    cd eggd && make && sudo make install

Set the variable <code>$INSTALL</code> to change the path.

TODO: upstart support.

Usage
-----

1.  In your project directory, make sure you have a Procfile in the root
    directory and a git remote set up to a bare repository on your EC2 instance.

2.  On your EC2 instance, let eggd track your remote repository:

        eggd add /path/to/remote/repo.git

3.  In your project repository, push your commits to your EC2 remote:

        git push your-remote your-branch

4.  eggd will run the Procfile on your instance. Profit!

Supported Platforms
-------------------

Currently tested on the Ubuntu Server 12.04 EC2 image.

FAQ
---

TODO

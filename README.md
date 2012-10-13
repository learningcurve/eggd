eggd
====

eggd is an automated git and foreman deployment daemon for Procfile apps on Amazon EC2 instances.
(Experimental, not for general use yet.)

Installation
------------

    git clone https://github.com/learningcurve/eggd.git
    cd eggd && make && make install

Usage
-----

1. In your project's git repository, make sure you have a Procfile in the project directory and remote set up to a bare repository on your EC2 instance.

2. On your EC2 instance, run
        eggd /path/to/remote/repo.git

3. In your project repository, your commits to your EC2 remote:
        git push your-remote your-branch

4. eggd will run the Procfile on your instance. Profit!

Platforms
---------

Currently tested on the Ubuntu Server 12.04 EC2 image.

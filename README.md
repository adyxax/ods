# ODS

ODS is a simple French crosswords dictionnary web application. It displays a
simple input form that users can use to submit words. The web application will
return whether the words are valid or not in the French scrabble.

## Dependencies

go is required. Though the code most certainly work on many operating systems,
only go version >= 1.24 on Linux amd64 is being regularly tested.

## Quick Start

There is a makefile with everything you need, just type `make help` (or `gmake
help` if running BSD) to see the possible actions.

Use `make build` to build a static binary of the service, embedding everything.

## The dictionary

The "Officiel Du Scrabble" (ODS for short) is what the official dictionary for
the scrabble game is called. One very sad thing is that this dictionary is not
free! You cannot download it digitally, which seems crazy for a simple list of
words. You might use your google-fu and maybe find it on some random GitHub
account if you look for it, but I certainly did not.

This repository relies on git-crypt to secure the content of my own dictionary
file, sorry for not sharing it.

## Systemd service

I use this simple systemd user service unit to run ODS as an unprivileged user:

``` ini
[Unit]
Description=ods.adyxax.org service

[Service]
Environment="ODS_PORT=8090"
ExecStart=/usr/local/bin/ods
Type=simple
```

## Nginx reverse proxy

I use this simple nginx configuration in front of ODS:

``` nginx
server {
	listen  80;
	listen  [::]:80;
	server_name  ods.adyxax.org;
	location / {
		return 308 https://$server_name$request_uri;
	}
}

server {
	listen  443       ssl;
	listen  [::]:443  ssl;
	server_name  ods.adyxax.org;

    error_page  404  /404.html;
	location / {
        include  headers_static.conf;
        proxy_pass  http://127.0.0.1:8090;
	}
	ssl_certificate      adyxax.org.fullchain;
	ssl_certificate_key  adyxax.org.key;
}
```

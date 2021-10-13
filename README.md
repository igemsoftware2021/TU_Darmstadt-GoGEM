# GoGEM

## About

GoGEM is a tool designed to make the deployment process of your iGEM-Wiki as easy as possible.
It is able to automatically fetch a page from a WordPress Instance hosted by your team. If you use this tool please give credit to me and link to this repo (https://github.com/Jackd4w/GoGEM).

## Installation

To install this program download clone the repo to your PC and install go (https://golang.org). After that change into the directory of the repo and execute the _go build_ command.

A pre-compiled version will be available later on.

## Usage

If installed the tool can be used by executing the _GoGEM_ command in your CLI.

## Examples

Upload: _GoGEM upload -u "[Your Username]" -y 2021 -t "TU_Darmstadt" -w "[Your WP Wiki]" -o "test"_

Save your WP Page: _GoGEM fetchWP [URL]_

Purge: _GoGEM purge -u "[Username]" -y [Wiki Year] -t "[Teamname]" -o "[Offset]"_

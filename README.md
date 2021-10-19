# GoGEM

![GoGEM Logo](/goGEM-Logo.png)

## About

GoGEM is a tool designed to make the deployment process of your iGEM-Wiki as easy as possible.
It is able to automatically fetch a page from a WordPress Instance hosted by your team. If you use this tool please give credit to me and link to this repo (<https://github.com/Jackd4w/GoGEM>).

**For detailed information about this tool visit our [iGEM Wiki](https://2021.igem.org/Team:TU_Darmstadt/software) and check out the [wiki guide](https://2021.igem.org/wiki/images/7/70/T--TU_Darmstadt--GoGEM-How-to-Wiki-the-Darmstadt-Way.pdf)!**

## Installation

To install this programm you can use the _go install github.com/Jackd4w/GoGEM_ command.
alternatively the cloning of this repo and _go run_ or _go build_ can be used.

A pre-compiled version will be available with each release.

## Usage

If installed the tool can be used by executing the _GoGEM_ command in your CLI.

If you want to use a pre-compiled version you will have to download the executables and place them in a folder that you can access with you command line of choice. Open the folder and run the downloaded executable with the commands stated below.

## Examples

**Help**: _GoGEM_ or _GoGEM --help_

A help message will be displayed if you run the tool without or with wrong input.

**Upload**: _GoGEM upload -u "[Username]" -y [year] -t "[Teamname]" -w "[WP URL]" -o "[offset]"_

This is the all-in-one command. It downloads your WordPress Page, uploads all the media files, replaces all the links and then uploads all the pages.

**Save your WP Page locally**: _GoGEM fetchWP [URL]_

**Purge**: _GoGEM purge -u "[Username]" -y [Wiki Year] -t "[Teamname]" -o "[Offset]"_

Purge overwrites **all** pages in the defined subspace with an empty one.

You will get a list with deleted pages beforehand and will have to enter your password a second time.
**BE SURE YOU KNOW WHAT THIS DOES BEFORE USING!**

## Issues

Please report Issues to this repo (<https://github.com/Jackd4w/GoGEM>), this is where the development will continue.

## Contribution

I tried to comment the code reasonably. Please try writing verbose comments when contributing, as this is intended to be a project a beginner programmer can understand. As stated above development will only be continued on this repo (<https://github.com/Jackd4w/GoGEM>).

## Required Packages

Directly required non standard packages:

-[GoGEM-WikiAPI](https://github.com/Jackd4w/GoGEM-WikiAPI)
-[OrderedMap](https://github.com/elliotchance/orderedmap)
-[Colly](https://github.com/gocolly/colly)
-[Cobra](https://github.com/spf13/cobra)
-[Viper](https://github.com/spf13/viper)
-[Term](https://golang.org/x/term)

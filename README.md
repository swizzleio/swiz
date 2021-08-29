# Why swiz

This is the amazing, fantabulous, swizzle command line. Connecting through bastion hosts generally sucks and requires
all sorts of command line-fu known only to the elite. Even the elites end up modifying .rc files or writing scripts so
they don't have to remember arcaine ssh incantations.

Worse yet, there are plenty of footguns which may create security problems for your team if the right incantation is not
used.

Well no more secret handshakes, hidden ceremonies, robes, and hazing rituals because the `swiz` command line is here
to help you. This tool will create a connection through a tunnel based on the list of (currently AWS) resources and
then open the necessary tool to interact with the resource.

# How do I use this?

TBD. I'm still writing it :)

# Getting started with development

### Dependencies

Install the following optional dependencies (from outside the project):

* `go install github.com/loov/goda@latest`

# Future roadmap

* Allow for a reverse tunnel that allows for explicit interactions with local resources (webhooks in your local dev
  environment anyone)
* Support for additional clouds (Azure, GCP, etc)
* Support for local datacenter environments
* Support for additional tunnel types
* Incorporate a command that orders artisnal grilled cheese sandwiches right to your developer workstation

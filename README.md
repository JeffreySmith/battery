# Battery
A Go library that gives you battery information, assuming one is present. 

# Platform Support
Currently, this library only supports OpenBSD. The plan is to eventually extend support to Linux and MacOS.

# Rationale
I built this because I got tired of doing the math when running the `apm` command in OpenBSD (it returns minutes). This way, I would always be able to see the estimated battery life in a more human friendly way.

This is a library so that I can eventually use this to display battery info in a cross platform way (for info that can be made cross platform)

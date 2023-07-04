![ShellWeGo-Banner](img/logo.jpg)


# ShellWeGo Reverse Shell

ShellWeGo is a lightweight reverse shell written in Go, which allows you to remotely execute commands on a compromised system, providing a nice tool for penetration testing, and security research.

The connection configuration isn't hard-coded or passed in as command-line arguments, but rather encoded in the filename of the executable itself. This has several advantages, such as making the reverse shell more flexible, evading certain detection mechanisms, and avoiding triggering sandbox environment warnings.

The file name of the executable is split into parts, each of them defining a specific configuration:

* IP Address: specifying the IP address of the server that the reverse shell should connect to.

* Port: indicates the port on the server that the reverse shell should connect to.

* Wait Time: an optional component that, when provided, defines the delay before the reverse shell attempts to establish a connection. This delay is specified in seconds.

* Execution Path: another optional component, when provided, the reverse shell will only run if its current working directory contains the specified string. If the condition isn't met, the program will exit immediately.

This mechanism provides an additional layer of evasion, allowing the reverse shell to blend in better with legitimate software and making it harder for sandbox environments and automated analysis tools to accurately assess its behavior. 

For example, the delay before the connection can help to avoid sandboxes that only observe the behavior of executable for a limited amount of time. Also, the execution path check can help to ensure that the shell is being run in a specific, intended environment.

Here's an example of how you might name the executable to use these features:

**192.168.50.129-4444-10-RandomName.exe**

In this case, the reverse shell will connect to 192.168.50.129 on port 4444, will wait 10 seconds before attempting to establish a connection, and will only run if its current working directory contains RandomName.

It's a simple, but effective way of adding some additional flexibility and stealth to reverse shell!

## Features

* Connection via TCP to a server
* Execution of remote commands
* User-friendly menu system 
* Persist in RunOnce registry key, 
* List processes with PIDs and kill processes
* Prank for fun
* Automatic reconnection if connection lost


## Build

To build from source, ensure you have Go installed and just use make !

## Note

This tool is created for educational and legal purposes only. Any misuse of this tool will not be the responsibility of the author.



**[`^        back to top        ^`](#)**
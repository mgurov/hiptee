## hiptee

tee to hipchat 

## tee
 
Tee mode is activated when the tool is invoked without extra arguments. In this mode, the standard input if forwarded to 
  the hipchat room as grey non-alerting messages, followed by a green "Done" notification (also non-alerting).  

````bash
$ echo Hello, World! | go run cmd/hiptee/main.go
Hello, World!
````

The stdin is also echoed to the stdout of the hiptee. 
 
### exec

Exec mode is activated when a command to execute is specified:  

````bash
go run cmd/hiptee/main.go echo exec mode!
exec mode!
````

the standard output of the command is treated the same way as the standard input in the tee mode: it is echoed to both 
standard output of the tool and sent as a grey non-alerting notification to hipchat. Standard error is shown as red alerts
in the hipchat. Green non-alerting "Done" notification is sent to hipchat upon completion.
  
### config

hipchat token and room can be specified either via command line parameters, `HIPTEE_TOKEN` and `HIPTEE_ROOM` environmental
variables or specified in a hiptee.yaml file of a following structure: 

````yaml
hipchat:
  token: tokentokentoken
  room:  999999
````

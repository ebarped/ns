# ns
Switch the namespace of your current k8s context, faster than kubens.

The main motivation was that kubens (part of https://github.com/ahmetb/kubectx) was very slow in my shell (order of seconds).

This alternative is executed in the order of milliseconds.

## "benchmark" (on my localhost, dont take it too seriously)

- ns:

```
time ns services
2024/02/06 09:13:52 switching to namespace "services"

________________________________________________________
Executed in  425.39 millis    fish           external
   usr time   19.76 millis  540.00 micros   19.22 millis
   sys time   11.15 millis  163.00 micros   10.98 millis
```

- kubens:

```
time kubens services

Context "mycontext" modified.
Active namespace is "services".

________________________________________________________
Executed in    2.91 secs      fish           external
   usr time  470.16 millis  550.00 micros  469.61 millis
   sys time  168.49 millis  142.00 micros  168.35 millis                                                                                                                 
```
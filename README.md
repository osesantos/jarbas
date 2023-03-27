# jarbas
jarbas is a cli chatgpt implementation that uses gpt APIs to answer questions

![image](https://user-images.githubusercontent.com/20876378/227887328-71d1a5e5-d041-42ae-9e88-396131193fdc.png)


## How to build it?

```bash
$ go build main.go
```

## How to run it?

```bash
$ go run main.go "QUESTION"
```
or after bulding it
```bash
$ ./main "QUESTION"
```

Note: create a link on to the `main` so that it will be accessible anywhere
```bash
$ ln main /usr/local/bin/jarbas
```


### Todo
- [ ] Create new subcommand to enter into chat mode, This mode will keep track of the messages and use the chatcompletion API to keep a chat session open. 

# jarbas
jarbas is a cli chatgpt implementation that uses gpt APIs to answer questions

![image](https://user-images.githubusercontent.com/20876378/227887438-f0d6b129-0c4c-4ca6-8be6-a180c08a32fd.png)

## Chat mode

```bash 
$ go run main.go chat
```

![image](https://user-images.githubusercontent.com/20876378/228389477-c64b037d-5cf4-41e1-9cc0-9764e742ed22.png)


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

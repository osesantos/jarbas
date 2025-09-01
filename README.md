[![Go](https://github.com/osesantos/jarbas/actions/workflows/go.yml/badge.svg)](https://github.com/osesantos/jarbas/actions/workflows/go.yml)

# jarbas
jarbas is a cli chatgpt implementation that uses gpt APIs to answer questions

![image](https://user-images.githubusercontent.com/20876378/227887438-f0d6b129-0c4c-4ca6-8be6-a180c08a32fd.png)

## Chat mode

```bash 
$ go run main.go chat
```

![image](https://user-images.githubusercontent.com/20876378/228389477-c64b037d-5cf4-41e1-9cc0-9764e742ed22.png)

### Editor mode

While running the chat, you can call the editor, which will open a Vim window. All input will be used for the chat question.

![image](https://github.com/user-attachments/assets/a3543fa7-53a8-4008-8893-df34bbda990f)

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
$ sudo ln -s $HOME/git/jarbas/main /usr/local/bin/jarbas
```

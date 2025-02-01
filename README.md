[![Go](https://github.com/osesantos/jarbas/actions/workflows/go.yml/badge.svg)](https://github.com/osesantos/jarbas/actions/workflows/go.yml)

# jarbas
jarbas is a cli chatgpt implementation that uses gpt APIs to answer questions

![image](https://user-images.githubusercontent.com/20876378/227887438-f0d6b129-0c4c-4ca6-8be6-a180c08a32fd.png)

## TODO:
- [x] Give the option to save all the chat messages in a .jarbas folder located in the home dir (filename = guid.txt) the file will contain the response from the api will all the messages.
- [x] Give the option to open the list of convesations and select one of them to continue.
- [x] Add the possibility to use claude api
- [x] Add the possibility to use the openai api
- [ ] Take the local list of conversations and show them in the chat
- [ ] take the local list of conversations and show them as a chat history without the need to open the chat
- [ ] Add agents concepts to the chat, so that the user can select the agent to use.
- [ ] Improve agent capabilities, by adding the possibility to scrape the web in real time, get an article and summarize it.
- [ ] Allow the user to select the api to use

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

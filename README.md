<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/katallaxie/autobot/blob/main/images/logo.png?raw=true" alt="Autobot"/>
</div>

# :robot: Autobot

[![Test & Build](https://github.com/katallaxie/autobot/actions/workflows/main.yml/badge.svg)](https://github.com/katallaxie/autobot/actions/workflows/main.yml)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

This is your happy :robot: chatbot for automating operations.

## Get Started

### Adapters

They ingest message from sources. The example uses Google Chats `cmd/hangouts-chat` to ingest messages. And also publish these messages back to the source.

### Plugins

Plugins represents features on top of message that are ingested from the adapters.

### Examples

* [`cmd/hello-world`](/cmd/hello-world/) returns a `Hello World` upon `hello me` command
* [`cmd/time`](/cmd/time) return the current time upon a `time now` command

Happy coding!

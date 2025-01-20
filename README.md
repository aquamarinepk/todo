# TODO Application

This is a TODO application. This simple app is an exploratory attempt to create a kit for developing microservices and self-contained apps.

This reference app is designed to showcase the features and patterns that the kit generator should later replicate for other specific projects.

<img src="docs/img/todo.png" alt="TODO Application" />

## Overview
This is a sample app using the common example of a todo list. The example is designed to express the desired outcome for the generator library.

Among other upcoming features, the main idea is that the generator will expose its features in two different ways. Related features involving one or more resources will be organized in separate packages under `feat`. The interface to handle its logic will be exposed through a command-query interface. For simple cases, where a client needs to manage basic CRUD operations, a RESTful interface can also be generated, with resource entities defined under a directory hanging from `res`. This is handy for simple cases or to interface with more traditional clients.

## Configuration
### Command-Line Flags
You can pass the following flags to the application:  
```txt
server.web.host: Host for the web server (default: localhost)
server.web.port: Port for the web server (default: 8080)
server.api.host: Host for the API server (default: localhost)
server.api.port: Port for the API server (default: 8081)
```

### Environment Variables
The same configuration can be set using environment variables. The environment variables should be prefixed with TODO_ and use underscores instead of dots. For example:  
```shell
TODO_SERVER_WEB_HOST: 127.0.0.1
TODO_SERVER_WEB_PORT: 8080
TODO_SERVER_API_HOST: 127.0.0.1
TODO_SERVER_API_PORT: 8081
```

## Usage
### Running the Application

Using environment variables:
```shell
export TODO_SERVER_WEB_HOST=127.0.0.1
export TODO_SERVER_WEB_PORT=8080
export TODO_SERVER_API_HOST=127.0.0.1
export TODO_SERVER_API_PORT=8081
./todo
```

Using command-line flags:

```shell
./todo -server.web.host=127.0.0.1 -server.web.port=8080 -server.api.host=127.0.0.1 -server.api.port=8081
```

## Notes

There is significant repetition due to the decision to use composition and delegation for providing core functionality to various entities. While this could be avoided by using embedding, the intention would not be as explicit. For now, we will stick to this approach. If embedding becomes a more sensible option in the future, we can test it out.

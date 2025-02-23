# go-mdapi (currently in beta)

## intention

go-mdapi is yet another cli client based on markdown declaration. Original idea is to use it togeather with nvim (to have syntax highlight).
At the moment supports built-in http client (aka simple http requests) + ability to extend the api via go templates: [samples](/samples/)

## Installation

To install go-mdapi, you need to have Go installed on your machine. Then, you can use the following command to install the CLI:

```
go install github.com/catmorte/go-mdapi@latest
```

## Usage

### CLI Commands

Here are the available CLI commands:

- `go-mdapi var_types`: Returns all available var types.
- `go-mdapi types`: Returns all available types declared in the `$HOME/.config/go-mdapi` folder.
- `go-mdapi generate [type]`: Generates an API of the specified type.
- `go-mdapi vars [var_name] [index]`: Shows all the vars in the format `name:type:count`.
- `go-mdapi run`: Runs the API.

### Examples

#### Example 1: Running a simple HTTP request

1. Create a markdown file `example.md` with the following content:

```
# Example API

## vars

## type[http]

### method

```
GET
```

### url

```text
http://localhost:3000
```

## after
```

2. Run the following command:

```
go-mdapi run -f example.md
```

Expected output:

```
result/example/status
result/example/headers
result/example/body
```

#### Example 2: Extending the API with custom templates

1. Create a new folder in `$HOME/.config/go-mdapi` with the name of your custom type (e.g., `custom`).
2. Create two files in the new folder: `run.tmpl` and `new_api.md`.

`run.tmpl`:

```
echo "Running custom API with method {{ index . "method" }} and URL {{ index . "url" }}"
```

`new_api.md`:

```
# Custom API

## vars

## type[custom]

### method

```
GET
```

### url

```text
http://localhost:3000
```

## after
```

3. Run the following command to generate the new API:

```
go-mdapi generate custom
```

4. Run the following command to execute the custom API:

```
go-mdapi run -f custom_api.md
```

Expected output:

```
Running custom API with method GET and URL http://localhost:3000
```

## Structure of Markdown Files

The markdown files used by the CLI have the following structure:

```
# API Name

## vars

## type[api_type]

### var_name

```
var_value
```

## after
```

## Built-in HTTP Client

The built-in HTTP client allows you to make simple HTTP requests by specifying the method and URL in the markdown file. You can also extend the API via Go templates to add custom functionality.

## Available Commands

- `go-mdapi var_types`: Returns all available var types.
- `go-mdapi types`: Returns all available types declared in the `$HOME/.config/go-mdapi` folder.
- `go-mdapi generate [type]`: Generates an API of the specified type.
- `go-mdapi vars [var_name] [index]`: Shows all the vars in the format `name:type:count`.
- `go-mdapi run`: Runs the API.

## Installation and Setup

To install go-mdapi, you need to have Go installed on your machine. Then, you can use the following command to install the CLI:

```
go install github.com/catmorte/go-mdapi@latest
```

To set up the project, create a folder in `$HOME/.config/go-mdapi` and add your custom templates in separate folders within this directory.

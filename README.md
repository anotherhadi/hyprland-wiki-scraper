# Hyprland Wiki Scraper

## Overview

The Hyprland Wiki Scraper is a tool designed to extract information from the Hyprland Wiki. It parses the Markdown content of the Variables page, organizes the information into a structured format, and outputs it as JSON.
It will read configuration variables from the "Variables", "Dwindle Layout" and "Master Layout" pages.

Used for [hyprsettings](https://github.com/anotherhadi/hyprsettings)

## Features

- **Markdown Parsing:** Utilizes regular expressions to extract data from the Markdown structure of the Variables page.
- **Data Structuring:** Organizes the parsed information into a structured format using Go data structures.
- **JSON Output:** Generates a JSON representation of the parsed data for easy consumption.

## Installation

Install it with `go`:

```bash
go install github.com/anotherhadi/hyprland-wiki-scraper@latest
```

## Usage

```bash
hyprland-wiki-scraper
```

**Output:**

```json
...
{
      "name": "General",
      "childs": [
        {
          "name": "Sensitivity",
          "variable": "sensitivity",
          "description": "Mouse sensitivity (legacy, may cause bugs if not 1, prefer `input:sensitivity`).",
          "type": "float",
          "default": "1.0"
        },
        {
          "name": "Border size",
          "variable": "border_size",
          "description": "Size of the border around windows.",
          "type": "int",
          "default": "1"
        },
        {
          "name": "No border on floating",
          "variable": "no_border_on_floating",
          "description": "Disable borders for floating windows.",
          "type": "bool",
          "default": "false"
        },
...
```

## Dependencies

- Go (Golang): The tool is developed using the Go programming language.

## License

This project is licensed under the [MIT License](LICENSE).

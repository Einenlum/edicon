# Edicon

A small project written in go to edit any configuration file from the terminal.

The goal is to be able to get and set values from any configuration file without having to rely on tools like `sed` or `awk`.
It should only change the matching lines and keep the rest of the file intact to have a clean diff.

This is a fun toy project to learn go. Does it make sense? Probably not. Is it well written or thought out? Definitely not.

## Installation

```bash
go get github.com/einenlum/edicon
```

## Usage

### Get the value of a key

```bash
edicon <config-type> get <key> <file>
```

Given the following PHP ini file:

```ini
; This is a comment
[Section1]
key1 = value1
key2.foo = value2
```

You can get the value of `key1` with:

```bash
edicon php get Section1.key1 file.ini
value1
```

If you want to get the value of a key or section that contains a dot (`.`) you can use the brackets notation:

```bash
edicon php get --brackets "Section1[key2.foo]" file.ini
value2
```

### Set the value of a key

```bash
edicon <config-type> set <key> <value> <file>
```

You can set the value of `key1` with:

```bash
edicon php set Section1.key1 newValue file.ini
; This is a comment
[Section1]
key1=newValue
key2.foo = value2
```

This will print the modified version of the config. If you want to save the changes to the file you can use the `-w` or `--write` flag:

```bash
edicon php set -w Section1.key1 newValue file.ini
```

If you want to set the value of a key or section that contains a dot (`.`) you can use the brackets notation:

```bash
edicon php set --brackets "Section1[key2.foo]" newValue file.ini
; This is a comment
[Section1]
key1 = value1
key2.foo=newValue
```

If you want to only print meaningful lines (without comments) you can use the `--values-only` flag:

```bash
edicon php set --values-only Section1.key1 newValue file.ini
[Section1]
key1=newValue
key2.foo = value2
```

## Currently supported configuration types

| Type       | config key | Misc                   | Get parameter      | Set existing parameter | Set new parameter |
| ---        | ---        | ---                    | ---                | ---                    | ---               |
| INI config | `ini`      |                        | :heavy_check_mark: | :heavy_check_mark:     | _missing_         |
| PHP Ini    | `php`      | Just an alias to `ini` | :heavy_check_mark: | :heavy_check_mark:     | _missing_         |

## Misc

- Why not use a parser like [go-ini](https://github.com/go-ini/ini)?

Because I want to keep the original formatting of the file. If I use a parser I would lose comments, empty lines, etc.

# sourceenv
Run command with environment read from .env file

Usage: sourceenv <.env file> <command> <arg0> .. <argN>

Whacky tool so far only intended for personal use.

## Expected syntax of .env file

- Lines starting immediately with # are ignored as comments
- Lines with key=value pairs: whitespace around key and value are ignored.
- Multi-line mode: Start with key<<<END_MARK where END_MARK can contain any non-whitespace characters. All lines, including those starting with # are then read until a line starting with END_MARK is found. Therefore all multi-line values end with a line break. Works for me ;-)


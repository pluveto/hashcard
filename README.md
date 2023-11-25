# Hashcard

Hashcard is a command-line tool that scans markdown (`.md`) files within a specified directory, identifies special card markers, and extracts content into a structured JSON format.

## Usage

Assert that the following markdown file exists at `./cards.md`:

```markdown
Front of Card 1
#card <!--2023/11/25/card1-->
Back of Card 1
---
Front of Card 2
#card <!--2023/11/25/card2-->
Back of Card 2`
```

Run the following command:

```bash
./hash_card --dir . --strategy md-json --out-dir ./out
```

The program will create a file named `cards.json` in the `./out` directory with the following contents:

```json
[
  {
    "id": "2023/11/25/card1",
    "front": "Front of Card 1",
    "back": "Back of Card 1",
  },
  {
    "id": "2023/11/25/card2",
    "front": "Front of Card 2",
    "back": "Back of Card 2",
  }
]
```

## Installation

To use hashcard, you need to have Go installed on your system. If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/).

After installing Go, you can build the program using the following command:

```bash
go build -o hash_card
```

This will compile the source code into an executable file named `hash_card`.

## Usage

You can run the program with the following command-line arguments:

```bash
./hash_card --dir <directory> --strategy md-json --out-dir <output-directory>
```

- `--dir`: The directory to scan for markdown files.
- `--strategy`: The strategy for processing markdown files (currently

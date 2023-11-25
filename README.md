# Hashcard

Hashcard is a command-line tool that scans markdown (`.md`) files within a specified directory, identifies special card markers, and extracts content into a structured JSON format.

Generated JSON files can be easily imported into [Anki](https://apps.ankiweb.net/) or other flashcard apps.

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

i.e. seperate cards with `---` and mark each card with `#card <!--id-->` to distinguish front and back.

> NOTE: id should be unique because it is used as the key of the card. This is important when you want to sync cards between different versions where the front text may be different but the id is the same.

Run the following command:

```bash
./hash_card
```

This is equivalent to:

```bash
./hash_card  --dir . --strategy md-json --out-dir ./out
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

You can download from the [releases page](https://github.com/pluveto/hashcard/releases).

Or you can build from source:

To use hashcard from source, you need to have Go installed on your system. If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/).

Also make sure that you have the `make` command installed on your system.

Then build, using the following command:

```bash
make build
```

Artifacts will be placed in the `./dist` directory.

# chub

> Run multiple commands simultaneously. Useful for dev servers.

## Installation

```bash
git clone https://github.com/megapctr/chub
cd chub
go install
```

## Usage

```
$ cat .chub.json
{
  "commands": {
    "three": ["sh", "three.sh"],
    "four": ["sh", "four.sh"]
  }
}

$ cat three.sh four.sh
while true; do
  echo three; sleep 1
done
while true; do
  echo against four; sleep 0.75
done

$ chub
three ▏three
four  ▏against four
three ▏three
four  ▏against four
```

# refill

Rehydrate JSON schemas from arbitrary text-based data using language models.

# Usage

```sh
$ go get github.com/joehewett/refill
```

```sh
$ refill -json structure.json -dir data -verbose
```

# Example

Given this empty JSON schema containing PII:

```json
{
  "name": "",
  "email": "",
  "phone": "",
  "address": "",
}
```

And this text file:
```
Hi Steve, it's Joe from the office. I'm calling to let you know that your order is ready for pickup. Please come by the office at your earliest convenience to pick it up. Thanks!

J. Hewett
Sales Manager

joe@company.net
(555) 555-5555
```

The `refill` command will produce this output:

```json
{
  "name": "Joe Hewett",
  "email": "joe@company.net",
  "phone": "(555) 555-5555",
  "address": "",
}
```

# License

[MIT](LICENSE)

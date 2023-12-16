# reval // llm prompt evalation using regular expressions

Evaluate the efficacy of your prompts based on a predefined set of regular expressions, and an arbitrary list of tasks.

# Usage

```sh
$ go get github.com/joehewett/refill
```

```sh
$ reval -prompt prompt.json -rules rules.json -tasks tasks.json

```

Reval iterates over tasks, substituting their content into prompts for a language model, with the response being run over the suite of rules specified in the ruleset. Aggregate stats are returned on the scores for each file for a given prompt.

The tasks file must be an array of strings that contain only the content to be substituted into the prompt, whilst the prompt must contain `{{ .task }}`` to designate where each tasks content is to be inserted.

The ruleset constitutes a list of regex rules in a designated JSON format specifying the rule, display name and weighting of each rule. A positive rating indicates that hitting the rule is a positive thing, whilst negative scores indicate regexes you do not want to match.

Each prompt is tested on every task in the tasks file. Every rule is run on each LLM response, giving a score for each. The overall prompt performance is then a function of the average score for a giving LLM response using that prompt.


### Example Rules

```json
[
  {
    "name": "Claims to be AI",
    "regex": "/(?i)\b(as a language model|as an ai)\b/",
    "weight": "-40",
  },
  {
    "name": "Uses name correctly",
    "regex": "(my|I\sam)\s*(name)?\s*(is)?\s*[Jj][Ee][Ee][Vv][Ee][Ss]",
    "weight": "25",
  }
]
```

# License

[MIT](LICENSE)

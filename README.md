sluk (symbol lookup) is a program to easiliy search and print unicode symbols from the command line

```bash

$ sluk white heavy check mark
✔

$ sluk automobile
🚗

$ sluk smiling face with open mouth --print-description
😃	SMILING FACE WITH OPEN MOUTH
😄	SMILING FACE WITH OPEN MOUTH AND SMILING EYES
😅	SMILING FACE WITH OPEN MOUTH AND COLD SWEAT
😆	SMILING FACE WITH OPEN MOUTH AND TIGHTLY-CLOSED EYES

$ sluk smiling face with open mouth --print-unicode --exact-match
😃	'\U0001F603'

```

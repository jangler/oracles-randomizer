# Owl statue hints

In the randomizer, owl statue messages are replaced with information about the
seed. The messages are of the form "[location] holds [item].", and follow these
rules:

- Checks that are already required to reach the owl statue are not hinted at.
- Each dungeon is a location.
- Groups of map tiles sharing a name are locations.
- Single map tiles with unique names are not locations, and use the name of the
  adjacent/containing group instead.
- Subrosia is treated as a single location in Seasons, as is Rolling Ridge in
  Ages.
- No two owls give the same hint.

Hints and their corresponding owls appear in the log file.

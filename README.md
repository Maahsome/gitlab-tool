# gitlab-tool

This cli tool sprang into life because switching to the browser to grab pipeline job logs, copy and pasting them into `vi` so I could interact with them, well, was quite annoying.  And just like everything annoying, there is a solution.

The interesting challenge was the propagation of the IDs from one command line to another.  Typically we run a command, copy the ID we need to pass along to the next command, and paste it in where we need on the next command line.  iTerm2 has a solution for this in Smart Selection and Actions.  Smart Selection is what happens when you hold down the COMMAND key and sections of text turn into hyperlinks.  

Once you have Smart Selection and Actions setup, the output of this CLI tool contains `bash:` links.  When you hold the COMMAND key, these links will be clickable and will send the matching text (after the `bash:`) to the command line.

## Installation

```bash
brew install maahsome/tap/gitlab-tool --formula
```

## Smart Selection and Action Setup

Smart Selection is in iTerm2 Preferences, Profiles, Advanced, Smart Selection.  Edit the Smart Selection list, click the + to add a new item.

- Notes: CMD URL
- Regex: `<bash:(.*)>`

Once you have that entered, click the `Edit Actions...` option

- Title: Send to Command Line
- Action: Send Text...
- Parameter: \1

### Dynamic Profile Entry

Here is an example of the configuration for a Dynamic Profile for iTerm2

```json
{
  "Profiles": [
    {
      "Custom Window Title" : "Smart Selection Test",
      "Name" : "SMART-SELECTION",
      "Use Custom Window Title" : true,
      "Guid" : "f3aabf14-aeac-433b-8495-c79fbf1689eb",
      "Smart Selection Rules" : [
        {
          "notes" : "CMD URL",
          "precision" : "very_low",
          "regex" : "<bash:(.*)>",
          "actions" : [
            {
              "title" : "Run Command",
              "action" : 4,
              "parameter" : "\\1"
            }
          ]
        }
      ]
    }
  ]
}
```

# hockey-stats
Go CLI to measure NHL stats as the season progresses

## Usage:
```./hockey-stats [-stat] <options>```

stats currently supported are

| stat | description |
--- | ---
| points | - |
| goals | - |
| assists | - |
| plusminus | - |
| PPG | Powerplay goals |
| PPP | Powerplay points |
| SHG | Shorthanded goals |
| SHP | Shorthanded points |
| shots | Shots on goal | 
| PIM | Penalties in minutes |
| GWG | Game winning goals |
| OTG | Overtime goals |

Options:

This can be an indefinite number of players IDs, these can be found on the nhl.com website. For example, the url for Nikita Kucherov is https://www.nhl.com/player/nikita-kucherov-8476453, his ID will be 8476453
  
To compare Nikita Kucherov and Evgeni Malkins power play points through the season, we'd type ```./hockey-stats -stat=PPG 8476453 8471215```

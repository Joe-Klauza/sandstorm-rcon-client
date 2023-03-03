# Sandstorm RCON Client

This tool properly interfaces with Insurgency: Sandstorm servers to issue RCON commands.

```log
] ./sandstorm-rcon-client --help
Usage of sandstorm-rcon-client:
  -c string
        RCON command to execute (short)
  -command string
        RCON command to execute. Omit this for REPL mode.
  -cronfile string
        Path to crontab file for scheduled commands
  -debug
        Enable debug logging
  -f string
        Path to crontab file for scheduled commands (short)
  -n    Disable REPL mode (e.g. when using --cronfile) (short)
  -norepl
        Disable REPL mode (e.g. when using --cronfile)
  -p string
        RCON password (short)
  -password string
        RCON password
  -s string
        Server IP:Port (short) (default "127.0.0.1:27015")
  -server string
        Server IP:Port (default "127.0.0.1:27015")
  -v    Print version and exit (short)
  -version
        Print version and exit
```

## REPL mode

Omit the command argument to enter REPL mode, which **reuses one port to avoid leaking threads per command** (Sandstorm RCON server bug since beta). Otherwise each connection causes the server to create one more thread which is never reaped and can eventually cause server degradation.

Use the tab key for tab completion of top-level command names:

```log
] ./sandstorm-rcon-client
2023-03-03T18:40:14.435 INFO    Connected to server:127.0.0.1:33331
2023-03-03T18:40:14.462 INFO    Authentication successful.
2023-03-03T18:40:14.462 INFO    Entering command REPL. Press ctrl+c or ctrl+d to exit.
RCON 127.0.0.1:27015] <tab>
help                     listplayers              kick                     permban                  travel                   ban                      banid
listbans                 unban                    say                      restartround             maps                     scenarios                travelscenario
gamemodeproperty         listgamemodeproperties
RCON 127.0.0.1:33331] ^C
2023-03-03T18:40:46.152 INFO    Exiting
```

Use the up and down arrow keys to cycle through command history.

See <https://github.com/chzyer/readline/blob/master/doc/shortcut.md> for other readline shortcuts.

## Crontab / Schedule mode

You can use a crontab file to schedule RCON commands to execute at specific times. There is an example file at [example/example.cron](example/example.cron). This can be combined with REPL or single-command mode as desired.

```log
] ./sandstorm-rcon-client --norepl --cronfile example/example.cron
2023-03-03T18:35:26.847 INFO    Connected to server:127.0.0.1:27015
2023-03-03T18:35:26.858 INFO    Authentication successful.
2023-03-03T18:35:26.859 WARN    [CRON] ignoring commented line: # Example crontab file for scheduled commands (e.g. MOTD)
2023-03-03T18:35:26.859 WARN    [CRON] ignoring commented line: # Documentation: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┏━ Second of the minute
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┏━ Minute of the hour
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┃     ┏━ Hour of the day
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┃     ┃     ┏━ Day of the month
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┃     ┃     ┃     ┏━ Month
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┃     ┃     ┃     ┃     ┏━ Day of the week
2023-03-03T18:35:26.860 WARN    [CRON] ignoring commented line: #┃     ┃     ┃     ┃     ┃     ┃     ┏━ RCON command
2023-03-03T18:35:26.860 INFO    [CRON] Added: [*/30 * * * * *] say Every 30 seconds
2023-03-03T18:35:26.860 INFO    [CRON] Added: [0 * * * * *] say Every minute
2023-03-03T18:35:26.860 INFO    [CRON] Added: [0 */10 * * * *] say Every 10 minutes
2023-03-03T18:35:26.860 INFO    Crontab running. Press ctrl+c to stop.
2023-03-03T18:35:30.003 INFO    [CRON] [*/30 * * * * *] say Every 30 seconds
2023-03-03T18:35:30.024 INFO    [CRON] Server response empty
2023-03-03T18:36:00.022 INFO    [CRON] [0 * * * * *] say Every minute
2023-03-03T18:36:00.037 INFO    [CRON] Server response empty
2023-03-03T18:36:00.037 INFO    [CRON] [*/30 * * * * *] say Every 30 seconds
2023-03-03T18:36:00.058 INFO    [CRON] Server response empty
```

## Donate

If you'd like to show your appreciation of `Sandstorm RCON Client`, please [donate via PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=QZDY3PPUMH5TU&item_name=Sandstorm%20RCON%20Client&currency_code=USD) (or suggest other methods).

[![](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=QZDY3PPUMH5TU&item_name=Sandstorm%20RCON%20Client&currency_code=USD)

## Contact

Join the unofficial [Insurgency: Sandstorm Community Server Hosts Discord](https://discord.gg/DSwnmyA)! We'd love to help you with any server hosting questions/issues you may have.

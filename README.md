# Sandstorm RCON Client

This tool properly interfaces with Insurgency: Sandstorm servers to issue RCON commands.

```txt
Usage of sandstorm-rcon-client:
  -c string
        RCON command to execute (short)
  -command string
        RCON command to execute. Omit this for REPL mode.
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

```
sandstorm-rcon-client
Connected to server: 127.0.0.1:27015
Authentication successful.
Entering command REPL. Press ctrl+c or ctrl+d to exit.
RCON 127.0.0.1:27015] <tab>
help                     listplayers              kick                     permban                  travel                   ban                      banid
listbans                 unban                    say                      restartround             maps                     scenarios                travelscenario
gamemodeproperty         listgamemodeproperties
```

Use the up and down arrow keys to cycle through command history.

See <https://github.com/chzyer/readline/blob/master/doc/shortcut.md> for other readline shortcuts.

## Donate

If you'd like to show your appreciation of `Sandstorm RCON Client`, please [donate via PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=QZDY3PPUMH5TU&item_name=Sandstorm%20RCON%20Client&currency_code=USD) (or suggest other methods).

[![](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=QZDY3PPUMH5TU&item_name=Sandstorm%20RCON%20Client&currency_code=USD)

## Contact

Join the unofficial [Insurgency: Sandstorm Community Server Hosts Discord](https://discord.gg/DSwnmyA)! We'd love to help you with any server hosting questions/issues you may have.

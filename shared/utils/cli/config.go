package cli

import (
    "github.com/urfave/cli"
    "github.com/urfave/cli/altsrc"
)


// Configure the application options
func Configure(app *cli.App) {

    // Register global application options & defaults
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "config",
            Usage: "Rocket Pool CLI config file absolute `path`",
            Value: "/.rocketpool/config.yml",
        },
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "database",
            Usage: "Rocket Pool CLI database absolute `path`",
            Value: "/.rocketpool/rocketpool.db",
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "password",
            Usage: "Rocket Pool CLI keystore password `path`",
            Value: "/.rocketpool/password",
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "keychainPow",
            Usage: "PoW chain account keychain absolute `path`",
            Value: "/.rocketpool/accounts",
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "keychainBeacon",
            Usage: "Beacon chain account keychain absolute `path`",
            Value: "/.rocketpool/validators",
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "providerPow",
            Usage: "PoW chain provider `url`",
            Value: "http://127.0.0.1:8545", // Local node
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "providerBeacon",
            Usage: "Beacon chain provider `url`",
            Value: "ws://127.0.0.1:9545", // Local simulator
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "storageAddress",
            Usage: "PoW chain Rocket Pool storage contract `address`",
            Value: "0x2CaBe1D447092CB2A0740c30F43c68624B2cbaD4", // RP2 Beta network
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "uniswapAddress",
            Usage: "PoW chain Uniswap factory contract `address`",
            Value: "0x9c83dCE8CA20E9aAF9D3efc003b2ea62aBC08351", // Ropsten
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "input",
            Usage: "Rocket Pool CLI input file `path` (advanced use only)",
        }),
        altsrc.NewStringFlag(cli.StringFlag{
            Name:  "output",
            Usage: "Rocket Pool CLI output file `path` (advanced use only)",
        }),
    }

    // Load external config, squelch load errors
    app.Before = func(c *cli.Context) error {
        altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
        return nil
    }

}


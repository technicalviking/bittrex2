# README #

This is (my second attempt at) a Go Library for the Bittrex API.  The previous version can be found at https://github.com/technicalviking/bittrex.

### Why make a second repo?  Is the original abandoned? ###

There's nothing wrong with the original repo at all, outside of missing significant code coverage in tests.  Bittrex2 is designed primarily to take advantage of the websocket features of Bittrex's beta api released in march (https://github.com/Bittrex/beta).  

I'm not abandoning the previous version, and eventually the codebases will be merged (once the new bittrex api is out of beta.)

### Installation and Usage ###

    go get github.com/technicalviking/bittrex2

To perform operations such as placing orders or retrieving account balances, you will need to have set up an API Key within your account settings on Bittrex.  If you only want the public methods, you'll need to use version 1 of my lib for now, as this one automatically tries to authenticate against the websocket api with the provided credentials.

    /*the renaming of the package isn't really necessary, but it does call out the mismatch between repo name and package name*/
    import bittrex "github.com/technicalviking/bittrex2"
    ...
    client, err := bittrex.New("YOUR_API_KEY", "YOUR_API_SECRET")

Once your client object has been created, you can call any existing v1.1 endpoint, two of the v2.0 endpoints, or any of the websocket hubcalls described on their github.

####Calling a V1.1 endpoint

Methods for V1.1 endpoint calls are named consistently based on the path.  For example, to make a call to bittrex.com/api/v1.1/acount/getbalance, the method used would be:

    client.AccountGetBalance("BTC") //replace "BTC" with your desired currency.

All v1.1 methods can be found in clientAccount.go, clientMarket.go,  or clientPublic.go .

####Calling a v2.0 endpoint

The only V2.0 endpoints are the important ones with no equivalent in v1.1, used for getting historical data.  They're found in clientUndocumented.go.  The valid intervals that can be used for the second parameter are provided as exported const values within the same file.

    //bittrex.com/api/v2.0/pub/market/getticks
    client.PubMarketGetTicks("USDT-BTC", bittrex.TickIntervalFiveMin)

	//bittrex.com/api/v2.0/pub/market/getlatesttick
	client.PubMarketGetLatestTick("USDT-BTC", "fiveMin"

####Subscribe to Market Summary Deltas

An example of subscribing to data from the websocket connection:

    summaryChan := client.SubscribeToMarketSummary("USDT-BTC")

If you look at the beta api readme, you'll see that the subscription channel actually gives ALL market data once you've subscribed to the channel.  However to make it easier to filter that data by market, I'm providing individual channels for summary changes by market.   The same applies to Summary Lite Delta.  Exchange Deltas actually require a seperate subscription per market, but as this is handled under the hood, each subscription gives you a new channel for the market data from that channel, so the code between these three is consistent.  the code is in socketSubscriptions.go

####Query Socket Data

Unlike subscriptions, these query methods return their replies in a synchronous manner to the user of this sdk.  Code is in socketQueries.go

    state, err := client.QueryExchangeState("USDT-BTC")


### Questions? ###

* Why are you using Float64 for the decimal values?  

Even though the API claims that the decimal values are "string-formatted decimal with 18 significant digits and 8 digit precision", within the json they are not, in fact, encoded as strings.  However the number of significant digits and the precision are well within the safe (ish) bounds of a Float64.  I hope they eventually change this (https://github.com/Bittrex/beta/issues/15) Converting this value to a safer type is left as an exercise for the consuming application.

* Why not just use shopspring/decimal?

First: I don't feel comfortable making presumptive data changes that could result in precision or data loss before the consuming application even sees the data.   Second, on a more personal note, my application that consumes this library does simulations on historical data, and the performance issues of shopspring/decimal come directly to the forefront in that scenario  (look at how many operations call "rescale" in that codebase).

* Your SignalR libary looks familiar....

Not a question, exactly, but I'll answer it anyway.  Yes, it's copy/pasted from either github.com/hweom or from github.com/thebotguys, I can't remember which.  I had to do some significant debugging within that lib when trying to get the socket stuff working, and it was easier to just include that code since it was just one file.  99.9% of the credit goes to those individuals (and most of THAT goes to hweom, the originator of the code).

* Can I make contributions to this repo?

Yes!  Please!  I'd really like to see a golang bittrex api library not get old and stale like others I've seen.

* If I raise an issue, how quickly will you reply?  

Within 24 hours during the week, or 72 hours if the issue is raised on the weekend.

* This code is garbage!

I'm sorry.  Feel free to make a Pull Request into this repo, or if the crimes are too heinous, fork it and modify your own version!

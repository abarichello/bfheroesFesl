[![HitCount](http://hits.dwyl.io/Synaxis/bfheroesFesl.svg)](http://hits.dwyl.io/Synaxis/bfheroesFesl)
# Open Heroes Backend (FESL)
```UNFINISHED CODE```

Remember to configure your GOPATH and type
==>```go build main.go && ./main.go```

## Configuration

Enviroment (.env) variables You can look in `./config/config.go` for more details

| String Name           | Default value        |
|-----------------------|----------------------|
| `LOG_LEVEL`           | `INFO`               |
| `HTTP_BIND`           | `0.0.0.0:8080`       |
| `HTTPS_BIND`          | `0.0.0.0:443`        |
| `GAMESPY_IP`          | `0.0.0.0`(auto bind) |
| `FESL_CLIENT_PORT`    | `18270`              |//cannot be changed
| `FESL_SERVER_PORT`    | `18051`              |//cannot be changed
| `THEATER_CLIENT_PORT` | `18275`              |//cannot be changed
| `THEATER_SERVER_PORT` | `18056`              |//cannot be changed
| `THEATER_ADDR`        | `127.0.0.1`          |
| `LEVEL_DB_PATH`       | `_data/lvl.db`       |
| `DATABASE_USERNAME`   | `root`               |
| `DATABASE_PASSWORD`   |                      |
| `DATABASE_HOST`       | `127.0.0.1`          |
| `DATABASE_PORT`       | `3306`               |
| `DATABASE_NAME`       | `tutorialDB`         |

WARNING for testing environment! Use Safe values in Production!

### Example `.env` file
```ini
DATABASE_NAME=tutorialDB
DATABASE_HOST=127.0.0.1
DATABASE_PASSWORD=dbPass
LOG_LEVEL=DEBUG /INFO.
=================================================================================================================================
# FESL PROTOCOL
This provides the info about the Backend/FESL . Between the Master server , Fesl Server and Theather Server

## Overview

Battlefield Heroes has a network structure similar to many other online games. It is based on previous games that also used the gamespy protocol, such as Battlefield 2 or Battlefield 2142. Need For Speed Carbon , and others

The general  consists of the following components:
1. gameClient.exe: the front-end software that runs on the player's computer. Consists mainly of a graphical userinterface and some game-logic.
2. gameServer.exe: the back-end server that acts as a central game coordinator for the players in a match. Consists mainly of game logic and connections to game clients.
3. Master server: the back-end server that stores player and server data and does match-making. This server provides persistance in between matches.

This specification provides details on the communication between the game client and master server, and the game server and the master server. 
##This document does not specify the protocol between game server and game client.##

## Master server overview
The MASTER server has 3 components:
1. Two FESL servers: a message based protocol server that handles authentication, quering account info, ...
2. The Magma server: a HTTPS based server for more account info and addons(Store,Entitlements,friend list ,Bookmark)
3. Two Theater servers: a message based protocol server that handles querying, joining, leaving, ... For Both game servers and clients

A game client will first connect to the FESL server, then the HTTP server, then the Theater server and finally the game server.
-> game server  = BFHeroes_w32ded.exe

## FESL

### TLS
On Start, both the game client and the game server will first connect to their respective, seperate FESL server. 
The address of the FESL server is inside the game client/server exe HEX. 
They are changed with an Hex editor.

The default value is "bfwest-server.fesl.ea.com".

Communication is encrypted with TLS. The game checks if the TLS certificate is Valid and disconnects if it doesn't match the EA certificate.
We use the FESL Patch , from Aluigi , to patch this Check
http://aluigi.altervista.org/patches/fesl.lpatch

After the TLS handshake, FESL format these messages like this:

|Offset (bytes) |Length (bytes)     |Data type                          |Field name    |
|---------------|-------------------|-----------------------------------|--------------|
|0x0            |4                  |ASCII string (no null terminator)  |Type          |
|0x4            |4                  |32-bit big-endian unsigned integer |ID            |
|0x8            |4                  |32-bit big-endian unsigned integer |              |
|0xC            |Pkt length - 12    |ASCII string (no null terminator)  |FESLData      |

### Packet structure
http://old.zenhax.com/post10292.html

The FESLData field is a key-value map where each pair is seperated by a newline (\n), and the key and value are seperated by '='.
For example:
```
Key1=Value1
Key2=Value2
```
### Message types
One of the keys in the FESLData key-value store is 'TXN'. This entry determines the message type , and Order
Depending on the message type, and whether the message is to or from the FESL server, other fields may be present in the FESLData.
Response Packets are always sent with the same Type and ID values as the query Packet.

#### TXN = Hello, game client/server => FESL server
This is the first Pkt that is sent when a FESL connection is made.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|SDKVersion                |5.0.0.0.0                  |GameSpy SDK version            |
|clientPlatform            |PC                         |                               |
|clientString              |bfwest-pc                  |                               |
|clientType                |server                     |GameServer.exe                 |
|clientVersion             |1.46.222034                |Static only                    |
|locale                    |en_US                      |                               |
|sku                       |125170                     |                               |
|protocolVersion           |2.0                        |TLS version                    |
|fragmentSize              |8096                       |Max buffer size                |

#### TXN = Hello, FESL server => game client/server

|Key                       |Example value              |Note                                             |
|--------------------------|---------------------------|-------------------------------------------------|
|domainPartition.domain    |eagames                    |                                                 |
|domainPartition.subDomain |bfwest-server              |bfwest-server if it's gameServer.exe             |
|                          |                           |bfwest-dedicated for gameClient.exe              |
|curTime                   |Nov-02-2017 22:29:00 UTC   |                                                 |
|activityTimeoutSecs       |3600                       |AFK timeout                                      |
|messengerIp               |messaging.ea.com           |This is not required to play the game.           |
|messengerPort             |13505                      |this was used by EA in the past (friendlist?)    |
|theaterIp                 |bfwest-pc.theater.ea.com   |this needs to be redirected                      |
|theaterPort               |18056                      |By default, 18056 is for game servers and        |
|                          |                           |18275 for game clients                           |

#### TXN = MemCheck, FESL server => game client/server
This message is sent every 10 seconds, and acts a heartbeat. 
If either stops receiving the MemCheck messages, connection loss is assumed.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|memcheck.[]               |0  (guessed response)      |                               |
|salt                      |5  (guessed response)      |                               |

#### TXN = MemCheck, game client/server => FESL server
This message is always a response to a MemCheck query message by the FESL server.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|result                    |*empty*                    |                               |

#### TXN = NuLogin, game client/server => FESL server
This message is sent by clients/servers to authenticate.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|returnEncryptedInfo       |0                          |                               |
|nuid                      |XxX_b3stP1ayer_XxX         |                               |
|password                  |thisIsAPassword            |                               |
|macAddr                   |$31dc51d43797              |                               |

#### TXN = NuLogin, FESL server => game client/server, on error

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|localizedMessage          |Incorrect password.        |                               |
|errorContainer.[]         |0                          |custom error message           |
|errorCode                 |122                        |                               |

#### TXN = NuLogin, FESL server => game client/server, on success

|Key                       |Example value              |Note                                    |
|--------------------------|---------------------------|----------------------------------------|
|profileId                 |1                          |PID                                     |
|userId                    |1                          |uID                                     |
|nuid                      |XxX_b3stP1ayer_XxX         |                                        |
|lkey                      |OwPcFq[xA338SppTjx0Ybw4c   |A 24 character BF2Random (see Appendix) |


#### TXN = NuGetPersonas, game client/server => FESL server
This message is a query to lookup all characters owned by a user.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|                          |                           |                               |

#### TXN = NuGetPersonas, FESL server => game client/server

|Key                       |Example value              |Note                                             |
|--------------------------|---------------------------|-------------------------------------------------|
|personas.*i*              |xXx_1337Sn1per_xXx         |One entry for every character owned by the user. |
|                          |                           |Contains the character name.                     |
|                          |                           |*i* is the zero-based index of the character.    |
|personas.[]               |1                          |The total amount of characters.                  |


#### TXN = NuGetAccount, game client/server => FESL server
This message retrieves general account information, based on the parameters sent.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|                          |                           |                               |

#### TXN = NuGetAccount, FESL server => game client/server

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|heroName                  |xXx_1337Sn1per_xXx         |                               |
|nuid                      |email@account.com          |                               |
|DOBDay                    |1                          |Date Of Birth                  |
|DOBMonth                  |1                          |                               |
|DOBYear                   |1992                       |                               |
|userId                    |1                          |                               |
|globalOptin               |0                          |always 0                       |
|thidPartyOptin            |0                          |always 0                       |
|language                  |enUS                       |                               |
|country                   |US                         |                               |


#### TXN = NuLoginPersona, game client/server => FESL server
This message is sent to login to a character/server.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|name                      |Hero2                      |                               |

#### TXN = NuLoginPersona, FESL server => game client/server

|Key                       |Example value              |Note                                    |
|--------------------------|---------------------------|----------------------------------------|
|lkey                      |OwPcFq[xA338SppTjx0Ybw4c   |A 24 character BF2Random (see Appendix) |
|profileId                 |1                          |                                        |
|userId                    |1                          |                                        |

#### TXN = GetStatsForOwners, game client/server => FESL server
This message is sent to retrieve info for the character selection screen.

|Key                       |Example value              |Note                                     |
|--------------------------|---------------------------|-----------------------------------------|
|keys.*i*                  |c_ltm                      |One entry for every stat to be retrieved |
|keys.[]                   |1                          |Amount of stats to be retrieved          |
|owner                     |2                          |                                         |
|ownerType                 |1                          |                                         |
|periodId                  |0                          |                                         |
|periodPast                |0                          |                                         |

#### TXN = GetStatsForOwners, FESL server => game client/server

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|stats.*i*.ownerId         |1                          |                               |
|stats.*i*.ownerType       |1                          |                               |
|stats.*i*.stats.*j*.key   |level                      |                               |
|stats.*i*.stats.*j*.value |3.0000                     |                               |
|stats.*i*.stats.*j*.text  |3.0000                     |                               |
|stats.*i*.stats.[]        |1                          |                               |
|stats.[]                  |1                          |                               |

#### TXN = GetStats, game client/server => FESL server
This message is sent to retrieve info about a character/user.

|Key                       |Example value              |Note                                     |
|--------------------------|---------------------------|-----------------------------------------|
|owner                     |35                         |                                         |
|ownerType                 |1                          |                                         |
|periodId                  |0                          |                                         |
|periodPast                |0                          |                                         |
|keys.*i*                  |c_items(abilities)         |One entry for every stat to be retrieved |
|keys.[]                   |1                          |Amount of stats to be retrieved          |

#### TXN = GetStats, FESL server => game client/server

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|ownerId                   |2                          |                               |
|ownerType                 |1                          |                               |
|stats.*i*.key             |edm                        |                               |
|stats.*i*.value           |*empty*                    |                               |
|stats.*i*.text            |*empty*                    |                               |
|stats.[]                  |2                          |                               |


#### TXN = NuLookupUserInfo, game client/server => FESL server
This message is sent to retrieve basic information about a user.

|Key                       |Example value              |Note                              |
|--------------------------|---------------------------|----------------------------------|
|userInfo.*i*.userName     |xXx_1337Sn1per_xXx         |Names of the characters to lookup |
|userInfo.[]               |1                          |Amount of characters to lookup    |

#### TXN = NuLookupUserInfo, FESL server => game client/server

|Key                       |Example value              |Note                              |
|--------------------------|---------------------------|----------------------------------|
|userInfo.*i*.userName     |xXx_1337Sn1per_xXx         |                                  |
|userInfo.*i*.userId       |1                          |                                  |
|userInfo.*i*.masterUserId |1                          |                                  |
|userInfo.*i*.namespace    |MAIN                       |                                  |
|userInfo.*i*.xuid         |24                         |                                  |
|userInfo.*i*.cid          |1                          |client ID(not implemented         |
|userInfo.[]               |3                          |Amount of users to lookup info of |

This is a query for a list of endpoints to test for the lowest latency on a game client.
This is not working at the moment, maybe EA used it to change from DataCenters , based on Ping
like LoadBalancer
#### TXN = GetPingSites, FESL server => game client/server

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|minPingSitesToPing        |1                          |this was used in the past tho  |
|pingSites.*i*.addr        |8.8.8.8                    |it doesnt seem to work. check  | 
|pingSites.*i*.name        |iad                        | valid response ?              |
|pingSites.*i*.type        |0                          | or it's just telemetric shit  |
|pingSites.[]              |1                          |                               |


#### TXN = UpdateStats, game client/server => FESL server
This message is sent to update character stats.

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|u.*i*.o                   |2                          |owner                          |
|u.*i*.ot                  |1                          |ownerType                      |
|u.*i*.s.*j*.k             |c_slm                      |key                            |
|u.*i*.s.*j*.ut            |0                          |updateType                     |
|u.*i*.s.*j*.t             |*empty*                    |text                           |
|u.*i*.s.*j*.v             |1.0000                     |value                          |
|u.*i*.s.*j*.pt            |0                          |                               |
|u.*i*.s.[]                |1                          |Amount of stats to query       |
|u.[]                      |1                          |Amount of character to query   |

#### TXN = UpdateStats, FESL server => game client/server

|Key                       |Example value              |Note                             |
|--------------------------|---------------------------|---------------------------------|
|u.*i*.o                   |                           |Values are copied from the query |
|u.*i*.s.*j*.k             |                           |                                 |
|u.*i*.s.*j*.ut            |                           |                                 |
|u.*i*.s.*j*.t             |                           |                                 |
|u.*i*.s.*j*.v             |                           |                                 |
|u.*i*.s.*j*.pt            |                           |                                 |
|u.*i*.s.[]                |                           |                                 |
|u.[]                      |                           |                                 |


#### TXN = GetTelemetryToken, game client/server => FESL server
Returns a unique token for game telemetry.
This is only used in 2009 client ?

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|                          |                           |                               |

#### TXN = GetTelemetryToken, FESL server => game client/server
#### only requested in 2009 client
|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|telemetryToken            |MTU5LjE1My4yMzUuMjYsOTk0Nix|                               |
|                          |lblVTLF7ZmajcnLfGpKSJk53K/4|                               |
|                          |WQj7LRw9asjLHvxLGhgoaMsrDE3|                               |
|                          |bGWhsyb4e6woYKGjJiw4MCBg4bM|                               |
|                          |srnKibuDppiWxYKditSp0amvhJm|                               |
|                          |StMiMlrHk4IGzhoyYsO7A4dLM26|                               |
|                          |rTgAo%3d                   |                               |
|enabled                   |US                         |                               |
|filters                   |*empty*                    |                               |
|disabled                  |*empty*                    |                               |

#### TXN = Start, game client/server => FESL server
This message is sent to initiate a "playnow".

|Key                               |Example value              |Note                                  |
|----------------------------------|---------------------------|--------------------------------------|
|partition.partition               |/eagames/bfwest-dedicated  |                                      |
|debugLevel                        |high                       |                                      |
|players.*i*.ownerId               |2                          |                                      |
|players.*i*.ownerType             |1                          |                                      |
|players.*i*.props.{*propertykey*} |3                          |Example *propertykey* is pref-lvl_avg |
|players.[]                        |1                          |                                      |

#### TXN = Start, FESL server => game client/server

|Key                       |Example value              |Note                           |
|--------------------------|---------------------------|-------------------------------|
|id.id                     |1                          |                               |
|id.partition              |/eagames/bfwest-dedicated  |                               |

## Magma server

On start, the game will connect to an HTTPS server. called MAGMA, this connection is encrypted using TLS. Multiple Magma server domains are in the game executable.

The following HTTPS paths are listed in the game executable:
* /dc/submit
* /nucleus/authToken
* /nucleus/check/%s/%I64d
* /nucleus/entitlement/%s/status/%s
* /nucleus/entitlement/%s/useCount/%d
* /nucleus/entitlements/%I64d
* /nucleus/entitlements/%I64d
* /nucleus/entitlements/%I64d?entitlementTag=%s
* /nucleus/name/%I64d
* /nucleus/personas/%s
* /nucleus/refundAbilities/%I64d
* /nucleus/wallets/%I64d
* /nucleus/wallets/%I64d/%s/%d/%s
* /ofb/products
* /ofb/purchase/%I64d/%s
* /relationships/acknowledge/nucleus:%I64d/%I64d
* /relationships/acknowledge/server:%s/%I64d
* /relationships/decrease/nucleus:%I64d/nucleus:%I64d/%s
* /relationships/decrease/nucleus:%I64d/server:%s/%s
* /relationships/decrease/server:%s/nucleus:%I64d/%s
* /relationships/increase/nucleus:%I64d/nucleus:%I64d/%s
* /relationships/increase/nucleus:%I64d/server:%s/%s
* /relationships/increase/server:%s/nucleus:%I64d/%s
* /relationships/roster/nucleus:%I64d
* /relationships/roster/nucleus:%I64d
* /relationships/roster/server:%s
* /relationships/roster/server:%s/bvip/1,3
* /relationships/status/nucleus:%I64d
* /relationships/status/server:%s
* /user
* /user/updateUserProfile/%I64d

### /nucleus/authToken
  For game servers: 
    ```<success><token>$serverSECRET$</token></success>```
 
  `$SECRET$` equal to the value of the `X-SERVER-KEY` and "+Secret" given in start parameter
 
  For game clients:
 
    ```<success><token code="NEW_TOKEN">$userKey$</token></success>```
 
  with `$userKey$` equal to the value of the `magma` cookie.

### /relationships/roster/{type}:{id}
    ```
    <update>
        <id>1</id>
        <name>Test</name>
        <state>ACTIVE</state>
        <type>server</type>
        <status>Online</status>
        <realid>$id$</realid>
    </update>
    ```
  with $id$ equal to the id in the URL.

### /nucleus/entitlements/{heroID}  // this is not a valid response

    ```
    <?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
    <entitlements>
        <entitlement>
            <entitlementId>1</entitlementId>
            <entitlementTag>WEST_Custom_Item_142</entitlementTag>
            <status>ACTIVE</status>
            <userId>$heroID$</userId>
        </entitlement>
        <entitlement>
            <entitlementId>1253</entitlementId>
            <entitlementTag>WEST_Custom_Item_142</entitlementTag>
            <status>ACTIVE</status>
            <userId>$heroID$</userId>
        </entitlement>
    </entitlements>
    ```
  with $heroID$ equal to heroID in the URL.

### /nucleus/wallets/{heroID}
    ```
    <?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
    <billingAccounts>
        <walletAccount>
            <currency>hp</currency>
            <balance>1</balance>
        </walletAccount>
    </billingAccounts>
    ```

### /ofb/products // this is not a valid response
    ```
    <?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
    <products>
        <product overwrite="false" default_locale="en_US" status="Published" productId="">
            <status>Published</status>
            <productId>142</productId>
            <productName>WEST_Custom_Item_142</productName>
            <attributes>
                <attribute value="WEST_Custom_Item_142" name="productName"/>
                <attribute value="WEST_Custom_Item_142" name="Product Name"/>
                <attribute value="WEST_Custom_Item_142" name="Long Description"/>
                <attribute value="WEST_Custom_Item_142" name="Short Description"/>
                <attribute value="BFHPC_Neck" name="groupName"/>
                <attribute value="WEST_Custom_Item_142" name="entitlementTag"/>
                <attribute value="142" name="entitlementId"/>
                <attribute value="1" name="sortKey"/>
                <attribute value="1" name="duration" />
                <attribute value="MONTH" name="durationType"/>
            </attributes>
        </product>
        <product overwrite="false" default_locale="en_US" status="Published" productId="">
            <status>Published</status>
            <productId>142</productId>
            <productName>WEST_Custom_Item_142</productName>
            <attributes>
                <attribute value="WEST_Custom_Item_141" name="productName"/>
                <attribute value="WEST_Custom_Item_141" name="Product Name"/>
                <attribute value="WEST_Custom_Item_141" name="Long Description"/>
                <attribute value="WEST_Custom_Item_141" name="Short Description"/>
                <attribute value="BFHPC_Neck" name="groupName"/>
                <attribute value="WEST_Custom_Item_141" name="entitlementTag"/>
                <attribute value="141" name="entitlementId"/>
                <attribute value="1" name="sortKey"/>
                <attribute value="1" name="duration" />
                <attribute value="MONTH" name="durationType"/>
            </attributes>
        </product>
        <product overwrite="false" default_locale="en_US" status="Published" productId="">
            <status>Published</status>
            <productId>142</productId>
            <productName>WEST_Custom_Item_142</productName>
            <attributes>
                <attribute value="WEST_Custom_Item_140" name="productName"/>
                <attribute value="WEST_Custom_Item_140" name="Product Name"/>
                <attribute value="WEST_Custom_Item_140" name="Long Description"/>
                <attribute value="WEST_Custom_Item_140" name="Short Description"/>
                <attribute value="BFHPC_Neck" name="groupName"/>
                <attribute value="WEST_Custom_Item_140" name="entitlementTag"/>
                <attribute value="140" name="entitlementId"/>
                <attribute value="1" name="sortKey"/>
                <attribute value="1" name="duration" />
                <attribute value="MONTH" name="durationType"/>
            </attributes>
        </product>
    </products>
    ``` 
### Generating a BF2Random

A BF2Random of length `n` consists of `n` characters chosen randomly from the following string:
`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ][`
```
Theater Protocol
## Theater
The 3rd type of connection is the Theater connection and runs over TCP and UDP 
The ports can be found inside the Readme.txt inside the original files
A set of sockets is used by gameServer.exe and the gameClient.exe
Packets received or sent from the UDP port are decoded/encoded using the "gamespy XOR"

Copyright Disclaimer Under Section 107 of the Copyright Act 1976, allowance is made for "fair use" for purposes such as criticism, comment, news reporting, teaching, scholarship, and research. Fair use is a use permitted by copyright statute that might otherwise be infringing. Non-profit, educational or personal use tips the balance in favor of fair use.



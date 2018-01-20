# nextWave/backend (fesl)

`backend` is an implementation of the `GameSpy` network adapted for Battlefield Heroes use.

## Configuration

Below there is table with all enviroment variables which are used by the `nextWave/backend`. You can refer to `config/config.go` file if you need more information about specific variable.


| Name                  | Default value        |
|-----------------------|----------------------|
| `LOG_LEVEL`           | `INFO`               |
| `HTTP_BIND`           | `0.0.0.0:8080`       |
| `HTTPS_BIND`          | `0.0.0.0:443`        |
| `GAMESPY_IP`          | `0.0.0.0`(auto bind) |
| `FESL_CLIENT_PORT`    | `18270`              |
| `FESL_SERVER_PORT`    | `18051`              |
| `THEATER_CLIENT_PORT` | `18275`              |
| `THEATER_SERVER_PORT` | `18056`              |
| `THEATER_ADDR`        | `127.0.0.1`          |
| `TELEMETRICS_IP`      | `127.0.0.1`          |
| `TELEMETRICS_PORT`    | `13505`              |
| `LEVEL_DB_PATH`       | `_data/lvl.db`       |
| `DATABASE_USERNAME`   | `root`               |
| `DATABASE_PASSWORD`   |                      |
| `DATABASE_HOST`       | `127.0.0.1`          |
| `DATABASE_PORT`       | `3306`               |
| `DATABASE_NAME`       | `open-heroes`        |
| `CERT_PATH`           | `_fixtures/cert.pem` |
| `PRIVATE_KEY_PATH`    | `_fixtures/key.pem`  |

Note: It is okay to run on default configuration if you run server on your local PC for testing and development. But if you are thinking about exposing your server you probably should change some variables (i.e. `THEATER_ADDR`).

### Example `.env` file

```ini
DATABASE_NAME=open-heroes
DATABASE_HOST=192.168.33.10
DATABASE_PASSWORD=test
```

`openheroes/backend`  currently only uses `MySQL` as a backing services. If you are on platform where `docker` is available, you may use following command to quickly download and start container with a MySQL database:

```bash
sudo docker-compose start
```

### Start

===WINDOWS===
go to root folder and  -> ```go build main.go```

Note: You Must Set your GOPATH at Windows Environment

LINUX
```bash
make run```
Which is alias to:
```bash
go build -o main cmd/backend/main.go && sudo -H ./main`

## Credits ##
All The Idea/Project/Prototy Behind Bringing Back Battlefield Heroes was by #Synaxis
Credits to #MakaHost For being able to translate the code from BF2BC-emulator to golang
Credits to #Freeze-18, #Spencer and #mDawg From Revive Network.
Credits to #piotr and #Temp
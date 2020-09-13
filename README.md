 Go DMX controller with ArtNet discovery + Angular Web UI with WebSockets

## Web UI

![qmsk::dmx Web UI](https://raw.githubusercontent.com/qmsk/dmx/master/docs/web-main.png)

## Usage

#### Docker
The recommended way to build and run the project is using Docker:

    $ git clone https://github.com/qmsk/dmx.git qmsk-dmx && cd qmsk-dmx
    $ docker build -t qmsk/dmx .
    $ docker run --rm --name qmsk-dmx -v $PWD/config:/etc/qmsk-dmx:ro -e ARTNET_DISCOVERY=2.255.255.255 -p 8000:8000 qmsk/dmx

The `-v $PWD/config:/go/src/github.com/qmsk/dmx/config:ro` allows editing the config and reloading it use `docker restart qmsk-dmx`.

The `-e ARTNET_DISCOVERY=2.255.255.255` allows configuring a comma-separated list of broadcast/unicast addresses for ArtNet discovery.

The `-p 8000:8000` allows accessing the API/UI at `http://localhost:8000` on the machine running the Docker container.

### `github.com/qmsk/dmx/cmd/qmsk-dmx`

```

Usage:
  qmsk-dmx [OPTIONS] HeadsConfig

Application Options:
      --log=
      --demo                                  Demo Effect

ArtNet:
      --artnet-listen=ADDR
      --artnet-discovery=ADDR
      --artnet-discovery-interval=DURATION
      --artnet-discovery-timeout=DURATION
      --artnet-dmx-refresh=DURATION
      --log.artnet=

Heads:
      --log.heads=
      --heads-library=PATH

Web:
      --http-listen=[HOST]:PORT
      --http-static=PATH

Help Options:
  -h, --help                                  Show this help message
```

## `github.com/qmsk/dmx/artnet`

Go package supporting Art-Net discovery and DMX output.

* Uses a single `--artnet-listen=:6454` to send and receive UDP packets
* Supports multiple `--artnet-discovery=192.168.2.102` targets for unicast or broadcast use
* Supports dynamic [Device Discovery](http://art-net.org.uk/?page_id=454) of nodes and output ports using [`ArtPoll`](http://art-net.org.uk/?page_id=575) and [`ArtPollReply`](http://art-net.org.uk/?page_id=570) packets
* Supports outgoing [DMX Streaming](http://art-net.org.uk/?page_id=456) using [Broadcast/Unicast Subscription](http://art-net.org.uk/?page_id=649)
  * Sequence numbers for outgoing packets
  * Periodic output refresh

## `github.com/qmsk/dmx/heads`

DMX controller with support for multi-channel heads.
Configured using TOML configuration files, outputs DMX over Art-Net, controlled using a REST API and WebSocket event stream.

## Configuration
Can be configured using a single `toml` file, or a structured directory of configuration files:

#### `config/colors.toml`
```toml
[colors.red]
Red     = 1.0

[colors.green]
Green   = 1.0

[colors.blue]
Blue    = 1.0
```

#### `config/groups.toml`
```
[led-par]
Name        = "LED-Par"

[tri-bar]
Name        = "Tri-Bar"
```

#### `config/heads.toml`
```toml
[led-par]
Type        = "stairville/ledpar56-5ch"
Universe    = 1
Address     = 1
Count       = 6
Name        = "LED-Par"
Groups      = ["led-par"]

[tribar-1]
Type        = "american-dj/megatri60_mode2"
Universe    = 1
Address     = 20
Name        = "TriBar @ floor"
Groups      = ["tri-bar"]

[tribar-2]
Type        = "american-dj/megatri60_mode2"
Universe    = 1
Address     = 30
Name        = "TriBar @ wall"
Groups      = ["tri-bar"]
```

#### `config/presets/test.toml`
```toml
[Groups.led-par.Color]
Red = 1.0
Blue = 0.5

[Groups.tri-bar.Intensity]
Intensity = 1.0
[Groups.tri-bar.Color]
Red = 1.0
Blue = 0.5
```

## Web API

#### `GET /api/`
```json
{
   "Heads" : {
      "tribar-1" : {
         "Channels" : {
            "color:green" : {
               "Address" : 21,
               "DMX" : 0,
               "ID" : "color:green",
               "Value" : 0,
               "Type" : {
                  "Color" : "green"
               },
               "Index" : 1
            },
            "intensity" : {
               "ID" : "intensity",
               "DMX" : 0,
               "Address" : 24,
               "Index" : 4,
               "Type" : {
                  "Intensity" : true
               },
               "Value" : 0
            },
            "color:red" : {
               "ID" : "color:red",
               "DMX" : 0,
               "Address" : 20,
               "Type" : {
                  "Color" : "red"
               },
               "Index" : 0,
               "Value" : 0
            },
            "control:control" : {
               "Value" : 0,
               "Index" : 3,
               "Type" : {
                  "Control" : "control"
               },
               "Address" : 23,
               "DMX" : 0,
               "ID" : "control:control"
            },
            "color:blue" : {
               "Index" : 2,
               "Type" : {
                  "Color" : "blue"
               },
               "Value" : 0,
               "ID" : "color:blue",
               "DMX" : 0,
               "Address" : 22
            }
         },
         "Intensity" : {
            "ScaleIntensity" : null,
            "Intensity" : 0
         },
         "Config" : {
            "Name" : "TriBar @ floor",
            "Type" : "american-dj/megatri60_mode2",
            "Universe" : 1,
            "Address" : 20,
            "Groups" : [
               "tri-bar"
            ],
            "Count" : 0
         },
         "Type" : {
            "URL" : "",
            "Model" : "Mega Tri 60",
            "Mode" : "2",
            "Channels" : [
               {
                  "Color" : "red"
               },
               {
                  "Color" : "green"
               },
               {
                  "Color" : "blue"
               },
               {
                  "Control" : "control"
               },
               {
                  "Intensity" : true
               }
            ],
            "Vendor" : "American DJ",
            "Colors" : {
               "red" : {
                  "Green" : 0,
                  "Blue" : 0,
                  "Red" : 1
               },
               "magenta" : {
                  "Red" : 1,
                  "Blue" : 1,
                  "Green" : 0
               },
               "blue" : {
                  "Red" : 0,
                  "Green" : 0,
                  "Blue" : 1
               },
               "cyan" : {
                  "Blue" : 1,
                  "Green" : 1,
                  "Red" : 0
               },
               "amber" : {
                  "Blue" : 0,
                  "Green" : 0.5,
                  "Red" : 1
               },
               "green" : {
                  "Green" : 1,
                  "Blue" : 0,
                  "Red" : 0
               }
            }
         },
         "Color" : {
            "Green" : 0,
            "Blue" : 0,
            "Red" : 0,
            "ScaleIntensity" : null
         },
         "ID" : "tribar-1"
      },
   },
   "Outputs" : [
      {
         "Universe" : 1,
         "ArtNetNode" : {
            "OutputPorts" : [
               {
                  "Address" : {
                     "Net" : 0,
                     "SubUni" : 1
                  },
                  "Type" : 0,
                  "Status" : 128
               },
               {
                  "Address" : {
                     "Net" : 0,
                     "SubUni" : 2
                  },
                  "Status" : 128,
                  "Type" : 0
               }
            ],
            "Version" : 1,
            "Description" : "",
            "BaseAddress" : {
               "Net" : 0,
               "SubUni" : 0
            },
            "OEM" : 0,
            "Name" : "NodeMCU-ARTNET",
            "Report" : "",
            "InputPorts" : null,
            "Ethernet" : "00:00:00:00:00:00"
         }
      }
   ],
   "Groups" : {
      "tri-bar" : {
         "Color" : {
            "Red" : 0,
            "ScaleIntensity" : null,
            "Green" : 0,
            "Blue" : 0
         },
         "ID" : "tri-bar",
         "Name" : "Tri-Bar",
         "Intensity" : {
            "Intensity" : 0,
            "ScaleIntensity" : null
         },
         "Colors" : {
            "green" : {
               "Red" : 0,
               "Green" : 1,
               "Blue" : 0
            },
            "amber" : {
               "Green" : 0.5,
               "Blue" : 0,
               "Red" : 1
            },
            "cyan" : {
               "Red" : 0,
               "Blue" : 1,
               "Green" : 1
            },
            "red" : {
               "Green" : 0,
               "Blue" : 0,
               "Red" : 1
            },
            "magenta" : {
               "Green" : 0,
               "Blue" : 1,
               "Red" : 1
            },
            "blue" : {
               "Red" : 0,
               "Blue" : 1,
               "Green" : 0
            }
         },
         "Heads" : [
            "tribar-1",
            "tribar-2"
         ]
      },
   },
   "Presets" : {
      "test" : {
         "ID" : "test",
         "Groups" : {
            "led-par" : {
               "Intensity" : null,
               "Color" : {
                  "Blue" : 0.5,
                  "Green" : 0,
                  "ScaleIntensity" : null,
                  "Red" : 1
               }
            },
            "tri-bar" : {
               "Intensity" : {
                  "Intensity" : 1,
                  "ScaleIntensity" : null
               },
               "Color" : {
                  "Red" : 1,
                  "ScaleIntensity" : null,
                  "Green" : 0,
                  "Blue" : 0.5
               }
            }
         },
         "Config" : {
            "Name" : "",
            "All" : null,
            "Heads" : null,
            "Groups" : {
               "led-par" : {
                  "Intensity" : null,
                  "Color" : {
                     "Green" : 0,
                     "Blue" : 0.5,
                     "Red" : 1,
                     "ScaleIntensity" : null
                  }
               },
               "tri-bar" : {
                  "Intensity" : {
                     "Intensity" : 1,
                     "ScaleIntensity" : null
                  },
                  "Color" : {
                     "Green" : 0,
                     "Blue" : 0.5,
                     "Red" : 1,
                     "ScaleIntensity" : null
                  }
               }
            }
         },
         "Heads" : {}
      }
   }
}
```

#### `POST /api/heads/tribar-1`
```json
{ "Color": { "Red": 0.517, "Green": 0.0, "Blue": 0.0}}
```

#### `POST /api/groups/dimmer`
```json
{ "Intensity": 0.21 }
```

#### `POST /api/presets/test`
```json
{ "Intensity": 0.69 }
```

## Concepts
#### Heads

Top-level object that binds together everything else.

#### Output

A DMX universe (512 channels). Refreshed on every update. Connected to something like a `github.com/qmsk/dmx/artnet` `Universe` output.

#### Channel

Each ***Head** has multiple ***Channels***, which are patched to some ***Output*** universe.

###

#### Head

#### Group

A number of ***Heads***.

#### Preset

### Parameters

Both ***Heads***, ***Groups*** and ***Presets*** use **Parameters**.

#### Intensity

A single `Intensity` value.

#### Color

A combination of `Red`, `Green` and `Blue` values.

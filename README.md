# Pacific
Pacific turns on collecting urls from different providers . These providers are:

* urlscan.io
* Common Crawl
* Alienvault OTX
* Wayback Machine
* GrayHatWarfare    [Api key required]
* Hybrid Analysis   [Api key required]

## Usage

`$ pacific -d example.com -o output.txt`


## Configuration
You can use api keys by editing  ~/.config/pacific/config.yaml

Example :
```yaml
grayhatwarfare:
  - b5823caee729973eb55de59fc9a54f89
hybridanalysis:
  - tooy0y114b5823caeepd59oy03eb55de53ea54f89fvu1bbgpd5be1e3e
```

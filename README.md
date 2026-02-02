# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a conf file to define the mapping of the osm tags to tiles.

In addition, you can add SRTM files to the program, to be able to handle different tiles for different altitudes (not supported by OSM only).

## Usage

```bash
./osm2tmx -conf <conf.yaml> [-out <my.osm.tmx>] [-options...] <my.osm.pbf>

- conf: configuration file for tileset, see below
- out: default to my.osm.tmx
```

### Examples

- Simple

```bash
./osm2tmx -out example/example01/centre_les_gets.osm.tmx data/centre_les_gets.osm.pbf
```

- With downscale 4 and drawing

```bash
./osm2tmx -downscale 4 -draw -out example/example01/centre_les_gets.osm.tmx data/centre_les_gets.osm.pbf
```

- With downscale 4, srtm file and drawing

```bash
./osm2tmx -downscale 10 -srtm-tif data/N46E006.tif -draw -out example/example01/les_gets.osm.tmx data/les_gets.osm.pbf
```

### Get OSM files

You have multiple options to get the OSM files you want, to generate the desired tmx map:

- Go to https://www.openstreetmap.org , search for your location, adapt your zoom level, and then export the data (supported only to a certain zoom level, because the data cannot be too large). The data will be a `.osm` file, you must convert it to `.osm.pbf` (protobuf) file, with `osmium` for instance (see below), so it can be read by this program.
- If you want larger maps, you can visit https://download.geofabrik.de/ to get full planet, countries, regions or cities within the world. The data will already be in `.osm.pbf` format, so you won't need any conversion.

### Convert OSM files

You can use `osmium`: https://wiki.openstreetmap.org/wiki/Osmium , or `osmosis`: https://wiki.openstreetmap.org/wiki/Osmosis#Usage (deprecated in favor of osmium)

- Osmium usage: https://osmcode.org/osmium-tool/manual.html

- Osmosis usage:

```bash
# convert .osm to .osm.pbf files
osmosis --read-xml my_region.osm --write-pbf my_region.osm.pbf
# convert .osm.pbf to .osm files
osmosis --read-pbf my_region.osm.pbf --write-osm my_region.osm
```

### Altitudes handling (SRTM files)

This program is able to parse SRTM .tif files, to be able to detect the altitude for each lat,lon included in OSM files (evelation is not included in OSM files, except some points, see https://wiki.openstreetmap.org/wiki/Altitude).

Expected SRTM files format: `[N|S]xx[W|E]yyy.tif`

For example, `N26W080.tif` gives the altitude for all the (lat,lon) couples between (26.0,-80.0) and (26.9999999,-80.9999999).

And `S56E158.tif` handles (lat,lon) from (-56.0,1580.0) to (-56.9999999,1580.9999999).

To get your SRTM files, we can go for instance visit: https://portal.opentopography.org/datasetMetadata?otCollectionID=OT.042013.4326.1 (you'll need a free account to download the data).

### Conf file format: YAML

- List of all possible tags here: https://wiki.openstreetmap.org/wiki/Map_features

Example:

```yaml
# TODO: replace with real example
tileset:
  source: "tileset/basechip_pipo.png"
  tile_width: 16
  tile_height: 16

layers: 3

tags:
  default:
    0:
      tile: 2
  building:
    1:
      tile: 38
      values:
        roof:
          2:
            tile: 21
  highway:
    1:
      tile: 42
```

## TODO

- Conf & handle the mapping.yaml file
- More examples with other tilesets
- Generate a tileset like OSM stylesheets to have a pretty map in tmx :)
- Optimisation
  - See simd/archisimd in go1.26

- Explore a direct integration in 2d video games, without tilesets, using directional svg ?

## Perfs

```bash
./osm2tmx -downscale 4 -out example/example01/les_gets_div4.osm.tmx data/les_gets.osm.pbf
Number of CPUs: 12
Number of workers: 11
2026/02/01 15:27:21 will write output to example/example01/les_gets_div4.osm.tmx
2026/02/01 15:27:22 osm.Bounds{MinLat:46.144219999, MaxLat:46.167969999, MinLon:6.64407, MaxLon:6.69458}
2026/02/01 15:27:22 Max: UTM: [east:745237.240000,north:5807307.510000]
2026/02/01 15:27:22 Min: UTM: [east:739614.480000,north:5803490.770000]
2026/02/01 15:27:22 Map size: (1405,954) meters (4x)
2026/02/01 15:27:22 Altitude: 0 -> 0
2026/02/01 15:27:22 Nodes: 40081
2026/02/01 15:27:22 Ways: 5468
2026/02/01 15:27:22 Relations: 86
2026/02/01 15:27:22 Generated map: height: 954, width: 1405
2026/02/01 15:27:22 Number of points out of bounds: 9047
```

## Resources

- https://en.wikipedia.org/wiki/Geographic_coordinate_system

- OpenStreetMap:
  - https://planet.openstreetmap.org/
  - https://download.geofabrik.de/
  - https://wiki.openstreetmap.org/wiki/Map_features
  - https://wiki.openstreetmap.org/wiki/Osmosis#Usage
  - https://github.com/maguro/pbf/blob/master/model/elements.go
  - https://wiki.openstreetmap.org/wiki/PBF_Format
  - https://wiki.openstreetmap.org/wiki/Altitude

- UI:
  - https://github.com/hajimehoshi/ebiten/blob/main/examples/vector/main.go
  - https://github.com/systemed/tilemaker/tree/master
  - https://wiki.openstreetmap.org/wiki/Vector_tiles

- Inspirations:
  - https://github.com/mapnik/mapnik/wiki/XMLConfigReference
  - https://github.com/mapnik/mapnik/wiki/OsmPlugin
  - https://github.com/maplibre/maplibre-gl-js?tab=readme-ov-file

- 3D
  - https://maplibre.org/maplibre-gl-js/docs/examples/animate-map-camera-around-a-point/

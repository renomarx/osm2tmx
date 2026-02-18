# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a conf file to define the mapping of the osm tags to tiles.

In addition, you can add SRTM files to the program, to be able to handle different tiles for different altitudes (not supported by OSM only).

## Usage

```bash
./osm2tmx -mapping <mapping.yaml> [-out <my.osm.tmx>] [-options...] <my.osm.pbf>

- mapping: mapping file for tileset, see below
- out: default to my.osm.tmx
```

### Examples

- Simple

```bash
./osm2tmx -mapping example/example01/mapping.yaml -out example/example01/centre_les_gets.osm.tmx example/centre_les_gets.osm.pbf
```

- With downscale 4 and drawing

```bash
./osm2tmx -downscale 4 -draw -mapping example/example01/mapping.yaml -out example/example01/centre_les_gets.osm.tmx example/centre_les_gets.osm.pbf
```

- With downscale 10, srtm file and drawing

```bash
./osm2tmx -downscale 10 -srtm-tif example/N46E006.tif -draw -mapping example/example01/mapping.yaml -out example/example01/les_gets.osm.tmx example/les_gets.osm.pbf
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

### Mapping file format: YAML

- List of all possible tags here: https://wiki.openstreetmap.org/wiki/Map_features
- See the [go model file](pkg/mapper/model.go) to have a full description of supported features

Examples:

- [example01](example/example01/mapping.yaml)

## TODO

- More examples with other tilesets

- Explore a direct integration in 2d video games, without tilesets, using directional svg ?

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

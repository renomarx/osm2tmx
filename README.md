# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a conf file to define the mapping of the osm tags to tiles.

## Usage

```bash
./osm2tmx -conf <my_mapping_file.yaml> [-out <my.osm.tmx>] <my.osm.pbf>

- conf: configuration file for tileset, see below
- out: default to my.osm.tmx
```

## Examples

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
./osm2tmx -downscale 4 -srtm-tif data/N46E006.tif -draw -out example/example01/centre_les_gets.osm.tmx data/centre_les_gets.osm.pbf
```

### Conf file format: YAML

Example:

```yaml
tileset:
  source: "tileset/basechip_pipo.png"
  tile_width: 16
  tile_height: 16

layers: 3

tags:
  default:
    layer: 0
    tile: 2
  building:
    layer: 1
    tile: 38
    # values:
    #   roof:
    #     layer: 2
    #     tile: 21
  highway:
    layer: 1
    tile: 42
```

## TODO

- Handle  altitudes with SRTM data files
- Trace "primary" roads with 5 meters diameter
- Conf & handle the mapping.yaml file
- More examples with other tilesets
- Generate a tileset like OSM stylesheets to have a pretty map in tmx :)
- Optimisation
  - See simd/archisimd in go1.26

- Explore a direct integration in 2d video games, without tilesets, using directional svg ?

## Perfs

```bash
./osm2tmx -downscale 4 -out example/example01/les_gets_div4.osm.tmx data/les_gets.osm.pbf 
Number of CPUs: 8
Number of workers: 7
2026/01/26 18:42:26 main.go:74: will write output to example/example01/les_gets_div4.osm.tmx
2026/01/26 18:42:35 main.go:91: osm.Bounds{MinLat:46.144219999, MaxLat:46.167969999, MinLon:6.64407, MaxLon:6.69458}
2026/01/26 18:42:35 main.go:92: Max: UTM: [east:745237.240000,north:5807307.510000]
2026/01/26 18:42:35 main.go:93: Min: UTM: [east:739614.480000,north:5803490.770000]
2026/01/26 18:42:35 main.go:94: Map size: (1405,954) meters (4x)
2026/01/26 18:42:35 main.go:96: Nodes: 40085
2026/01/26 18:42:35 main.go:97: Ways: 5468
2026/01/26 18:42:35 main.go:98: Relations: 86
2026/01/26 18:42:35 main.go:100: Generated map: height: 954, width: 1405
2026/01/26 18:42:35 main.go:102: Number of points out of bounds: 9043
```

## Resources

Les valeurs de latitude sont mesurées par rapport à l'équateur et à une plage comprise entre -90° au pôle Sud et +90° au pôle Nord. Les valeurs de longitude sont mesurées par rapport au premier méridien. Elles sont comprises entre -180° en allant vers l'ouest et 180° vers l'est.

- OpenStreetMap specifications:
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

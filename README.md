# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a conf file to define the mapping of the osm tags to tiles.

## Usage

```bash
./osm2tmx -conf <my_mapping_file.yaml> [-out <my.osm.tmx>] <my.osm.pbf>

- conf: configuration file for tileset, see below
- out: default to my.osm.tmx
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

- Draw way corners, area corners and buildings with different tiles:
  - evenodd.IsInsidePolygon return a struct{Left,Right,Top,Bottom} representing the distance between the point and the nearest boundary
  - then the mapper can use these values to select the right tile
- Add arguments and options to handle bounded exports

- POC drawer 2d (ebiten, draw2d ?)

- How to handle z-index and level lines ?
  - Not included in OSM data, we'll have to find another way
- Conf & handle the mapping.yaml file
- Optimisation
- More examples with other tilesets
- Generate a tileset like OSM stylesheets to have a pretty map in tmx :)

- Explore a direct integration in 2d video games, without tilesets, using directional svg ?

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

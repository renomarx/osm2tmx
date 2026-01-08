# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a tileset mapping file.

## Usage

```bash
# Ile de la reunion
./osm2tmx --mapping=<my_mapping_file.yaml> <my.osm.pbf> [--out=<my.osm.tmx>]

- mapping: mapping file of osm tags <-> tileset pos, see below
- out: default to my.osm.tmx
```

### Tags mapping file format: YAML

Example:

```yaml
tileset: "tileset/basechip_pipo.png"

layers: 3

tags:
  default:
    layer: 0
    x: 0
    y: 0
  building:
    layer: 1
    x: 0
    y: 493
    # values:
    #   roof:
    #     layer: 2
    #     x: 0
    #     y: 232
  highway:
    layer: 1
    x: 120
    y: 0
```

## TODO

- Refacto & tests
- Handle relations, fillign areas of relation
- Handle all relevant tags
- Add arguments and options to handle bounded exports

- POC drawer 2d (ebiten, draw2d ?)

- How to handle multiple tags ?
- How to handle z-index and level lines ?
- Handle the mapping.yaml file
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

- UI:
    - https://github.com/hajimehoshi/ebiten/blob/main/examples/vector/main.go
    - https://github.com/systemed/tilemaker/tree/master
    - https://wiki.openstreetmap.org/wiki/Vector_tiles

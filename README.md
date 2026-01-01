# osm2tmx

Util to convert osm.pbf files to tmx files, using a tileset and a tileset mapping file.

## Usage

```bash
# Ile de la reunion
./osm2tmx --tileset=<my_tileset.tsx> --tags=<my_mapping_file.csv> --tl_lat=-21.129278 --tl_long=55.4150827 [--out=<my.osm.tmx>] <my.osm.pbf>

- tr_lat, tr_long: top-left latitude, longitude of the map. The corresponding position will be 0,0
- tags: mapping file of osm tags <-> tmx pos, see below
- out: default to my.osm.tmx
```

### Tags mapping file format: CSV

```csv
key,value,<pos_in_tmx_file>
key,value,<pos2_in_tmx_file>
```

Example:

```csv
building,pagoda,99
building,roof,1442
```

## TODO

- How to handle multiple tags ?
- OSM file is a polygon, tmx file is a square. How to handle the translation ? Should we auto-detect, taking the max lat and min lon and extrapolate to have the corresponding top_left corner of the square ?
- What default tile to give for positions without corresponding lat,long in the osm file ? Should we begin the tmx file, at pos 0, with conventionnaly the default tile (ex: water ?) ?
- What precision to give for the translation lat,lon->pos ? We'd need to calculate the tile size in meters to allow a realistic translation
- How to handle z-index and level lines ?

## Resources

Les valeurs de latitude sont mesurées par rapport à l'équateur et à une plage comprise entre -90° au pôle Sud et +90° au pôle Nord. Les valeurs de longitude sont mesurées par rapport au premier méridien. Elles sont comprises entre -180° en allant vers l'ouest et 180° vers l'est.

- OpenStreetMap specifications:
    - https://planet.openstreetmap.org/
    - https://wiki.openstreetmap.org/wiki/Map_features
    - https://wiki.openstreetmap.org/wiki/Osmosis#Usage
    - https://github.com/maguro/pbf/blob/master/model/elements.go
    - https://wiki.openstreetmap.org/wiki/PBF_Format

- UI:
    - https://github.com/hajimehoshi/ebiten/blob/main/examples/vector/main.go
    - https://github.com/systemed/tilemaker/tree/master
    - https://wiki.openstreetmap.org/wiki/Vector_tiles

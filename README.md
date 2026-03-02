# osm2tmx

A command-line utility to convert OpenStreetMap PBF files to Tiled TMX map files. The tool uses a configurable YAML mapping file to define how OSM tags are translated to tileset tiles, enabling custom game map generation from real-world geographic data.

## Key Features

- **OSM to TMX Conversion**: Convert OpenStreetMap data (`.osm.pbf` format) to Tiled Map Editor format (`.tmx`)
- **Flexible Mapping System**: YAML-based configuration for mapping OSM tags to tileset tiles
- **Multi-layer Support**: Generate maps with multiple layers for complex tile arrangements
- **Altitude/Topography Support**: Integrate SRTM elevation data to generate different tiles based on altitude
- **Custom Tile Shapes**: Support for custom tile patterns (corners, borders, lines, rectangles) based on tile positions
- **Polygon Filling**: Automatic filling of areas defined by ways and relations
- **Downscaling**: Reduce map size by applying a downscale factor
- **Built-in Viewer**: Optional rendering of the generated map using Ebiten game engine
- **Parallel Processing**: Multi-threaded processing for improved performance
- **Automatic Layer Management**: Automatically adds layers when rectangles overlap to maintain proper rendering order

## Installation

```bash
go build -o osm2tmx
```

## Usage

```bash
./osm2tmx -mapping <mapping.yaml> [options] <input.osm.pbf>
```

### Command-Line Options

| Option | Type | Description |
|--------|------|-------------|
| `-mapping` | string | Path to YAML mapping configuration file (required) |
| `-out` | string | Output TMX file path (default: replaces `.osm.pbf` with `.tmx`) |
| `-downscale` | int | Downscale factor to reduce map size (e.g., 10 = 1/10th size) (default: 1) |
| `-offset-x` | int | X-axis offset in meters, applied after downscaling (default: 0) |
| `-offset-y` | int | Y-axis offset in meters, applied after downscaling (default: 0) |
| `-limit-x` | int | Maximum width in tiles, applied after downscaling (default: unlimited) |
| `-limit-y` | int | Maximum height in tiles, applied after downscaling (default: unlimited) |
| `-workers` | int | Number of parallel workers (default: CPU count - 1) |
| `-draw` | bool | Display the generated map in a game UI window (default: false) |
| `-srtm-tif` | string | Path to SRTM .tif file (can be specified multiple times) |
| `-srtm-dir` | string | Directory containing SRTM .tif files (recursive search) |
| `-help` | bool | Display usage information |

## Examples

### Basic Conversion

```bash
./osm2tmx -mapping example/example01/mapping.yaml \
  -out example/example01/centre_les_gets.osm.tmx \
  example/centre_les_gets.osm.pbf
```

### With Downscaling and Viewer

```bash
./osm2tmx -downscale 4 -draw \
  -mapping example/example01/mapping.yaml \
  -out example/example01/centre_les_gets.osm.tmx \
  example/centre_les_gets.osm.pbf
```

### With Elevation Data (SRTM)

```bash
./osm2tmx -downscale 4 \
  -srtm-tif example/N46E006.tif \
  -draw \
  -mapping example/example01/mapping.yaml \
  -out example/example01/les_gets.osm.tmx \
  example/les_gets.osm.pbf
```

### Using Multiple SRTM Files

```bash
./osm2tmx -mapping mapping.yaml \
  -srtm-tif data/N45E006.tif \
  -srtm-tif data/N46E006.tif \
  -srtm-tif data/N46E007.tif \
  input.osm.pbf
```

### Using SRTM Directory

```bash
./osm2tmx -mapping mapping.yaml \
  -srtm-dir data/srtm/ \
  input.osm.pbf
```

## Working with OpenStreetMap Data

### Obtaining OSM Files

#### Small Areas (Direct Export)

1. Visit [openstreetmap.org](https://www.openstreetmap.org)
2. Navigate to your desired location
3. Zoom to the appropriate level
4. Click "Export" to download `.osm` XML file
5. Convert to `.osm.pbf` format (see below)

**Note**: Direct exports are limited in size. For larger areas, use Geofabrik extracts.

#### Large Areas (Geofabrik Extracts)

Download pre-made extracts in `.osm.pbf` format from [Geofabrik](https://download.geofabrik.de/):

- Full planet file
- Continent extracts
- Country extracts
- Regional/city extracts

These files are already in PBF format and ready to use.

### Converting OSM Formats

#### Using Osmium (Recommended)

[Osmium Tool](https://osmcode.org/osmium-tool/) is the modern, fast way to work with OSM data.

```bash
# Install
sudo apt-get install osmium-tool  # Ubuntu/Debian
brew install osmium-tool          # macOS

# Convert .osm to .osm.pbf
osmium cat my_region.osm -o my_region.osm.pbf

# Convert .osm.pbf to .osm
osmium cat my_region.osm.pbf -o my_region.osm

# Extract a bounding box
osmium extract --bbox 6.0,46.0,7.0,47.0 source.osm.pbf -o area.osm.pbf
```

#### Using Osmosis (Legacy)

[Osmosis](https://wiki.openstreetmap.org/wiki/Osmosis) is deprecated but still functional:

```bash
# Convert .osm to .osm.pbf
osmosis --read-xml my_region.osm --write-pbf my_region.osm.pbf

# Convert .osm.pbf to .osm
osmosis --read-pbf my_region.osm.pbf --write-xml my_region.osm
```

## Elevation Data (SRTM)

This tool can parse SRTM (Shuttle Radar Topography Mission) `.tif` files to add elevation data to map generation. Since OSM files generally don't include elevation data (see [OSM Altitude wiki](https://wiki.openstreetmap.org/wiki/Altitude)), SRTM files enable altitude-based tile selection.

### SRTM File Format

Expected filename format: `[N|S]dd[E|W]ddd.tif`

Examples:
- `N46E006.tif` - Covers coordinates from (46.0°N, 6.0°E) to (46.999°N, 6.999°E)
- `S26W080.tif` - Covers coordinates from (26.0°S, 80.0°W) to (26.999°S, 80.999°W)

### Obtaining SRTM Data

Download SRTM elevation data from:
- [OpenTopography](https://portal.opentopography.org/datasetMetadata?otCollectionID=OT.042013.4326.1) (free account required)
- Other SRTM data providers

### How It Works

1. The tool automatically loads SRTM tiles covering your OSM data bounds
2. For each map position, it looks up the elevation from SRTM data
3. The mapping configuration can specify different tiles based on altitude thresholds
4. Elevation precision is automatically adjusted based on the downscale factor

## Mapping Configuration

The mapping file is a YAML configuration that defines how OpenStreetMap tags are converted to tileset tiles.

### Configuration Structure

```yaml
# Tileset definitions
tilesets:
  - source: path/to/tileset.tsx
    tile_width: 32
    tile_height: 32
    first_gid: 1

# Default tile for unmapped areas
default:
  tile: 1

# Layer-based tag mappings
layers:
  0:  # Base layer (terrain)
    tags:
      natural:
        values:
          water: { tile: 10 }
          grass: { tile: 5 }
  1:  # Feature layer (buildings, roads)
    tags:
      highway:
        values:
          primary: { tile: 20 }
          secondary: { tile: 21 }

# Custom tile patterns (optional)
custom_tiles:
  20:  # Tile ID to customize
    shapes:
      square9: [[1,2,3], [4,5,6], [7,8,9]]
```

### Mapping Features

See [pkg/mapper/mapping.go](pkg/mapper/mapping.go) for complete field documentation.

#### Basic Tile Mapping

Map OSM tags directly to tile IDs:

```yaml
layers:
  0:
    tags:
      landuse:
        values:
          forest: { tile: 15 }
          residential: { tile: 20 }
```

#### Altitude-Based Mapping

Different tiles based on elevation (requires SRTM data):

```yaml
layers:
  0:
    tags:
      natural:
        values:
          terrain:
            tile: 5  # Default tile
            altitude:
              min: 1500  # Meters
              tile: 8    # Tile for altitude >= 1500m
```

#### Random Tile Selection

Add variety with probabilistic tile selection:

```yaml
layers:
  0:
    tags:
      natural:
        values:
          grass:
            random:
              - probability: 70  # 70% chance
                tile: 5
              - probability: 20  # 20% chance
                tile: 6
              - probability: 10  # 10% chance
                tile: 7
```

#### Custom Tile Shapes

Define tiles that change based on adjacent tiles (auto-tiling):

```yaml
custom_tiles:
  10:  # Road tile
    shapes:
      line: [101, 102, 103]  # [start, middle, end] tiles
      column: [111, 112, 113]  # [top, middle, bottom]
      square9:  # 3x3 pattern for corners and edges
        - [201, 202, 203]  # top-left, top, top-right
        - [204, 205, 206]  # left, center, right
        - [207, 208, 209]  # bottom-left, bottom, bottom-right
      square4:  # 2x2 pattern for diagonal corners
        - [301, 302]  # top-left, top-right
        - [303, 304]  # bottom-left, bottom-right
      standalone: 999  # isolated tile
```

#### Rectangle Objects

Place multi-tile objects within polygons:

```yaml
custom_tiles:
  50:  # Polygon tile (e.g., forest)
    rectangle:
      tiles:
        - [0, 81, 82, 83, 0]  # Tree object rows (0 = transparent)
        - [0, 84, 85, 86, 0]
        - [0, 87, 88, 89, 0]
      inside_polygon:
        density: 2  # Overlap factor (0=max density, higher=more spacing)
        overflow: "ORTHOGONAL"  # ALWAYS, ORTHOGONAL, QUARTER, HALF, or ""
```

**Overflow modes** for rectangle placement:
- `""` (empty): Only draw if rectangle fully inside polygon
- `ORTHOGONAL`: Draw if the last column and last row are inside polygon
- `QUARTER`: Draw if at least 25% of rectangle is inside
- `HALF`: Draw if at least 50% of rectangle is inside
- `ALWAYS`: Always draw regardless of polygon boundaries

#### Random Custom Tiles

Combine randomization with custom tiles:

```yaml
custom_tiles:
  50:
    random:
      - probability: 60
        rectangle:
          tiles: [[71, 72], [73, 74]]
      - probability: 40
        rectangle:
          tiles: [[75, 76], [77, 78]]
```

### Example Configurations

- [example01](example/example01/mapping.yaml) - Basic mapping example
- [example02](example/example02/mapping.yaml) - Advanced features

## Project Architecture

### Package Overview

```
pkg/
├── bresenham/      Line drawing algorithm for connecting OSM way nodes
├── draw/           Built-in map viewer using Ebiten game engine
├── evenodd/        Point-in-polygon testing for area filling
├── floodfill/      Flood fill algorithm for polygon areas
├── mapper/         Tag-to-tile mapping logic and configuration
├── mercator/       Mercator projection coordinate conversions
├── model/          Core data structures (Map, Layer, Cell, Tile, etc.)
├── raster/         Main OSM parsing and rasterization engine
├── tmx/            TMX file format encoding/decoding
└── topography/
    └── srtm/       SRTM elevation data parsing
```

### Processing Pipeline

1. **Parse OSM PBF File** (`pkg/raster`)
   - Read nodes, ways, and relations from the input file
   - Convert lat/lon coordinates to map X,Y positions using Mercator projection
   - Calculate map bounds and initialize layers

2. **Map Nodes** (`pkg/mapper`)
   - For each OSM node, look up tags in mapping configuration
   - Apply altitude-based rules if SRTM data available
   - Set appropriate tile on each layer

3. **Process Ways** (parallel)
   - Draw lines for linear features (roads, paths, etc.) using Bresenham algorithm
   - Fill polygons for areas (buildings, water bodies, forests)
   - Use even-odd rule for point-in-polygon testing

4. **Process Relations** (parallel)
   - Handle multipolygon relations (complex areas)
   - Merge member ways into complete polygons
   - Fill with appropriate tiles

5. **Apply Custom Tiles**
   - Re-scan the map to detect tile patterns (corners, edges, etc.)
   - Replace tiles with custom variants based on adjacent tiles
   - Place rectangle objects within polygons
   - Automatically add layers when rectangles overlap

6. **Generate TMX Output** (`pkg/tmx`)
   - Convert internal map structure to TMX XML format
   - Write CSV-encoded tile data for each layer
   - Reference external tileset files

### Key Algorithms

- **Mercator Projection**: Convert geographic coordinates (lat/lon) to planar coordinates (x/y)
- **Bresenham Line Algorithm**: Draw lines between way nodes without gaps
- **Even-Odd Rule**: Determine if a point is inside a polygon for area filling
- **Auto-tiling**: Detect tile patterns (corners, borders) and apply custom tiles
- **Overlap Resolution**: Automatically create new layers when tiles overlap

## Built-in Viewer

When using the `-draw` flag, the tool opens an interactive viewer window:

### Controls

- **Arrow Keys**: Pan the camera (Up, Down, Left, Right)
- **Ctrl + Up**: Zoom in
- **Ctrl + Down**: Zoom out

The viewer automatically centers on the map and displays the current zoom level. Large maps are automatically cropped to 16000x16000 pixels for performance.

## Use Cases

- **Game Development**: Generate realistic game maps from real-world geography
- **Educational Tools**: Create interactive maps for geography or urban planning education
- **Prototyping**: Quickly generate test maps for game engines or mapping applications
- **Visualization**: Convert geographic data to tile-based visualizations
- **Custom Cartography**: Create stylized maps with custom tilesets

## Performance Considerations

- **Downscaling**: Use `-downscale` to reduce processing time and output size for large areas
- **Workers**: Adjust `-workers` based on CPU cores (default is CPU count - 1)
- **Bounds**: Use `-offset-x/y` and `-limit-x/y` to process only a subset of the data
- **SRTM Loading**: SRTM files are loaded on-demand per coordinate; preloading optimizes based on area coverage

## Limitations

- Maximum drawable image size: 16000x16000 pixels (automatically cropped in viewer)
- TMX format supports up to 2^31 tiles per layer
- SRTM precision decreases with higher downscale values
- Memory usage scales with map size and number of OSM elements

## Resources

### Geographic & Coordinate Systems

- [Geographic Coordinate System](https://en.wikipedia.org/wiki/Geographic_coordinate_system)
- [Mercator Projection](https://en.wikipedia.org/wiki/Mercator_projection)

### OpenStreetMap

- [OSM Main Site](https://www.openstreetmap.org/)
- [OSM Planet Dumps](https://planet.openstreetmap.org/)
- [Geofabrik Downloads](https://download.geofabrik.de/)
- [OSM Map Features (Tag List)](https://wiki.openstreetmap.org/wiki/Map_features)
- [OSM PBF Format](https://wiki.openstreetmap.org/wiki/PBF_Format)
- [OSM Altitude Data](https://wiki.openstreetmap.org/wiki/Altitude)
- [Osmium Tool](https://wiki.openstreetmap.org/wiki/Osmium)
- [Osmosis Tool](https://wiki.openstreetmap.org/wiki/Osmosis)

### Tiled Map Editor

- [Tiled Map Editor](https://www.mapeditor.org/)
- [TMX Format Documentation](https://doc.mapeditor.org/en/stable/reference/tmx-map-format/)

### Libraries & Tools

- [Ebiten Game Engine](https://ebiten.org/) (used for viewer)
- [go-tiled](https://github.com/lafriks/go-tiled) (TMX parsing)
- [paulmach/osm](https://github.com/paulmach/osm) (OSM PBF parsing)

### Inspiration

- [Mapnik](https://github.com/mapnik/mapnik) - Map rendering toolkit
- [Tilemaker](https://github.com/systemed/tilemaker) - Vector tile generator
- [MapLibre GL JS](https://github.com/maplibre/maplibre-gl-js) - Interactive map renderer

## Troubleshooting

### Common Issues

**Issue**: "ERROR: no mapping file in parameter"
- **Solution**: Ensure you specify the `-mapping` flag with a valid YAML file path

**Issue**: Output TMX appears empty or has only default tiles
- **Solution**: Check that your mapping configuration matches the OSM tags in your data
- **Tip**: Use [taginfo](https://taginfo.openstreetmap.org/) to explore available tags in OSM data

**Issue**: "tif not found" errors with SRTM files
- **Solution**: Verify SRTM filenames follow the `[N|S]dd[E|W]ddd.tif` format
- **Solution**: Ensure SRTM files cover the geographic bounds of your OSM data

**Issue**: Tiles appear incorrectly positioned or sized
- **Solution**: Verify your tileset `tile_width` and `tile_height` match your actual tileset
- **Solution**: Check that `first_gid` is correctly set (usually 1 for the first tileset)

**Issue**: Custom tiles not appearing
- **Solution**: Ensure the tile ID in `custom_tiles` matches the base tile ID from your mapping
- **Solution**: Verify adjacent tiles meet the pattern requirements (e.g., corners need specific neighbors)

**Issue**: Map size too large or out of memory
- **Solution**: Increase `-downscale` factor to reduce map resolution
- **Solution**: Use `-limit-x` and `-limit-y` to process only a portion of the data
- **Solution**: Use `-offset-x` and `-offset-y` to shift the processing area

**Issue**: Viewer window shows cropped image
- **Solution**: This is expected for maps larger than 16000x16000 pixels; reduce `-downscale` or use bounds flags

### Debug Tips

1. Start with a small area to test your mapping configuration
2. Use `-draw` to visually inspect the output immediately
3. Examine the generated TMX file in Tiled Map Editor to debug layer issues
4. Check the console output for statistics (nodes, ways, relations, bounds)
5. Verify OSM tag names match exactly (case-sensitive)

## Contributing

Contributions are welcome! Areas for improvement:

- Additional auto-tiling patterns
- Performance optimizations
- Support for more TMX features (object layers, image layers)
- Extended SRTM data sources
- Better error messages and validation
- Documentation and examples

Please open issues or pull requests on the project repository.

## License

This project structure and code organization follows standard Go practices. See project files for specific licensing information.

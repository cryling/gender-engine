# gender-engine

This project provides an easy-to-use API for the [World Gender Name Dictionary (WGND 2.0)](https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/MSEGSJ).

## Getting Started

### Docker Image

If you want to use the pre-built Docker image, you can pull it from the GitHub Container Registry.

```sh
docker pull ghcr.io/cryling/gender-engine:latest
```

Then run it with:

```sh
docker run -it --rm -p 8080:8080 ghcr.io/cryling/gender-engine
```

### Building from Source

#### Prerequisites

- Docker

#### Installation

1. Build the Docker image. The dataset is downloaded automatically during the build.

   ```sh
   docker build --tag gender-engine .
   ```

2. Run the Docker container.

   ```sh
   docker run -it --rm -p 8080:8080 gender-engine
   ```

## Configuration

When creating a container from the image, you can modify the rate limit parameters. The default rate limit settings are:

```sh
ENV RATE_LIMIT_ENABLED=true
ENV RATE_LIMIT=50
ENV RATE_BURST=500
```

To disable rate limiting, set RATE_LIMIT_ENABLED to any other value than 'true' (e.g., 'false') by adding the -e flag to the docker run command. For example:

```sh
docker run -it --rm -p 8080:8080 -e RATE_LIMIT_ENABLED=false gender-engine
```

## API Documentation

The API consists of a single endpoint that returns the gender of a given name.

### Endpoint

- GET /gender
  - Query Parameters:
    - name: The name to be queried.
    - country: The country code for the name.

### Response

A typical response would look like this:

```sh
curl -X GET "http://localhost:8080/api/v1/gender?name=tom&country=US"
```

```json
{
  "message": "Tom could be found in US",
  "result": {
    "Name": "tom",
    "Gender": "M",
    "Country": "US",
    "Probability": "0.99560356"
  }
}
```

If the name could not be found, the response would look like this:

```json
{
  "message": "Tom could not be found"
}
```

Valid countries are two-letter country codes. You can find a list of valid country codes in the [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) Wikipedia article. Not all countries are available in the dataset. Check [api/domain/countries.go](https://github.com/cryling/gender-engine/blob/main/api/domain/countries.go) for a list of available countries.

## Dataset Acknowledgements

This project utilizes the World Gender Name Dictionary (WGND 2.0) dataset. The dataset is licensed under the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication. The dataset can be accessed and downloaded from the [Harvard Dataverse](https://dataverse.harvard.edu), DOI: 10.7910/DVN/MSEGSJ.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

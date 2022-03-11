# md5_hasher

![example workflow](https://github.com/Moxxx1e/md5_hasher/actions/workflows/workflow.yml/badge.svg)

A tool which makes http requests and prints the address of the request along with the MD5 hash of the response.

## Build

```
make build
```

## Usage

The -parallel option specifies the number of goroutines used by the program. Default value is 10.

```
./myhttp -parallel 3 google.com yandex.com twitter.com reddit.com/r/funny
```

## Example
![example workflow](https://github.com/Moxxx1e/md5_hasher/blob/main/img/example.png?raw=true)

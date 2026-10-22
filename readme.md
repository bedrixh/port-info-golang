# lightweight tool to get infromation about port

## Cli usage
`port 80`
will return information about the port 80

## How does that work?
 - Data are being fetched from IANA as csv file, cached in {cache directory}/port-info-golang/
 - To force reload port list remove the cache file and the data will refresh durring next run
